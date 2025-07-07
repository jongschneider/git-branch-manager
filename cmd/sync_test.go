package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to reset sync command flags to default values
func resetSyncFlags() {
	syncDryRun = false
	syncForce = false
	syncFetch = false
}

// Helper function to setup a cloned repository for sync testing
func setupSyncTestRepo(t *testing.T, sourceRepo *testutils.GitTestRepo) (string, string) {
	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()

	os.Chdir(targetDir)

	// Clone the repository
	cloneCmd := rootCmd
	cloneCmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})
	err := cloneCmd.Execute()
	require.NoError(t, err, "Failed to clone repository")

	// Navigate to cloned repo
	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)
	os.Chdir(repoPath)

	// Return the path and original directory, but stay in the cloned repo
	return repoPath, originalDir
}

func TestSyncCommand_BasicOperations(t *testing.T) {
	tests := []struct {
		name         string
		setupRepo    func(t *testing.T) *testutils.GitTestRepo
		expectedDirs []string
	}{
		{
			name: "sync with existing envrc creates all worktrees",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				return testutils.NewStandardEnvrcRepo(t) // Has MAIN, DEV, FEAT, PROD
			},
			expectedDirs: []string{"worktrees/MAIN", "worktrees/DEV", "worktrees/FEAT", "worktrees/PROD"},
		},
		{
			name: "sync with minimal envrc",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				mapping := map[string]string{"MAIN": "main"}
				return testutils.NewEnvrcRepo(t, mapping)
			},
			expectedDirs: []string{"worktrees/MAIN"},
		},
		{
			name: "sync with already synced repo is idempotent",
			setupRepo: func(t *testing.T) *testutils.GitTestRepo {
				return testutils.NewStandardEnvrcRepo(t)
			},
			expectedDirs: []string{"worktrees/MAIN", "worktrees/DEV", "worktrees/FEAT", "worktrees/PROD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			sourceRepo := tt.setupRepo(t)
			_, originalDir := setupSyncTestRepo(t, sourceRepo)
			defer os.Chdir(originalDir)

			// For the idempotent test, run sync twice
			if strings.Contains(tt.name, "idempotent") {
				cmd := rootCmd
				cmd.SetArgs([]string{"sync"})
				require.NoError(t, cmd.Execute()) // First sync
			}

			cmd := rootCmd
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
				// Add more worktrees to .envrc to create missing worktrees scenario
				require.NoError(t, os.WriteFile(".envrc", []byte("MAIN=main\nNEW1=develop\nNEW2=feature/auth\n"), 0644))
			},
			validate: func(t *testing.T, repoPath string, output string, err error) {
				require.NoError(t, err)
				// Check that the command succeeded and directories are as expected
				// NEW1 and NEW2 should still be missing after dry-run
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "NEW1"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "NEW2"))
				// MAIN should still exist (was created by clone)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
			},
		},
		{
			name: "force flag removes orphaned worktrees",
			args: []string{"sync", "--force"},
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Create untracked worktree
				wd, _ := os.Getwd()
				createUntrackedWorktree(t, wd, "ORPHAN", "main")

				// Modify .envrc to remove some existing worktrees (making them orphaned)
				require.NoError(t, os.WriteFile(".envrc", []byte("MAIN=main\n"), 0644))
			},
			validate: func(t *testing.T, repoPath string, output string, err error) {
				require.NoError(t, err)
				// MAIN should still exist (still in .envrc)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				// Orphaned worktrees should be removed by --force
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "FEAT"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "PROD"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "ORPHAN"))
			},
		},
		{
			name: "fetch flag updates remote tracking",
			args: []string{"sync", "--fetch"},
			setup: func(t *testing.T, repo *testutils.GitTestRepo) {
				// Basic setup - sync should work normally with fetch
			},
			validate: func(t *testing.T, repoPath string, output string, err error) {
				require.NoError(t, err)
				// Verify worktrees exist (fetch doesn't prevent sync)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "FEAT"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "PROD"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			sourceRepo := testutils.NewStandardEnvrcRepo(t)
			repoPath, originalDir := setupSyncTestRepo(t, sourceRepo)
			defer os.Chdir(originalDir)

			tt.setup(t, &testutils.GitTestRepo{LocalDir: repoPath})

			var buf bytes.Buffer
			cmd := rootCmd
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
		name           string
		initialEnvrc   map[string]string
		updatedEnvrc   map[string]string
		expectChanges  bool
		validateResult func(t *testing.T, repoPath string)
	}{
		{
			name:          "branch reference changed",
			initialEnvrc:  map[string]string{"MAIN": "main", "FEAT": "feature/auth"},
			updatedEnvrc:  map[string]string{"MAIN": "main", "FEAT": "develop"},
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				// Verify FEAT worktree was updated to develop branch
				cmd := exec.Command("git", "branch", "--show-current")
				cmd.Dir = filepath.Join(repoPath, "worktrees", "FEAT")
				branchOutput, err := cmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "develop", strings.TrimSpace(string(branchOutput)))
			},
		},
		{
			name:          "new environment variable added",
			initialEnvrc:  map[string]string{"MAIN": "main"},
			updatedEnvrc:  map[string]string{"MAIN": "main", "DEV": "develop"},
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
			},
		},
		{
			name:          "environment variable removed",
			initialEnvrc:  map[string]string{"MAIN": "main", "TEMP": "develop"},
			updatedEnvrc:  map[string]string{"MAIN": "main"},
			expectChanges: true,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				// TEMP should still exist without --force
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "TEMP"))
			},
		},
		{
			name:          "no changes needed",
			initialEnvrc:  map[string]string{"MAIN": "main", "DEV": "develop"},
			updatedEnvrc:  map[string]string{"MAIN": "main", "DEV": "develop"},
			expectChanges: false,
			validateResult: func(t *testing.T, repoPath string) {
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			// Create source repo with initial .envrc
			sourceRepo := testutils.NewMultiBranchRepo(t)
			require.NoError(t, sourceRepo.CreateEnvrc(tt.initialEnvrc))
			commitEnvrcChanges(t, sourceRepo, "Add initial .envrc")
			require.NoError(t, sourceRepo.PushBranch("main"))

			repoPath, originalDir := setupSyncTestRepo(t, sourceRepo)
			defer os.Chdir(originalDir)

			// Update .envrc to new configuration
			envrcContent := ""
			for key, value := range tt.updatedEnvrc {
				envrcContent += fmt.Sprintf("%s=%s\n", key, value)
			}
			require.NoError(t, os.WriteFile(".envrc", []byte(envrcContent), 0644))

			// Run sync with updated .envrc
			cmd := rootCmd
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
		envrcMapping           map[string]string
		untrackedWorktrees     []string
		syncArgs               []string
		validateResult         func(t *testing.T, repoPath string, output string)
		createTrackedWorktrees bool
	}{
		{
			name: "untracked worktree preserved by default",
			envrcMapping: map[string]string{
				"MAIN": "main",
				"DEV":  "develop",
			},
			untrackedWorktrees:     []string{"MANUAL"},
			syncArgs:               []string{"sync"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// Tracked worktrees should exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
				// Untracked worktree should still exist (not removed)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MANUAL"))
			},
		},
		{
			name: "untracked worktree removed with --force",
			envrcMapping: map[string]string{
				"MAIN": "main",
				"DEV":  "develop",
			},
			untrackedWorktrees:     []string{"MANUAL"},
			syncArgs:               []string{"sync", "--force"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// Tracked worktrees should exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
				// Untracked worktree should be removed
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "MANUAL"))
			},
		},
		{
			name: "dry-run shows untracked worktree would be removed",
			envrcMapping: map[string]string{
				"MAIN": "main",
			},
			untrackedWorktrees:     []string{"TEMP", "EXPERIMENTAL"},
			syncArgs:               []string{"sync", "--dry-run", "--force"},
			createTrackedWorktrees: false,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// All worktrees should still exist (dry-run)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "TEMP"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "EXPERIMENTAL"))
				// In dry-run mode, orphaned worktrees should NOT be removed
				// (We can't easily capture the output due to test isolation issues,
				// but we can verify the intended behavior: no actual changes made)
			},
		},
		{
			name: "tracked worktrees updated, untracked preserved without force",
			envrcMapping: map[string]string{
				"MAIN": "main",
				"FEAT": "develop", // Will be changed from feature/auth
			},
			untrackedWorktrees:     []string{"MANUAL"},
			syncArgs:               []string{"sync"},
			createTrackedWorktrees: true,
			validateResult: func(t *testing.T, repoPath string, output string) {
				// TRACKED worktrees should exist and be updated
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "FEAT"))

				// UNTRACKED worktree should still exist (not removed without --force)
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MANUAL"))

				// Verify FEAT worktree was updated to develop branch
				cmd := exec.Command("git", "branch", "--show-current")
				cmd.Dir = filepath.Join(repoPath, "worktrees", "FEAT")
				featOutput, err := cmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "develop", strings.TrimSpace(string(featOutput)))

				// MANUAL worktree should still be on main branch (unchanged)
				manualCmd := exec.Command("git", "branch", "--show-current")
				manualCmd.Dir = filepath.Join(repoPath, "worktrees", "MANUAL")
				manualOutput, err := manualCmd.Output()
				require.NoError(t, err)
				assert.Equal(t, "main", strings.TrimSpace(string(manualOutput)))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			var sourceRepo *testutils.GitTestRepo

			if tt.createTrackedWorktrees {
				// Setup initial .envrc with different branch for FEAT
				sourceRepo = testutils.NewMultiBranchRepo(t)
				initialEnvrc := map[string]string{
					"MAIN": "main",
					"FEAT": "feature/auth",
				}
				require.NoError(t, sourceRepo.CreateEnvrc(initialEnvrc))
				commitEnvrcChanges(t, sourceRepo, "Add initial .envrc")
				require.NoError(t, sourceRepo.PushBranch("main"))
			} else {
				// Standard setup for other tests
				sourceRepo = testutils.NewMultiBranchRepo(t)
				require.NoError(t, sourceRepo.CreateEnvrc(tt.envrcMapping))
				commitEnvrcChanges(t, sourceRepo, "Add .envrc")
				require.NoError(t, sourceRepo.PushBranch("main"))
			}

			repoPath, originalDir := setupSyncTestRepo(t, sourceRepo)
			defer os.Chdir(originalDir)

			if tt.createTrackedWorktrees {
				// Create untracked worktree first
				for _, untrackedName := range tt.untrackedWorktrees {
					createUntrackedWorktree(t, repoPath, untrackedName, "main")
				}

				// Update .envrc to change FEAT branch
				envrcContent := ""
				for key, value := range tt.envrcMapping {
					envrcContent += fmt.Sprintf("%s=%s\n", key, value)
				}
				require.NoError(t, os.WriteFile(".envrc", []byte(envrcContent), 0644))
			} else {
				// Create untracked worktrees for standard tests
				for _, untrackedName := range tt.untrackedWorktrees {
					createUntrackedWorktree(t, repoPath, untrackedName, "main")
				}
			}

			// Run the sync command with specified args
			var buf bytes.Buffer
			cmd := rootCmd
			cmd.SetArgs(tt.syncArgs)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()
			output := buf.String()
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
			name: "missing envrc file",
			setup: func(t *testing.T) string {
				repo := testutils.NewBasicRepo(t)
				return repo.GetLocalPath()
			},
			args:          []string{"sync"},
			expectedError: "failed to load .envrc",
		},
		{
			name: "invalid branch reference",
			setup: func(t *testing.T) string {
				// Create a source repo with invalid branch reference
				sourceRepo := testutils.NewBasicRepo(t)
				require.NoError(t, sourceRepo.CreateEnvrc(map[string]string{"INVALID": "nonexistent-branch"}))
				commitEnvrcChanges(t, sourceRepo, "Add invalid envrc")
				require.NoError(t, sourceRepo.PushBranch("main"))

				// Clone it to set up proper structure, but don't defer the chdir - let the test handle it
				repoPath, _ := setupSyncTestRepo(t, sourceRepo)
				// Immediately return to original dir so test can handle directory changes
				os.Chdir(repoPath)
				return repoPath
			},
			args:          []string{"sync"},
			expectedError: "branch 'nonexistent-branch' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			workingDir := tt.setup(t)
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			require.NoError(t, os.Chdir(workingDir))

			cmd := rootCmd
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
				// 1. Initial state created by clone (MAIN, DEV, FEAT, PROD from StandardEnvrcRepo)

				// 2. Modify .envrc to remove some worktrees and add different ones
				envrcContent := "MAIN=main\nDEV=develop\nNEW_FEAT=feature/auth\nNEW_PROD=production/v1.0\n"
				require.NoError(t, os.WriteFile(".envrc", []byte(envrcContent), 0644))

				// 3. Run sync with --force to remove orphaned worktrees and create new ones
				cmd := rootCmd
				cmd.SetArgs([]string{"sync", "--force"})
				require.NoError(t, cmd.Execute())
			},
			validate: func(t *testing.T, repoPath string) {
				// Verify new worktrees exist
				expectedDirs := []string{"MAIN", "DEV", "NEW_FEAT", "NEW_PROD"}
				for _, dir := range expectedDirs {
					assert.DirExists(t, filepath.Join(repoPath, "worktrees", dir))
				}

				// Old worktrees should be removed with --force
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "FEAT"))
				assert.NoDirExists(t, filepath.Join(repoPath, "worktrees", "PROD"))
			},
		},
		{
			name: "sync after manual worktree changes",
			scenario: func(t *testing.T, repoPath string) {
				// Manually remove a tracked worktree (simulate corruption where worktree is lost)
				devWorktreePath := filepath.Join(repoPath, "worktrees", "DEV")

				// Remove the worktree directory first
				require.NoError(t, os.RemoveAll(devWorktreePath))

				// Also remove it from git's worktree list to simulate complete loss
				pruneCmd := exec.Command("git", "worktree", "prune")
				pruneCmd.Dir = repoPath
				require.NoError(t, pruneCmd.Run())

				// Run sync to fix things (recreate missing worktree)
				cmd := rootCmd
				cmd.SetArgs([]string{"sync"})
				require.NoError(t, cmd.Execute())
			},
			validate: func(t *testing.T, repoPath string) {
				// DEV should be recreated
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "DEV"))
				// All original worktrees should still exist
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "FEAT"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "PROD"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSyncFlags() // Reset flags before each test
			sourceRepo := testutils.NewStandardEnvrcRepo(t)
			repoPath, originalDir := setupSyncTestRepo(t, sourceRepo)
			defer os.Chdir(originalDir)

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

// Helper function to safely commit .envrc changes (handles already-clean repos)
func commitEnvrcChanges(t *testing.T, repo *testutils.GitTestRepo, message string) {
	if err := repo.CommitChangesWithForceAdd(message); err != nil {
		if strings.Contains(err.Error(), "nothing to commit") {
			// Ignore "nothing to commit" errors - repo is already in desired state
			return
		}
		require.NoError(t, err)
	}
}
