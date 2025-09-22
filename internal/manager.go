package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
)

type Manager struct {
	config     *Config
	state      *State
	gitManager *GitManager
	gbmConfig  *GBMConfig
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
	InSync             bool
	MissingWorktrees   []string
	OrphanedWorktrees  []string
	BranchChanges      map[string]BranchChange
	WorktreePromotions []WorktreePromotion
}

type BranchChange struct {
	OldBranch string
	NewBranch string
}

type WorktreePromotion struct {
	SourceWorktree string
	TargetWorktree string
	Branch         string
	SourceBranch   string
	TargetBranch   string
}

type ConfirmationFunc func(message string) bool

func NewManager(repoPath string) (*Manager, error) {
	gbmDir := filepath.Join(repoPath, DefaultConfigDirname)
	config, err := LoadConfig(gbmDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	state, err := LoadState(gbmDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	gitManager, err := NewGitManager(repoPath, config.Settings.WorktreePrefix)
	if err != nil {
		return nil, err
	}

	// Initialize the global icon manager with the loaded config
	iconManager := NewIconManager(config)
	SetGlobalIconManager(iconManager)

	return &Manager{
		config:     config,
		state:      state,
		gitManager: gitManager,
		repoPath:   repoPath,
		gbmDir:     gbmDir,
	}, nil
}

func (m *Manager) LoadGBMConfig(configPath string) error {
	if configPath == "" {
		configPath = DefaultBranchConfigFilename
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(m.repoPath, configPath)
	}

	// Resolve symlinks to handle macOS /var -> /private/var and similar cases
	if resolved, err := filepath.EvalSymlinks(configPath); err == nil {
		configPath = resolved
	}

	gbmConfig, err := ParseGBMConfig(configPath)
	if err != nil {
		return err
	}

	m.gbmConfig = gbmConfig
	return nil
}

func (m *Manager) GetSyncStatus() (*SyncStatus, error) {
	if m.gbmConfig == nil {
		if err := m.LoadGBMConfig(""); err != nil {
			return nil, fmt.Errorf("no %s loaded", DefaultBranchConfigFilename)
		}
	}

	status := &SyncStatus{
		InSync:             true,
		MissingWorktrees:   []string{},
		OrphanedWorktrees:  []string{},
		BranchChanges:      make(map[string]BranchChange),
		WorktreePromotions: []WorktreePromotion{},
	}

	worktrees, err := m.gitManager.GetWorktrees()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	worktreeMap := make(map[string]*WorktreeInfo)
	// Resolve symlinks for robust prefix checks on systems where /var -> /private/var
	prefix := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix)
	resolvedPrefix, err := filepath.EvalSymlinks(prefix)
	if err != nil {
		resolvedPrefix = prefix
	}
	for _, wt := range worktrees {
		resolvedPath, err := filepath.EvalSymlinks(wt.Path)
		if err != nil {
			resolvedPath = wt.Path
		}
		if strings.HasPrefix(resolvedPath, resolvedPrefix) {
			worktreeMap[wt.Name] = wt
		}
	}

	for worktreeName, worktreeConfig := range m.gbmConfig.Worktrees {
		if wt, exists := worktreeMap[worktreeName]; exists {
			if wt.Branch != worktreeConfig.Branch {
				status.BranchChanges[worktreeName] = BranchChange{
					OldBranch: wt.Branch,
					NewBranch: worktreeConfig.Branch,
				}
				status.InSync = false
			}
			delete(worktreeMap, worktreeName)
		} else {
			status.MissingWorktrees = append(status.MissingWorktrees, worktreeName)
			status.InSync = false
		}
	}

	for worktreeName := range worktreeMap {
		status.OrphanedWorktrees = append(status.OrphanedWorktrees, worktreeName)
		status.InSync = false
	}

	// Detect worktree promotions: when a branch moves from one worktree to another
	status.WorktreePromotions = m.detectWorktreePromotions(status.BranchChanges, worktrees)

	return status, nil
}

func (m *Manager) detectWorktreePromotions(branchChanges map[string]BranchChange, allWorktrees []*WorktreeInfo) []WorktreePromotion {
	var promotions []WorktreePromotion

	// Create a map of branch -> worktree for existing worktrees
	branchToWorktree := make(map[string]string)
	expectedPrefix := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix)

	// Resolve symlinks for consistent path comparison (e.g., /var vs /private/var on macOS)
	resolvedExpectedPrefix, err := filepath.EvalSymlinks(expectedPrefix)
	if err != nil {
		resolvedExpectedPrefix = expectedPrefix // fallback to original if resolution fails
	}

	for _, wt := range allWorktrees {
		resolvedWorktreePath, err := filepath.EvalSymlinks(wt.Path)
		if err != nil {
			resolvedWorktreePath = wt.Path // fallback to original if resolution fails
		}

		if strings.HasPrefix(resolvedWorktreePath, resolvedExpectedPrefix) {
			branchToWorktree[wt.Branch] = wt.Name
		}
	}

	// Check each branch change to see if the new branch is currently checked out elsewhere
	for targetWorktree, change := range branchChanges {
		if sourceWorktree, exists := branchToWorktree[change.NewBranch]; exists {
			// This is a promotion: the new branch is currently in another worktree
			promotion := WorktreePromotion{
				SourceWorktree: sourceWorktree,
				TargetWorktree: targetWorktree,
				Branch:         change.NewBranch,
				SourceBranch:   change.NewBranch,
				TargetBranch:   change.OldBranch,
			}
			promotions = append(promotions, promotion)
		}
	}

	return promotions
}

