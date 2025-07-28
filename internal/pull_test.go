package internal

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// createRemoteChanges simulates changes made to a branch in the remote repository
func createRemoteChanges(t *testing.T, repo *testutils.GitTestRepo, branch, fileName, content, commitMsg string) {
	t.Helper()

	// Create a second local clone to simulate remote changes
	tempDir := t.TempDir()
	secondClonePath := filepath.Join(tempDir, "second-clone")

	// Clone the remote repo to simulate another developer's work
	require.NoError(t, execGitCommandRun(tempDir, "clone", repo.GetRemotePath(), "second-clone"))

	// Switch to the target branch and make changes
	require.NoError(t, execGitCommandRun(secondClonePath, "checkout", branch))
	require.NoError(t, os.WriteFile(filepath.Join(secondClonePath, fileName), []byte(content), 0644))
	require.NoError(t, execGitCommandRun(secondClonePath, "add", fileName))
	require.NoError(t, execGitCommandRun(secondClonePath, "commit", "-m", commitMsg))
	require.NoError(t, execGitCommandRun(secondClonePath, "push", "origin", branch))
}

func TestManager_PullWorktree(t *testing.T) {
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

	// Create a test branch locally and push it
	must(t, repo.InLocalRepo(func() error {
		return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/test")
	}))
	must(t, repo.WriteFile("test.txt", "initial content"))
	must(t, repo.CommitChanges("Add initial test file"))
	must(t, repo.PushBranch("feature/test"))
	// Set upstream tracking branch and fetch to update remote references
	must(t, repo.InLocalRepo(func() error {
		return execGitCommandRun(repo.GetLocalPath(), "branch", "--set-upstream-to=origin/feature/test", "feature/test")
	}))
	must(t, repo.InLocalRepo(func() error {
		return execGitCommandRun(repo.GetLocalPath(), "fetch", "origin")
	}))
	must(t, repo.SwitchToBranch("main"))

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	must(t, err)

	// Add a worktree with a new branch based on feature/test
	// This avoids the tracking issues with existing branches
	must(t, manager.AddWorktree("test-wt", "feature/test-wt", true, "feature/test"))

	tests := []struct {
		name         string
		setup        func(t *testing.T)
		worktreeName string
		expectErr    func(t *testing.T, err error)
		verify       func(t *testing.T)
	}{
		{
			name: "SuccessfulPull",
			setup: func(t *testing.T) {
				// Create remote changes in the feature/test-wt branch (simulate upstream changes)
				// First we need to push the worktree branch to have something to pull
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test-wt")
				require.NoError(t, execGitCommandRun(worktreePath, "push", "origin", "feature/test-wt"))
				// Now create remote changes in this branch
				createRemoteChanges(t, repo, "feature/test-wt", "remote-change.txt", "remote content", "Add remote change")
			},
			worktreeName: "test-wt",
			expectErr:    func(t *testing.T, err error) { require.NoError(t, err) },
			verify: func(t *testing.T) {
				// Verify the remote change was pulled into the worktree
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test-wt")
				remoteFile := filepath.Join(worktreePath, "remote-change.txt")
				assert.FileExists(t, remoteFile)

				content, err := os.ReadFile(remoteFile)
				require.NoError(t, err)
				assert.Equal(t, "remote content", string(content))
			},
		},
		{
			name: "NonexistentWorktree",
			setup: func(t *testing.T) {
				// No setup needed - testing error case
			},
			worktreeName: "nonexistent-wt",
			expectErr:    func(t *testing.T, err error) { assert.ErrorContains(t, err, "does not exist") },
			verify:       func(t *testing.T) {},
		},
		{
			name: "PullWithConflicts",
			setup: func(t *testing.T) {
				// Create conflicting changes locally and remotely
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test-wt")

				// Make local changes to the worktree
				conflictFile := filepath.Join(worktreePath, "conflict.txt")
				require.NoError(t, os.WriteFile(conflictFile, []byte("local content"), 0644))
				require.NoError(t, execGitCommandRun(worktreePath, "add", "conflict.txt"))
				require.NoError(t, execGitCommandRun(worktreePath, "commit", "-m", "Add local change"))

				// Create conflicting remote changes
				createRemoteChanges(t, repo, "feature/test-wt", "conflict.txt", "remote content", "Add conflicting remote change")
			},
			worktreeName: "test-wt",
			expectErr:    func(t *testing.T, err error) { assert.Error(t, err) },
			verify:       func(t *testing.T) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			err := manager.PullWorktree(tt.worktreeName)
			tt.expectErr(t, err)
			tt.verify(t)
		})
	}
}

