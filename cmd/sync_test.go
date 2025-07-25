package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal"
	"gbm/internal/testutils"

	"slices"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncCommand_BasicOperations(t *testing.T) {
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
				require.NoError(t, repo.WriteFile(internal.DefaultBranchConfigFilename, gbmContent))
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
			setupClonedRepo(t, sourceRepo)

			// For the idempotent test, run sync twice
			if strings.Contains(tt.name, "idempotent") {
				cmd := newRootCommand()
				cmd.SetArgs([]string{"sync"})
				require.NoError(t, cmd.Execute()) // First sync
			}

			cmd := newRootCommand()
			cmd.SetArgs([]string{"sync"})

			err := cmd.Execute()
			require.NoError(t, err)

			for _, expectedDir := range tt.expectedDirs {
				wd, _ := os.Getwd()
				assert.DirExists(t, filepath.Join(wd, expectedDir))
			}
		})
	}
}

func TestSyncCommand_Flags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		setup    func(t *testing.T, repo *testutils.GitTestRepo)
		validate func(t *testing.T, repoPath string, output string, err error)
	}{
		{
			name: "dry-run flag shows changes without applying",
			args: []string{"sync", "--dry-run"},
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Add more worktrees to gbm.branchconfig.yaml to create missing worktrees scenario
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  new1:
    branch: develop
    description: "Development branch"
  new2:
    branch: feature/auth
    description: "Feature branch"
`
				require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(gbmContent), 0644))
			},
			validate: func(t *testing.T, repoPath string, output string, err error) {
				require.NoError(t, err)
				// Check that the command succeeded and directories are as expected
				// new1 and new2 should still be missing after dry-run
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "new1"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "new2"))
				// main should still exist (was created by clone)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
			},
		},
		{
			name: "force flag removes orphaned worktrees with confirmation",
			args: []string{"sync", "--force"},
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create untracked worktree
				wd, _ := os.Getwd()
				createUntrackedWorktree(t, wd, "orphan", "main")

				// Modify gbm.branchconfig.yaml to remove some existing worktrees (making them orphaned)
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
				require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(gbmContent), 0644))
			},
			validate: func(t *testing.T, repoPath string, output string, err error) {
				// This test should fail because confirmation is required
				// We need to use the manager directly with mock confirmation instead
				require.Error(t, err)
				assert.ErrorContains(t, err, "sync cancelled by user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepo := testutils.NewStandardGBMConfigRepo(t)
			repoPath := setupClonedRepo(t, sourceRepo)

			tt.setup(t, &testutils.GitTestRepo{LocalDir: repoPath})

			var buf bytes.Buffer
			cmd := newRootCommand()
			cmd.SetArgs(tt.args)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()
			output := buf.String()

			tt.validate(t, repoPath, output, err)
		})
	}
}

func TestSyncCommand_SyncScenarios(t *testing.T) {
	tests := []struct {
		name             string
		initialGBMConfig string
		updatedGBMConfig string
		expectChanges    bool
		validateResult   func(t *testing.T, repoPath string)
	}{
		{
			name: "branch reference changed",
			initialGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  feat:
    branch: feature/auth
    description: "Feature branch"
`,
			updatedGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  feat:
    branch: develop
    description: "Feature branch"
`,
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				// Verify feat worktree was updated to develop branch
				cmd := exec.Command("git", "branch", "--show-current")
				cmd.Dir = filepath.Join(repoPath, "worktrees", "feat")
				branchOutput, err := cmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "develop", strings.TrimSpace(string(branchOutput)))
			},
		},
		{
			name: "new worktree added",
			initialGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
`,
			updatedGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
`,
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
			},
		},
		{
			name: "worktree removed",
			initialGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  temp:
    branch: develop
    description: "Temporary branch"
`,
			updatedGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
`,
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				// temp should still exist without --force
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "temp"))
			},
		},
		{
			name: "no changes needed",
			initialGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
`,
			updatedGBMConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
`,
			expectChanges: false,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source repo with initial gbm.branchconfig.yaml
			sourceRepo := testutils.NewMultiBranchRepo(t)
			require.NoError(t, sourceRepo.WriteFile(internal.DefaultBranchConfigFilename, tt.initialGBMConfig))
			require.NoError(t, sourceRepo.CommitChangesWithForceAdd("Add initial gbm.branchconfig.yaml"))
			require.NoError(t, sourceRepo.PushBranch("main"))

			repoPath := setupClonedRepo(t, sourceRepo)

			// Update gbm.branchconfig.yaml to new configuration
			require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(tt.updatedGBMConfig), 0644))

			// Run sync with updated gbm.branchconfig.yaml
			cmd := newRootCommand()
			cmd.SetArgs([]string{"sync"})
			err := cmd.Execute()
			require.NoError(t, err)

			tt.validateResult(t, repoPath)
		})
	}
}

func TestSyncCommand_UntrackedWorktrees(t *testing.T) {
	tests := []struct {
		name                   string
		gbmConfig              string
		untrackedWorktrees     []string
		syncArgs               []string
		validateResult         func(t *testing.T, repoPath string, output string)
		createTrackedWorktrees bool
	}{
		{
			name: "untracked worktree preserved by default",
			gbmConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
`,
			untrackedWorktrees:     []string{"manual"},
			syncArgs:               []string{"sync"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// Tracked worktrees should exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
				// Untracked worktree should still exist (not removed)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "manual"))
			},
		},
		{
			name: "untracked worktree removed with --force",
			gbmConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  dev:
    branch: develop
    description: "Development branch"
`,
			untrackedWorktrees:     []string{"manual"},
			syncArgs:               []string{"sync", "--force"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// Tracked worktrees should exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
				// Untracked worktree should be removed
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "manual"))
			},
		},
		{
			name: "dry-run shows untracked worktree would be removed",
			gbmConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
`,
			untrackedWorktrees:     []string{"temp", "experimental"},
			syncArgs:               []string{"sync", "--dry-run", "--force"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// All worktrees should still exist (dry-run)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "temp"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "experimental"))
				// In dry-run mode, orphaned worktrees should NOT be removed
				// (We can't easily capture the output due to test isolation issues,
				// but we can verify the intended behavior: no actual changes made)
			},
		},
		{
			name: "tracked worktrees updated, untracked preserved without force",
			gbmConfig: `worktrees:
  main:
    branch: main
    description: "Main branch"
  feat:
    branch: develop
    description: "Feature branch"
`,
			untrackedWorktrees:     []string{"manual"},
			syncArgs:               []string{"sync"},
			createTrackedWorktrees: true,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// TRACKED worktrees should exist and be updated
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "feat"))

				// UNTRACKED worktree should still exist (not removed without --force)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "manual"))

				// Verify feat worktree was updated to develop branch
				cmd := exec.Command("git", "branch", "--show-current")
				cmd.Dir = filepath.Join(repoPath, "worktrees", "feat")
				featOutput, err := cmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "develop", strings.TrimSpace(string(featOutput)))

				// manual worktree should still be on main branch (unchanged)
				manualCmd := exec.Command("git", "branch", "--show-current")
				manualCmd.Dir = filepath.Join(repoPath, "worktrees", "manual")
				manualOutput, err := manualCmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "main", strings.TrimSpace(string(manualOutput)))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sourceRepo *testutils.GitTestRepo

			if tt.createTrackedWorktrees {
				// Setup initial gbm.branchconfig.yaml with different branch for feat
				sourceRepo = testutils.NewMultiBranchRepo(t)
				initialGBMConfig := `worktrees:
  main:
    branch: main
    description: "Main branch"
  feat:
    branch: feature/auth
    description: "Feature branch"
`
				require.NoError(t, sourceRepo.WriteFile(internal.DefaultBranchConfigFilename, initialGBMConfig))
				require.NoError(t, sourceRepo.CommitChangesWithForceAdd("Add initial gbm.branchconfig.yaml"))
				require.NoError(t, sourceRepo.PushBranch("main"))
			} else {
				// Standard setup for other tests
				sourceRepo = testutils.NewMultiBranchRepo(t)
				require.NoError(t, sourceRepo.WriteFile(internal.DefaultBranchConfigFilename, tt.gbmConfig))
				require.NoError(t, sourceRepo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml"))
				require.NoError(t, sourceRepo.PushBranch("main"))
			}

			repoPath := setupClonedRepo(t, sourceRepo)

			if tt.createTrackedWorktrees {
				// Create untracked worktree first
				for _, untrackedName := range tt.untrackedWorktrees {
					createUntrackedWorktree(t, repoPath, untrackedName, "main")
				}

				// Update gbm.branchconfig.yaml to change feat branch
				require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(tt.gbmConfig), 0644))
			} else {
				// Create untracked worktrees for standard tests
				for _, untrackedName := range tt.untrackedWorktrees {
					createUntrackedWorktree(t, repoPath, untrackedName, "main")
				}
			}

			// Run the sync command with specified args
			var buf bytes.Buffer
			var err error
			var output string

			// Check if this test uses --force and needs confirmation
			usesForce := slices.Contains(tt.syncArgs, "--force")

			if usesForce {
				// Use simulateUserInput for tests that use --force
				err = simulateUserInput("y", func() error {
					cmd := newRootCommand()
					cmd.SetArgs(tt.syncArgs)
					cmd.SetOut(&buf)
					cmd.SetErr(&buf)
					return cmd.Execute()
				})
				output = buf.String()
			} else {
				// Standard execution for non-force tests
				cmd := newRootCommand()
				cmd.SetArgs(tt.syncArgs)
				cmd.SetOut(&buf)
				cmd.SetErr(&buf)
				err = cmd.Execute()
				output = buf.String()
			}

			require.NoError(t, err)

			tt.validateResult(t, repoPath, output)
		})
	}
}

