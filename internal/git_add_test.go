package internal

import (
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// must is a test helper that wraps error-returning functions and fails the test on error
func must(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

// verifyWorktreeLinked verifies that a worktree exists and is correctly linked to the expected branch
func verifyWorktreeLinked(t *testing.T, gitManager *GitManager, worktreeName, branchName string) {
	t.Helper()
	worktrees, err := gitManager.GetWorktrees()
	require.NoError(t, err)

	var foundWorktree *WorktreeInfo
	for _, wt := range worktrees {
		if wt.Name == worktreeName {
			foundWorktree = wt
			break
		}
	}
	require.NotNil(t, foundWorktree)
	assert.Equal(t, branchName, foundWorktree.Branch)
}

func TestGitManager_AddWorktree(t *testing.T) {
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
		worktreeName string
		branchName   string
		createBranch bool
		baseBranch   string
		expect       func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string)
		expectErr    func(t *testing.T, err error)
	}{
		{
			name:         "CreateWorktreeWithNewBranch",
			setup:        func(t *testing.T, repo *testutils.GitTestRepo) {},
			worktreeName: "feature-test",
			branchName:   "feature/new-feature",
			createBranch: true,
			baseBranch:   "",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string) {
				// Verify worktree directory exists
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", worktreeName)
				assert.DirExists(t, worktreePath)

				// Verify worktree is correctly linked
				verifyWorktreeLinked(t, gitManager, worktreeName, branchName)

				// Verify branch was created
				exists, err := gitManager.BranchExists(branchName)
				require.NoError(t, err)
				assert.True(t, exists)
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name: "CreateWorktreeFromExistingBranch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a local-only existing branch
				must(t, repo.InLocalRepo(func() error {
					return ExecGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "existing-branch")
				}))
				must(t, repo.WriteFile("test.txt", "test content"))
				must(t, repo.CommitChanges("Add test file"))
				must(t, repo.SwitchToBranch("main"))
			},
			worktreeName: "existing-test",
			branchName:   "existing-branch",
			createBranch: false,
			baseBranch:   "",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string) {
				// Verify worktree directory exists
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", worktreeName)
				assert.DirExists(t, worktreePath)

				// Verify worktree is correctly linked
				verifyWorktreeLinked(t, gitManager, worktreeName, branchName)

				// Verify test file exists (from the existing branch)
				assert.FileExists(t, filepath.Join(worktreePath, "test.txt"))
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name: "CreateWorktreeWithNewBranchFromBaseBranch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a base branch
				must(t, repo.CreateSynchronizedBranch("develop"))
				must(t, repo.WriteFile("develop.txt", "develop content"))
				must(t, repo.CommitChanges("Add develop file"))
				must(t, repo.SwitchToBranch("main"))
			},
			worktreeName: "hotfix-test",
			branchName:   "hotfix/urgent-fix",
			createBranch: true,
			baseBranch:   "develop",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string) {
				// Verify worktree directory exists
				worktreePath := filepath.Join(repo.GetLocalPath(), "worktrees", worktreeName)
				assert.DirExists(t, worktreePath)

				// Verify worktree is correctly linked
				verifyWorktreeLinked(t, gitManager, worktreeName, branchName)

				// Verify branch was created
				exists, err := gitManager.BranchExists(branchName)
				require.NoError(t, err)
				assert.True(t, exists)

				// Verify develop file exists (inherited from base branch)
				assert.FileExists(t, filepath.Join(worktreePath, "develop.txt"))
			},
			expectErr: func(t *testing.T, err error) { require.NoError(t, err) },
		},
		{
			name:         "ErrorNonexistentBranch",
			setup:        func(t *testing.T, repo *testutils.GitTestRepo) {},
			worktreeName: "error-test",
			branchName:   "nonexistent-branch",
			createBranch: false,
			baseBranch:   "",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string) {
			},
			expectErr: func(t *testing.T, err error) { assert.ErrorContains(t, err, "does not exist") },
		},
		{
			name: "ErrorDuplicateWorktreeName",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create a worktree first
				must(t, gitManager.AddWorktree("duplicate-test", "feature/first", true, ""))
			},
			worktreeName: "duplicate-test",
			branchName:   "feature/second",
			createBranch: true,
			baseBranch:   "",
			expect: func(t *testing.T, repo *testutils.GitTestRepo, gitManager *GitManager, worktreeName, branchName string) {
			},
			expectErr: func(t *testing.T, err error) { assert.ErrorContains(t, err, "already exists") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)
			err := gitManager.AddWorktree(tt.worktreeName, tt.branchName, tt.createBranch, tt.baseBranch)
			tt.expectErr(t, err)
			tt.expect(t, repo, gitManager, tt.worktreeName, tt.branchName)
		})
	}
}
