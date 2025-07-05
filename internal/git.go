package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type GitManager struct {
	repo     *git.Repository
	repoPath string
}

type WorktreeInfo struct {
	Name       string
	Path       string
	Branch     string
	IsOrphaned bool
	NeedsSync  bool
	GitStatus  *GitStatus
}

type GitStatus struct {
	IsDirty   bool
	Ahead     int
	Behind    int
	Untracked int
	Modified  int
	Staged    int
}

func (gs *GitStatus) HasChanges() bool {
	return gs.IsDirty || gs.Untracked > 0 || gs.Modified > 0 || gs.Staged > 0
}

// execCommand executes a command with debug output
func execCommand(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.Output()
	return output, err
}

// execCommandRun executes a command using Run() instead of Output() with debug output
func execCommandRun(cmd *exec.Cmd) error {
	err := cmd.Run()
	return err
}

// FindGitRoot finds the root directory of the git repository
func FindGitRoot(startPath string) (string, error) {
	// First, try direct git commands from the current directory
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = startPath
	gitDirOutput, err := cmd.Output()
	if err == nil {
		gitDir := strings.TrimSpace(string(gitDirOutput))

		// If .git-dir contains "worktrees", we're in a worktree
		if strings.Contains(gitDir, "worktrees") {
			// Extract the main repository path from the worktree git dir
			// Format: /path/to/main/repo/.git/worktrees/worktree-name
			parts := strings.Split(gitDir, "/.git/worktrees/")
			if len(parts) >= 2 {
				return parts[0], nil
			}
		}

		// Check if this is a bare repository
		cmd = exec.Command("git", "rev-parse", "--is-bare-repository")
		cmd.Dir = startPath
		output, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(output)) == "true" {
			// For bare repositories, the git directory is the repository root
			if filepath.IsAbs(gitDir) {
				return filepath.Dir(gitDir), nil
			} else {
				// gitDir is relative (e.g., ".git"), so repository root is startPath
				return startPath, nil
			}
		}

		// Standard git root detection if we're in a regular repo
		cmd = exec.Command("git", "rev-parse", "--show-toplevel")
		cmd.Dir = startPath
		output, err = cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}

		// If show-toplevel failed, try alternative detection methods
		// Check if we have a .git directory or file
		gitPath := filepath.Join(startPath, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			// We have a .git entry, so this is likely the repository root
			return startPath, nil
		}

		// If gitDir is absolute, we can derive the repository root
		if filepath.IsAbs(gitDir) {
			// Remove .git suffix to get repository root
			if strings.HasSuffix(gitDir, ".git") {
				return strings.TrimSuffix(gitDir, ".git"), nil
			}
			return filepath.Dir(gitDir), nil
		}
	}

	// If we're not in a git repository, check for git repositories in subdirectories
	// This handles the case where we're in a directory that contains worktrees
	// but isn't itself a git repository
	entries, err := os.ReadDir(startPath)
	if err == nil {
		var foundRepos []string

		for _, entry := range entries {
			if entry.IsDir() {
				subPath := filepath.Join(startPath, entry.Name())

				// Check if this subdirectory is a git repository
				cmd := exec.Command("git", "rev-parse", "--git-dir")
				cmd.Dir = subPath
				gitDirOutput, err := cmd.Output()
				if err == nil {
					gitDir := strings.TrimSpace(string(gitDirOutput))

					// If this is a worktree, get the main repository path
					if strings.Contains(gitDir, "worktrees") {
						parts := strings.Split(gitDir, "/.git/worktrees/")
						if len(parts) >= 2 {
							foundRepos = append(foundRepos, parts[0])
						}
					} else {
						// Check if this is a bare repository
						cmd = exec.Command("git", "rev-parse", "--is-bare-repository")
						cmd.Dir = subPath
						output, err := cmd.Output()
						if err == nil && strings.TrimSpace(string(output)) == "true" {
							// For bare repositories, the git directory is the repository root
							if filepath.IsAbs(gitDir) {
								foundRepos = append(foundRepos, filepath.Dir(gitDir))
							} else {
								// gitDir is relative (e.g., ".git"), so repository root is subPath
								foundRepos = append(foundRepos, subPath)
							}
						} else {
							// If this is a regular git repository, return its root
							cmd = exec.Command("git", "rev-parse", "--show-toplevel")
							cmd.Dir = subPath
							output, err := cmd.Output()
							if err == nil {
								foundRepos = append(foundRepos, strings.TrimSpace(string(output)))
							}
						}
					}
				}
			}
		}

		// If we found repositories in subdirectories, use the first one
		if len(foundRepos) > 0 {
			return foundRepos[0], nil
		}
	}

	return "", fmt.Errorf("not in a git repository and no git repositories found in subdirectories")
}