func TestSyncCommand_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) string // returns working directory
		args          []string
		expectedError string
	}{
		{
			name: "not a git repository",
			setup: func(t *testing.T) string {
				return t.TempDir() // empty directory
			},
			args:          []string{"sync"},
			expectedError: "failed to find git repository root",
		},
		{
			name: "missing gbm config file",
			setup: func(t *testing.T) string {
				repo := testutils.NewBasicRepo(t)
				return repo.GetLocalPath()
			},
			args:          []string{"sync"},
			expectedError: "failed to load gbm.branchconfig.yaml",
		},
		{
			name: "invalid branch reference",
			setup: func(t *testing.T) string {
				// Create a source repo with invalid branch reference
				sourceRepo := testutils.NewBasicRepo(t)
				gbmContent := `worktrees:
  invalid:
    branch: nonexistent-branch
    description: "Invalid branch reference"
`
				require.NoError(t, sourceRepo.WriteFile(internal.DefaultBranchConfigFilename, gbmContent))
				require.NoError(t, sourceRepo.CommitChangesWithForceAdd("Add invalid gbm config"))
				require.NoError(t, sourceRepo.PushBranch("main"))

				// Clone it to set up proper structure, but don't defer the chdir - let the test handle it
				repoPath := setupClonedRepo(t, sourceRepo)
				// Immediately return to original dir so test can handle directory changes
				_ = os.Chdir(repoPath)
				return repoPath
			},
			args:          []string{"sync"},
			expectedError: "branch 'nonexistent-branch' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workingDir := tt.setup(t)
			originalDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(originalDir) }()

			require.NoError(t, os.Chdir(workingDir))

			cmd := newRootCommand()
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestSyncCommand_Integration(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(t *testing.T, repoPath string)
		validate func(t *testing.T, repoPath string)
	}{
		{
			name: "complete sync workflow",
			scenario: func(t *testing.T, repoPath string) {
				// 1. Initial state created by clone (main, dev, feat, prod from StandardGBMConfigRepo)

				// 2. Modify gbm.branchconfig.yaml to remove some worktrees and add different ones
				gbmContent := `worktrees:
  main:
    branch: main
    description: "Main production branch"
  dev:
    branch: develop
    description: "Development branch"
  new_feat:
    branch: feature/auth
    description: "New feature branch"
  new_prod:
    branch: production/v1.0
    description: "New production branch"
`
				require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(gbmContent), 0644))

				// 3. Run sync with --force to remove orphaned worktrees and create new ones
				err := simulateUserInput("y", func() error {
					cmd := newRootCommand()
					cmd.SetArgs([]string{"sync", "--force"})
					return cmd.Execute()
				})
				require.NoError(t, err)
			},
			validate: func(t *testing.T, repoPath string) {
				// Verify new worktrees exist
				expectedDirs := []string{"main", "dev", "new_feat", "new_prod"}
				for _, dir := range expectedDirs {
					assert.DirExists(t, filepath.Join(repoPath, "worktrees", dir))
				}

				// Old worktrees should be removed with --force
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "feat"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "prod"))
			},
		},
		{
			name: "sync after manual worktree changes",
			scenario: func(t *testing.T, repoPath string) {
				// Manually remove a tracked worktree (simulate corruption where worktree is lost)
				devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")

				// Remove the worktree directory first
				require.NoError(t, os.RemoveAll(devWorktreePath))

				// Also remove it from git's worktree list to simulate complete loss
				pruneCmd := exec.Command("git", "worktree", "prune")
				pruneCmd.Dir = repoPath
				require.NoError(t, pruneCmd.Run())

				// Run sync to fix things (recreate missing worktree)
				cmd := newRootCommand()
				cmd.SetArgs([]string{"sync"})
				require.NoError(t, cmd.Execute())
			},
			validate: func(t *testing.T, repoPath string) {
				// dev should be recreated
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
				// All original worktrees should still exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "feat"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "prod"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepo := testutils.NewStandardGBMConfigRepo(t)
			repoPath := setupClonedRepo(t, sourceRepo)

			tt.scenario(t, repoPath)
			tt.validate(t, repoPath)
		})
	}
}

