package internal

import (
	"path/filepath"
	"strings"
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

func TestGitManager_VerifyRef(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *testutils.GitTestRepo)
		ref         string
		expectExist bool
		expectError bool
	}{
		{
			name: "existing branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// main branch should exist by default
			},
			ref:         "main",
			expectExist: true,
			expectError: false,
		},
		{
			name: "non-existing branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			ref:         "nonexistent-branch",
			expectExist: false,
			expectError: false,
		},
		{
			name: "existing remote branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create and push a branch to create remote ref
				err := repo.InLocalRepo(func() error {
					err := execGitCommandRun(repo.GetLocalPath(), "checkout", "-b", "feature/remote-test")
					if err != nil {
						return err
					}
					return execGitCommandRun(repo.GetLocalPath(), "push", "-u", "origin", "feature/remote-test")
				})
				require.NoError(t, err)
			},
			ref:         "origin/feature/remote-test",
			expectExist: true,
			expectError: false,
		},
		{
			name: "existing commit hash",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Get the current commit hash - it should exist
			},
			ref:         "HEAD",
			expectExist: true,
			expectError: false,
		},
		{
			name: "invalid ref format",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			ref:         "invalid..ref",
			expectExist: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)

			exists, err := gitManager.VerifyRef(tt.ref)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectExist, exists)
			}
		})
	}
}

func TestGitManager_VerifyRefInPath(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *testutils.GitTestRepo)
		testPath    func(repo *testutils.GitTestRepo) string
		ref         string
		expectExist bool
		expectError bool
	}{
		{
			name: "existing branch in main repo",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// main branch should exist by default
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			ref:         "main",
			expectExist: true,
			expectError: false,
		},
		{
			name: "non-existing branch in main repo",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			ref:         "nonexistent-branch",
			expectExist: false,
			expectError: false,
		},
		{
			name: "invalid path",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// No setup needed
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return "/nonexistent/path"
			},
			ref:         "main",
			expectExist: false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, repo)

			testPath := tt.testPath(repo)
			exists, err := gitManager.VerifyRefInPath(testPath, tt.ref)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectExist, exists)
			}
		})
	}
}

func TestGitManager_GetCommitHash(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *testutils.GitTestRepo) string // returns expected hash
		ref         string
		expectError bool
	}{
		{
			name: "HEAD reference",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				// Get the current HEAD commit hash to compare with
				output, err := ExecGitCommand(repo.GetLocalPath(), "rev-parse", "HEAD")
				require.NoError(t, err)
				return strings.TrimSpace(string(output))
			},
			ref:         "HEAD",
			expectError: false,
		},
		{
			name: "main branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				// main branch should point to same commit as HEAD initially
				output, err := ExecGitCommand(repo.GetLocalPath(), "rev-parse", "main")
				require.NoError(t, err)
				return strings.TrimSpace(string(output))
			},
			ref:         "main",
			expectError: false,
		},
		{
			name: "non-existent branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				return "" // No expected hash for error case
			},
			ref:         "non-existent-branch",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHash := tt.setup(t, repo)

			hash, err := gitManager.GetCommitHash(tt.ref)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedHash, hash)
				assert.Len(t, hash, 40) // Full SHA-1 hash should be 40 characters
			}
		})
	}
}

func TestGitManager_GetCommitHashInPath(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
		testutils.WithRemoteName("origin"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	tests := []struct {
		name        string
		setup       func(t *testing.T, repo *testutils.GitTestRepo) string // returns expected hash
		testPath    func(repo *testutils.GitTestRepo) string
		ref         string
		expectError bool
	}{
		{
			name: "HEAD from repository root",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				output, err := ExecGitCommand(repo.GetLocalPath(), "rev-parse", "HEAD")
				require.NoError(t, err)
				return strings.TrimSpace(string(output))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			ref:         "HEAD",
			expectError: false,
		},
		{
			name: "main branch from repository root",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				output, err := ExecGitCommand(repo.GetLocalPath(), "rev-parse", "main")
				require.NoError(t, err)
				return strings.TrimSpace(string(output))
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			ref:         "main",
			expectError: false,
		},
		{
			name: "non-existent branch",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				return ""
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return repo.GetLocalPath()
			},
			ref:         "non-existent-branch",
			expectError: true,
		},
		{
			name: "invalid path",
			setup: func(t *testing.T, repo *testutils.GitTestRepo) string {
				return ""
			},
			testPath: func(repo *testutils.GitTestRepo) string {
				return "/non/existent/path"
			},
			ref:         "HEAD",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHash := tt.setup(t, repo)
			testPath := tt.testPath(repo)

			hash, err := gitManager.GetCommitHashInPath(testPath, tt.ref)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedHash, hash)
				assert.Len(t, hash, 40) // Full SHA-1 hash should be 40 characters
			}
		})
	}
}