func NewGitManager(repoPath string) (*GitManager, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &GitManager{
		repo:     repo,
		repoPath: repoPath,
	}, nil
}

func (gm *GitManager) IsGitRepository() bool {
	_, err := git.PlainOpen(gm.repoPath)
	return err == nil
}

func (gm *GitManager) BranchExists(branchName string) (bool, error) {
	refs, err := gm.repo.References()
	if err != nil {
		return false, err
	}

	var found bool
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			if ref.Name().Short() == branchName {
				found = true
				return storer.ErrStop
			}
		}
		return nil
	})

	if err != nil && err != storer.ErrStop {
		return false, err
	}

	return found, nil
}

func (gm *GitManager) GetWorktrees() ([]*WorktreeInfo, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = gm.repoPath
	output, err := execCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	var infos []*WorktreeInfo
	lines := strings.Split(string(output), "\n")

	var currentWorktree *WorktreeInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentWorktree != nil {
				infos = append(infos, currentWorktree)
				currentWorktree = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			path := strings.TrimPrefix(line, "worktree ")
			currentWorktree = &WorktreeInfo{
				Name: filepath.Base(path),
				Path: path,
			}
		} else if strings.HasPrefix(line, "branch ") && currentWorktree != nil {
			branch := strings.TrimPrefix(line, "branch ")
			branch = strings.TrimPrefix(branch, "refs/heads/")
			currentWorktree.Branch = branch
		}
	}

	if currentWorktree != nil {
		infos = append(infos, currentWorktree)
	}

	return infos, nil
}

func (gm *GitManager) CreateWorktree(envVar, branchName, worktreeDir string) error {
	worktreePath := filepath.Join(gm.repoPath, worktreeDir, envVar)

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree directory already exists: %s", worktreePath)
	}

	branchExists, err := gm.BranchExists(branchName)
	if err != nil {
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}

	if !branchExists {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	// Check if remote tracking branch exists
	remoteBranch := fmt.Sprintf("origin/%s", branchName)
	cmd := exec.Command("git", "rev-parse", "--verify", remoteBranch)
	cmd.Dir = gm.repoPath
	_, err = execCommand(cmd)

	if err == nil {
		// Remote tracking branch exists, use --track
		cmd = exec.Command("git", "worktree", "add", "--track", "-b", branchName, worktreePath, remoteBranch)
	} else {
		// No remote tracking branch, create worktree without tracking
		cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
	}

	cmd.Dir = gm.repoPath
	if err := execCommandRun(cmd); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	return nil
}

func (gm *GitManager) RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
	cmd.Dir = gm.repoPath
	if err := execCommandRun(cmd); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

func (gm *GitManager) UpdateWorktree(worktreePath, newBranch string) error {
	if err := gm.RemoveWorktree(worktreePath); err != nil {
		return fmt.Errorf("failed to remove old worktree: %w", err)
	}

	worktreeDir := filepath.Dir(worktreePath)
	envVar := filepath.Base(worktreePath)
	relativeWorktreeDir := strings.TrimPrefix(worktreeDir, gm.repoPath+string(filepath.Separator))

	return gm.CreateWorktree(envVar, newBranch, relativeWorktreeDir)
}