func (m *Manager) Sync(dryRun, force bool) error {
	return m.SyncWithConfirmation(dryRun, force, false, nil)
}

func (m *Manager) SyncWithConfirmation(dryRun, force bool, removeOrphans bool, confirmFunc ConfirmationFunc) error {
	// Validate all branches exist before performing any operations
	if err := m.ValidateConfig(); err != nil {
		return err
	}

	// Always fetch from remote before sync
	if err := m.gitManager.FetchAll(); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
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

	// Ensure worktrees directory exists before creating/moving worktrees
	worktreesDir := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix)
	if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
		return fmt.Errorf("failed to ensure worktrees directory: %w", err)
	}

	// Remove orphaned worktrees first (if --remove-orphans is used) to free up branches
	if removeOrphans && len(status.OrphanedWorktrees) > 0 {
		// Ask for confirmation unless --force is used
		if !force && confirmFunc != nil {
			message := "The following worktrees will be PERMANENTLY DELETED:\n"
			for _, envVar := range status.OrphanedWorktrees {
				worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, envVar)
				message += fmt.Sprintf("  • %s (%s)\n", envVar, worktreePath)
			}
			message += "Do you want to continue?"

			if !confirmFunc(message) {
				return fmt.Errorf("orphan removal cancelled by user")
			}
		}

		for _, envVar := range status.OrphanedWorktrees {
			worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, envVar)
			err := m.gitManager.RemoveWorktree(worktreePath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned worktree %s: %w", envVar, err)
			}
		}
	}

	for _, worktreeName := range status.MissingWorktrees {
		worktreeConfig := m.gbmConfig.Worktrees[worktreeName]
		// If the directory exists but is empty (e.g., created by .gitignore), remove it first
		worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)
		if stat, err := os.Stat(worktreePath); err == nil && stat.IsDir() {
			// Check if directory is empty
			entries, readErr := os.ReadDir(worktreePath)
			if readErr == nil && len(entries) == 0 {
				_ = os.Remove(worktreePath)
			}
		}

		err := m.gitManager.CreateWorktree(worktreeName, worktreeConfig.Branch, m.config.Settings.WorktreePrefix)
		if err != nil {
			// Special case: if creating a worktree fails because directory already exists,
			// check if this is the main worktree already present in repository root
			if errors.Is(err, ErrWorktreeDirectoryExists) {
				if worktreeName == worktreeConfig.Branch {
					// Skip creating this worktree since it already exists as the main repository
					continue
				}
			}
			return fmt.Errorf("failed to create worktree for %s: %w", worktreeName, err)
		}
	}

	// Handle worktree promotions with confirmation (always required for destructive operations)
	if len(status.WorktreePromotions) > 0 {
		if confirmFunc != nil {
			for _, promotion := range status.WorktreePromotions {
				message := fmt.Sprintf("Worktree %s (%s) will be promoted to %s.\nThis is a destructive action:\n  1. Worktree %s (%s) will be removed.\n  2. Worktree %s (%s) will be moved to %s.\nContinue?",
					promotion.SourceWorktree, promotion.Branch, promotion.TargetWorktree,
					promotion.TargetWorktree, promotion.TargetBranch,
					promotion.SourceWorktree, promotion.Branch, promotion.TargetWorktree)

				if !confirmFunc(message) {
					return fmt.Errorf("worktree promotion cancelled by user")
				}
			}
		} else {
			return fmt.Errorf("worktree promotions require confirmation, but no confirmation function provided")
		}
	}

	// Process worktree promotions first by removing both worktrees, then recreating with swapped branches
	for _, promotion := range status.WorktreePromotions {
		sourceWorktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, promotion.SourceWorktree)
		targetWorktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, promotion.TargetWorktree)

		// Remove both worktrees to free up branches
		if err := m.gitManager.RemoveWorktree(sourceWorktreePath); err != nil {
			return fmt.Errorf("failed to remove source worktree %s: %w", promotion.SourceWorktree, err)
		}
		if err := m.gitManager.RemoveWorktree(targetWorktreePath); err != nil {
			return fmt.Errorf("failed to remove target worktree %s: %w", promotion.TargetWorktree, err)
		}

		// Recreate worktrees with swapped branches
		if err := m.gitManager.CreateWorktree(promotion.TargetWorktree, promotion.SourceBranch, m.config.Settings.WorktreePrefix); err != nil {
			return fmt.Errorf("failed to create target worktree %s with branch %s: %w", promotion.TargetWorktree, promotion.SourceBranch, err)
		}
		if err := m.gitManager.CreateWorktree(promotion.SourceWorktree, promotion.TargetBranch, m.config.Settings.WorktreePrefix); err != nil {
			return fmt.Errorf("failed to create source worktree %s with branch %s: %w", promotion.SourceWorktree, promotion.TargetBranch, err)
		}

		// Remove from regular branch changes since already handled
		delete(status.BranchChanges, promotion.TargetWorktree)
		delete(status.BranchChanges, promotion.SourceWorktree)
	}

	for worktreeName, change := range status.BranchChanges {
		worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)
		err := m.gitManager.UpdateWorktree(worktreePath, change.NewBranch)
		if err != nil {
			return fmt.Errorf("failed to update worktree for %s: %w", worktreeName, err)
		}
	}

	var trackedWorktrees []string
	for worktreeName := range m.gbmConfig.Worktrees {
		trackedWorktrees = append(trackedWorktrees, worktreeName)
	}
	m.state.TrackedVars = trackedWorktrees
	m.state.LastSync = time.Now()

	return m.SaveState()
}

