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

func setupValidateTest(t *testing.T, repo *testutils.GitTestRepo, envrcMapping map[string]string) {
	if envrcMapping != nil {
		err := repo.CreateEnvrc(envrcMapping)
		require.NoError(t, err, "Failed to create .envrc")

		err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
		require.NoError(t, err, "Failed to commit .envrc")
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
		if strings.Contains(line, "ENV VARIABLE") ||
			strings.Contains(line, "┌") ||
			strings.Contains(line, "├") ||
			strings.Contains(line, "└") ||
			strings.TrimSpace(line) == "" {
			continue
		}

		// Parse table rows (format: "│ ENV_VAR │ BRANCH │ STATUS │")
		if strings.Contains(line, "│") {
			parts := strings.Split(line, "│")
			if len(parts) >= 4 {
				envVar := strings.TrimSpace(parts[1])
				status := strings.TrimSpace(parts[3])
				if envVar != "" && status != "" {
					results[envVar] = status
				}
			}
		}
	}

	return results
}

func assertValidationResults(t *testing.T, output string, expectedResults map[string]string) {
	results := parseValidationTable(output)

	for envVar, expectedStatus := range expectedResults {
		actualStatus, exists := results[envVar]
		assert.True(t, exists, "Environment variable %s should appear in validation results", envVar)
		assert.Contains(t, actualStatus, expectedStatus, "Environment variable %s should have status %s", envVar, expectedStatus)
	}
}

// Test functions

func TestValidateCommand_AllBranchesValid(t *testing.T) {
	repo := testutils.NewStandardEnvrcRepo(t)

	setupValidateTest(t, repo, nil) // .envrc already exists in StandardEnvrcRepo

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Validate command should succeed when all branches are valid")

	expectedResults := map[string]string{
		"MAIN": "VALID",
		"DEV":  "VALID",
		"FEAT": "VALID",
		"PROD": "VALID",
	}

	assertValidationResults(t, output, expectedResults)
}

func TestValidateCommand_SomeBranchesInvalid(t *testing.T) {
	repo := testutils.NewMultiBranchRepo(t)

	// Create .envrc with some valid and some invalid branch references
	envrcMapping := map[string]string{
		"MAIN":    "main",               // valid - exists
		"DEV":     "develop",            // valid - exists
		"FEAT":    "feature/auth",       // valid - exists
		"INVALID": "nonexistent-branch", // invalid - doesn't exist
		"MISSING": "another-missing",    // invalid - doesn't exist
	}

	setupValidateTest(t, repo, envrcMapping)

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when some branches are invalid")

	expectedResults := map[string]string{
		"MAIN":    "VALID",
		"DEV":     "VALID",
		"FEAT":    "VALID",
		"INVALID": "NOT FOUND",
		"MISSING": "NOT FOUND",
	}

	assertValidationResults(t, output, expectedResults)
}

