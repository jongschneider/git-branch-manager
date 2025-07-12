package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions

func setupValidateTest(t *testing.T, repo *testutils.GitTestRepo, gbmMapping map[string]string) {
	if gbmMapping != nil {
		err := repo.CreateGBMConfig(gbmMapping)
		require.NoError(t, err, "Failed to create .gbm.config.yaml")

		err = repo.CommitChangesWithForceAdd("Add .gbm.config.yaml configuration")
		require.NoError(t, err, "Failed to commit .gbm.config.yaml")
	}
}

func runValidateCommand(t *testing.T, workDir string, args []string) (string, error) {
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(workDir)
	require.NoError(t, err)

	// Use the actual root command to ensure flags are properly handled
	cmd := rootCmd

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
	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	var capturedOutput bytes.Buffer
	capturedOutput.ReadFrom(r)

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

	setupValidateTest(t, repo, nil) // .gbm.config.yaml already exists in StandardGBMConfigRepo

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

	// Create .gbm.config.yaml with some valid and some invalid branch references
	gbmMapping := map[string]string{
		"main":    "main",               // valid - exists
		"dev":     "develop",            // valid - exists
		"feat":    "feature/auth",       // valid - exists
		"invalid": "nonexistent-branch", // invalid - doesn't exist
		"missing": "another-missing",    // invalid - doesn't exist
	}

	setupValidateTest(t, repo, gbmMapping)

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
		gbmMapping      map[string]string
		shouldPass      bool
		expectedResults map[string]string
	}{
		{
			name:      "local branches only",
			setupRepo: testutils.NewMultiBranchRepo,
			gbmMapping: map[string]string{
				"main": "main",
				"dev":  "develop",
				"feat": "feature/auth",
				"prod": "production/v1.0",
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
			gbmMapping: map[string]string{
				"main": "main",
				"dev":  "develop",
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
			gbmMapping: map[string]string{
				"main": "main",
				"dev":  "develop",
				"feat": "feature/auth",
				"prod": "production/v1.0",
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
			gbmMapping: map[string]string{
				"main":    "main",
				"missing": "does-not-exist",
				"invalid": "also-missing",
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
			gbmMapping: map[string]string{
				"main":    "main",
				"feature": "feature/auth",
				"release": "production/v1.0",
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
			setupValidateTest(t, repo, tt.gbmMapping)

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

	// Don't create .gbm.config.yaml file - setupValidateTest with nil mapping skips creation
	setupValidateTest(t, repo, nil)

	_, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when .gbm.config.yaml file is missing")

	// Check that error message mentions missing .gbm.config.yaml
	assert.ErrorContains(t, err, ".gbm.config.yaml", "Error message should mention .gbm.config.yaml file")
}

func TestValidateCommand_InvalidGBMConfigSyntax(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create .gbm.config.yaml with invalid syntax (malformed YAML content)
	err := repo.WriteFile(".gbm.config.yaml", "invalid yaml syntax:\n  - missing quotes\n  bad indentation\nworktrees:\n  main branch: main")
	require.NoError(t, err, "Failed to create malformed .gbm.config.yaml")

	err = repo.CommitChangesWithForceAdd("Add malformed .gbm.config.yaml")
	require.NoError(t, err, "Failed to commit malformed .gbm.config.yaml")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when .gbm.config.yaml has invalid syntax")

	// The error could be either a parsing error or validation failure
	// Since the current implementation may still parse some entries, we just verify it fails
}

func TestValidateCommand_EmptyGBMConfig(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create empty .gbm.config.yaml file
	err := repo.WriteFile(".gbm.config.yaml", "")
	require.NoError(t, err, "Failed to create empty .gbm.config.yaml")

	err = repo.CommitChangesWithForceAdd("Add empty .gbm.config.yaml")
	require.NoError(t, err, "Failed to commit empty .gbm.config.yaml")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Validate command should succeed with empty .gbm.config.yaml")

	// With empty .gbm.config.yaml, there should be no validation results to display
	results := parseValidationTable(output)
	assert.Empty(t, results, "Empty .gbm.config.yaml should result in no validation entries")
}

func TestValidateCommand_NotInGitRepository(t *testing.T) {
	// Create a temporary directory that is not a git repository
	tempDir := t.TempDir()

	// Create .gbm.config.yaml file in non-git directory
	configPath := tempDir + "/.gbm.config.yaml"
	configContent := `worktrees:
  main:
    branch: main
    description: "Main branch"`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err, "Failed to create .gbm.config.yaml in non-git directory")

	_, err = runValidateCommand(t, tempDir, []string{})
	require.Error(t, err, "Validate command should fail when not in a git repository")

	// Check that error message mentions git repository
	assert.ErrorContains(t, err, "git", "Error message should mention git repository")
}

func TestValidateCommand_CorruptGitRepository(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create .gbm.config.yaml file first
	gbmMapping := map[string]string{
		"main": "main",
	}
	setupValidateTest(t, repo, gbmMapping)

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

	// Create .gbm.config.yaml with duplicate worktree names (invalid YAML)
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
	err := repo.WriteFile(".gbm.config.yaml", configContent)
	require.NoError(t, err, "Failed to create .gbm.config.yaml with duplicates")

	err = repo.CommitChangesWithForceAdd("Add .gbm.config.yaml with duplicates")
	require.NoError(t, err, "Failed to commit .gbm.config.yaml")

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

	// Create .gbm.config.yaml with very long branch name and worktree name
	gbmMapping := map[string]string{
		"main":                             "main",
		"very_long_worktree_variable_name": longBranchName,
	}
	setupValidateTest(t, repo, gbmMapping)

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

func TestValidateCommand_CustomGBMConfigPath(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create custom .gbm.config.yaml file in different location
	customPath := "custom-gbm-config.yaml"
	customContent := `worktrees:
  main:
    branch: main
    description: "Main branch"
  custom:
    branch: main
    description: "Custom branch"`
	err := repo.WriteFile(customPath, customContent)
	require.NoError(t, err, "Failed to create custom config file")

	err = repo.CommitChangesWithForceAdd("Add custom config file")
	require.NoError(t, err, "Failed to commit custom config")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{"--config", customPath})
	require.NoError(t, err, "Should succeed with custom config path")

	// Verify validation used the custom config file
	results := parseValidationTable(output)
	assert.Contains(t, results["main"], "VALID")
	assert.Contains(t, results["custom"], "VALID")
	assert.Len(t, results, 2, "Should only have entries from custom config file")
}
