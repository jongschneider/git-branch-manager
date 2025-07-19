package cmd

import (
	"os"
	"os/exec"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindMergeTargetWithTreeStructure(t *testing.T) {
	tests := []struct {
		name           string
		worktrees      map[string]testutils.WorktreeConfig
		expectBranch   string
		expectWorktree string
		expectError    bool
		hasCommits     bool
		commitSetup    func(*testutils.GitTestRepo) error
	}{
		{
			name: "Production to Master mergeback needed",
			worktrees: map[string]testutils.WorktreeConfig{
				"master": {
					Branch:      "master",
					Description: "Master branch",
				},
				"production": {
					Branch:      "production-2025-05-1",
					MergeInto:   "master",
					Description: "Production branch",
				},
			},
			expectBranch:   "master", // Target branch (master)
			expectWorktree: "master", // Target worktree
			hasCommits:     true,
			commitSetup: func(repo *testutils.GitTestRepo) error {
				// Create commits on production that need to be merged to master
				err := repo.SwitchToBranch("production-2025-05-1")
				if err != nil {
					return err
				}
				err = repo.WriteFile("production-change.txt", "Production change")
				if err != nil {
					return err
				}
				return repo.CommitChangesWithForceAdd("Add production change")
			},
		},
		{
			name: "Preview to Master mergeback needed when production is up to date",
			worktrees: map[string]testutils.WorktreeConfig{
				"master": {
					Branch:      "master",
					Description: "Master branch",
				},
				"preview": {
					Branch:      "production-2025-07-1",
					MergeInto:   "master",
					Description: "Preview branch",
				},
				"production": {
					Branch:      "production-2025-05-1",
					MergeInto:   "preview",
					Description: "Production branch",
				},
			},
			expectBranch:   "master", // Target branch (master)
			expectWorktree: "master", // Target worktree
			hasCommits:     true,
			commitSetup: func(repo *testutils.GitTestRepo) error {
				// Create commits on preview that need to be merged to master
				err := repo.SwitchToBranch("production-2025-07-1")
				if err != nil {
					return err
				}
				err = repo.WriteFile("preview-change.txt", "Preview change")
				if err != nil {
					return err
				}
				return repo.CommitChangesWithForceAdd("Add preview change")
			},
		},
		{
			name: "Simple two-level chain",
			worktrees: map[string]testutils.WorktreeConfig{
				"master": {
					Branch:      "master",
					Description: "Master branch",
				},
				"production": {
					Branch:      "production-branch",
					MergeInto:   "master",
					Description: "Production branch",
				},
			},
			expectBranch:   "master",
			expectWorktree: "master",
			hasCommits:     true,
			commitSetup: func(repo *testutils.GitTestRepo) error {
				err := repo.SwitchToBranch("production-branch")
				if err != nil {
					return err
				}
				err = repo.WriteFile("prod-change.txt", "Production change")
				if err != nil {
					return err
				}
				return repo.CommitChangesWithForceAdd("Add production change")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test repository
			repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("master"))
			defer repo.Cleanup()

			// Create GBM config
			err := repo.CreateGBMConfig(tt.worktrees)
			require.NoError(t, err)

			// Create branches manually in the correct order to avoid divergent histories
			if tt.name == "Preview to Master mergeback needed when production is up to date" {
				// Create preview branch from master
				err := repo.CreateBranch("production-2025-07-1", "Preview content")
				require.NoError(t, err)

				// Create production branch from preview WITHOUT additional commits
				err = repo.SwitchToBranch("production-2025-07-1")
				require.NoError(t, err)
				// Use exec.Command to create branch without commits
				cmd := exec.Command("git", "checkout", "-b", "production-2025-05-1")
				cmd.Dir = repo.GetLocalPath()
				err = cmd.Run()
				require.NoError(t, err)
				// Push the branch
				cmd = exec.Command("git", "push", "origin", "production-2025-05-1")
				cmd.Dir = repo.GetLocalPath()
				err = cmd.Run()
				require.NoError(t, err)
			} else {
				// Default behavior for other tests
				for worktreeName, config := range tt.worktrees {
					if worktreeName == "master" && config.Branch == "master" {
						continue // Skip master as it already exists
					}
					err := repo.CreateBranch(config.Branch, "Initial "+worktreeName+" content")
					require.NoError(t, err)
				}
			}

			// Set up commits if needed
			if tt.hasCommits && tt.commitSetup != nil {
				err := tt.commitSetup(repo)
				require.NoError(t, err)
			}

			// Switch to the repo directory
			err = os.Chdir(repo.GetLocalPath())
			require.NoError(t, err)

			// Create manager
			manager, err := createInitializedManager()
			require.NoError(t, err)

			// Test findMergeTargetBranchAndWorktree
			branch, worktree, err := findMergeTargetBranchAndWorktree(manager)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectBranch, branch, "Should return correct target branch")
				assert.Equal(t, tt.expectWorktree, worktree, "Should return correct target worktree")
			}
		})
	}
}

