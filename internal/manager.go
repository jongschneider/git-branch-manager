package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Manager struct {
	config     *Config
	gitManager *GitManager
	envMapping *EnvMapping
	repoPath   string
	gbmDir     string
}

type WorktreeListInfo struct {
	Path           string
	ExpectedBranch string
	CurrentBranch  string
	GitStatus      *GitStatus
}

type SyncStatus struct {
	InSync            bool
	MissingWorktrees  []string
	OrphanedWorktrees []string
	BranchChanges     map[string]BranchChange
}

type BranchChange struct {
	OldBranch string
	NewBranch string
}

func NewManager(repoPath string) (*Manager, error) {
	gitManager, err := NewGitManager(repoPath)
	if err != nil {
		return nil, err
	}

	gbmDir := filepath.Join(repoPath, ".gbm")
	config, err := LoadConfig(gbmDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize the global icon manager with the loaded config
	iconManager := NewIconManager(config)
	SetGlobalIconManager(iconManager)

	return &Manager{
		config:     config,
		gitManager: gitManager,
		repoPath:   repoPath,
		gbmDir:     gbmDir,
	}, nil
}

func (m *Manager) LoadEnvMapping(envrcPath string) error {
	if envrcPath == "" {
		envrcPath = filepath.Join(m.repoPath, ".envrc")
	}

	if !filepath.IsAbs(envrcPath) {
		envrcPath = filepath.Join(m.repoPath, envrcPath)
	}

	mapping, err := ParseEnvrc(envrcPath)
	if err != nil {
		return err
	}

	m.envMapping = mapping
	return nil
}

func (m *Manager) GetSyncStatus() (*SyncStatus, error) {
	if m.envMapping == nil {
		return nil, fmt.Errorf("no .envrc mapping loaded")
	}

	status := &SyncStatus{
		InSync:            true,
		MissingWorktrees:  []string{},
		OrphanedWorktrees: []string{},
		BranchChanges:     make(map[string]BranchChange),
	}

	worktrees, err := m.gitManager.GetWorktrees()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	worktreeMap := make(map[string]*WorktreeInfo)
	for _, wt := range worktrees {
		if strings.HasPrefix(wt.Path, filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix)) {
			worktreeMap[wt.Name] = wt
		}
	}

	for envVar, expectedBranch := range m.envMapping.Variables {
		if wt, exists := worktreeMap[envVar]; exists {
			if wt.Branch != expectedBranch {
				status.BranchChanges[envVar] = BranchChange{
					OldBranch: wt.Branch,
					NewBranch: expectedBranch,
				}
				status.InSync = false
			}
			delete(worktreeMap, envVar)
		} else {
			status.MissingWorktrees = append(status.MissingWorktrees, envVar)
			status.InSync = false
		}
	}

	for envVar := range worktreeMap {
		status.OrphanedWorktrees = append(status.OrphanedWorktrees, envVar)
		status.InSync = false
	}

	return status, nil
}

func (m *Manager) Sync(dryRun, force, fetch bool) error {
	// if m.envMapping == nil {
	// 	return fmt.Errorf("no .envrc mapping loaded")
	// }

	// Validate all branches exist before performing any operations
	if err := m.ValidateEnvrc(); err != nil {
		return err
	}

	if fetch {
		if err := m.gitManager.FetchAll(); err != nil {
			return fmt.Errorf("failed to fetch: %w", err)
		}
	}

	status, err := m.GetSyncStatus()
	if err != nil {
		return err
	}

	if status.InSync {
		return nil
	}

	if dryRun {
		return nil
	}

	// Remove orphaned worktrees first (if --force is used) to free up branches
	if force {
		for _, envVar := range status.OrphanedWorktrees {
			worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, envVar)
			err := m.gitManager.RemoveWorktree(worktreePath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned worktree %s: %w", envVar, err)
			}
		}
	}

	for _, envVar := range status.MissingWorktrees {
		branchName := m.envMapping.Variables[envVar]
		err := m.gitManager.CreateWorktree(envVar, branchName, m.config.Settings.WorktreePrefix)
		if err != nil {
			return fmt.Errorf("failed to create worktree for %s: %w", envVar, err)
		}
	}

	for envVar, change := range status.BranchChanges {
		worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, envVar)
		err := m.gitManager.UpdateWorktree(worktreePath, change.NewBranch)
		if err != nil {
			return fmt.Errorf("failed to update worktree for %s: %w", envVar, err)
		}
	}

	var trackedVars []string
	for envVar := range m.envMapping.Variables {
		trackedVars = append(trackedVars, envVar)
	}
	m.config.State.TrackedVars = trackedVars
	m.config.State.LastSync = time.Now()

	return m.config.Save(m.gbmDir)
}

