package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type GitManager struct {
	repo           *git.Repository
	repoPath       string
	worktreePrefix string
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

// ExecGitCommand executes a git command in the specified directory with optional output capture
// This unified function replaces multiple duplicate git execution patterns across the codebase
func ExecGitCommand(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Output()
}

// execGitCommandRun executes a git command in the specified directory without capturing output
func execGitCommandRun(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// ExecGitCommandCombined executes a git command and returns combined stdout/stderr output
func ExecGitCommandCombined(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.CombinedOutput()
}

// ExecGitCommandInteractive executes a git command in the specified directory with live output to terminal
func ExecGitCommandInteractive(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// enhanceGitError provides more specific error messages for common git failures
func enhanceGitError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Extract exit code from error
	if exitError, ok := err.(*exec.ExitError); ok {
		stderr := string(exitError.Stderr)

		switch exitError.ExitCode() {
		case 128:
			if strings.Contains(stderr, "already checked out") {
				return fmt.Errorf("branch is already checked out in another worktree: %w", err)
			}
			if strings.Contains(stderr, "not a git repository") {
				return fmt.Errorf("not a git repository: %w", err)
			}
			if strings.Contains(stderr, "does not exist") && strings.Contains(operation, "worktree") {
				return fmt.Errorf("branch or worktree does not exist: %w", err)
			}
			if strings.Contains(stderr, "destination path") && strings.Contains(stderr, "already exists") {
				return fmt.Errorf("worktree directory already exists: %w", err)
			}
			return fmt.Errorf("git %s failed (exit 128): %s", operation, stderr)
		case 1:
			if strings.Contains(stderr, "not found") {
				return fmt.Errorf("branch or reference not found: %w", err)
			}
			return fmt.Errorf("git %s failed: %s", operation, stderr)
		default:
			return fmt.Errorf("git %s failed (exit %d): %s", operation, exitError.ExitCode(), stderr)
		}
	}

	return fmt.Errorf("git %s failed: %w", operation, err)
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

func NewGitManager(repoPath string, worktreePrefix string) (*GitManager, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &GitManager{
		repo:           repo,
		repoPath:       repoPath,
		worktreePrefix: worktreePrefix,
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
		// Also check remote branches
		if ref.Name().IsRemote() {
			remoteBranch := ref.Name().Short()
			if strings.HasPrefix(remoteBranch, "origin/") {
				localBranch := strings.TrimPrefix(remoteBranch, "origin/")
				if localBranch == branchName {
					found = true
					return storer.ErrStop
				}
			}
		}
		return nil
	})

	if err != nil && err != storer.ErrStop {
		return false, err
	}

	return found, nil
}

// Remote returns the remote branch name for a given branch (e.g., "main" -> "origin/main")
func Remote(branchName string) string {
	return fmt.Sprintf("origin/%s", branchName)
}

// BranchExistsLocalOrRemote checks if a branch exists either locally or remotely
func (gm *GitManager) BranchExistsLocalOrRemote(branchName string) (bool, error) {
	// // Check if local branch exists
	// _, err := ExecGitCommand(gm.repoPath, "rev-parse", "--verify", branchName)
	// if err == nil {
	// 	return true, nil
	// }

	// Check if remote branch exists
	remoteBranch := Remote(branchName)
	_, err := ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)
	return err == nil, nil
}

func (gm *GitManager) IsBranchAvailable(branchName string) (bool, error) {
	// First check if branch exists
	exists, err := gm.BranchExists(branchName)
	if err != nil {
		return false, fmt.Errorf("failed to check if branch exists: %w", err)
	}
	if !exists {
		return false, fmt.Errorf("branch %s does not exist", branchName)
	}

	// Get all worktrees and check if branch is checked out elsewhere
	worktrees, err := gm.GetWorktrees()
	if err != nil {
		return false, fmt.Errorf("failed to get worktrees: %w", err)
	}

	for _, wt := range worktrees {
		if wt.Branch == branchName {
			return false, fmt.Errorf("branch %s is already checked out in worktree %s", branchName, wt.Name)
		}
	}

	return true, nil
}

func (gm *GitManager) GetWorktrees() ([]*WorktreeInfo, error) {
	output, err := ExecGitCommand(gm.repoPath, "worktree", "list", "--porcelain")
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

var ErrWorktreeDirectoryExists = fmt.Errorf("worktree directory already exists")

func (gm *GitManager) CreateWorktree(envVar, branchName, worktreeDir string) error {
	worktreePath := filepath.Join(gm.repoPath, worktreeDir, envVar)

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrWorktreeDirectoryExists, worktreePath)
	}

	branchExists, err := gm.BranchExists(branchName)
	if err != nil {
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}

	if !branchExists {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	// Check if remote tracking branch exists
	remoteBranch := Remote(branchName)
	_, err = ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)

	if err == nil {
		// Remote tracking branch exists, create worktree and set up tracking
		if err := execGitCommandRun(gm.repoPath, "worktree", "add", worktreePath, branchName); err != nil {
			return enhanceGitError(err, "worktree add")
		}

		// Set up tracking for the remote branch
		if err := execGitCommandRun(worktreePath, "branch", "--set-upstream-to", remoteBranch, branchName); err != nil {
			return fmt.Errorf("failed to set up tracking: %w", err)
		}

		return nil
	} else {
		// No remote tracking branch, create worktree normally
		if err := execGitCommandRun(gm.repoPath, "worktree", "add", worktreePath, branchName); err != nil {
			return enhanceGitError(err, "worktree add")
		}
	}

	return nil
}

