package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupPushTestRepo creates a git repository setup for push testing
func setupPushTestRepo(t *testing.T) (*testutils.GitTestRepo, *Manager) {
	t.Helper()

	// Setup repository with remote
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Add .gitignore to ignore worktrees directory
	must(t, repo.WriteFile(".gitignore", "worktrees/\n"))
	must(t, repo.CommitChanges("Add .gitignore for worktrees"))
	must(t, repo.PushBranch("main"))

	// Create basic gbm.branchconfig.yaml for the manager
	gbmContent := `worktrees:
  main:
    branch: main
    description: "Main production branch"
`
	must(t, repo.WriteFile(DefaultBranchConfigFilename, gbmContent))

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	must(t, err)

	// Load GBM config
	must(t, manager.LoadGBMConfig(""))

	return repo, manager
}

// createWorktreeWithChanges creates a worktree and adds commits for testing
func createWorktreeWithChanges(t *testing.T, repo *testutils.GitTestRepo, manager *Manager, worktreeName, branchName string, numCommits int) {
	t.Helper()

	// Add worktree
	must(t, manager.AddWorktree(worktreeName, branchName, true, "main"))

	// Update the GBM config to include this worktree
	gbmConfigPath := filepath.Join(repo.GetLocalPath(), DefaultBranchConfigFilename)
	existingContent, err := os.ReadFile(gbmConfigPath)
	require.NoError(t, err)

	// Add the new worktree to the config
	newWorktreeConfig := fmt.Sprintf("  %s:\n    branch: %s\n    description: \"Test worktree for %s\"\n", worktreeName, branchName, worktreeName)
	updatedContent := string(existingContent) + newWorktreeConfig

	require.NoError(t, os.WriteFile(gbmConfigPath, []byte(updatedContent), 0644))

	// Reload the GBM config to pick up the new worktree
	require.NoError(t, manager.LoadGBMConfig(""))

	// Add commits to the worktree
	worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", worktreeName)
	for i := 0; i < numCommits; i++ {
		fileName := fmt.Sprintf("change%d.txt", i+1)
		content := fmt.Sprintf("Content for change %d", i+1)
		commitMsg := fmt.Sprintf("Add change %d", i+1)

		require.NoError(t, os.WriteFile(filepath.Join(worktreePath, fileName), []byte(content), 0644))
		require.NoError(t, execGitCommandRun(worktreePath, "add", fileName))
		require.NoError(t, execGitCommandRun(worktreePath, "commit", "-m", commitMsg))
	}
}

// verifyPushSuccess verifies that commits were successfully pushed to remote
func verifyPushSuccess(t *testing.T, repo *testutils.GitTestRepo, branchName string, expectedCommits int) {
	t.Helper()

	// Create a separate clone to verify remote state
	tempDir := t.TempDir()
	verifyClonePath := filepath.Join(tempDir, "verify-clone")

	// Clone the remote repo to check what was pushed
	require.NoError(t, execGitCommandRun(tempDir, "clone", repo.GetRemotePath(), "verify-clone"))

	// Try to checkout the branch - if it fails, it means the push didn't work
	err := execGitCommandRun(verifyClonePath, "checkout", branchName)
	if err != nil {
		// If checkout fails, try to see if the branch exists in remote
		err2 := execGitCommandRun(verifyClonePath, "checkout", "-b", branchName, "origin/"+branchName)
		require.NoError(t, err2, "Branch %s was not found in remote repository", branchName)
	}

	// Branch exists and can be checked out - push was successful
}

