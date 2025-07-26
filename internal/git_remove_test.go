package internal

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// verifyWorktreeRemoved verifies that a worktree no longer exists
func verifyWorktreeRemoved(t *testing.T, gitManager *GitManager, worktreeName string) {
	t.Helper()
	worktrees, err := gitManager.GetWorktrees()
	require.NoError(t, err)

	for _, wt := range worktrees {
		if wt.Name == worktreeName {
			t.Errorf("worktree '%s' still exists after removal", worktreeName)
		}
	}
}

func TestGitManager_RemoveWorktree(t *testing.T) {
	// Setup shared repository and git manager
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)

	// Cleanup test infrastructure when test completes
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Add .gitignore to ignore worktrees directory
	must(t, repo.WriteFile(".gitignore", "worktrees/\n"))
	must(t, repo.CommitChanges("Add .gitignore for worktrees"))

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	must(t, err)

	tests := []struct {
		name         string
		setup        func(t *testing.T, repo *testutils.GitTestRepo)
		worktreePath string
		expect       func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string)
		expectErr    func(t *testing.T, err error)
	}{
		{
			name: "RemoveExistingWorktree",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a worktree to remove
				must(t, gitManager.AddWorktree("test-remove", "feature/test-remove", true, ""))
			},
			worktreePath: filepath.Join(repo.GetLocalPath(), "worktrees", "test-remove"),
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string) {
				// Verify worktree directory no longer exists
				assert.NoDirExists(t, worktreePath)

				// Verify worktree is no longer linked
				verifyWorktreeRemoved(t, gitManager, "test-remove")
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name: "RemoveWorktreeWithAbsolutePath",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a worktree to remove
				must(t, gitManager.AddWorktree("absolute-test", "feature/absolute-test", true, ""))
			},
			worktreePath: filepath.Join(repo.GetLocalPath(), "worktrees", "absolute-test"),
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string) {
				// Verify worktree directory no longer exists
				assert.NoDirExists(t, worktreePath)

				// Verify worktree is no longer linked
				verifyWorktreeRemoved(t, gitManager, "absolute-test")
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name: "RemoveWorktreeWithUncommittedChanges",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a worktree and add uncommitted changes
				must(t, gitManager.AddWorktree("dirty-test", "feature/dirty-test", true, ""))
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", "dirty-test")
				err := os.WriteFile(filepath.Join(worktreePath, "uncommitted.txt"), []byte("uncommitted content"), 0o644)
				must(t, err)
			},
			worktreePath: filepath.Join(repo.GetLocalPath(), "worktrees", "dirty-test"),
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string) {
				// Verify worktree directory no longer exists (--force should handle uncommitted changes)
				assert.NoDirExists(t, worktreePath)

				// Verify worktree is no longer linked
				verifyWorktreeRemoved(t, gitManager, "dirty-test")
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name: "ErrorNonexistentWorktree",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup - worktree doesn't exist
			},
			worktreePath: filepath.Join(repo.GetLocalPath(), "worktrees", "nonexistent"),
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string) {
				// Nothing to verify since we expect an error
			},
			expectErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "worktree remove")
			},
		},
		{
			name: "ErrorInvalidPath",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			worktreePath: "/path/that/does/not/exist",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreePath string) {
				// Nothing to verify since we expect an error
			},
			expectErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "worktree remove")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)
			err := gitManager.RemoveWorktree(tt.worktreePath)
			tt.expectErr(t, err)
			tt.expect(t, repo, gitManager, tt.worktreePath)
		})
	}
}