func (gm *GitManager) FetchAll() error {
	remote, err := gm.repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("failed to get origin remote: %w", err)
	}

	err = remote.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch from remote: %w", err)
	}

	return nil
}

func (gm *GitManager) GetWorktreeStatus(worktreePath string) (*GitStatus, error) {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree path does not exist: %s", worktreePath)
	}

	status := &GitStatus{}

	// Get git status
	cmd := exec.Command("git", "status", "--porcelain", "--ahead-behind")
	cmd.Dir = worktreePath
	output, err := execCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	statusLines := strings.Split(string(output), "\n")
	for _, line := range statusLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		status.IsDirty = true

		// Parse git status output
		if len(line) >= 2 {
			indexStatus := line[0]
			worktreeStatus := line[1]

			switch indexStatus {
			case 'A', 'M', 'D', 'R', 'C':
				status.Staged++
			}

			switch worktreeStatus {
			case 'M', 'D':
				status.Modified++
			}

			if indexStatus == '?' && worktreeStatus == '?' {
				status.Untracked++
			}
		}
	}

	// Get ahead/behind info
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = worktreePath
	output, err = execCommand(cmd)
	if err == nil {
		parts := strings.Fields(string(output))
		if len(parts) == 2 {
			if _, err := fmt.Sscanf(parts[0], "%d", &status.Ahead); err == nil {
				fmt.Sscanf(parts[1], "%d", &status.Behind)
			}
		}
	}

	return status, nil
}

func (gm *GitManager) GetStatusIcon(gitStatus *GitStatus) string {
	iconManager := GetGlobalIconManager()

	if gitStatus == nil {
		return iconManager.Error()
	}

	var icons []string

	// Check ahead/behind first
	if gitStatus.Ahead > 0 && gitStatus.Behind > 0 {
		icons = append(icons, iconManager.GitDiverged())
	} else if gitStatus.Ahead > 0 {
		icons = append(icons, iconManager.GitAhead())
	} else if gitStatus.Behind > 0 {
		icons = append(icons, iconManager.GitBehind())
	}

	// Check dirty status
	if gitStatus.IsDirty {
		if gitStatus.Staged > 0 {
			icons = append(icons, "●")
		}
		if gitStatus.Modified > 0 {
			icons = append(icons, "✚")
		}
		if gitStatus.Untracked > 0 {
			icons = append(icons, "?")
		}
	}

	if len(icons) == 0 {
		return iconManager.Success()
	}

	return strings.Join(icons, "")
}