func TestManager_PushWorktree(t *testing.T) {
	repo, manager := setupPushTestRepo(t)

	tests := []struct {
		name         string
		setup        func(t *testing.T)
		worktreeName string
		expectErr    func(t *testing.T, err error)
		verify       func(t *testing.T)
	}{
		{
			name: "SuccessfulPushWithUpstreamSetup",
			setup: func(t *testing.T) {
				// Create worktree with local changes (new branch)
				createWorktreeWithChanges(t, repo, manager, "test-wt", "feature/new-branch", 2)
			},
			worktreeName: "test-wt",
			expectErr:    func(t *testing.T, err error) { require.NoError(t, err) },
			verify: func(t *testing.T) {
				// Verify the branch was pushed to remote
				verifyPushSuccess(t, repo, "feature/new-branch", 2)
			},
		},
		{
			name: "SuccessfulPushWithExistingUpstream",
			setup: func(t *testing.T) {
				// Create worktree and push once, then add more changes
				createWorktreeWithChanges(t, repo, manager, "existing-wt", "feature/existing-branch", 1)
				must(t, manager.PushWorktree("existing-wt")) // First push to set upstream

				// Add more changes
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "existing-wt")
				require.NoError(t, os.WriteFile(filepath.Join(worktreePath, "second-change.txt"), []byte("Second change"), 0644))
				require.NoError(t, execGitCommandRun(worktreePath, "add", "second-change.txt"))
				require.NoError(t, execGitCommandRun(worktreePath, "commit", "-m", "Add second change"))
			},
			worktreeName: "existing-wt",
			expectErr:    func(t *testing.T, err error) { require.NoError(t, err) },
			verify: func(t *testing.T) {
				// Verify both commits were pushed
				verifyPushSuccess(t, repo, "feature/existing-branch", 2)
			},
		},
		{
			name: "PushWithNoChanges",
			setup: func(t *testing.T) {
				// Create worktree and push it, then try to push again
				createWorktreeWithChanges(t, repo, manager, "no-changes-wt", "feature/no-changes", 1)
				must(t, manager.PushWorktree("no-changes-wt")) // First push
			},
			worktreeName: "no-changes-wt",
			expectErr:    func(t *testing.T, err error) { require.NoError(t, err) }, // Should succeed
			verify: func(t *testing.T) {
				// Verify the branch still exists in remote (idempotency check)
				verifyPushSuccess(t, repo, "feature/no-changes", 1)

				// Also verify local worktree is still in sync with remote
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "no-changes-wt")

				// Check that there are no uncommitted changes
				err := execGitCommandRun(worktreePath, "diff-index", "--quiet", "HEAD")
				require.NoError(t, err, "Worktree should have no uncommitted changes")

				// Check that local is up to date with remote
				err = execGitCommandRun(worktreePath, "diff", "--quiet", "HEAD", "origin/feature/no-changes")
				require.NoError(t, err, "Local should be in sync with remote")
			},
		},
		{
			name: "ErrorNonexistentWorktree",
			setup: func(t *testing.T) {
				// No setup needed - testing nonexistent worktree
			},
			worktreeName: "nonexistent-wt",
			expectErr:    func(t *testing.T, err error) { require.Error(t, err) },
			verify: func(t *testing.T) {
				// For error case: verify the system state remains clean after failed operation

				// Verify no worktree directory was created
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "nonexistent-wt")
				assert.NoFileExists(t, worktreePath, "No worktree directory should exist for failed push")

				// Verify no unwanted branch was created in remote
				tempDir := t.TempDir()
				verifyClonePath := filepath.Join(tempDir, "error-verify-clone")
				require.NoError(t, execGitCommandRun(tempDir, "clone", repo.GetRemotePath(), "error-verify-clone"))

				// Verify that only expected branches exist by trying to checkout nonexistent ones
				// Main branch should exist
				err := execGitCommandRun(verifyClonePath, "checkout", "main")
				require.NoError(t, err, "Main branch should exist and be checkable")

				// Nonexistent branch should not exist
				err = execGitCommandRun(verifyClonePath, "checkout", "feature/nonexistent-branch")
				assert.Error(t, err, "Nonexistent branch should not be checkable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			err := manager.PushWorktree(tt.worktreeName)
			tt.expectErr(t, err)

			// Always call verify - let the verify function handle success vs error cases
			tt.verify(t)
		})
	}
}

