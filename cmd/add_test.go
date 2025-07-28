package cmd

import (
	"testing"

	"gbm/internal"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// UNIT TESTS (Using mocks - these are fast and don't require real git operations)
// ============================================================================
// These tests use the worktreeAdder interface with mocks to test business logic
// without requiring real git operations. They run in milliseconds.

// Unit tests using table-driven approach with mocks for ArgsResolver
func TestArgsResolver_ResolveArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		newBranch bool
		mockSetup func() *worktreeAdderMock
		expectErr func(t *testing.T, err error)
		expect    func(t *testing.T, result *WorktreeArgs)
	}{
		{
			name:      "missing worktree name",
			args:      []string{},
			newBranch: false,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "worktree name is required")
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "new branch with default base",
			args:      []string{"test-worktree"},
			newBranch: true,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetDefaultBranchFunc: func() (string, error) {
						return "main", nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Equal(t, "test-worktree", result.WorktreeName)
				assert.Equal(t, "feature/test-worktree", result.BranchName)
				assert.True(t, result.NewBranch)
				assert.Equal(t, "main", result.ResolvedBaseBranch)
			},
		},
		{
			name:      "existing branch",
			args:      []string{"test-worktree", "existing-branch"},
			newBranch: false,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Equal(t, "test-worktree", result.WorktreeName)
				assert.Equal(t, "existing-branch", result.BranchName)
				assert.False(t, result.NewBranch)
				assert.Equal(t, "", result.ResolvedBaseBranch)
			},
		},
		{
			name:      "new branch with valid base branch",
			args:      []string{"test-worktree", "new-branch", "develop"},
			newBranch: true,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					BranchExistsFunc: func(branch string) (bool, error) {
						return branch == "develop", nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Equal(t, "test-worktree", result.WorktreeName)
				assert.Equal(t, "new-branch", result.BranchName)
				assert.True(t, result.NewBranch)
				assert.Equal(t, "develop", result.ResolvedBaseBranch)
			},
		},
		{
			name:      "new branch with invalid base branch",
			args:      []string{"test-worktree", "new-branch", "invalid-base"},
			newBranch: true,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					BranchExistsFunc: func(branch string) (bool, error) {
						return false, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "base branch 'invalid-base' does not exist")
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "JIRA key without branch name should suggest",
			args:      []string{"PROJ-123"},
			newBranch: false,
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GenerateBranchFromJiraFunc: func(jiraKey string) (string, error) {
						return "feature/PROJ-123-implement-feature", nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "branch name required")
				assert.Contains(t, err.Error(), "feature/PROJ-123-implement-feature")
				assert.Contains(t, err.Error(), "gbm add PROJ-123 feature/PROJ-123-implement-feature -b")
			},
			expect: func(t *testing.T, result *WorktreeArgs) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := tt.mockSetup()
			resolver := &ArgsResolver{manager: mockManager}

			result, err := resolver.ResolveArgs(tt.args, tt.newBranch)

			tt.expectErr(t, err)
			tt.expect(t, result)
		})
	}
}

func TestGenerateBranchName(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		mockSetup    func() *worktreeAdderMock
		expect       func(t *testing.T, result string)
	}{
		{
			name:         "JIRA key success",
			worktreeName: "PROJ-456",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GenerateBranchFromJiraFunc: func(jiraKey string) (string, error) {
						return "feature/PROJ-456-add-new-feature", nil
					},
				}
			},
			expect: func(t *testing.T, result string) {
				assert.Equal(t, "feature/PROJ-456-add-new-feature", result)
			},
		},
		{
			name:         "JIRA key with error fallback",
			worktreeName: "PROJ-789",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GenerateBranchFromJiraFunc: func(jiraKey string) (string, error) {
						return "", assert.AnError
					},
				}
			},
			expect: func(t *testing.T, result string) {
				assert.Equal(t, "feature/proj-789", result)
			},
		},
		{
			name:         "non-JIRA key",
			worktreeName: "my feature",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{}
			},
			expect: func(t *testing.T, result string) {
				assert.Equal(t, "feature/my-feature", result)
			},
		},
		{
			name:         "already has prefix",
			worktreeName: "bugfix/urgent-fix",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{}
			},
			expect: func(t *testing.T, result string) {
				assert.Equal(t, "bugfix/urgent-fix", result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := tt.mockSetup()

			result := generateBranchName(tt.worktreeName, mockManager)

			tt.expect(t, result)
		})
	}
}