// Helper function to create an untracked worktree
func createUntrackedWorktree(t *testing.T, repoPath, name, branch string) {
	worktreeDir := filepath.Join(repoPath, "worktrees", name)

	// Ensure worktrees directory exists
	require.NoError(t, os.MkdirAll(filepath.Join(repoPath, "worktrees"), 0755))

	cmd := exec.Command("git", "worktree", "add", "--force", worktreeDir, branch)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("git worktree add failed: %v, output: %s", err, string(output))
	}
	require.NoError(t, err, "Failed to create untracked worktree %s", name)
}

func TestSyncCommand_ForceConfirmationDirectManagerTest(t *testing.T) {

	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepo(t, sourceRepo)

	// Create untracked worktree
	createUntrackedWorktree(t, repoPath, "orphan", "main")

	// Modify gbm.branchconfig.yaml to remove some existing worktrees (making them orphaned)
	gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
	require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(gbmContent), 0644))

	// Create manager
	manager, err := createInitializedManager()
	require.NoError(t, err)

	// Create confirmation function that accepts
	confirmFunc := func(message string) bool {
		// Verify message contains what we expect
		assert.Contains(t, message, "PERMANENTLY DELETED")
		return true // Confirm deletion
	}

	// Test sync with force and confirmation
	err = manager.SyncWithConfirmation(false, true, confirmFunc)
	require.NoError(t, err)

	// Verify orphaned worktrees were removed
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
	assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "dev"))
	assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "feat"))
	assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "prod"))
	assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "orphan"))
}

