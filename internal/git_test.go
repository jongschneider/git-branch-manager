package internal

import (
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitManager_GetCurrentBranchInPath(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name           string
		setup          func(t *testing.T, repo *testutils.GitTestRepo)
		testPath       func(repo *testutils.GitTestRepo) string
		expectedBranch string
		expectError    bool
	}{
		{
			name: "get current branch from main repo",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Ensure we're on main branch
				must(t, repo.InLocalRepo(func() error {
					return execGitCommandRun(repo.GetLocalPath(), "checkout", "main")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedBranch: "main",
			expectError:    false,
		},
		{
			name: "get current branch from feature branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create and checkout a feature branch
				must(t, repo.InLocalRepo(func() error {
					return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/test-branch")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedBranch: "feature/test-branch",
			expectError:    false,
		},
		{
			name: "get current branch from worktree",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a worktree with a specific branch
				must(t, gitManager.AddWorktree("test-worktree", "feature/worktree-branch", true, ""))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return filepath.Join(repo.GetLocalPath(), "worktrees", "test-worktree")
			},
			expectedBranch: "feature/worktree-branch",
			expectError:    false,
		},
		{
			name: "error when path is not a git repository",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return "/tmp/not-a-git-repo"
			},
			expectedBranch: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)

			testPath := tt.testPath(repo)
			branch, err := gitManager.GetCurrentBranchInPath(testPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, branch)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBranch, branch)
			}
		})
	}
}