// Table-driven unit tests for command functions using mocks
func TestAddCommand_ValidArgsFunction(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		toComplete string
		mockSetup  func() *worktreeAdderMock
		expectErr  func(t *testing.T, err error)
		expect     func(t *testing.T, completions []string, directive cobra.ShellCompDirective)
	}{
		{
			name:       "JIRA completion success",
			args:       []string{},
			toComplete: "PROJ",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetJiraIssuesFunc: func() ([]internal.JiraIssue, error) {
						return []internal.JiraIssue{
							{Key: "PROJ-123", Summary: "Implement new feature"},
							{Key: "PROJ-124", Summary: "Fix critical bug"},
						}, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				// No error expected for completion
			},
			expect: func(t *testing.T, completions []string, directive cobra.ShellCompDirective) {
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
				assert.Len(t, completions, 2)
				assert.Contains(t, completions, "PROJ-123\tImplement new feature")
				assert.Contains(t, completions, "PROJ-124\tFix critical bug")
			},
		},
		{
			name:       "JIRA completion error",
			args:       []string{},
			toComplete: "PROJ",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetJiraIssuesFunc: func() ([]internal.JiraIssue, error) {
						return nil, assert.AnError
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				// No error expected for completion failure
			},
			expect: func(t *testing.T, completions []string, directive cobra.ShellCompDirective) {
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
				assert.Nil(t, completions)
			},
		},
		{
			name:       "branch name completion for JIRA key",
			args:       []string{"PROJ-456"},
			toComplete: "",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GenerateBranchFromJiraFunc: func(jiraKey string) (string, error) {
						return "feature/PROJ-456-implement-feature", nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				// No error expected for completion
			},
			expect: func(t *testing.T, completions []string, directive cobra.ShellCompDirective) {
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
				assert.Len(t, completions, 1)
				assert.Equal(t, "feature/PROJ-456-implement-feature", completions[0])
			},
		},
		{
			name:       "branch name completion error fallback",
			args:       []string{"PROJ-789"},
			toComplete: "",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GenerateBranchFromJiraFunc: func(jiraKey string) (string, error) {
						return "", assert.AnError
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				// No error expected for completion
			},
			expect: func(t *testing.T, completions []string, directive cobra.ShellCompDirective) {
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
				assert.Len(t, completions, 1)
				assert.Equal(t, "feature/proj-789", completions[0])
			},
		},
		{
			name:       "non-JIRA key no completion",
			args:       []string{"regular-name"},
			toComplete: "",
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{}
			},
			expectErr: func(t *testing.T, err error) {
				// No error expected for completion
			},
			expect: func(t *testing.T, completions []string, directive cobra.ShellCompDirective) {
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
				assert.Nil(t, completions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := tt.mockSetup()
			cmd := newAddCommand(mockManager)

			completions, directive := cmd.ValidArgsFunction(cmd, tt.args, tt.toComplete)

			tt.expect(t, completions, directive)
		})
	}
}

func TestAddCommand_Execute(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mockSetup func() *worktreeAdderMock
		expectErr func(t *testing.T, err error)
		expect    func(t *testing.T, mock *worktreeAdderMock)
	}{
		{
			name: "successful execution",
			args: []string{"test-worktree", "-b"},
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetDefaultBranchFunc: func() (string, error) {
						return "main", nil
					},
					AddWorktreeFunc: func(worktreeName, branchName string, newBranch bool, baseBranch string) error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expect: func(t *testing.T, mock *worktreeAdderMock) {
				assert.Len(t, mock.GetDefaultBranchCalls(), 1)
				assert.Len(t, mock.AddWorktreeCalls(), 1)

				addCall := mock.AddWorktreeCalls()[0]
				assert.Equal(t, "test-worktree", addCall.WorktreeName)
				assert.Equal(t, "feature/test-worktree", addCall.BranchName)
				assert.True(t, addCall.NewBranch)
				assert.Equal(t, "main", addCall.BaseBranch)
			},
		},
		{
			name: "AddWorktree error",
			args: []string{"test-worktree", "-b"},
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetDefaultBranchFunc: func() (string, error) {
						return "main", nil
					},
					AddWorktreeFunc: func(worktreeName, branchName string, newBranch bool, baseBranch string) error {
						return assert.AnError
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to add worktree")
			},
			expect: func(t *testing.T, mock *worktreeAdderMock) {
				assert.Len(t, mock.GetDefaultBranchCalls(), 1)
				assert.Len(t, mock.AddWorktreeCalls(), 1)
			},
		},
		{
			name: "GetDefaultBranch error",
			args: []string{"test-worktree", "-b"},
			mockSetup: func() *worktreeAdderMock {
				return &worktreeAdderMock{
					GetDefaultBranchFunc: func() (string, error) {
						return "", assert.AnError
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
			expect: func(t *testing.T, mock *worktreeAdderMock) {
				assert.Len(t, mock.GetDefaultBranchCalls(), 1)
				assert.Len(t, mock.AddWorktreeCalls(), 0) // Should not reach AddWorktree
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := tt.mockSetup()
			cmd := newAddCommand(mockManager)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			tt.expectErr(t, err)
			tt.expect(t, mockManager)
		})
	}
}