func TestGitManager_GetCommitHistory(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	// Create some test commits using git harness utilities
	var secondCommitHash string
	must(t, repo.InLocalRepo(func() error {
		// Create first commit
		must(t, repo.WriteFile("test1.txt", "content1"))
		must(t, repo.CommitChanges("First test commit"))

		// Create second commit
		must(t, repo.WriteFile("test2.txt", "content2"))
		must(t, repo.CommitChanges("Second test commit"))

		// Get second commit hash using GitManager utility
		hash, err := gitManager.GetCommitHash("HEAD")
		if err != nil {
			return err
		}
		secondCommitHash = hash

		return nil
	}))

	tests := []struct {
		name          string
		options       CommitHistoryOptions
		expectedCount int
		expectedFirst string // Hash of first commit in result
		expectError   bool
	}{
		{
			name: "get recent commits with limit",
			options: CommitHistoryOptions{
				Limit: 1,
			},
			expectedCount: 1,
			expectedFirst: secondCommitHash, // Most recent commit
			expectError:   false,
		},
		{
			name: "get all commits",
			options: CommitHistoryOptions{
				Limit: 10, // More than we have
			},
			expectedCount: 3, // main initial + first + second
			expectedFirst: secondCommitHash,
			expectError:   false,
		},
		{
			name: "get commits with custom format",
			options: CommitHistoryOptions{
				Limit:        1,
				CustomFormat: "%H|%s|%an",
			},
			expectedCount: 1,
			expectedFirst: secondCommitHash,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commits, err := gitManager.GetCommitHistory("", tt.options)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, commits)
			} else {
				assert.NoError(t, err)
				assert.Len(t, commits, tt.expectedCount)
				if tt.expectedCount > 0 {
					assert.Equal(t, tt.expectedFirst, commits[0].Hash)
					assert.NotEmpty(t, commits[0].Message)
					assert.NotEmpty(t, commits[0].Author)
				}
			}
		})
	}
}

func TestGitManager_GetFileChanges(t *testing.T) {
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	defer repo.Cleanup()

	gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
	require.NoError(t, err)

	// Create some test file changes
	must(t, repo.InLocalRepo(func() error {
		// Create and stage a new file
		must(t, repo.WriteFile("new.txt", "new content"))
		if err := execGitCommandRun(repo.GetLocalPath(), "add", "new.txt"); err != nil {
			return err
		}

		// Modify existing file (unstaged)
		must(t, repo.WriteFile("README.md", "# Modified README\nAdded content"))

		return nil
	}))

	tests := []struct {
		name           string
		options        FileChangeOptions
		expectStaged   bool
		expectUnstaged bool
		expectError    bool
	}{
		{
			name: "get unstaged changes only",
			options: FileChangeOptions{
				Unstaged: true,
			},
			expectStaged:   false,
			expectUnstaged: true,
			expectError:    false,
		},
		{
			name: "get staged changes only",
			options: FileChangeOptions{
				Staged: true,
			},
			expectStaged:   true,
			expectUnstaged: false,
			expectError:    false,
		},
		{
			name: "get both staged and unstaged changes",
			options: FileChangeOptions{
				Staged:   true,
				Unstaged: true,
			},
			expectStaged:   true,
			expectUnstaged: true,
			expectError:    false,
		},
		{
			name:           "default behavior (unstaged only)",
			options:        FileChangeOptions{},
			expectStaged:   false,
			expectUnstaged: true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes, err := gitManager.GetFileChanges("", tt.options)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, changes)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, changes)

				// Check if we have the expected file changes
				var hasNewFile, hasReadme bool
				for _, change := range changes {
					if change.Path == "new.txt" {
						hasNewFile = true
						assert.Equal(t, "added", change.Status)
						assert.Greater(t, change.Additions, 0)
					}
					if change.Path == "README.md" {
						hasReadme = true
						assert.Equal(t, "modified", change.Status)
					}
				}

				if tt.expectStaged {
					assert.True(t, hasNewFile, "Should find staged new.txt file")
				}
				if tt.expectUnstaged {
					assert.True(t, hasReadme, "Should find unstaged README.md file")
				}
			}
		})
	}
}
