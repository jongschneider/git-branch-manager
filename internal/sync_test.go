package internal

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_SyncBasicOperations(t *testing.T) {
	tests := []struct {
		name         string
		setupRepo    func(t *testing.T) *testutils.GitTestRepo
		expectedDirs []string
	}{
		{
			name: "sync with existing gbm config creates all worktrees",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				return testutils.NewStandardGBMConfigRepo(t) // Has main, dev, feat, prod
			},
			expectedDirs: []string{"worktrees/main", "worktrees/dev", "worktrees/feat", "worktrees/prod"},
		},
		{
			name: "sync with minimal gbm config",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				repo := testutils.NewBasicRepo(t)
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
				require.NoError(t, repo.WriteFile(DefaultBranchConfigFilename, gbmContent))
				require.NoError(t, repo.CommitChangesWithForceAdd("Add gbm config"))
				require.NoError(t, repo.PushBranch("main"))
				return repo
			},
			expectedDirs: []string{"worktrees/main"},
		},
		{
			name: "sync with already synced repo is idempotent",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				return testutils.NewStandardGBMConfigRepo(t)
			},
			expectedDirs: []string{"worktrees/main", "worktrees/dev", "worktrees/feat", "worktrees/prod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepo := tt.setupRepo(t)
			defer sourceRepo.Cleanup()

			// Clone into a new directory to simulate working in a cloned repo
			wd := t.TempDir()
			require.NoError(t, os.Chdir(wd))
			require.NoError(t, execGitCommandRun(wd, "clone", sourceRepo.GetRemotePath(), "."))

			// Create manager and test sync operations
			manager, err := NewManager(wd)
			require.NoError(t, err)
			// Load gbm.branchconfig.yaml before sync operations
			require.NoError(t, manager.LoadGBMConfig(""))
			// Load gbm.branchconfig.yaml before sync operations
			require.NoError(t, manager.LoadGBMConfig(""))

			// For the idempotent test, run sync twice
			if len(tt.expectedDirs) == 4 { // Standard config test
				err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
				require.NoError(t, err) // First sync for idempotent test
			}

			err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
			require.NoError(t, err)

			for _, expectedDir := range tt.expectedDirs {
				assert.DirExists(t, filepath.Join(wd, expectedDir))
			}
		})
	}
}