func (m *Manager) ValidateEnvrc() error {
	if m.envMapping == nil {
		return fmt.Errorf("no .envrc mapping loaded")
	}

	for envVar, branchName := range m.envMapping.Variables {
		exists, err := m.gitManager.BranchExists(branchName)
		if err != nil {
			return fmt.Errorf("failed to check branch %s for %s: %w", branchName, envVar, err)
		}
		if !exists {
			return fmt.Errorf("branch '%s' does not exist", branchName)
		}
	}

	return nil
}

func (m *Manager) GetEnvMapping() (map[string]string, error) {
	if m.envMapping == nil {
		return nil, fmt.Errorf("no .envrc mapping loaded")
	}
	return m.envMapping.Variables, nil
}

func (m *Manager) BranchExists(branchName string) (bool, error) {
	return m.gitManager.BranchExists(branchName)
}

func (m *Manager) GetWorktreeList() (map[string]*WorktreeListInfo, error) {
	if m.envMapping == nil {
		return nil, fmt.Errorf("no .envrc mapping loaded")
	}

	result := make(map[string]*WorktreeListInfo)

	for envVar, expectedBranch := range m.envMapping.Variables {
		worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, envVar)

		info := &WorktreeListInfo{
			Path:           worktreePath,
			ExpectedBranch: expectedBranch,
			CurrentBranch:  "",
		}

		if _, err := os.Stat(worktreePath); err == nil {
			worktrees, err := m.gitManager.GetWorktrees()
			if err == nil {
				for _, wt := range worktrees {
					if wt.Path == worktreePath {
						info.CurrentBranch = wt.Branch
						break
					}
				}
			}

			// Get git status for the worktree
			if gitStatus, err := m.gitManager.GetWorktreeStatus(worktreePath); err == nil {
				info.GitStatus = gitStatus
			}
		}

		result[envVar] = info
	}

	return result, nil
}

func (m *Manager) GetStatusIcon(gitStatus *GitStatus) string {
	return m.gitManager.GetStatusIcon(gitStatus)
}

func (m *Manager) GetWorktreePath(worktreeName string) (string, error) {
	worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return "", fmt.Errorf("worktree directory '%s' does not exist", worktreeName)
	}

	return worktreePath, nil
}

func (m *Manager) GetAllWorktrees() (map[string]*WorktreeListInfo, error) {
	result := make(map[string]*WorktreeListInfo)

	// Get all actual worktrees from git
	worktrees, err := m.gitManager.GetWorktrees()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	worktreePrefix := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix)

	for _, wt := range worktrees {
		if strings.HasPrefix(wt.Path, worktreePrefix) {
			// Extract worktree name from path
			worktreeName := filepath.Base(wt.Path)

			info := &WorktreeListInfo{
				Path:          wt.Path,
				CurrentBranch: wt.Branch,
			}

			// Set expected branch if it's tracked in .envrc
			if m.envMapping != nil {
				if expectedBranch, exists := m.envMapping.Variables[worktreeName]; exists {
					info.ExpectedBranch = expectedBranch
				} else {
					info.ExpectedBranch = wt.Branch // Use current branch as expected for ad hoc worktrees
				}
			} else {
				info.ExpectedBranch = wt.Branch
			}

			// Get git status for the worktree
			if gitStatus, err := m.gitManager.GetWorktreeStatus(wt.Path); err == nil {
				info.GitStatus = gitStatus
			}

			result[worktreeName] = info
		}
	}

	return result, nil
}

func (m *Manager) AddWorktree(worktreeName, branchName string, createBranch bool) error {
	err := m.gitManager.AddWorktree(worktreeName, branchName, createBranch)
	if err != nil {
		return err
	}

	// Track this worktree as ad hoc if it's not in .envrc
	if m.envMapping != nil {
		if _, exists := m.envMapping.Variables[worktreeName]; !exists {
			// Add to ad hoc worktrees list if not already there
			if !contains(m.config.State.AdHocWorktrees, worktreeName) {
				m.config.State.AdHocWorktrees = append(m.config.State.AdHocWorktrees, worktreeName)
				// Save the updated config
				if saveErr := m.config.Save(m.gbmDir); saveErr != nil {
					// Log warning but don't fail the operation
					fmt.Printf("Warning: failed to save config: %v\n", saveErr)
				}
			}
		}
	}

	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (m *Manager) GetRemoteBranches() ([]string, error) {
	return m.gitManager.GetRemoteBranches()
}

func (m *Manager) GetCurrentBranch() (string, error) {
	return m.gitManager.GetCurrentBranch()
}

func (m *Manager) PushWorktree(worktreeName string) error {
	worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)
	return m.gitManager.PushWorktree(worktreePath)
}