func (gm *GitManager) CreateBranch(branchName, baseBranch string) error {
	if baseBranch == "" {
		baseBranch = "HEAD"
	}

	cmd := exec.Command("git", "branch", branchName, baseBranch)
	cmd.Dir = gm.repoPath
	if err := execCommandRun(cmd); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func (gm *GitManager) AddWorktree(worktreeName, branchName string, createBranch bool) error {
	worktreeDir := filepath.Join(gm.repoPath, "worktrees")
	worktreePath := filepath.Join(worktreeDir, worktreeName)

	// Check if worktree already exists
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree '%s' already exists", worktreeName)
	}

	// Create worktrees directory if it doesn't exist
	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	var cmd *exec.Cmd
	if createBranch {
		// Check if branch already exists
		branchExists, err := gm.BranchExists(branchName)
		if err != nil {
			return fmt.Errorf("failed to check if branch exists: %w", err)
		}

		if branchExists {
			// Branch exists, create worktree on existing branch
			cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
		} else {
			// Create new branch and worktree
			cmd = exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
		}
	} else {
		// Create worktree on existing branch
		branchExists, err := gm.BranchExists(branchName)
		if err != nil {
			return fmt.Errorf("failed to check if branch exists: %w", err)
		}

		if !branchExists {
			return fmt.Errorf("branch '%s' does not exist", branchName)
		}

		// Check if remote tracking branch exists
		remoteBranch := fmt.Sprintf("origin/%s", branchName)
		checkCmd := exec.Command("git", "rev-parse", "--verify", remoteBranch)
		checkCmd.Dir = gm.repoPath
		_, err = execCommand(checkCmd)

		if err == nil {
			// Remote tracking branch exists, use --track
			cmd = exec.Command("git", "worktree", "add", "--track", "-b", branchName, worktreePath, remoteBranch)
		} else {
			// No remote tracking branch, create worktree without tracking
			cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
		}
	}

	cmd.Dir = gm.repoPath
	if err := execCommandRun(cmd); err != nil {
		return fmt.Errorf("failed to add worktree (command: %s): %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (gm *GitManager) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = gm.repoPath
	output, err := execCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (gm *GitManager) GetRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r")
	cmd.Dir = gm.repoPath
	output, err := execCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	var branches []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "HEAD") {
			continue
		}

		// Remove "origin/" prefix
		if strings.HasPrefix(line, "origin/") {
			branch := strings.TrimPrefix(line, "origin/")
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

func (gm *GitManager) PushWorktree(worktreePath string) error {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree path does not exist: %s", worktreePath)
	}

	// Get the current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = worktreePath
	output, err := execCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	currentBranch := strings.TrimSpace(string(output))

	// Check if upstream is set
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	cmd.Dir = worktreePath
	_, err = execCommand(cmd)

	if err != nil {
		// No upstream set, push with -u flag
		cmd = exec.Command("git", "push", "-u", "origin", currentBranch)
	} else {
		// Upstream is set, simple push
		cmd = exec.Command("git", "push")
	}

	cmd.Dir = worktreePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (gm *GitManager) PullWorktree(worktreePath string) error {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree path does not exist: %s", worktreePath)
	}

	// Get the current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = worktreePath
	output, err := execCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	currentBranch := strings.TrimSpace(string(output))

	// Check if upstream is set
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	cmd.Dir = worktreePath
	_, err = execCommand(cmd)

	if err != nil {
		// No upstream set, try to set it and pull
		remoteBranch := fmt.Sprintf("origin/%s", currentBranch)

		// Check if remote branch exists
		cmd = exec.Command("git", "rev-parse", "--verify", remoteBranch)
		cmd.Dir = worktreePath
		_, err = execCommand(cmd)

		if err == nil {
			// Remote branch exists, set upstream and pull
			cmd = exec.Command("git", "branch", "--set-upstream-to", remoteBranch)
			cmd.Dir = worktreePath
			if err := execCommandRun(cmd); err != nil {
				return fmt.Errorf("failed to set upstream: %w", err)
			}

			// Now pull
			cmd = exec.Command("git", "pull")
		} else {
			// Remote branch doesn't exist, try to pull with explicit remote and branch
			cmd = exec.Command("git", "pull", "origin", currentBranch)
		}
	} else {
		// Upstream is set, simple pull
		cmd = exec.Command("git", "pull")
	}

	cmd.Dir = worktreePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (gm *GitManager) IsInWorktree(currentPath string) (bool, string, error) {
	// Check if we're in a worktree
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = currentPath
	output, err := execCommand(cmd)
	if err != nil {
		return false, "", fmt.Errorf("not in a git repository: %w", err)
	}

	worktreePath := strings.TrimSpace(string(output))

	// Check if this is a worktree by checking if it's under the worktrees directory
	worktreePrefix := filepath.Join(gm.repoPath, "worktrees")
	if strings.HasPrefix(worktreePath, worktreePrefix) {
		worktreeName := filepath.Base(worktreePath)
		return true, worktreeName, nil
	}

	return false, "", nil
}
