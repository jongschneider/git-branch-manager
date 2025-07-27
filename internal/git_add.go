package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (gm *GitManager) CreateBranch(branchName, baseBranch string) error {
	if baseBranch == "" {
		baseBranch = "HEAD"
	}

	if err := execGitCommandRun(gm.repoPath, "branch", branchName, baseBranch); err != nil {
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

	var finalArgs []string
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
				baseCommitHash, err := gm.GetCommitHash(baseBranch)
				if err != nil {
					return fmt.Errorf("failed to get base branch commit: %w", err)
				}
				baseCommit := []byte(baseCommitHash)

				// Check if the existing branch is based on the correct base branch
				if strings.TrimSpace(string(mergeBase)) != strings.TrimSpace(string(baseCommit)) {
					return fmt.Errorf("branch '%s' exists but is not based on '%s'. Please delete the branch and try again, or use a different branch name", branchName, baseBranch)
				}
			}

			// Branch exists and is based on correct base (or no base specified), create worktree on existing branch
			finalArgs = append(finalArgs, "worktree", "add", worktreePath, branchName)
		} else {
			// Create new branch and worktree with base branch
			if baseBranch != "" {
				finalArgs = append(finalArgs, "worktree", "add", "-b", branchName, worktreePath, baseBranch)
			} else {
				// Use default behavior (current HEAD)
				finalArgs = append(finalArgs, "worktree", "add", "-b", branchName, worktreePath)
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
		exists, err := gm.VerifyRef(remoteBranch)
		if err != nil {
			return fmt.Errorf("failed to verify remote branch: %w", err)
		}

		if exists {
			// Remote tracking branch exists, use --track but don't create new branch
			finalArgs = append(finalArgs, "worktree", "add", "--track", worktreePath, remoteBranch)
		} else {
			// No remote tracking branch, create worktree without tracking
			finalArgs = append(finalArgs, "worktree", "add", worktreePath, branchName)
		}
	}

	if err := execGitCommandRun(gm.repoPath, finalArgs...); err != nil {
		return fmt.Errorf("failed to add worktree (command: git %s): %w", strings.Join(finalArgs, " "), err)
	}

	return nil
}