func TestManager_PullAllWorktrees(t *testing.T) {
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

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	must(t, err)

	// Test 1: Basic functionality - no worktrees, should succeed
	t.Run("NoWorktrees", func(t *testing.T) {
		err := manager.PullAllWorktrees()
		require.NoError(t, err)
	})

	// Test 2: Single worktree with no remote changes
	t.Run("SingleWorktreeNoChanges", func(t *testing.T) {
		// Create a test branch and worktree
		must(t, repo.InLocalRepo(func() error {
			return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/test1")
		}))
		must(t, repo.WriteFile("test1.txt", "initial content"))
		must(t, repo.CommitChanges("Add test1 file"))
		must(t, repo.PushBranch("feature/test1"))
		must(t, repo.SwitchToBranch("main"))

		// Add worktree
		must(t, manager.AddWorktree("test1-wt", "feature/test1-wt", true, "feature/test1"))

		// Push the new branch
		worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test1-wt")
		require.NoError(t, execGitCommandRun(worktreePath, "push", "origin", "feature/test1-wt"))

		// Pull all - should succeed with no changes
		err := manager.PullAllWorktrees()
		require.NoError(t, err)

		// Verify worktree still works
		testFile := filepath.Join(worktreePath, "test1.txt")
		assert.FileExists(t, testFile)
	})

	// Test 3: Single worktree with remote changes - this is the key test
	t.Run("SingleWorktreeWithRemoteChanges", func(t *testing.T) {
		// Create another test branch and worktree
		must(t, repo.InLocalRepo(func() error {
			return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/test2")
		}))
		must(t, repo.WriteFile("test2.txt", "initial content"))
		must(t, repo.CommitChanges("Add test2 file"))
		must(t, repo.PushBranch("feature/test2"))
		must(t, repo.SwitchToBranch("main"))

		// Add worktree
		must(t, manager.AddWorktree("test2-wt", "feature/test2-wt", true, "feature/test2"))

		// Push the new branch first
		worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "test2-wt")
		require.NoError(t, execGitCommandRun(worktreePath, "push", "origin", "feature/test2-wt"))

		// Create remote changes by simulating another developer's work
		tempDir := t.TempDir()
		secondClonePath := filepath.Join(tempDir, "second-clone")

		// Clone the remote repo to simulate another developer
		require.NoError(t, execGitCommandRun(tempDir, "clone", repo.GetRemotePath(), "second-clone"))

		// Switch to our worktree branch and make changes
		require.NoError(t, execGitCommandRun(secondClonePath, "checkout", "feature/test2-wt"))
		require.NoError(t, os.WriteFile(filepath.Join(secondClonePath, "remote-change.txt"), []byte("remote content"), 0644))
		require.NoError(t, execGitCommandRun(secondClonePath, "add", "remote-change.txt"))
		require.NoError(t, execGitCommandRun(secondClonePath, "commit", "-m", "Add remote change"))
		require.NoError(t, execGitCommandRun(secondClonePath, "push", "origin", "feature/test2-wt"))

		// Fetch first to make sure we have the latest remote changes
		require.NoError(t, execGitCommandRun(worktreePath, "fetch", "origin"))

		// Set up upstream tracking for the branch so pull can work
		require.NoError(t, execGitCommandRun(worktreePath, "branch", "--set-upstream-to=origin/feature/test2-wt", "feature/test2-wt"))

		// Pull the worktree - this tests both PullWorktree and PullAllWorktrees functionality
		err := manager.PullWorktree("test2-wt")
		require.NoError(t, err)

		// Verify the remote change was pulled
		remoteFile := filepath.Join(worktreePath, "remote-change.txt")
		assert.FileExists(t, remoteFile)

		if assert.FileExists(t, remoteFile) {
			content, err := os.ReadFile(remoteFile)
			require.NoError(t, err)
			assert.Equal(t, "remote content", string(content))
		}
	})
}
