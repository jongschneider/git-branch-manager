package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions

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
	// Unit-test path with mock: all branches exist
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
				"dev":  "develop",
				"feat": "feature/auth",
				"prod": "production/v1.0",
			}, nil
		},
		BranchExistsFunc: func(branch string) (bool, error) { return true, nil },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	require.NoError(t, err)
	expectedResults := map[string]string{
		"main": "VALID",
		"dev":  "VALID",
		"feat": "VALID",
		"prod": "VALID",
	}
	assertValidationResults(t, output, expectedResults)
}

func TestValidateCommand_SomeBranchesInvalid(t *testing.T) {
	// Unit-test path with mock: some missing
	missing := map[string]bool{"nonexistent-branch": true, "another-missing": true}
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main":    "main",
				"dev":     "develop",
				"feat":    "feature/auth",
				"invalid": "nonexistent-branch",
				"missing": "another-missing",
			}, nil
		},
		BranchExistsFunc: func(branch string) (bool, error) { return !missing[branch], nil },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	require.Error(t, err)
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
	// Focused unit tests with mock covering combinations
	cases := []struct {
		name            string
		mapping         map[string]string
		existsSet       map[string]bool
		wantErr         bool
		expectedResults map[string]string
	}{
		{
			name: "all valid",
			mapping: map[string]string{
				"main": "main", "dev": "develop", "feat": "feature/auth", "prod": "production/v1.0",
			},
			existsSet:       map[string]bool{"main": true, "develop": true, "feature/auth": true, "production/v1.0": true},
			wantErr:         false,
			expectedResults: map[string]string{"main": "VALID", "dev": "VALID", "feat": "VALID", "prod": "VALID"},
		},
		{
			name: "some missing",
			mapping: map[string]string{
				"main": "main", "missing": "does-not-exist", "invalid": "also-missing",
			},
			existsSet:       map[string]bool{"main": true},
			wantErr:         true,
			expectedResults: map[string]string{"main": "VALID", "missing": "NOT FOUND", "invalid": "NOT FOUND"},
		},
		{
			name:            "error on branch check",
			mapping:         map[string]string{"main": "main"},
			existsSet:       map[string]bool{},
			wantErr:         true,
			expectedResults: map[string]string{"main": "ERROR"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &worktreeValidatorMock{
				GetWorktreeMappingFunc: func() (map[string]string, error) { return tc.mapping, nil },
				BranchExistsFunc: func(branch string) (bool, error) {
					if tc.name == "error on branch check" {
						return false, errors.New("boom")
					}
					return tc.existsSet[branch], nil
				},
			}

			var buf bytes.Buffer
			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := handleValidate(mock)

			_ = w.Close()
			os.Stdout = stdout
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assertValidationResults(t, output, tc.expectedResults)
		})
	}
}

func TestValidateCommand_MissingGBMConfig(t *testing.T) {
	// Unit-test path: propagate error from GetWorktreeMapping
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return nil, errors.New("missing " + internal.DefaultBranchConfigFilename)
		},
		BranchExistsFunc: func(string) (bool, error) { return false, nil },
	}

	err := handleValidate(mock)
	require.Error(t, err)
	assert.ErrorContains(t, err, internal.DefaultBranchConfigFilename)
}

func TestValidateCommand_InvalidGBMConfigSyntax(t *testing.T) {
	// Unit-test path: malformed config error surfaced from GetWorktreeMapping
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) { return nil, errors.New("invalid yaml syntax") },
		BranchExistsFunc:       func(string) (bool, error) { return false, nil },
	}

	err := handleValidate(mock)
	require.Error(t, err)
}

func TestValidateCommand_EmptyGBMConfig(t *testing.T) {
	// Unit-test path: empty mapping should succeed and print nothing
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) { return map[string]string{}, nil },
		BranchExistsFunc:       func(string) (bool, error) { return true, nil },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	require.NoError(t, err)
	results := parseValidationTable(output)
	assert.Empty(t, results)
}

func TestValidateCommand_NotInGitRepository(t *testing.T) {
	// Unit-test path: simulate git error bubbling from BranchExists on any branch
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{"main": "main"}, nil
		},
		BranchExistsFunc: func(string) (bool, error) { return false, errors.New("git error: not a repository") },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)

	require.Error(t, err)
}

func TestValidateCommand_CorruptGitRepository(t *testing.T) {
	// Unit-test path: propagate a git-related error during branch check
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) { return map[string]string{"main": "main"}, nil },
		BranchExistsFunc:       func(string) (bool, error) { return false, errors.New("git: corrupted repo") },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)

	require.Error(t, err)
}

func TestValidateCommand_DuplicateWorktrees(t *testing.T) {
	// Unit-test path: duplicate mapping error surfaced from GetWorktreeMapping
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) { return nil, errors.New("mapping key \"main\" already defined") },
		BranchExistsFunc:       func(string) (bool, error) { return false, nil },
	}

	err := handleValidate(mock)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mapping key \"main\" already defined")
}

func TestValidateCommand_VeryLongBranchNames(t *testing.T) {
	// Unit-test path: long names still produce a readable table
	longBranchName := "feature/very-long-branch-name-that-exceeds-normal-length-limits-and-tests-table-formatting"
	mock := &worktreeValidatorMock{
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main":                             "main",
				"very_long_worktree_variable_name": longBranchName,
			}, nil
		},
		BranchExistsFunc: func(string) (bool, error) { return true, nil },
	}

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handleValidate(mock)

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	require.NoError(t, err)
	results := parseValidationTable(output)
	assert.Contains(t, results["main"], "VALID")
	assert.Contains(t, results["very_long_worktree_variable_name"], "VALID")
	assert.Contains(t, output, "┌")
	assert.Contains(t, output, "│")
}
