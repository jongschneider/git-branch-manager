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

	cmd := exec.Command("git", "worktree", "add", worktreePath, branchName)
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
	if gitStatus == nil {
		return "❌"
	}

	var icons []string

	// Check ahead/behind first
	if gitStatus.Ahead > 0 && gitStatus.Behind > 0 {
		icons = append(icons, "⇕")
	} else if gitStatus.Ahead > 0 {
		icons = append(icons, "↑")
	} else if gitStatus.Behind > 0 {
		icons = append(icons, "↓")
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
		return "✅"
	}

	return strings.Join(icons, "")
}