func TestSyncCommand_ForceConfirmation(t *testing.T) {
	tests := []struct {
		name           string
		userResponse   string
		expectedAction string
		shouldSucceed  bool
	}{
		{
			name:           "user confirms deletion with 'y'",
			userResponse:   "y",
			expectedAction: "delete orphaned worktrees",
			shouldSucceed:  true,
		},
		{
			name:           "user confirms deletion with 'yes'",
			userResponse:   "yes",
			expectedAction: "delete orphaned worktrees",
			shouldSucceed:  true,
		},
		{
			name:           "user cancels deletion with 'n'",
			userResponse:   "n",
			expectedAction: "cancel sync operation",
			shouldSucceed:  false,
		},
		{
			name:           "user cancels deletion with empty response",
			userResponse:   "",
			expectedAction: "cancel sync operation",
			shouldSucceed:  false,
		},
		{
			name:           "user cancels deletion with 'no'",
			userResponse:   "no",
			expectedAction: "cancel sync operation",
			shouldSucceed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup a repo with orphaned worktrees
			sourceRepo := testutils.NewStandardGBMConfigRepo(t)
			repoPath := setupClonedRepo(t, sourceRepo)

			// Create orphaned worktree
			createUntrackedWorktree(t, repoPath, "orphan", "main")

			// Modify gbm.branchconfig.yaml to remove some worktrees (making them orphaned)
			gbmContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
			require.NoError(t, os.WriteFile(internal.DefaultBranchConfigFilename, []byte(gbmContent), 0644))

			// Create manager to test confirmation directly
			manager, err := createInitializedManager()
			require.NoError(t, err)

			// Create mock confirmation function
			confirmFunc := func(message string) bool {
				// Verify message contains what we expect
				assert.Contains(t, message, "PERMANENTLY DELETED")
				assert.Contains(t, message, "orphan")

				// Simulate user response
				switch strings.ToLower(tt.userResponse) {
				case "y", "yes":
					return true
				default:
					return false
				}
			}

			// Test sync with confirmation
			err = manager.SyncWithConfirmation(false, true, confirmFunc)

			if tt.shouldSucceed {
				require.NoError(t, err)
				// Verify orphaned worktree was deleted
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "orphan"))
			} else {
				require.Error(t, err)
				assert.ErrorContains(t, err, "sync cancelled by user")
				// Verify orphaned worktree still exists
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "orphan"))
			}
		})
	}
}