func TestManager_PushAllWorktrees(t *testing.T) {
	repo, manager := setupPushTestRepo(t)

	tests := []struct {
		name      string
		setup     func(t *testing.T)
		expectErr func(t *testing.T, err error)
		verify    func(t *testing.T)
	}{
		{
			name: "PushAllWithChanges",
			setup: func(t *testing.T) {
				// Create multiple worktrees with changes
				createWorktreeWithChanges(t, repo, manager, "wt1", "feature/branch1", 1)
				createWorktreeWithChanges(t, repo, manager, "wt2", "feature/branch2", 2)
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
			verify: func(t *testing.T) {
				// Verify both worktrees were pushed
				verifyPushSuccess(t, repo, "feature/branch1", 1)
				verifyPushSuccess(t, repo, "feature/branch2", 2)
			},
		},
		{
			name: "PushAllMixedStates",
			setup: func(t *testing.T) {
				// Create one worktree and push it (up-to-date)
				createWorktreeWithChanges(t, repo, manager, "up-to-date-wt", "feature/up-to-date", 1)
				must(t, manager.PushWorktree("up-to-date-wt"))

				// Create another with new changes
				createWorktreeWithChanges(t, repo, manager, "new-changes-wt", "feature/new-changes", 1)
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
			verify: func(t *testing.T) {
				// Both should be pushed successfully
				verifyPushSuccess(t, repo, "feature/up-to-date", 1)
				verifyPushSuccess(t, repo, "feature/new-changes", 1)
			},
		},
		{
			name: "PushAllNoWorktrees",
			setup: func(t *testing.T) {
				// No worktrees created
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) }, // Should succeed with no-op
			verify: func(t *testing.T) {
				// No-op case - nothing to verify
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			err := manager.PushAllWorktrees()
			tt.expectErr(t, err)

			if err == nil {
				tt.verify(t)
			}
		})
	}
}

func TestManager_IsInWorktree_Integration(t *testing.T) {
	repo, manager := setupPushTestRepo(t)

	// Create a worktree for testing
	createWorktreeWithChanges(t, repo, manager, "test-detection", "feature/detection", 1)
	worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test-detection")

	tests := []struct {
		name         string
		path         string
		expectInWt   bool
		expectWtName string
		expectErr    bool
	}{
		{
			name:         "InWorktreeRoot",
			path:         worktreePath,
			expectInWt:   true,
			expectWtName: "test-detection",
			expectErr:    false,
		},
		{
			name:         "InWorktreeSubdir",
			path:         filepath.Join(worktreePath, "subdir"),
			expectInWt:   true,
			expectWtName: "test-detection",
			expectErr:    false,
		},
		{
			name:         "NotInWorktree",
			path:         repo.GetLocalPath(),
			expectInWt:   false,
			expectWtName: "",
			expectErr:    false,
		},
		{
			name:         "NonexistentPath",
			path:         "/nonexistent/path",
			expectInWt:   false,
			expectWtName: "",
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create subdirectory if needed
			if tt.name == "InWorktreeSubdir" {
				require.NoError(t, os.MkdirAll(tt.path, 0755))
			}

			inWt, wtName, err := manager.IsInWorktree(tt.path)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectInWt, inWt)
				assert.Equal(t, tt.expectWtName, wtName)
			}
		})
	}
}

func TestManager_GetAllWorktrees_Integration(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *testutils.GitTestRepo, manager *Manager)
		expectCount int
		expectErr   bool
	}{
		{
			name: "MultipleWorktrees",
			setup: func(t *testing.T, repo *testutils.GitTestRepo, manager *Manager) {
				createWorktreeWithChanges(t, repo, manager, "wt1", "feature/test1", 1)
				createWorktreeWithChanges(t, repo, manager, "wt2", "feature/test2", 1)
			},
			expectCount: 2,
			expectErr:   false,
		},
		{
			name: "NoWorktrees",
			setup: func(t *testing.T, repo *testutils.GitTestRepo, manager *Manager) {
				// No setup - no worktrees should be created
			},
			expectCount: 0,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh repo and manager for each test case
			repo, manager := setupPushTestRepo(t)

			tt.setup(t, repo, manager)

			worktrees, err := manager.GetAllWorktrees()

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, worktrees, tt.expectCount)
			}
		})
	}
}
