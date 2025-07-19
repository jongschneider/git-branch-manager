package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gbm/internal"
	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions

func setupValidateTest(t *testing.T, repo *testutils.GitTestRepo, worktrees map[string]testutils.WorktreeConfig) {
	if worktrees != nil {
		err := repo.CreateGBMConfig(worktrees)
		require.NoError(t, err, "Failed to create gbm.branchconfig.yaml")

		err = repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml configuration")
		require.NoError(t, err, "Failed to commit gbm.branchconfig.yaml")
	}
}

func runValidateCommand(t *testing.T, workDir string, args []string) (string, error) {
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(workDir)
	require.NoError(t, err)

	// Use the actual root command to ensure flags are properly handled
	cmd := newRootCommand()

	// Capture both stdout and stderr since validate command uses fmt.Println
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Redirect os.Stdout temporarily to capture fmt.Println output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.SetArgs(append([]string{"validate"}, args...))
	err = cmd.Execute()

	// Restore stdout
	_ = w.Close()
	os.Stdout = originalStdout

	// Read captured output
	var capturedOutput bytes.Buffer
	_, _ = capturedOutput.ReadFrom(r)

	// Combine both outputs
	combined := buf.String() + capturedOutput.String()

	return combined, err
}

func parseValidationTable(output string) map[string]string {
	results := make(map[string]string)

	for line := range strings.SplitSeq(output, "\n") {
		// Skip header lines, separators, and empty lines
		if strings.Contains(line, "WORKTREE") ||
			strings.Contains(line, "┌") ||
			strings.Contains(line, "├") ||
			strings.Contains(line, "└") ||
			strings.TrimSpace(line) == "" {
			continue
		}

		// Parse table rows (format: "│ WORKTREE │ BRANCH │ STATUS │")
		if strings.Contains(line, "│") {
			parts := strings.Split(line, "│")
			if len(parts) >= 4 {
				worktree := strings.TrimSpace(parts[1])
				status := strings.TrimSpace(parts[3])
				if worktree != "" && status != "" {
					results[worktree] = status
				}
			}
		}
	}

	return results
}

func assertValidationResults(t *testing.T, output string, expectedResults map[string]string) {
	results := parseValidationTable(output)

	for worktree, expectedStatus := range expectedResults {
		actualStatus, exists := results[worktree]
		assert.True(t, exists, "Worktree %s should appear in validation results", worktree)
		assert.Contains(t, actualStatus, expectedStatus, "Worktree %s should have status %s", worktree, expectedStatus)
	}
}

// Test functions

func TestValidateCommand_AllBranchesValid(t *testing.T) {
	repo := testutils.NewStandardGBMConfigRepo(t)

	setupValidateTest(t, repo, nil) // gbm.branchconfig.yaml already exists in StandardGBMConfigRepo

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Validate command should succeed when all branches are valid")

	expectedResults := map[string]string{
		"main": "VALID",
		"dev":  "VALID",
		"feat": "VALID",
		"prod": "VALID",
	}

	assertValidationResults(t, output, expectedResults)
}

func TestValidateCommand_SomeBranchesInvalid(t *testing.T) {
	repo := testutils.NewMultiBranchRepo(t)

	// Create gbm.branchconfig.yaml with some valid and some invalid branch references
	worktrees := map[string]testutils.WorktreeConfig{
		"main": {
			Branch:      "main",               // valid - exists
			Description: "Main branch",
		},
		"dev": {
			Branch:      "develop",            // valid - exists
			MergeInto:   "main",
			Description: "Dev branch",
		},
		"feat": {
			Branch:      "feature/auth",       // valid - exists
			MergeInto:   "dev",
			Description: "Feat branch",
		},
		"invalid": {
			Branch:      "nonexistent-branch", // invalid - doesn't exist
			MergeInto:   "feat",
			Description: "Invalid branch",
		},
		"missing": {
			Branch:      "another-missing",    // invalid - doesn't exist
			MergeInto:   "invalid",
			Description: "Missing branch",
		},
	}

	setupValidateTest(t, repo, worktrees)

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when some branches are invalid")

	expectedResults := map[string]string{
		"main":    "VALID",
		"dev":     "VALID",
		"feat":    "VALID",
		"invalid": "NOT FOUND",
		"missing": "NOT FOUND",
	}

	assertValidationResults(t, output, expectedResults)
}