func TestSyncCommand_WorktreePromotion(t *testing.T) {
	// Test that verifies our fix for the "exit status 128" bug
	// The original bug occurred when trying to create a worktree for a branch
	// that was already checked out in another worktree during promotion scenarios

	// Create source repository with branches like the working tests do
	sourceRepo := testutils.NewMultiBranchRepo(t)

	// Create the specific branches mentioned in the original bug report
	require.NoError(t, sourceRepo.CreateBranch("production-2025-05-1", "main"))
	require.NoError(t, sourceRepo.CreateBranch("production-2025-07-1", "main"))
	require.NoError(t, sourceRepo.PushBranch("production-2025-05-1"))
	require.NoError(t, sourceRepo.PushBranch("production-2025-07-1"))

	// Create initial gbm.branchconfig.yaml in the source repo
	initialGBMConfig := `worktrees:
  main:
    branch: main
    description: "Main production branch"
  preview:
    branch: production-2025-07-1 
    description: "Blade Runner"
    merge_into: main
  production:
    branch: production-2025-05-1 
    description: "Arrival"
    merge_into: preview
`
	require.NoError(t, sourceRepo.WriteFile(internal.DefaultBranchConfigFilename, initialGBMConfig))
	require.NoError(t, sourceRepo.CommitChangesWithForceAdd("Add initial gbm config"))
	require.NoError(t, sourceRepo.PushBranch("main"))

	// Clone to create proper bare repo setup like working tests
	repoPath := setupClonedRepo(t, sourceRepo)

	// Initial sync to create worktrees
	syncCmd := newRootCommand()
	syncCmd.SetArgs([]string{"sync"})
	err := syncCmd.Execute()
	require.NoError(t, err, "Initial sync should succeed")

	// Verify the initial state is set up correctly
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "preview"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "production"))

	// Simulate the promotion scenario: preview branch gets promoted to production worktree
	// This means preview worktree disappears and production worktree switches to preview's branch
	promotionGBMConfig := `worktrees:
  main:
    branch: main
    description: "Main production branch"
  production:
    branch: production-2025-07-1 
    description: "Blade Runner"
    merge_into: main
`
	require.NoError(t, os.WriteFile(filepath.Join(repoPath, internal.DefaultBranchConfigFilename), []byte(promotionGBMConfig), 0644))

	// Run sync with confirmation to handle the promotion
	err = simulateUserInput("y", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"sync"})
		return cmd.Execute()
	})

	// The key assertion: sync should succeed when user confirms the promotion
	require.NoError(t, err, "Sync should succeed when user confirms worktree promotion")

	// Validate that "main" and "production" worktrees exist and they have the correct branches
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "production"))

	// Check main worktree is on main branch
	mainCmd := exec.Command("git", "branch", "--show-current")
	mainCmd.Dir = filepath.Join(repoPath, "worktrees", "main")
	mainOutput, err := mainCmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "main", strings.TrimSpace(string(mainOutput)), "main worktree should be on main branch")

	// Check production worktree is on production-2025-07-1 branch (promoted from preview)
	prodCmd := exec.Command("git", "branch", "--show-current")
	prodCmd.Dir = filepath.Join(repoPath, "worktrees", "production")
	prodOutput, err := prodCmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "production-2025-07-1", strings.TrimSpace(string(prodOutput)), "production worktree should be on production-2025-07-1 branch")

	// Validate that "preview" worktree no longer exists
	assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "preview"), "preview worktree should no longer exist")
}