func TestValidateCommand_BranchExistence(t *testing.T) {
	tests := []struct {
		name            string
		setupRepo       func(t *testing.T) *testutils.GitTestRepo
		envMapping      map[string]string
		shouldPass      bool
		expectedResults map[string]string
	}{
		{
			name:      "local branches only",
			setupRepo: testutils.NewMultiBranchRepo,
			envMapping: map[string]string{
				"MAIN": "main",
				"DEV":  "develop",
				"FEAT": "feature/auth",
				"PROD": "production/v1.0",
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"MAIN": "VALID",
				"DEV":  "VALID",
				"FEAT": "VALID",
				"PROD": "VALID",
			},
		},
		{
			name:      "remote branches only",
			setupRepo: testutils.NewMultiBranchRepo,
			envMapping: map[string]string{
				"MAIN": "main",
				"DEV":  "develop",
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"MAIN": "VALID",
				"DEV":  "VALID",
			},
		},
		{
			name:      "both local and remote branches",
			setupRepo: testutils.NewMultiBranchRepo,
			envMapping: map[string]string{
				"MAIN": "main",
				"DEV":  "develop",
				"FEAT": "feature/auth",
				"PROD": "production/v1.0",
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"MAIN": "VALID",
				"DEV":  "VALID",
				"FEAT": "VALID",
				"PROD": "VALID",
			},
		},
		{
			name:      "non-existent branches",
			setupRepo: testutils.NewBasicRepo,
			envMapping: map[string]string{
				"MAIN":    "main",
				"MISSING": "does-not-exist",
				"INVALID": "also-missing",
			},
			shouldPass: false,
			expectedResults: map[string]string{
				"MAIN":    "VALID",
				"MISSING": "NOT FOUND",
				"INVALID": "NOT FOUND",
			},
		},
		{
			name:      "branches with special characters/slashes",
			setupRepo: testutils.NewMultiBranchRepo,
			envMapping: map[string]string{
				"MAIN":    "main",
				"FEATURE": "feature/auth",
				"RELEASE": "production/v1.0",
			},
			shouldPass: true,
			expectedResults: map[string]string{
				"MAIN":    "VALID",
				"FEATURE": "VALID",
				"RELEASE": "VALID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo(t)
			setupValidateTest(t, repo, tt.envMapping)

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

func TestValidateCommand_MissingEnvrc(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Don't create .envrc file - setupValidateTest with nil mapping skips creation
	setupValidateTest(t, repo, nil)

	_, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when .envrc file is missing")

	// Check that error message mentions missing .envrc
	assert.ErrorContains(t, err, ".envrc", "Error message should mention .envrc file")
}

func TestValidateCommand_InvalidEnvrcSyntax(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create .envrc with invalid syntax (malformed content)
	err := repo.WriteFile(".envrc", "INVALID SYNTAX WITHOUT EQUALS\nALSO=MISSING=QUOTES\nVALID=branch")
	require.NoError(t, err, "Failed to create malformed .envrc")

	err = repo.CommitChangesWithForceAdd("Add malformed .envrc")
	require.NoError(t, err, "Failed to commit malformed .envrc")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail when .envrc has invalid syntax")

	// The error could be either a parsing error or validation failure
	// Since the current implementation may still parse some entries, we just verify it fails
}

func TestValidateCommand_EmptyEnvrc(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create empty .envrc file
	err := repo.WriteFile(".envrc", "")
	require.NoError(t, err, "Failed to create empty .envrc")

	err = repo.CommitChangesWithForceAdd("Add empty .envrc")
	require.NoError(t, err, "Failed to commit empty .envrc")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Validate command should succeed with empty .envrc")

	// With empty .envrc, there should be no validation results to display
	results := parseValidationTable(output)
	assert.Empty(t, results, "Empty .envrc should result in no validation entries")
}

func TestValidateCommand_NotInGitRepository(t *testing.T) {
	// Create a temporary directory that is not a git repository
	tempDir := t.TempDir()

	// Create .envrc file in non-git directory
	envrcPath := tempDir + "/.envrc"
	err := os.WriteFile(envrcPath, []byte("MAIN=main\n"), 0644)
	require.NoError(t, err, "Failed to create .envrc in non-git directory")

	_, err = runValidateCommand(t, tempDir, []string{})
	require.Error(t, err, "Validate command should fail when not in a git repository")

	// Check that error message mentions git repository
	assert.ErrorContains(t, err, "git", "Error message should mention git repository")
}

func TestValidateCommand_CorruptGitRepository(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create .envrc file first
	envrcMapping := map[string]string{
		"MAIN": "main",
	}
	setupValidateTest(t, repo, envrcMapping)

	// Corrupt the .git directory by removing essential files
	gitDir := repo.GetLocalPath() + "/.git"
	err := os.RemoveAll(gitDir + "/refs")
	require.NoError(t, err, "Failed to corrupt git repository")

	_, err = runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Validate command should fail with corrupted git repository")

	// The error should be related to git operations
	assert.Contains(t, err.Error(), "git", "Error message should mention git-related issue")
}

func TestValidateCommand_DuplicateEnvVars(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create .envrc with duplicate environment variables (last one wins)
	err := repo.WriteFile(".envrc", "MAIN=main\nMAIN=develop\nTEST=feature/auth")
	require.NoError(t, err, "Failed to create .envrc with duplicates")

	err = repo.CommitChangesWithForceAdd("Add .envrc with duplicates")
	require.NoError(t, err, "Failed to commit .envrc")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.Error(t, err, "Should fail because 'develop' branch doesn't exist")

	// Verify that the last occurrence wins (MAIN=develop, not MAIN=main)
	results := parseValidationTable(output)
	assert.Contains(t, results["MAIN"], "NOT FOUND", "Should use the last value 'develop' which doesn't exist")
	assert.Contains(t, results["TEST"], "NOT FOUND", "TEST branch should also not be found")
}

func TestValidateCommand_VeryLongBranchNames(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create branch with very long name
	longBranchName := "feature/very-long-branch-name-that-exceeds-normal-length-limits-and-tests-table-formatting"
	err := repo.CreateBranch(longBranchName, "Long branch content")
	require.NoError(t, err, "Failed to create long branch")

	// Create .envrc with very long branch name and environment variable
	envrcMapping := map[string]string{
		"MAIN":                                "main",
		"VERY_LONG_ENVIRONMENT_VARIABLE_NAME": longBranchName,
	}
	setupValidateTest(t, repo, envrcMapping)

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{})
	require.NoError(t, err, "Should succeed with long branch names")

	// Verify table formatting remains readable with long names
	results := parseValidationTable(output)
	assert.Contains(t, results["MAIN"], "VALID")
	assert.Contains(t, results["VERY_LONG_ENVIRONMENT_VARIABLE_NAME"], "VALID")

	// Verify table structure is maintained
	assert.Contains(t, output, "┌", "Table should have proper borders")
	assert.Contains(t, output, "│", "Table should have column separators")
}

func TestValidateCommand_CustomEnvrcPath(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create custom .envrc file in different location
	customPath := "custom-env-config"
	err := repo.WriteFile(customPath, "MAIN=main\nCUSTOM=main")
	require.NoError(t, err, "Failed to create custom config file")

	err = repo.CommitChangesWithForceAdd("Add custom config file")
	require.NoError(t, err, "Failed to commit custom config")

	output, err := runValidateCommand(t, repo.GetLocalPath(), []string{"--config", customPath})
	require.NoError(t, err, "Should succeed with custom config path")

	// Verify validation used the custom config file
	results := parseValidationTable(output)
	assert.Contains(t, results["MAIN"], "VALID")
	assert.Contains(t, results["CUSTOM"], "VALID")
	assert.Len(t, results, 2, "Should only have entries from custom config file")
}