func TestValidateCommand_BranchExistence(t *testing.T) {
	tests := []struct {
		name            string
		setupRepo       func(t *testing.T) *testutils.GitTestRepo
		worktrees       map[string]testutils.WorktreeConfig
		shouldPass      bool
		expectedResults map[string]string
	}{
		{
			name:      "local branches only",
			setupRepo: testutils.NewMultiBranchRepo,
			worktrees: map[string]testutils.WorktreeConfig{
				"main": {
					Branch:      "main",
					Description: "Main branch",
				},
				"dev": {
					Branch:      "develop",
					MergeInto:   "main",
					Description: "Dev branch",
				},
				"feat": {
					Branch:      "feature/auth",
					MergeInto:   "dev",
					Description: "Feat branch",
				},
				"prod": {
					Branch:      "production/v1.0",
					MergeInto:   "feat",
					Description: "Prod branch",
				},
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"main": "VALID",
				"dev":  "VALID",
				"feat": "VALID",
				"prod": "VALID",
			},
		},
		{
			name:      "remote branches only",
			setupRepo: testutils.NewMultiBranchRepo,
			worktrees: map[string]testutils.WorktreeConfig{
				"main": {
					Branch:      "main",
					Description: "Main branch",
				},
				"dev": {
					Branch:      "develop",
					MergeInto:   "main",
					Description: "Dev branch",
				},
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"main": "VALID",
				"dev":  "VALID",
			},
		},
		{
			name:      "both local and remote branches",
			setupRepo: testutils.NewMultiBranchRepo,
			worktrees: map[string]testutils.WorktreeConfig{
				"main": {
					Branch:      "main",
					Description: "Main branch",
				},
				"dev": {
					Branch:      "develop",
					MergeInto:   "main",
					Description: "Dev branch",
				},
				"feat": {
					Branch:      "feature/auth",
					MergeInto:   "dev",
					Description: "Feat branch",
				},
				"prod": {
					Branch:      "production/v1.0",
					MergeInto:   "feat",
					Description: "Prod branch",
				},
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"main": "VALID",
				"dev":  "VALID",
				"feat": "VALID",
				"prod": "VALID",
			},
		},
		{
			name:      "non-existent branches",
			setupRepo: testutils.NewBasicRepo,
			worktrees: map[string]testutils.WorktreeConfig{
				"main": {
					Branch:      "main",
					Description: "Main branch",
				},
				"missing": {
					Branch:      "does-not-exist",
					MergeInto:   "main",
					Description: "Missing branch",
				},
				"invalid": {
					Branch:      "also-missing",
					MergeInto:   "missing",
					Description: "Invalid branch",
				},
			},
			shouldPass: false,
			expectedResults: map[string]string{
				"main":    "VALID",
				"missing": "NOT FOUND",
				"invalid": "NOT FOUND",
			},
		},
		{
			name:      "branches with special characters/slashes",
			setupRepo: testutils.NewMultiBranchRepo,
			worktrees: map[string]testutils.WorktreeConfig{
				"main": {
					Branch:      "main",
					Description: "Main branch",
				},
				"feature": {
					Branch:      "feature/auth",
					MergeInto:   "main",
					Description: "Feature branch",
				},
				"release": {
					Branch:      "production/v1.0",
					MergeInto:   "feature",
					Description: "Release branch",
				},
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"main":    "VALID",
				"feature": "VALID",
				"release": "VALID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo(t)
			setupValidateTest(t, repo, tt.worktrees)

			output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})

			if tt.shouldPass {
				require.NoError(t, err, "Validate command should succeed for test case: %s", tt.name)
			} else {
				require.Error(t, err, "Validate command should fail for test case: %s", tt.name)
			}

			assertValidationResults(t, output, tt.expectedResults)
		})
	}
}

func TestValidateCommand_MissingGBMConfig(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Don't create gbm.branchconfig.yaml file - setupValidateTest with nil mapping skips creation
	setupValidateTest(t, repo, nil)

	_, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when gbm.branchconfig.yaml file is missing")

	// Check that error message mentions missing gbm.branchconfig.yaml
	assert.ErrorContains(t, err, internal.DefaultBranchConfigFilename, "Error message should mention gbm.branchconfig.yaml file")
}

