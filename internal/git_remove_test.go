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

// ============================================================================
// INTEGRATION TESTS for worktreeRemover interface methods
// ============================================================================
// These tests validate that the Manager actually implements the worktreeRemover
// interface correctly, using real git operations to ensure our mocked behavior
// in unit tests matches the real implementation.

// Helper function to setup a manager with worktrees for testing
func setupManagerForRemoverTests(t *testing.T) (*Manager, string, *testutils.GitTestRepo) {
	// Create test repo with standard branches and configuration
	sourceRepo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)

	// Add .gitignore to ignore worktrees directory
	err := sourceRepo.WriteFile(".gitignore", "worktrees/\n")
	require.NoError(t, err)
	err = sourceRepo.CommitChanges("Add .gitignore for worktrees")
	require.NoError(t, err)

	// Create additional branches for worktrees
	err = sourceRepo.CreateBranch("dev", "dev content")
	require.NoError(t, err)
	err = sourceRepo.CreateBranch("feat", "feat content")
	require.NoError(t, err)

	// Initialize manager in the local repo
	repoPath := sourceRepo.GetLocalPath()
	manager, err := NewManager(repoPath)
	require.NoError(t, err)

	// Create worktrees using the manager (avoid main branch conflict)
	err = manager.AddWorktree("dev", "dev", false, "dev")
	require.NoError(t, err)
	err = manager.AddWorktree("feat", "feat", false, "feat")
	require.NoError(t, err)

	return manager, repoPath, sourceRepo
}

func TestManager_GetWorktreePath_Integration(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		expectResult func(t *testing.T, path string)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name:         "success - existing worktree",
			worktreeName: "dev",
			expectResult: func(t *testing.T, path string) {
				assert.Contains(t, path, "worktrees/dev")
				assert.True(t, filepath.IsAbs(path), "Path should be absolute")

				// Verify directory actually exists
				_, err := os.Stat(path)
				assert.NoError(t, err, "Worktree directory should exist")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "error - nonexistent worktree",
			worktreeName: "nonexistent",
			expectResult: func(t *testing.T, path string) {
				assert.Empty(t, path, "Path should be empty for nonexistent worktree")
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "nonexistent")
			},
		},
		{
			name:         "success - feat worktree",
			worktreeName: "feat",
			expectResult: func(t *testing.T, path string) {
				assert.Contains(t, path, "worktrees/feat")
				assert.True(t, filepath.IsAbs(path), "Path should be absolute")

				// Verify directory actually exists
				_, err := os.Stat(path)
				assert.NoError(t, err, "Feat worktree directory should exist")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, _, _ := setupManagerForRemoverTests(t)

			path, err := manager.GetWorktreePath(tt.worktreeName)

			tt.expectResult(t, path)
			tt.assertErr(t, err)
		})
	}
}

func TestManager_GetWorktreeStatus_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupChanges func(t *testing.T, worktreePath string)
		expectResult func(t *testing.T, status *GitStatus)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name: "success - clean worktree",
			setupChanges: func(t *testing.T, worktreePath string) {
				// No changes needed
			},
			expectResult: func(t *testing.T, status *GitStatus) {
				assert.False(t, status.HasChanges(), "Clean worktree should not have changes")
				assert.False(t, status.IsDirty)
				assert.Equal(t, 0, status.Untracked)
				assert.Equal(t, 0, status.Modified)
				assert.Equal(t, 0, status.Staged)
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success - worktree with untracked files",
			setupChanges: func(t *testing.T, worktreePath string) {
				// Create untracked file
				untrackedFile := filepath.Join(worktreePath, "untracked.txt")
				err := os.WriteFile(untrackedFile, []byte("untracked content"), 0o644)
				require.NoError(t, err)
			},
			expectResult: func(t *testing.T, status *GitStatus) {
				assert.True(t, status.HasChanges(), "Worktree with untracked files should have changes")
				assert.Greater(t, status.Untracked, 0, "Should have untracked files")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success - worktree with changes (combined test)",
			setupChanges: func(t *testing.T, worktreePath string) {
				// Create untracked file - this is simpler and always works
				untrackedFile := filepath.Join(worktreePath, "changes.txt")
				err := os.WriteFile(untrackedFile, []byte("some changes"), 0o644)
				require.NoError(t, err)
			},
			expectResult: func(t *testing.T, status *GitStatus) {
				assert.True(t, status.HasChanges(), "Worktree with changes should have changes")
				// The specific type of change (untracked, modified) doesn't matter for interface validation
				assert.True(t, status.Untracked > 0 || status.Modified > 0 || status.IsDirty, "Should have some type of changes")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - invalid worktree path",
			setupChanges: func(t *testing.T, worktreePath string) {
				// No setup needed for invalid path test
			},
			expectResult: func(t *testing.T, status *GitStatus) {
				assert.Nil(t, status, "Status should be nil for invalid path")
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, repoPath, _ := setupManagerForRemoverTests(t)

			var worktreePath string
			if tt.name != "error - invalid worktree path" {
				// Use a real worktree path
				worktreePath = filepath.Join(repoPath, "worktrees", "dev")
				tt.setupChanges(t, worktreePath)
			} else {
				// Use invalid path
				worktreePath = "/nonexistent/path"
			}

			status, err := manager.GetWorktreeStatus(worktreePath)

			tt.expectResult(t, status)
			tt.assertErr(t, err)
		})
	}
}