func TestMergebackNamingWithTreeStructure(t *testing.T) {
	// Create test repository
	repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("master"))
	defer repo.Cleanup()

	// Create GBM config: production -> preview -> master
	worktrees := map[string]testutils.WorktreeConfig{
		"master": {
			Branch:      "master",
			Description: "Master branch",
		},
		"preview": {
			Branch:      "production-2025-07-1",
			MergeInto:   "master",
			Description: "Preview branch",
		},
		"production": {
			Branch:      "production-2025-05-1",
			MergeInto:   "preview",
			Description: "Production branch",
		},
	}

	err := repo.CreateGBMConfig(worktrees)
	require.NoError(t, err)

	// Create branches
	err = repo.CreateBranch("production-2025-07-1", "Preview content")
	require.NoError(t, err)

	err = repo.CreateBranch("production-2025-05-1", "Production content")
	require.NoError(t, err)

	// Add commits to production that need mergeback
	err = repo.SwitchToBranch("production-2025-05-1")
	require.NoError(t, err)

	err = repo.WriteFile("production-change.txt", "Production change")
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add production change")
	require.NoError(t, err)

	// Switch to repo directory
	err = os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	// Create manager
	manager, err := createInitializedManager()
	require.NoError(t, err)

	// Find merge target (should be production -> preview)
	targetBranch, targetWorktree, err := findMergeTargetBranchAndWorktree(manager)
	require.NoError(t, err)

	// Should target preview branch/worktree (immediate parent of production)
	assert.Equal(t, "production-2025-07-1", targetBranch)
	assert.Equal(t, "preview", targetWorktree)

	// Test mergeback branch name generation
	worktreeName := "INGSVC-5638"
	jiraTicket := "INGSVC-5638"

	branchName, err := generateMergebackBranchName(worktreeName, jiraTicket, targetWorktree, manager)
	require.NoError(t, err)

	// Should include target worktree name in branch name
	assert.Contains(t, branchName, "preview", "Branch name should contain target worktree name")
	assert.Contains(t, branchName, "INGSVC-5638", "Branch name should contain JIRA ticket")

	// Test worktree name generation
	mergebackPrefix := manager.GetConfig().Settings.MergebackPrefix
	expectedWorktreeName := mergebackPrefix + "_" + worktreeName + "_" + targetWorktree

	// Should be MERGE_INGSVC-5638_preview (target worktree, not source)
	assert.Equal(t, "MERGE_INGSVC-5638_preview", expectedWorktreeName)
}

func TestMergebackNamingProductionToMaster(t *testing.T) {
	// Test the specific case mentioned in the issue: production -> master should use "_master" suffix

	// Create test repository
	repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("master"))
	defer repo.Cleanup()

	// Create simple config: production -> master
	worktrees := map[string]testutils.WorktreeConfig{
		"master": {
			Branch:      "master",
			Description: "Master branch",
		},
		"production": {
			Branch:      "production-2025-05-1",
			MergeInto:   "master",
			Description: "Production branch",
		},
	}

	err := repo.CreateGBMConfig(worktrees)
	require.NoError(t, err)

	// Create production branch
	err = repo.CreateBranch("production-2025-05-1", "Production content")
	require.NoError(t, err)

	// Add commits to production that need mergeback to master
	err = repo.SwitchToBranch("production-2025-05-1")
	require.NoError(t, err)

	err = repo.WriteFile("production-change.txt", "Production change for master")
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add production change for master")
	require.NoError(t, err)

	// Switch to repo directory
	err = os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	// Create manager
	manager, err := createInitializedManager()
	require.NoError(t, err)

	// Find merge target (should be production -> master)
	targetBranch, targetWorktree, err := findMergeTargetBranchAndWorktree(manager)
	require.NoError(t, err)

	// Should target master branch/worktree
	assert.Equal(t, "master", targetBranch)
	assert.Equal(t, "master", targetWorktree)

	// Test the specific naming case from the issue
	worktreeName := "INGSVC-5638"
	jiraTicket := "INGSVC-5638"

	branchName, err := generateMergebackBranchName(worktreeName, jiraTicket, targetWorktree, manager)
	require.NoError(t, err)

	// Branch name should end with "_master" not "_production"
	assert.Contains(t, branchName, "master", "Branch name should contain target 'master'")
	assert.NotContains(t, branchName, "production", "Branch name should NOT contain source 'production'")

	// Test worktree name
	mergebackPrefix := manager.GetConfig().Settings.MergebackPrefix
	expectedWorktreeName := mergebackPrefix + "_" + worktreeName + "_" + targetWorktree

	// Should be MERGE_INGSVC-5638_master (target worktree "master", not source "production")
	assert.Equal(t, "MERGE_INGSVC-5638_master", expectedWorktreeName)

	// This is the fix for the issue: it should NOT be MERGE_INGSVC-5638_production
	assert.NotEqual(t, "MERGE_INGSVC-5638_production", expectedWorktreeName)
}