func TestManager_SyncScenarios(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(t *testing.T) *testutils.GitTestRepo
		modifyGbm func(t *testing.T, repoPath string)
		validate  func(t *testing.T, manager *Manager, repoPath string)
	}{
		{
			name: "branch reference change updates existing worktree",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				repo := testutils.NewMultiBranchRepo(t)
				// Seed initial GBM config
				initial := map[string]testutils.WorktreeConfig{
					"main": {Branch: "main", Description: "Main branch"},
					"dev":  {Branch: "develop", Description: "Development branch"},
				}
				require.NoError(t, repo.CreateGBMConfig(initial))
				require.NoError(t, repo.CommitChangesWithForceAdd("Add initial gbm config"))
				require.NoError(t, repo.PushBranch("main"))
				return repo
			},
			modifyGbm: func(t *testing.T, repoPath string) {
				// Change dev worktree to point to feature-1 instead of dev
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: feature/auth
    description: "Development branch"
`
				require.NoError(t, os.WriteFile(filepath.Join(repoPath, DefaultBranchConfigFilename), []byte(gbmContent), 0o644))
				require.NoError(t, execGitCommandRun(repoPath, "add", DefaultBranchConfigFilename))
				require.NoError(t, execGitCommandRun(repoPath, "commit", "-m", "Update gbm config"))
				require.NoError(t, execGitCommandRun(repoPath, "push", "origin", "main"))
			},
			validate: func(t *testing.T, manager *Manager, repoPath string) {
				// Check that dev worktree now has feature-1 branch
				branch, err := manager.GetGitManager().GetCurrentBranchInPath(filepath.Join(repoPath, "worktrees/dev"))
				require.NoError(t, err)
				assert.Equal(t, "feature/auth", branch)
			},
		},
		{
			name: "new worktree added to config creates worktree",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				repo := testutils.NewMultiBranchRepo(t)
				// Seed initial GBM config with main and dev
				initial := map[string]testutils.WorktreeConfig{
					"main": {Branch: "main", Description: "Main branch"},
					"dev":  {Branch: "develop", Description: "Development branch"},
				}
				require.NoError(t, repo.CreateGBMConfig(initial))
				require.NoError(t, repo.CommitChangesWithForceAdd("Add initial gbm config"))
				require.NoError(t, repo.PushBranch("main"))
				return repo
			},
			modifyGbm: func(t *testing.T, repoPath string) {
				// Add new worktree for feature-1
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
  feature:
    branch: feature/auth
    description: "Feature branch"
`
				require.NoError(t, os.WriteFile(filepath.Join(repoPath, DefaultBranchConfigFilename), []byte(gbmContent), 0o644))
				require.NoError(t, execGitCommandRun(repoPath, "add", DefaultBranchConfigFilename))
				require.NoError(t, execGitCommandRun(repoPath, "commit", "-m", "Add feature worktree"))
				require.NoError(t, execGitCommandRun(repoPath, "push", "origin", "main"))
			},
			validate: func(t *testing.T, manager *Manager, repoPath string) {
				// Check that feature worktree exists
				assert.DirExists(t, filepath.Join(repoPath, "worktrees/feature"))
				branch, err := manager.GetGitManager().GetCurrentBranchInPath(filepath.Join(repoPath, "worktrees/feature"))
				require.NoError(t, err)
				assert.Equal(t, "feature/auth", branch)
			},
		},
		{
			name: "worktree removed from config (no-op without force)",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				repo := testutils.NewMultiBranchRepo(t)
				// Seed initial GBM config with main and dev
				initial := map[string]testutils.WorktreeConfig{
					"main": {Branch: "main", Description: "Main branch"},
					"dev":  {Branch: "develop", Description: "Development branch"},
				}
				require.NoError(t, repo.CreateGBMConfig(initial))
				require.NoError(t, repo.CommitChangesWithForceAdd("Add initial gbm config"))
				require.NoError(t, repo.PushBranch("main"))
				return repo
			},
			modifyGbm: func(t *testing.T, repoPath string) {
				// Remove dev worktree from config
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
				require.NoError(t, os.WriteFile(filepath.Join(repoPath, DefaultBranchConfigFilename), []byte(gbmContent), 0o644))
				require.NoError(t, execGitCommandRun(repoPath, "add", DefaultBranchConfigFilename))
				require.NoError(t, execGitCommandRun(repoPath, "commit", "-m", "Remove dev worktree"))
				require.NoError(t, execGitCommandRun(repoPath, "push", "origin", "main"))
			},
			validate: func(t *testing.T, manager *Manager, repoPath string) {
				// Orphaned worktree should still exist (not removed without force)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees/dev"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepo := tt.setupRepo(t)
			defer sourceRepo.Cleanup()

			// Clone and do initial sync
			wd := t.TempDir()
			require.NoError(t, os.Chdir(wd))
			require.NoError(t, execGitCommandRun(wd, "clone", sourceRepo.GetRemotePath(), "."))

			manager, err := NewManager(wd)
			require.NoError(t, err)
			// Load gbm.branchconfig.yaml before sync operations
			require.NoError(t, manager.LoadGBMConfig(""))

			// Initial sync to create worktrees
			err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
			require.NoError(t, err)

			// Modify gbm config as per test (in the source repo), then push and pull in clone
			tt.modifyGbm(t, sourceRepo.GetLocalPath())
			// Ensure remote has the new commit
			require.NoError(t, execGitCommandRun(sourceRepo.GetLocalPath(), "push", "origin", "main"))

			// Pull changes and sync again (specify remote and branch since main might be in worktree)
			if output, err := ExecGitCommandCombined(wd, "pull", "origin", "main"); err != nil {
				t.Fatalf("git pull failed: %s", string(output))
			}
			// Reload gbm.branchconfig.yaml after pulling updates
			require.NoError(t, manager.LoadGBMConfig(""))
			err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
			require.NoError(t, err)

			// Validate results
			tt.validate(t, manager, wd)
		})
	}
}

func TestManager_SyncIntegration(t *testing.T) {
	t.Run("complete sync workflow with force operations", func(t *testing.T) {
		sourceRepo := testutils.NewStandardGBMConfigRepo(t)
		defer sourceRepo.Cleanup()

		// Clone into working directory
		wd := t.TempDir()
		require.NoError(t, os.Chdir(wd))
		require.NoError(t, execGitCommandRun(wd, "clone", sourceRepo.GetRemotePath(), "."))

		manager, err := NewManager(wd)
		require.NoError(t, err)
		// Load gbm.branchconfig.yaml before sync operations
		require.NoError(t, manager.LoadGBMConfig(""))

		// Initial sync
		err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
		require.NoError(t, err)

		// Manually corrupt worktrees by removing dev worktree directory but keeping git worktree entry
		devWorktreePath := filepath.Join(wd, "worktrees/dev")
		require.NoError(t, os.RemoveAll(devWorktreePath))

		// Prune worktree to clean up git worktree list
		require.NoError(t, execGitCommandRun(wd, "worktree", "prune"))

		// Sync with force should recreate the removed worktree
		err = manager.SyncWithConfirmation(false, true, func(string) bool { return true })
		require.NoError(t, err)

		// Verify dev worktree was recreated
		assert.DirExists(t, devWorktreePath)
		branch, err := manager.GetGitManager().GetCurrentBranchInPath(devWorktreePath)
		require.NoError(t, err)
		// The standard config maps worktree 'dev' to branch 'develop'
		assert.Equal(t, "develop", branch)
	})
}