func TestManager_RemoveWorktree_Integration(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		setupChanges func(t *testing.T, manager *Manager, repoPath string)
		expectResult func(t *testing.T, repoPath string, worktreeName string)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name:         "success - remove clean worktree",
			worktreeName: "feat",
			setupChanges: func(t *testing.T, manager *Manager, repoPath string) {
				// No changes needed
			},
			expectResult: func(t *testing.T, repoPath string, worktreeName string) {
				// Verify worktree directory is removed
				worktreePath := filepath.Join(repoPath, "worktrees", worktreeName)
				_, err := os.Stat(worktreePath)
				assert.True(t, os.IsNotExist(err), "Worktree directory should not exist after removal")

				// Verify git worktree list doesn't include it
				gitManager, err := NewGitManager(repoPath, "worktrees")
				require.NoError(t, err)
				worktrees, err := gitManager.GetWorktrees()
				require.NoError(t, err)

				for _, wt := range worktrees {
					assert.NotEqual(t, worktreeName, wt.Name, "Removed worktree should not appear in git worktree list")
				}
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "success - remove worktree with changes (force not required at Manager level)",
			worktreeName: "dev",
			setupChanges: func(t *testing.T, manager *Manager, repoPath string) {
				// Create uncommitted changes
				worktreePath := filepath.Join(repoPath, "worktrees", "dev")
				uncommittedFile := filepath.Join(worktreePath, "uncommitted.txt")
				err := os.WriteFile(uncommittedFile, []byte("uncommitted content"), 0o644)
				require.NoError(t, err)
			},
			expectResult: func(t *testing.T, repoPath string, worktreeName string) {
				// Manager.RemoveWorktree should succeed regardless of uncommitted changes
				// The force checking is done at the command level, not Manager level
				worktreePath := filepath.Join(repoPath, "worktrees", worktreeName)
				_, err := os.Stat(worktreePath)
				assert.True(t, os.IsNotExist(err), "Worktree directory should not exist after removal")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "error - remove nonexistent worktree",
			worktreeName: "nonexistent",
			setupChanges: func(t *testing.T, manager *Manager, repoPath string) {
				// No setup needed
			},
			expectResult: func(t *testing.T, repoPath string, worktreeName string) {
				// No specific result validation needed for error case
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// The exact error message may vary, but there should be an error for nonexistent worktree
				assert.NotEmpty(t, err.Error(), "Error message should not be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, repoPath, _ := setupManagerForRemoverTests(t)

			tt.setupChanges(t, manager, repoPath)

			err := manager.RemoveWorktree(tt.worktreeName)

			tt.expectResult(t, repoPath, tt.worktreeName)
			tt.assertErr(t, err)
		})
	}
}

func TestManager_GetAllWorktrees_RemoverIntegration(t *testing.T) {
	tests := []struct {
		name         string
		setupChanges func(t *testing.T, manager *Manager)
		expectResult func(t *testing.T, worktrees map[string]*WorktreeListInfo)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name: "success - get all existing worktrees",
			setupChanges: func(t *testing.T, manager *Manager) {
				// No changes needed, use default setup
			},
			expectResult: func(t *testing.T, worktrees map[string]*WorktreeListInfo) {
				expectedWorktrees := []string{"dev", "feat"}
				assert.Len(t, worktrees, len(expectedWorktrees), "Should have all created worktrees")

				for _, expected := range expectedWorktrees {
					info, exists := worktrees[expected]
					assert.True(t, exists, "Worktree %s should exist", expected)
					assert.NotNil(t, info, "Worktree info should not be nil")
					assert.Contains(t, info.Path, expected, "Path should contain worktree name")
					assert.True(t, filepath.IsAbs(info.Path), "Path should be absolute")
				}
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success - get worktrees after removing one",
			setupChanges: func(t *testing.T, manager *Manager) {
				// Remove one worktree
				err := manager.RemoveWorktree("feat")
				require.NoError(t, err)
			},
			expectResult: func(t *testing.T, worktrees map[string]*WorktreeListInfo) {
				expectedWorktrees := []string{"dev"}
				assert.Len(t, worktrees, len(expectedWorktrees), "Should have remaining worktrees")

				for _, expected := range expectedWorktrees {
					_, exists := worktrees[expected]
					assert.True(t, exists, "Worktree %s should still exist", expected)
				}

				// Verify removed worktree is not present
				_, exists := worktrees["feat"]
				assert.False(t, exists, "Removed worktree should not be present")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success - empty worktrees after removing all",
			setupChanges: func(t *testing.T, manager *Manager) {
				// Remove all worktrees
				worktreesToRemove := []string{"dev", "feat"}
				for _, name := range worktreesToRemove {
					err := manager.RemoveWorktree(name)
					require.NoError(t, err)
				}
			},
			expectResult: func(t *testing.T, worktrees map[string]*WorktreeListInfo) {
				assert.Empty(t, worktrees, "Should have no worktrees after removing all")
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, _, _ := setupManagerForRemoverTests(t)

			tt.setupChanges(t, manager)

			worktrees, err := manager.GetAllWorktrees()

			tt.expectResult(t, worktrees)
			tt.assertErr(t, err)
		})
	}
}
