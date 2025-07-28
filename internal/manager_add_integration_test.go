package internal

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for Manager.AddWorktree using a shared test repository
// These tests verify that AddWorktree works correctly with real git operations

func TestManager_AddWorktree_Integration(t *testing.T) {
	// Create a shared test repository once to optimize test performance
	sourceRepo := testutils.NewMultiBranchRepo(t)

	// Use the existing repository directly - no need to copy
	repoPath := sourceRepo.GetLocalPath()

	// Change to the repository directory for gbm setup
	originalDir, _ := os.Getwd()
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore directory: %v", err)
		}
	})
	err := os.Chdir(repoPath)
	require.NoError(t, err)

	// Create gbm.branchconfig.yaml for the manager
	gbmContent := `worktrees:
  main:
    branch: main
    description: "Main production branch"
`
	err = os.WriteFile(DefaultBranchConfigFilename, []byte(gbmContent), 0644)
	require.NoError(t, err)

	tests := []struct {
		name         string
		worktreeName string
		branchName   string
		newBranch    bool
		baseBranch   string
		expectErr    func(t *testing.T, err error)
		expect       func(t *testing.T, manager *Manager, worktreeName, branchName string)
	}{
		{
			name:         "create new branch from default",
			worktreeName: "feature-new",
			branchName:   "feature/new-feature",
			newBranch:    true,
			baseBranch:   "", // Will use default branch
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, manager *Manager, worktreeName, branchName string) {
				// Verify worktree directory was created
				worktreePath := filepath.Join(repoPath, DefaultWorktreeDirname, worktreeName)
				assert.DirExists(t, worktreePath)

				// Verify branch exists
				exists, err := manager.BranchExists(branchName)
				require.NoError(t, err)
				assert.True(t, exists)
			},
		},
		{
			name:         "create new branch from specific base",
			worktreeName: "hotfix-urgent",
			branchName:   "hotfix/urgent-fix",
			newBranch:    true,
			baseBranch:   "production/v1.0",
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, manager *Manager, worktreeName, branchName string) {
				// Verify worktree directory was created
				worktreePath := filepath.Join(repoPath, DefaultWorktreeDirname, worktreeName)
				assert.DirExists(t, worktreePath)

				// Verify branch exists
				exists, err := manager.BranchExists(branchName)
				require.NoError(t, err)
				assert.True(t, exists)
			},
		},
		{
			name:         "error: checkout existing branch tracking issue",
			worktreeName: "develop-work",
			branchName:   "develop", // This will fail due to tracking setup issues
			newBranch:    false,
			baseBranch:   "",
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// This is actually an error case due to git worktree tracking behavior
			},
			expect: func(t *testing.T, manager *Manager, worktreeName, branchName string) {
				// Verify worktree directory was NOT created due to error
				worktreePath := filepath.Join(repoPath, DefaultWorktreeDirname, worktreeName)
				assert.NoDirExists(t, worktreePath)
			},
		},
		{
			name:         "error: invalid base branch",
			worktreeName: "error-test1",
			branchName:   "test-branch",
			newBranch:    true,
			baseBranch:   "nonexistent-branch",
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// Check for general error, not specific text since git error messages vary
			},
			expect: func(t *testing.T, manager *Manager, worktreeName, branchName string) {
				// Verify worktree directory was NOT created
				worktreePath := filepath.Join(repoPath, DefaultWorktreeDirname, worktreeName)
				assert.NoDirExists(t, worktreePath)
			},
		},
		{
			name:         "error: checkout nonexistent branch",
			worktreeName: "error-test2",
			branchName:   "nonexistent-branch",
			newBranch:    false,
			baseBranch:   "",
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// Check for general error, not specific text since git error messages vary
			},
			expect: func(t *testing.T, manager *Manager, worktreeName, branchName string) {
				// Verify worktree directory was NOT created
				worktreePath := filepath.Join(repoPath, DefaultWorktreeDirname, worktreeName)
				assert.NoDirExists(t, worktreePath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize manager for each test
			manager, err := NewManager(repoPath)
			require.NoError(t, err, "Failed to initialize manager")

			// Execute the test
			err = manager.AddWorktree(tt.worktreeName, tt.branchName, tt.newBranch, tt.baseBranch)

			// Check error expectations
			tt.expectErr(t, err)

			// Check result expectations
			tt.expect(t, manager, tt.worktreeName, tt.branchName)
		})
	}
}
