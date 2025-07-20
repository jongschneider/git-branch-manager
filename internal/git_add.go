package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (gm *GitManager) CreateBranch(branchName, baseBranch string) error {
	if baseBranch == "" {
		baseBranch = "HEAD"
	}

	if err := ExecGitCommandRun(gm.repoPath, "branch", branchName, baseBranch); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func (gm *GitManager) AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error {
	worktreeDir := filepath.Join(gm.repoPath, gm.worktreePrefix)
	worktreePath := filepath.Join(worktreeDir, worktreeName)

	// Check if worktree already exists
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree '%s' already exists", worktreeName)
	}

	// Create worktrees directory if it doesn't exist
	if err := os.MkdirAll(worktreeDir, 0o755); err != nil {
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
			// Branch exists, check if it's based on the correct base branch
			if baseBranch != "" {
				// Get the merge base between the existing branch and the base branch
				mergeBase, err := ExecGitCommand(gm.repoPath, "merge-base", branchName, baseBranch)
				if err != nil {
					return fmt.Errorf("failed to get merge base: %w", err)
				}

				// Get the commit hash of the base branch
				baseCommit, err := ExecGitCommand(gm.repoPath, "rev-parse", baseBranch)
				if err != nil {
					return fmt.Errorf("failed to get base branch commit: %w", err)
				}

				// Check if the existing branch is based on the correct base branch
				if strings.TrimSpace(string(mergeBase)) != strings.TrimSpace(string(baseCommit)) {
					return fmt.Errorf("branch '%s' exists but is not based on '%s'. Please delete the branch and try again, or use a different branch name", branchName, baseBranch)
				}
			}

			// Branch exists and is based on correct base (or no base specified), create worktree on existing branch
			cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
		} else {
			// Create new branch and worktree with base branch
			if baseBranch != "" {
				cmd = exec.Command("git", "worktree", "add", "-b", branchName, worktreePath, baseBranch)
			} else {
				// Use default behavior (current HEAD)
				cmd = exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
			}
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
		remoteBranch := Remote(branchName)
		checkCmd := exec.Command("git", "rev-parse", "--verify", remoteBranch)
		checkCmd.Dir = gm.repoPath
		_, err = execCommand(checkCmd)

		if err == nil {
			// Remote tracking branch exists, use --track but don't create new branch
			cmd = exec.Command("git", "worktree", "add", "--track", worktreePath, remoteBranch)
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