func (m *Manager) ValidateConfig() error {
	if m.gbmConfig == nil {
		if err := m.LoadGBMConfig(""); err != nil {
			return fmt.Errorf("no %s loaded", DefaultBranchConfigFilename)
		}
	}

	for worktreeName, worktreeConfig := range m.gbmConfig.Worktrees {
		exists, err := m.gitManager.BranchExistsLocalOrRemote(worktreeConfig.Branch)
		if err != nil {
			return fmt.Errorf("failed to check branch %s for %s: %w", worktreeConfig.Branch, worktreeName, err)
		}
		if !exists {
			return fmt.Errorf("branch '%s' does not exist", worktreeConfig.Branch)
		}
	}

	return nil
}

func (m *Manager) GetWorktreeMapping() (map[string]string, error) {
	if m.gbmConfig == nil {
		if err := m.LoadGBMConfig(""); err != nil {
			return nil, fmt.Errorf("no %s loaded", DefaultBranchConfigFilename)
		}
	}

	mapping := make(map[string]string)
	for worktreeName, worktreeConfig := range m.gbmConfig.Worktrees {
		mapping[worktreeName] = worktreeConfig.Branch
	}
	return mapping, nil
}

func (m *Manager) BranchExists(branchName string) (bool, error) {
	return m.gitManager.BranchExists(branchName)
}

// BranchExistsLocal checks if a branch exists locally only (not remote)
func (m *Manager) BranchExistsLocal(branchName string) (bool, error) {
	return m.gitManager.BranchExistsLocal(branchName)
}

// GetDefaultBranch returns the repository's default branch
func (m *Manager) GetDefaultBranch() (string, error) {
	return m.gitManager.GetDefaultBranch()
}

