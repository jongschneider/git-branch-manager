package cmd

import (
	"errors"
	"strings"
	"testing"

	"gbm/internal"
)

func TestGenerateHotfixBranchName(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		jiraTicket   string
		expected     string
	}{
		{
			name:         "simple worktree name",
			worktreeName: "critical-bug",
			jiraTicket:   "",
			expected:     "hotfix/critical-bug",
		},
		{
			name:         "worktree name with spaces",
			worktreeName: "auth bug fix",
			jiraTicket:   "",
			expected:     "hotfix/auth-bug-fix",
		},
		{
			name:         "worktree name with underscores",
			worktreeName: "api_timeout_fix",
			jiraTicket:   "",
			expected:     "hotfix/api-timeout-fix",
		},
		{
			name:         "JIRA ticket as worktree name",
			worktreeName: "PROJECT-123",
			jiraTicket:   "PROJECT-123",
			expected:     "hotfix/PROJECT-123", // Fallback when JIRA API fails
		},
		{
			name:         "mixed case worktree name",
			worktreeName: "CriticalBug",
			jiraTicket:   "",
			expected:     "hotfix/criticalbug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pass nil manager to force fallback behavior in JIRA integration
			result, err := generateHotfixBranchName(tt.worktreeName, tt.jiraTicket, nil)
			if err != nil {
				t.Errorf("generateHotfixBranchName() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("generateHotfixBranchName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHotfixWorktreeNaming(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		hotfixPrefix string
		expected     string
	}{
		{
			name:         "default prefix",
			worktreeName: "PROJECT-123",
			hotfixPrefix: "HOTFIX",
			expected:     "HOTFIX_PROJECT-123",
		},
		{
			name:         "custom prefix",
			worktreeName: "critical-bug",
			hotfixPrefix: "FIX",
			expected:     "FIX_critical-bug",
		},
		{
			name:         "empty prefix",
			worktreeName: "auth-issue",
			hotfixPrefix: "",
			expected:     "auth-issue",
		},
		{
			name:         "single char prefix",
			worktreeName: "bug-123",
			hotfixPrefix: "H",
			expected:     "H_bug-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build worktree name using the same logic as the command
			var result string
			if tt.hotfixPrefix != "" {
				result = tt.hotfixPrefix + "_" + tt.worktreeName
			} else {
				result = tt.worktreeName
			}

			if result != tt.expected {
				t.Errorf("Hotfix worktree naming = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsProductionBranchName(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expected   bool
	}{
		{
			name:       "production branch",
			branchName: "production",
			expected:   true,
		},
		{
			name:       "prod branch",
			branchName: "prod",
			expected:   true,
		},
		{
			name:       "main branch",
			branchName: "main",
			expected:   true,
		},
		{
			name:       "master branch",
			branchName: "master",
			expected:   true,
		},
		{
			name:       "release branch",
			branchName: "release-v1.0",
			expected:   true,
		},
		{
			name:       "production with version",
			branchName: "production-v2.1",
			expected:   true,
		},
		{
			name:       "feature branch",
			branchName: "feature/new-api",
			expected:   false,
		},
		{
			name:       "development branch",
			branchName: "development",
			expected:   false,
		},
		{
			name:       "staging branch",
			branchName: "staging",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isProductionBranchName(tt.branchName)
			if result != tt.expected {
				t.Errorf("isProductionBranchName(%s) = %v, want %v", tt.branchName, result, tt.expected)
			}
		})
	}
}

func TestHandleHotfix(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		mockSetup     func(*hotfixCreatorMock)
		expectError   bool
		errorContains string
	}{
		{
			name: "successful hotfix creation with simple worktree name",
			args: []string{"critical-bug"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "main", nil
				}
				mock.GetConfigFunc = func() *internal.Config {
					return &internal.Config{
						Settings: internal.ConfigSettings{
							HotfixPrefix: "HOTFIX",
						},
					}
				}
				mock.GetGBMConfigFunc = func() *internal.GBMConfig {
					return &internal.GBMConfig{
						Worktrees: map[string]internal.WorktreeConfig{
							"main": {Branch: "main", MergeInto: ""},
						},
					}
				}
				mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
					if worktreeName != "HOTFIX_critical-bug" {
						t.Errorf("Expected worktree name 'HOTFIX_critical-bug', got '%s'", worktreeName)
					}
					if branchName != "hotfix/critical-bug" {
						t.Errorf("Expected branch name 'hotfix/critical-bug', got '%s'", branchName)
					}
					if baseBranch != "main" {
						t.Errorf("Expected base branch 'main', got '%s'", baseBranch)
					}
					if !createBranch {
						t.Errorf("Expected createBranch to be true")
					}
					return nil
				}
			},
			expectError: false,
		},
		{
			name: "successful hotfix creation with JIRA ticket",
			args: []string{"PROJECT-123"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "production", nil
				}
				mock.GetConfigFunc = func() *internal.Config {
					return &internal.Config{
						Settings: internal.ConfigSettings{
							HotfixPrefix: "FIX",
						},
					}
				}
				mock.GetGBMConfigFunc = func() *internal.GBMConfig {
					return &internal.GBMConfig{
						Worktrees: map[string]internal.WorktreeConfig{
							"production": {Branch: "production", MergeInto: ""},
						},
					}
				}
				mock.GenerateBranchFromJiraFunc = func(jiraKey string) (string, error) {
					if jiraKey != "PROJECT-123" {
						t.Errorf("Expected JIRA key 'PROJECT-123', got '%s'", jiraKey)
					}
					return "feature/PROJECT-123_fix_auth_bug", nil
				}
				mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
					if worktreeName != "FIX_PROJECT-123" {
						t.Errorf("Expected worktree name 'FIX_PROJECT-123', got '%s'", worktreeName)
					}
					if branchName != "hotfix/PROJECT-123_fix_auth_bug" {
						t.Errorf("Expected branch name 'hotfix/PROJECT-123_fix_auth_bug', got '%s'", branchName)
					}
					if baseBranch != "production" {
						t.Errorf("Expected base branch 'production', got '%s'", baseBranch)
					}
					return nil
				}
			},
			expectError: false,
		},
		{
			name: "hotfix creation with empty prefix",
			args: []string{"auth-fix"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "main", nil
				}
				mock.GetConfigFunc = func() *internal.Config {
					return &internal.Config{
						Settings: internal.ConfigSettings{
							HotfixPrefix: "", // Empty prefix
						},
					}
				}
				mock.GetGBMConfigFunc = func() *internal.GBMConfig {
					return &internal.GBMConfig{}
				}
				mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
					if worktreeName != "auth-fix" {
						t.Errorf("Expected worktree name 'auth-fix', got '%s'", worktreeName)
					}
					return nil
				}
			},
			expectError: false,
		},
		{
			name: "error when production branch detection fails",
			args: []string{"test-fix"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "", errors.New("no production branch found")
				}
			},
			expectError:   true,
			errorContains: "failed to determine production branch",
		},
		{
			name: "error when JIRA branch generation fails",
			args: []string{"PROJECT-456"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "main", nil
				}
				mock.GenerateBranchFromJiraFunc = func(jiraKey string) (string, error) {
					return "", errors.New("JIRA API unavailable")
				}
				mock.GetConfigFunc = func() *internal.Config {
					return &internal.Config{
						Settings: internal.ConfigSettings{
							HotfixPrefix: "HOTFIX",
						},
					}
				}
				mock.GetGBMConfigFunc = func() *internal.GBMConfig {
					return &internal.GBMConfig{}
				}
				mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
					// Should use fallback branch name
					if branchName != "hotfix/PROJECT-456" {
						t.Errorf("Expected fallback branch name 'hotfix/PROJECT-456', got '%s'", branchName)
					}
					return nil
				}
			},
			expectError: false, // Should not error, fallback to simple branch name
		},
		{
			name: "error when worktree creation fails",
			args: []string{"failing-worktree"},
			mockSetup: func(mock *hotfixCreatorMock) {
				mock.FindProductionBranchFunc = func() (string, error) {
					return "main", nil
				}
				mock.GetConfigFunc = func() *internal.Config {
					return &internal.Config{
						Settings: internal.ConfigSettings{
							HotfixPrefix: "HOTFIX",
						},
					}
				}
				mock.GetGBMConfigFunc = func() *internal.GBMConfig {
					return &internal.GBMConfig{}
				}
				mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
					return errors.New("failed to create worktree")
				}
			},
			expectError:   true,
			errorContains: "failed to add hotfix worktree",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mock := &hotfixCreatorMock{}
			tt.mockSetup(mock)

			// Execute
			err := handleHotfix(mock, tt.args)

			// Verify
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHandleHotfixWithJiraInSecondArg(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		jiraBranchResult string
		jiraBranchError  error
		expectedBranch   string
	}{
		{
			name:           "second arg is hotfix branch (from completion)",
			args:           []string{"my-fix", "hotfix/existing-branch"},
			expectedBranch: "hotfix/existing-branch",
		},
		{
			name:             "second arg is JIRA ticket",
			args:             []string{"my-fix", "ABC-789"},
			jiraBranchResult: "feature/ABC-789_implement_new_feature",
			expectedBranch:   "hotfix/ABC-789_implement_new_feature",
		},
		{
			name:            "second arg is JIRA ticket with API error",
			args:            []string{"my-fix", "ABC-999"},
			jiraBranchError: errors.New("JIRA API down"),
			expectedBranch:  "hotfix/ABC-999", // fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &hotfixCreatorMock{}
			mock.FindProductionBranchFunc = func() (string, error) {
				return "main", nil
			}
			mock.GetConfigFunc = func() *internal.Config {
				return &internal.Config{
					Settings: internal.ConfigSettings{
						HotfixPrefix: "HOTFIX",
					},
				}
			}
			mock.GetGBMConfigFunc = func() *internal.GBMConfig {
				return &internal.GBMConfig{}
			}
			if tt.jiraBranchResult != "" || tt.jiraBranchError != nil {
				mock.GenerateBranchFromJiraFunc = func(jiraKey string) (string, error) {
					return tt.jiraBranchResult, tt.jiraBranchError
				}
			}

			var actualBranch string
			mock.AddWorktreeFunc = func(worktreeName, branchName string, createBranch bool, baseBranch string) error {
				actualBranch = branchName
				return nil
			}

			err := handleHotfix(mock, tt.args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if actualBranch != tt.expectedBranch {
				t.Errorf("Expected branch '%s', got '%s'", tt.expectedBranch, actualBranch)
			}
		})
	}
}