func (m *Manager) PullWorktree(worktreeName string) error {
	worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)
	return m.gitManager.PullWorktree(worktreePath)
}

func (m *Manager) IsInWorktree(currentPath string) (bool, string, error) {
	return m.gitManager.IsInWorktree(currentPath)
}

func (m *Manager) PushAllWorktrees() error {
	worktrees, err := m.GetAllWorktrees()
	if err != nil {
		return fmt.Errorf("failed to get worktrees: %w", err)
	}

	for name, info := range worktrees {
		fmt.Printf("Pushing worktree '%s'...\n", name)
		if err := m.gitManager.PushWorktree(info.Path); err != nil {
			fmt.Printf("Failed to push worktree '%s': %v\n", name, err)
			continue
		}
		fmt.Printf("Successfully pushed worktree '%s'\n", name)
	}

	return nil
}

func (m *Manager) PullAllWorktrees() error {
	worktrees, err := m.GetAllWorktrees()
	if err != nil {
		return fmt.Errorf("failed to get worktrees: %w", err)
	}

	for name, info := range worktrees {
		fmt.Printf("Pulling worktree '%s'...\n", name)
		if err := m.gitManager.PullWorktree(info.Path); err != nil {
			fmt.Printf("Failed to pull worktree '%s': %v\n", name, err)
			continue
		}
		fmt.Printf("Successfully pulled worktree '%s'\n", name)
	}

	return nil
}

func (m *Manager) RemoveWorktree(worktreeName string) error {
	worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)

	// Remove the worktree using git
	if err := m.gitManager.RemoveWorktree(worktreePath); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	// Remove from ad hoc worktrees list if it exists there
	for i, name := range m.config.State.AdHocWorktrees {
		if name == worktreeName {
			m.config.State.AdHocWorktrees = append(m.config.State.AdHocWorktrees[:i], m.config.State.AdHocWorktrees[i+1:]...)
			break
		}
	}

	// Save the updated config
	if err := m.config.Save(m.gbmDir); err != nil {
		// Log warning but don't fail the operation
		fmt.Printf("Warning: failed to save config: %v\n", err)
	}

	return nil
}

func (m *Manager) GetWorktreeStatus(worktreePath string) (*GitStatus, error) {
	return m.gitManager.GetWorktreeStatus(worktreePath)
}

func (m *Manager) SetCurrentWorktree(worktreeName string) error {
	// Update previous worktree before changing current
	if m.config.State.CurrentWorktree != "" && m.config.State.CurrentWorktree != worktreeName {
		m.config.State.PreviousWorktree = m.config.State.CurrentWorktree
	}
	m.config.State.CurrentWorktree = worktreeName
	return m.config.Save(m.gbmDir)
}

func (m *Manager) GetPreviousWorktree() string {
	return m.config.State.PreviousWorktree
}

func (m *Manager) GetCurrentWorktree() string {
	return m.config.State.CurrentWorktree
}

func (m *Manager) GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string {
	var envrcNames []string
	var adHocNames []string

	// Get .envrc mapping if available
	envMapping := make(map[string]string)
	if m.envMapping != nil {
		envMapping = m.envMapping.Variables
	}

	// Separate worktrees into .envrc tracked and ad hoc
	for name := range worktrees {
		if _, exists := envMapping[name]; exists {
			envrcNames = append(envrcNames, name)
		} else {
			adHocNames = append(adHocNames, name)
		}
	}

	// Sort .envrc names alphabetically
	sort.Strings(envrcNames)

	// Sort ad hoc names by creation time (directory modification time) descending
	sort.Slice(adHocNames, func(i, j int) bool {
		pathI := worktrees[adHocNames[i]].Path
		pathJ := worktrees[adHocNames[j]].Path

		statI, errI := os.Stat(pathI)
		statJ, errJ := os.Stat(pathJ)

		// If we can't get stats, fall back to alphabetical
		if errI != nil || errJ != nil {
			return adHocNames[i] < adHocNames[j]
		}

		// Sort by modification time descending (newer first)
		return statI.ModTime().After(statJ.ModTime())
	})

	// Return .envrc worktrees first, then ad hoc worktrees
	result := make([]string, 0, len(envrcNames)+len(adHocNames))
	result = append(result, envrcNames...)
	result = append(result, adHocNames...)

	return result
}
