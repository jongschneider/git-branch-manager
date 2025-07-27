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

func TestGitManager_GetUpstreamBranch(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name             string
		setup            func(t *testing.T, repo *testutils.GitTestRepo)
		testPath         func(repo *testutils.GitTestRepo) string
		expectedUpstream string
		expectError      bool
	}{
		{
			name: "branch with upstream set",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create and checkout a branch, then set upstream
				must(t, repo.InLocalRepo(func() error {
					err := execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/test")
					if err != nil {
						return err
					}
					// Push to set up tracking
					return execGitCommandRun(repo.GetLocalPath(), "push", "-u", "origin", "feature/test")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedUpstream: "origin/feature/test",
			expectError:      false,
		},
		{
			name: "branch without upstream",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a local-only branch
				must(t, repo.InLocalRepo(func() error {
					return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "local-only")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedUpstream: "",
			expectError:      false,
		},
		{
			name: "main branch with upstream",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Ensure main branch has upstream tracking
				must(t, repo.InLocalRepo(func() error {
					err := execGitCommandRun(repo.GetLocalPath(), "checkout", "main")
					if err != nil {
						return err
					}
					return execGitCommandRun(repo.GetLocalPath(), "branch", "--set-upstream-to", "origin/main")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedUpstream: "origin/main",
			expectError:      false,
		},
		{
			name: "invalid path",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return "/nonexistent/path"
			},
			expectedUpstream: "",
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)

			testPath := tt.testPath(repo)
			upstream, err := gitManager.GetUpstreamBranch(testPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, upstream)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUpstream, upstream)
			}
		})
	}
}

func TestGitManager_GetAheadBehindCount(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name           string
		setup          func(t *testing.T, repo *testutils.GitTestRepo)
		testPath       func(repo *testutils.GitTestRepo) string
		expectedAhead  int
		expectedBehind int
		expectError    bool
	}{
		{
			name: "branch with upstream and commits ahead",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create and checkout a branch, set upstream, then add commits
				must(t, repo.InLocalRepo(func() error {
					err := execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/ahead")
					if err != nil {
						return err
					}
					// Push to set up tracking
					err = execGitCommandRun(repo.GetLocalPath(), "push", "-u", "origin", "feature/ahead")
					if err != nil {
						return err
					}
					// Add a local commit to be ahead
					err = repo.WriteFile("ahead.txt", "ahead content")
					if err != nil {
						return err
					}
					return execGitCommandRun(repo.GetLocalPath(), "add", "ahead.txt")
				}))
				must(t, repo.CommitChanges("Add ahead commit"))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedAhead:  1,
			expectedBehind: 0,
			expectError:    false,
		},
		{
			name: "branch without upstream",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a local-only branch
				must(t, repo.InLocalRepo(func() error {
					return execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "local-only")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedAhead:  0,
			expectedBehind: 0,
			expectError:    false,
		},
		{
			name: "branch with upstream and no divergence",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create and checkout a branch, set upstream with no extra commits
				must(t, repo.InLocalRepo(func() error {
					err := execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/even")
					if err != nil {
						return err
					}
					// Push to set up tracking
					return execGitCommandRun(repo.GetLocalPath(), "push", "-u", "origin", "feature/even")
				}))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			expectedAhead:  0,
			expectedBehind: 0,
			expectError:    false,
		},
		{
			name: "invalid path",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return "/nonexistent/path"
			},
			expectedAhead:  0,
			expectedBehind: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)

			testPath := tt.testPath(repo)
			ahead, behind, err := gitManager.GetAheadBehindCount(testPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAhead, ahead)
				assert.Equal(t, tt.expectedBehind, behind)
			}
		})
	}
}