// GetJiraIssues returns JIRA issues for the current user
func (m *Manager) GetJiraIssues() ([]JiraIssue, error) {
	return GetJiraIssues(m)
}

// GenerateBranchFromJira generates a branch name from a JIRA key
func (m *Manager) GenerateBranchFromJira(jiraKey string) (string, error) {
	return GenerateBranchFromJira(jiraKey, m)
}

func (m *Manager) GetWorktreeList() (map[string]*WorktreeListInfo, error) {
	if m.gbmConfig == nil {
		if err := m.LoadGBMConfig(""); err != nil {
			return nil, fmt.Errorf("no %s loaded", DefaultBranchConfigFilename)
		}
	}

	result := make(map[string]*WorktreeListInfo)

	for worktreeName, worktreeConfig := range m.gbmConfig.Worktrees {
		worktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, worktreeName)

		info := &WorktreeListInfo{
			Path:           worktreePath,
			ExpectedBranch: worktreeConfig.Branch,
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

		result[worktreeName] = info
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

	// Resolve symlinks to handle cases like /var vs /private/var on macOS
	resolvedWorktreePrefix, err := filepath.EvalSymlinks(worktreePrefix)
	if err != nil {
		resolvedWorktreePrefix = worktreePrefix // fallback to original if resolution fails
	}

	for _, wt := range worktrees {
		// Resolve symlinks for worktree path as well
		resolvedWtPath, err := filepath.EvalSymlinks(wt.Path)
		if err != nil {
			resolvedWtPath = wt.Path // fallback to original if resolution fails
		}

		if strings.HasPrefix(resolvedWtPath, resolvedWorktreePrefix) {
			// Extract worktree name from path
			worktreeName := filepath.Base(wt.Path)

			info := &WorktreeListInfo{
				Path:          wt.Path,
				CurrentBranch: wt.Branch,
			}

			// Set expected branch if it's tracked in gbm.branchconfig.yaml
			if m.gbmConfig != nil {
				if worktreeConfig, exists := m.gbmConfig.Worktrees[worktreeName]; exists {
					info.ExpectedBranch = worktreeConfig.Branch
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

func (m *Manager) AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error {
	err := m.gitManager.AddWorktree(worktreeName, branchName, createBranch, baseBranch)
	if err != nil {
		return err
	}

	// Check if this is an ad-hoc worktree (not tracked in gbm.branchconfig.yaml)
	isAdHoc := true
	if m.gbmConfig != nil {
		if _, exists := m.gbmConfig.Worktrees[worktreeName]; exists {
			isAdHoc = false
		}
	}

	// Only copy files for ad-hoc worktrees
	if isAdHoc {
		if err := m.copyFilesToWorktree(worktreeName); err != nil {
			fmt.Printf("Warning: failed to copy files to worktree: %v\n", err)
		}
	}

	// Store the base branch information for this worktree
	m.state.SetWorktreeBaseBranch(worktreeName, baseBranch)

	// Track this worktree as ad hoc if it's not in gbm.branchconfig.yaml
	if m.gbmConfig != nil {
		if _, exists := m.gbmConfig.Worktrees[worktreeName]; !exists {
			// Add to ad hoc worktrees list if not already there
			if !contains(m.state.AdHocWorktrees, worktreeName) {
				m.state.AdHocWorktrees = append(m.state.AdHocWorktrees, worktreeName)
			}
		}
	}

	// Save the updated state
	if saveErr := m.SaveState(); saveErr != nil {
		// Log warning but don't fail the operation
		fmt.Printf("Warning: failed to save state: %v\n", saveErr)
	}

	return nil
}

// copyFilesToWorktree copies files from source worktrees to the newly created worktree
func (m *Manager) copyFilesToWorktree(targetWorktreeName string) error {
	if len(m.config.FileCopy.Rules) == 0 {
		return nil
	}

	targetWorktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, targetWorktreeName)

	for _, rule := range m.config.FileCopy.Rules {
		sourceWorktreePath := filepath.Join(m.repoPath, m.config.Settings.WorktreePrefix, rule.SourceWorktree)

		// Check if source worktree exists
		if _, err := os.Stat(sourceWorktreePath); os.IsNotExist(err) {
			fmt.Printf("Warning: source worktree '%s' does not exist, skipping file copy rule\n", rule.SourceWorktree)
			continue
		}

		for _, filePattern := range rule.Files {
			if err := m.copyFileOrDirectory(sourceWorktreePath, targetWorktreePath, filePattern); err != nil {
				fmt.Printf("Warning: failed to copy '%s' from '%s': %v\n", filePattern, rule.SourceWorktree, err)
			}
		}
	}

	return nil
}

// copyFileOrDirectory copies a file or directory from source to target
func (m *Manager) copyFileOrDirectory(sourceWorktreePath, targetWorktreePath, filePattern string) error {
	sourcePath := filepath.Join(sourceWorktreePath, filePattern)
	targetPath := filepath.Join(targetWorktreePath, filePattern)

	sourceInfo, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source file/directory '%s' does not exist", sourcePath)
	}
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}

	if sourceInfo.IsDir() {
		return m.copyDirectory(sourcePath, targetPath)
	}
	return m.copyFile(sourcePath, targetPath)
}

// copyFile copies a single file from source to target
func (m *Manager) copyFile(sourcePath, targetPath string) error {
	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Check if target file already exists
	if _, err := os.Stat(targetPath); err == nil {
		fmt.Printf("File '%s' already exists in target worktree, skipping\n", filepath.Base(targetPath))
		return nil
	}

	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()

	// Create target file
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer func() { _ = targetFile.Close() }()

	// Copy file contents
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}
	if err := os.Chmod(targetPath, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// copyDirectory recursively copies a directory from source to target
func (m *Manager) copyDirectory(sourcePath, targetPath string) error {
	// Create target directory
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Read source directory
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		sourceEntryPath := filepath.Join(sourcePath, entry.Name())
		targetEntryPath := filepath.Join(targetPath, entry.Name())

		if entry.IsDir() {
			if err := m.copyDirectory(sourceEntryPath, targetEntryPath); err != nil {
				return err
			}
		} else {
			if err := m.copyFile(sourceEntryPath, targetEntryPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func (m *Manager) GetRemoteBranches() ([]string, error) {
	return m.gitManager.GetRemoteBranches()
}

func (m *Manager) GetCurrentBranch() (string, error) {
	return m.gitManager.GetCurrentBranch()
}

func (m *Manager) GetGitManager() *GitManager {
	return m.gitManager
}

func (m *Manager) GetGBMConfig() *GBMConfig {
	return m.gbmConfig
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
	for i, name := range m.state.AdHocWorktrees {
		if name == worktreeName {
			m.state.AdHocWorktrees = slices.Delete(m.state.AdHocWorktrees, i, i+1)
			break
		}
	}

	// Remove base branch information
	m.state.RemoveWorktreeBaseBranch(worktreeName)

	// Save the updated state
	if err := m.SaveState(); err != nil {
		// Log warning but don't fail the operation
		fmt.Printf("Warning: failed to save state: %v\n", err)
	}

	return nil
}

func (m *Manager) GetWorktreeStatus(worktreePath string) (*GitStatus, error) {
	return m.gitManager.GetWorktreeStatus(worktreePath)
}

func (m *Manager) SetCurrentWorktree(worktreeName string) error {
	// Update previous worktree before changing current
	if m.state.CurrentWorktree != "" && m.state.CurrentWorktree != worktreeName {
		m.state.PreviousWorktree = m.state.CurrentWorktree
	}
	m.state.CurrentWorktree = worktreeName
	return m.SaveState()
}

func (m *Manager) GetPreviousWorktree() string {
	return m.state.PreviousWorktree
}

func (m *Manager) GetCurrentWorktree() string {
	return m.state.CurrentWorktree
}

func (m *Manager) GetConfig() *Config {
	return m.config
}

func (m *Manager) GetState() *State {
	return m.state
}

func (m *Manager) GetRepoPath() string {
	return m.repoPath
}

func (m *Manager) SaveConfig() error {
	return m.config.Save(m.gbmDir)
}

func (m *Manager) SaveState() error {
	return m.state.Save(m.gbmDir)
}

func (m *Manager) GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string {
	var trackedNames []string
	var adHocNames []string

	// Get gbm.branchconfig.yaml mapping if available
	trackedWorktrees := make(map[string]string)
	if m.gbmConfig != nil {
		for name, config := range m.gbmConfig.Worktrees {
			trackedWorktrees[name] = config.Branch
		}
	}

	// Separate worktrees into tracked and ad hoc
	for name := range worktrees {
		if _, exists := trackedWorktrees[name]; exists {
			trackedNames = append(trackedNames, name)
		} else {
			adHocNames = append(adHocNames, name)
		}
	}

	// Sort tracked names alphabetically
	sort.Strings(trackedNames)

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

	// Return tracked worktrees first, then ad hoc worktrees
	result := make([]string, 0, len(trackedNames)+len(adHocNames))
	result = append(result, trackedNames...)
	result = append(result, adHocNames...)

	return result
}

// Wrapper methods for GitManager operations used by cmd/info.go

// GetWorktreeCommitHistory retrieves commit history for a specific worktree
func (m *Manager) GetWorktreeCommitHistory(worktreePath string, limit int) ([]CommitInfo, error) {
	return m.gitManager.GetCommitHistory(worktreePath, CommitHistoryOptions{
		Limit: limit,
	})
}

// GetWorktreeFileChanges retrieves modified files for a specific worktree
func (m *Manager) GetWorktreeFileChanges(worktreePath string) ([]FileChange, error) {
	return m.gitManager.GetFileChanges(worktreePath, FileChangeOptions{
		Staged:   true,
		Unstaged: true,
	})
}

// GetWorktreeCurrentBranch gets the current branch for a specific worktree
func (m *Manager) GetWorktreeCurrentBranch(worktreePath string) (string, error) {
	return m.gitManager.GetCurrentBranchInPath(worktreePath)
}

// GetWorktreeUpstreamBranch gets the upstream branch for a specific worktree
func (m *Manager) GetWorktreeUpstreamBranch(worktreePath string) (string, error) {
	return m.gitManager.GetUpstreamBranch(worktreePath)
}

// GetWorktreeAheadBehindCount gets the ahead/behind count for a specific worktree
func (m *Manager) GetWorktreeAheadBehindCount(worktreePath string) (int, int, error) {
	return m.gitManager.GetAheadBehindCount(worktreePath)
}

// VerifyWorktreeRef verifies if a ref exists in a specific worktree
func (m *Manager) VerifyWorktreeRef(ref string, worktreePath string) (bool, error) {
	return m.gitManager.VerifyRefInPath(worktreePath, ref)
}

// GetWorktrees retrieves all worktrees from the git repository
func (m *Manager) GetWorktrees() ([]*WorktreeInfo, error) {
	return m.gitManager.GetWorktrees()
}

// JIRA interaction methods - delegate to jira package

// GetJiraTicketDetails retrieves detailed JIRA ticket information using the JIRA CLI
func (m *Manager) GetJiraTicketDetails(jiraKey string) (*JiraTicketDetails, error) {
	return GetJiraTicketDetails(jiraKey)
}

// FindProductionBranch finds the actual production deployment branch for hotfixes
// This finds the deepest leaf nodes in the worktree tree - branches that nothing merges into
func (m *Manager) FindProductionBranch() (string, error) {
	// Look for gbm.branchconfig.yaml file
	configPath := filepath.Join(m.repoPath, DefaultBranchConfigFilename)
	config, err := ParseGBMConfig(configPath)
	if err != nil {
		// If no config file, fall back to default branch
		return m.GetDefaultBranch()
	}

	// Create worktree tree manager
	wtManager, err := NewWorktreeManager(config)
	if err != nil {
		return "", fmt.Errorf("failed to create worktree tree: %w", err)
	}

	// Get the deepest leaf nodes - these are the production branches
	deepestLeaves := wtManager.GetAllDeepestLeafNodes()

	if len(deepestLeaves) == 0 {
		// No leaves found, fall back to root nodes
		roots := wtManager.GetRoots()
		if len(roots) == 0 {
			return "", fmt.Errorf("no production branch found in mergeback configuration")
		}
		// Return the first root branch
		return roots[0].Config.Branch, nil
	}

	// Return the first deepest leaf branch
	// Note: In case of multiple valid production branches (multiple deployment chains),
	// this returns the first one found. This is consistent with the original behavior.
	return deepestLeaves[0].Config.Branch, nil
}