func TestValidateCommand_InvalidGBMConfigSyntax(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create gbm.branchconfig.yaml with invalid syntax (malformed YAML content)
	err := repo.WriteFile(internal.DefaultBranchConfigFilename, "invalid yaml syntax:\n  - missing quotes\n  bad indentation\nworktrees:\n  main branch: main")
	require.NoError(t, err, "Failed to create malformed gbm.branchconfig.yaml")

	err = repo.CommitChangesWithForceAdd("Add malformed gbm.branchconfig.yaml")
	require.NoError(t, err, "Failed to commit malformed gbm.branchconfig.yaml")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when gbm.branchconfig.yaml has invalid syntax")

	// The error could be either a parsing error or validation failure
	// Since the current implementation may still parse some entries, we just verify it fails
}

func TestValidateCommand_EmptyGBMConfig(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create empty gbm.branchconfig.yaml file
	err := repo.WriteFile(internal.DefaultBranchConfigFilename, "")
	require.NoError(t, err, "Failed to create empty gbm.branchconfig.yaml")

	err = repo.CommitChangesWithForceAdd("Add empty gbm.branchconfig.yaml")
	require.NoError(t, err, "Failed to commit empty gbm.branchconfig.yaml")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Validate command should succeed with empty gbm.branchconfig.yaml")

	// With empty gbm.branchconfig.yaml, there should be no validation results to display
	results := parseValidationTable(output)
	assert.Empty(t, results, "Empty gbm.branchconfig.yaml should result in no validation entries")
}

func TestValidateCommand_NotInGitRepository(t *testing.T) {
	// Create a temporary directory that is not a git repository
	tempDir := t.TempDir()

	// Create gbm.branchconfig.yaml file in non-git directory
	configPath := tempDir + "/gbm.branchconfig.yaml"
	configContent := `worktrees:
  main:
    branch: main
    description: "Main branch"`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err, "Failed to create gbm.branchconfig.yaml in non-git directory")

	_, err = runValidateCommand(t, tempDir, []string{})
	require.Error(t, err, "Validate command should fail when not in a git repository")

	// Check that error message mentions git repository
	assert.ErrorContains(t, err, "git", "Error message should mention git repository")
}

func TestValidateCommand_CorruptGitRepository(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create gbm.branchconfig.yaml file first
	worktrees := map[string]testutils.WorktreeConfig{
		"main": {
			Branch:      "main",
			Description: "Main branch",
		},
	}
	setupValidateTest(t, repo, worktrees)

	// Corrupt the .git directory by removing essential files
	gitDir := repo.GetLocalPath() + "/.git"
	err := os.RemoveAll(gitDir + "/refs")
	require.NoError(t, err, "Failed to corrupt git repository")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail with corrupted git repository")

	// The error should be related to git operations
	assert.Contains(t, err.Error(), "git", "Error message should mention git-related issue")
}

func TestValidateCommand_DuplicateWorktrees(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create gbm.branchconfig.yaml with duplicate worktree names (invalid YAML)
	configContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  main:
    branch: develop
    description: "Duplicate main"
  test:
    branch: feature/auth
    description: "Test branch"`
	err := repo.WriteFile(internal.DefaultBranchConfigFilename, configContent)
	require.NoError(t, err, "Failed to create gbm.branchconfig.yaml with duplicates")

	err = repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml with duplicates")
	require.NoError(t, err, "Failed to commit gbm.branchconfig.yaml")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Should fail due to duplicate YAML keys")

	// Verify that the error mentions duplicate mapping keys
	assert.Contains(t, err.Error(), "mapping key \"main\" already defined", "Should fail with YAML duplicate key error")
}

func TestValidateCommand_VeryLongBranchNames(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create branch with very long name
	longBranchName := "feature/very-long-branch-name-that-exceeds-normal-length-limits-and-tests-table-formatting"
	err := repo.CreateBranch(longBranchName, "Long branch content")
	require.NoError(t, err, "Failed to create long branch")

	// Create gbm.branchconfig.yaml with very long branch name and worktree name
	worktrees := map[string]testutils.WorktreeConfig{
		"main": {
			Branch:      "main",
			Description: "Main branch",
		},
		"very_long_worktree_variable_name": {
			Branch:      longBranchName,
			MergeInto:   "main",
			Description: "Very long worktree variable name branch",
		},
	}
	setupValidateTest(t, repo, worktrees)

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Should succeed with long branch names")

	// Verify table formatting remains readable with long names
	results := parseValidationTable(output)
	assert.Contains(t, results["main"], "VALID")
	assert.Contains(t, results["very_long_worktree_variable_name"], "VALID")

	// Verify table structure is maintained
	assert.Contains(t, output, "┌", "Table should have proper borders")
	assert.Contains(t, output, "│", "Table should have column separators")
}