func (gm *GitManager) MoveWorktree(sourceWorktreePath, targetWorktreePath string) error {
	if err := execGitCommandRun(gm.repoPath, "worktree", "move", sourceWorktreePath, targetWorktreePath); err != nil {
		return enhanceGitError(err, "worktree move")
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

func (gm *GitManager) PromoteWorktree(sourceWorktreePath, targetWorktreePath string) error {
	// First, remove the target worktree
	if err := gm.RemoveWorktree(targetWorktreePath); err != nil {
		return fmt.Errorf("failed to remove target worktree: %w", err)
	}

	// Then, move the source worktree to the target location
	if err := gm.MoveWorktree(sourceWorktreePath, targetWorktreePath); err != nil {
		return fmt.Errorf("failed to move worktree from %s to %s: %w", sourceWorktreePath, targetWorktreePath, err)
	}

	return nil
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
	output, err := ExecGitCommand(worktreePath, "status", "--porcelain", "--ahead-behind")
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	statusLines := strings.SplitSeq(string(output), "\n")
	for line := range statusLines {
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
	output, err = ExecGitCommand(worktreePath, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err == nil {
		parts := strings.Fields(string(output))
		if len(parts) == 2 {
			if _, err := fmt.Sscanf(parts[0], "%d", &status.Ahead); err == nil {
				_, _ = fmt.Sscanf(parts[1], "%d", &status.Behind)
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

func (gm *GitManager) GetCurrentBranch() (string, error) {
	output, err := ExecGitCommand(gm.repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranchInPath gets the current branch name from a specified directory path
func (gm *GitManager) GetCurrentBranchInPath(path string) (string, error) {
	output, err := ExecGitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", enhanceGitError(err, "get current branch")
	}

	return strings.TrimSpace(string(output)), nil
}

func (gm *GitManager) GetDefaultBranch() (string, error) {
	// Try to get the default branch from remote HEAD
	output, err := ExecGitCommand(gm.repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		// Parse refs/remotes/origin/main -> main
		defaultRef := strings.TrimSpace(string(output))
		if strings.HasPrefix(defaultRef, "refs/remotes/origin/") {
			return strings.TrimPrefix(defaultRef, "refs/remotes/origin/"), nil
		}
	}

	// Fallback: try common default branch names
	commonDefaults := []string{"main", "master", "develop"}
	for _, branch := range commonDefaults {
		exists, err := gm.BranchExists(branch)
		if err == nil && exists {
			return branch, nil
		}
	}

	// Last resort: get the current branch
	return gm.GetCurrentBranch()
}

func (gm *GitManager) GetRemoteBranches() ([]string, error) {
	output, err := ExecGitCommand(gm.repoPath, "branch", "-r")
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	var branches []string
	lines := strings.SplitSeq(string(output), "\n")
	for line := range lines {
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
	output, err := ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	currentBranch := strings.TrimSpace(string(output))

	// Check if upstream is set
	_, err = ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")

	var cmd *exec.Cmd
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

	output, err := ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	currentBranch := strings.TrimSpace(string(output))

	finalArgs := []string{"pull"}

	// Check if upstream is set
	_, err = ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")
	if err != nil {
		// No upstream set, try to set it and pull
		remoteBranch := Remote(currentBranch)

		// Check if remote branch exists
		_, err = ExecGitCommand(worktreePath, "rev-parse", "--verify", remoteBranch)
		if err == nil {
			// Remote branch exists, set upstream and pull
			_, err = ExecGitCommand(worktreePath, "branch", "--set-upstream-to", remoteBranch)
			if err != nil {
				return fmt.Errorf("failed to set upstream: %w", err)
			}
		} else {
			// Remote branch doesn't exist, try to pull with explicit remote and branch
			finalArgs = append(finalArgs, "origin", currentBranch)
		}
	}

	return ExecGitCommandInteractive(worktreePath, finalArgs...)
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

// WorktreeInfoData represents comprehensive information about a worktree
type WorktreeInfoData struct {
	Name          string
	Path          string
	Branch        string
	CreatedAt     time.Time
	GitStatus     *GitStatus
	BaseInfo      *BranchInfo
	Commits       []CommitInfo
	ModifiedFiles []FileChange
	JiraTicket    *JiraTicketDetails
}

// BranchInfo represents information about the base branch
type BranchInfo struct {
	Name       string
	DivergedAt string
	DaysAgo    int
	Upstream   string
	AheadBy    int
	BehindBy   int
}

// CommitInfo represents information about a commit
type CommitInfo struct {
	Hash      string
	Message   string
	Author    string
	Timestamp time.Time
}

// FileChange represents a modified file
type FileChange struct {
	Path      string
	Status    string
	Additions int
	Deletions int
}

// JiraTicketDetails represents detailed JIRA ticket information
type JiraTicketDetails struct {
	Key           string
	Summary       string
	Status        string
	Assignee      string
	Priority      string
	Reporter      string
	Created       time.Time
	DueDate       *time.Time
	Epic          string
	URL           string
	LatestComment *Comment
}

// Comment represents a JIRA comment
type Comment struct {
	Author    string
	Content   string
	Timestamp time.Time
}

// RecentActivity represents recent git activity that might necessitate a mergeback
type RecentActivity struct {
	Type          string // "hotfix", "merge", "feature"
	WorktreeName  string
	BranchName    string
	SourceBranch  string // For merges, what was merged
	TargetBranch  string // For merges, what it was merged into
	CommitHash    string
	CommitMessage string
	Author        string
	Timestamp     time.Time
	JiraTicket    string // Extracted JIRA ticket if found
}

// GetRecentMergeableActivity analyzes recent git history to find hotfixes or merges
// that might necessitate a mergeback operation
func (gm *GitManager) GetRecentMergeableActivity(maxDays int) ([]RecentActivity, error) {
	if maxDays <= 0 {
		maxDays = 7 // Default to last 7 days
	}

	var activities []RecentActivity

	// Get recent commits that might indicate hotfix/merge activity
	since := fmt.Sprintf("--since=%d.days.ago", maxDays)

	// Look for merge commits first
	mergeCommits, err := gm.getRecentMergeCommits(since)
	if err == nil {
		activities = append(activities, mergeCommits...)
	}

	// Look for hotfix branches that were recently created or merged
	hotfixCommits, err := gm.getRecentHotfixActivity(since)
	if err == nil {
		activities = append(activities, hotfixCommits...)
	}

	// Note: Removed feature branch detection per user request
	// Only consider hotfix and merge commits for auto-detection

	return activities, nil
}

// getRecentMergeCommits finds recent merge commits
func (gm *GitManager) getRecentMergeCommits(since string) ([]RecentActivity, error) {
	var activities []RecentActivity

	// Get merge commits with format: hash|author|date|message
	output, err := ExecGitCommand(gm.repoPath, "log", "--merges", since, "--pretty=format:%H|%an|%at|%s")
	if err != nil {
		return activities, err
	}

	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		hash := parts[0]
		author := parts[1]
		timestampStr := parts[2]
		message := parts[3]

		timestamp, err := parseTimestamp(timestampStr)
		if err != nil {
			continue
		}

		activity := RecentActivity{
			Type:          "merge",
			CommitHash:    hash,
			CommitMessage: message,
			Author:        author,
			Timestamp:     timestamp,
			JiraTicket:    ExtractJiraTicket(message),
		}

		// Try to extract source and target branches from merge commit
		sourceBranch, targetBranch := gm.extractMergeBranches(hash)
		activity.SourceBranch = sourceBranch
		activity.TargetBranch = targetBranch

		// Generate worktree name from branch or JIRA ticket
		if activity.JiraTicket != "" {
			activity.WorktreeName = activity.JiraTicket
		} else if sourceBranch != "" {
			activity.WorktreeName = ExtractWorktreeNameFromBranch(sourceBranch)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// getRecentHotfixActivity finds recent hotfix branch activity
func (gm *GitManager) getRecentHotfixActivity(since string) ([]RecentActivity, error) {
	var activities []RecentActivity

	// Get commits on hotfix branches
	output, err := ExecGitCommand(gm.repoPath, "log", "--all", since, "--pretty=format:%H|%an|%at|%s|%D", "--grep=hotfix")
	if err != nil {
		return activities, err
	}

	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 4 {
			continue
		}

		hash := parts[0]
		author := parts[1]
		timestampStr := parts[2]
		message := parts[3]
		refs := ""
		if len(parts) > 4 {
			refs = parts[4]
		}

		timestamp, err := parseTimestamp(timestampStr)
		if err != nil {
			continue
		}

		// Look for hotfix patterns in refs or message
		if !strings.Contains(refs, "hotfix") && !strings.Contains(strings.ToLower(message), "hotfix") {
			continue
		}

		activity := RecentActivity{
			Type:          "hotfix",
			CommitHash:    hash,
			CommitMessage: message,
			Author:        author,
			Timestamp:     timestamp,
			JiraTicket:    ExtractJiraTicket(message),
		}

		// Extract branch name from refs
		if refs != "" {
			for ref := range strings.SplitSeq(refs, ", ") {
				if strings.Contains(ref, "hotfix/") {
					activity.BranchName = extractBranchFromRef(ref)
					activity.WorktreeName = ExtractWorktreeNameFromBranch(activity.BranchName)
					break
				}
			}
		}

		// Fallback: extract from JIRA ticket or commit message
		if activity.WorktreeName == "" {
			if activity.JiraTicket != "" {
				activity.WorktreeName = activity.JiraTicket
			} else {
				activity.WorktreeName = ExtractWorktreeNameFromMessage(message)
			}
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// Helper functions for parsing git data
func parseTimestamp(timestampStr string) (time.Time, error) {
	// Parse as Unix timestamp
	var unixTime int64
	if _, err := fmt.Sscanf(timestampStr, "%d", &unixTime); err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp format: %w", err)
	}
	return time.Unix(unixTime, 0), nil
}

func ExtractJiraTicket(message string) string {
	// Common JIRA patterns: PROJECT-123, ABC-456, etc.
	jiraPattern := `[A-Z]{2,}-\d+`
	re := regexp.MustCompile(jiraPattern)
	match := re.FindString(message)
	return match
}

func ExtractWorktreeNameFromBranch(branchName string) string {
	// Extract meaningful part from branch names like:
	// hotfix/PROJECT-123_fix_auth -> PROJECT-123
	// feature/PROJECT-456_new_api -> PROJECT-456
	// hotfix/critical-bug -> critical-bug

	if branchName == "" {
		return ""
	}

	// Remove common prefixes
	prefixes := []string{"hotfix/", "feature/", "bugfix/", "merge/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(branchName, prefix) {
			branchName = strings.TrimPrefix(branchName, prefix)
			break
		}
	}

	// If it looks like a JIRA ticket, extract just that
	if jira := ExtractJiraTicket(branchName); jira != "" {
		return jira
	}

	// Otherwise use the cleaned branch name
	return branchName
}

func ExtractWorktreeNameFromMessage(message string) string {
	// Try to extract JIRA ticket first
	if jira := ExtractJiraTicket(message); jira != "" {
		return jira
	}

	// Fallback: use first few words of commit message
	words := strings.Fields(strings.ToLower(message))
	if len(words) > 0 {
		// Remove common commit prefixes
		commonPrefixes := []string{"feat:", "fix:", "hotfix:", "merge:", "add:", "update:"}
		firstWord := words[0]
		for _, prefix := range commonPrefixes {
			if firstWord == prefix && len(words) > 1 {
				firstWord = words[1]
				break
			}
		}
		return firstWord
	}

	return "unknown"
}

func extractBranchFromRef(ref string) string {
	// Parse refs like "origin/hotfix/PROJECT-123" -> "hotfix/PROJECT-123"
	parts := strings.Split(ref, "/")
	if len(parts) >= 2 {
		// Skip "origin" if present
		if parts[0] == "origin" {
			return strings.Join(parts[1:], "/")
		}
		return ref
	}
	return ref
}

func (gm *GitManager) extractMergeBranches(commitHash string) (string, string) {
	// Get the merge commit details to extract source and target branches
	output, err := ExecGitCommand(gm.repoPath, "show", "--format=%P %s", "--no-patch", commitHash)
	if err != nil {
		return "", ""
	}

	line := strings.TrimSpace(string(output))
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return "", ""
	}

	// For merge commits, try to extract from commit message
	message := parts[1]

	// Look for patterns like "Merge branch 'feature/xyz' into main"
	mergePattern := `Merge branch '([^']+)' into (.+)`
	re := regexp.MustCompile(mergePattern)
	matches := re.FindStringSubmatch(message)
	if len(matches) >= 3 {
		return matches[1], matches[2] // source, target
	}

	return "", ""
}