func TestManager_SyncWorkreePromotion(t *testing.T) {
	t.Run("worktree promotion workflow", func(t *testing.T) {
		// Create repo with multiple production branches scenario
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		defer repo.Cleanup()

		// Create production branches
		must(t, repo.CreateBranch("production", "Initial production content"))
		must(t, repo.CreateBranch("production-v2", "Initial production-v2 content"))
		must(t, repo.PushBranch("main"))
		must(t, repo.PushBranch("production"))
		must(t, repo.PushBranch("production-v2"))

		// Create initial GBM config with production worktree pointing to production branch
		gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  production:
    branch: production
    description: "Production branch"
  production-v2:
    branch: production-v2
    description: "Production v2 branch"
`
		must(t, repo.WriteFile(DefaultBranchConfigFilename, gbmContent))
		must(t, repo.CommitChangesWithForceAdd("Add initial gbm config"))
		must(t, repo.PushBranch("main"))

		// Clone and do initial sync
		wd := t.TempDir()
		require.NoError(t, os.Chdir(wd))
		require.NoError(t, execGitCommandRun(wd, "clone", repo.GetRemotePath(), "."))

		manager, err := NewManager(wd)
		require.NoError(t, err)
		// Load gbm.branchconfig.yaml before sync operations
		require.NoError(t, manager.LoadGBMConfig(""))

		// Initial sync creates worktrees
		err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
		require.NoError(t, err)

		// Modify config to cause promotion in source repo: production worktree should now point to production-v2
		// and production-v2 worktree should point to production
		newGbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  production:
    branch: production-v2
    description: "Production branch (promoted)"
  production-v2:
    branch: production
    description: "Production v2 branch (demoted)"
`
		require.NoError(t, repo.WriteFile(DefaultBranchConfigFilename, newGbmContent))
		require.NoError(t, repo.CommitChangesWithForceAdd("Update gbm config for promotion"))
		require.NoError(t, repo.PushBranch("main"))

		// Pull changes and sync again (specify remote and branch since main might be in worktree)
		if output, err := ExecGitCommandCombined(wd, "pull", "origin", "main"); err != nil {
			t.Fatalf("git pull failed: %s", string(output))
		}
		// Reload gbm.branchconfig.yaml after pulling updates
		require.NoError(t, manager.LoadGBMConfig(""))
		err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
		require.NoError(t, err)

		// Validate promotion occurred correctly
		prodBranch, err := manager.GetGitManager().GetCurrentBranchInPath(filepath.Join(wd, "worktrees/production"))
		require.NoError(t, err)
		assert.Equal(t, "production-v2", prodBranch)

		prodV2Branch, err := manager.GetGitManager().GetCurrentBranchInPath(filepath.Join(wd, "worktrees/production-v2"))
		require.NoError(t, err)
		assert.Equal(t, "production", prodV2Branch)
	})
}

func TestManager_GetSyncStatus(t *testing.T) {
	t.Run("sync status analysis", func(t *testing.T) {
		sourceRepo := testutils.NewStandardGBMConfigRepo(t)
		defer sourceRepo.Cleanup()

		wd := t.TempDir()
		require.NoError(t, os.Chdir(wd))
		require.NoError(t, execGitCommandRun(wd, "clone", sourceRepo.GetRemotePath(), "."))

		manager, err := NewManager(wd)
		require.NoError(t, err)

		// Before sync, all worktrees should be missing
		status, err := manager.GetSyncStatus()
		require.NoError(t, err)

		assert.False(t, status.InSync)
		assert.Contains(t, status.MissingWorktrees, "main")
		assert.Contains(t, status.MissingWorktrees, "dev")
		assert.Contains(t, status.MissingWorktrees, "feat")
		assert.Contains(t, status.MissingWorktrees, "prod")

		// After sync, should be in sync
		err = manager.SyncWithConfirmation(false, false, func(string) bool { return true })
		require.NoError(t, err)

		status, err = manager.GetSyncStatus()
		require.NoError(t, err)
		assert.True(t, status.InSync)
		assert.Empty(t, status.MissingWorktrees)
		assert.Empty(t, status.BranchChanges)
		assert.Empty(t, status.OrphanedWorktrees)
	})
}
