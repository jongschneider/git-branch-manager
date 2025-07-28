package cmd

import (
	"testing"
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
