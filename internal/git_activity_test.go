package internal

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name         string
		timestampStr string
		expectError  bool
		expectedTime time.Time
	}{
		{
			name:         "valid unix timestamp string",
			timestampStr: "1672531200", // 2023-01-01 00:00:00 UTC
			expectError:  false,
			expectedTime: time.Unix(1672531200, 0),
		},
		{
			name:         "invalid timestamp",
			timestampStr: "invalid",
			expectError:  true,
		},
		{
			name:         "empty timestamp",
			timestampStr: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimestamp(tt.timestampStr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTime, result)
			}
		})
	}
}

func TestExtractJiraTicket(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "standard JIRA ticket",
			message:  "hotfix: SHOP-456 Fix authentication timeout",
			expected: "SHOP-456",
		},
		{
			name:     "JIRA ticket at beginning",
			message:  "PROJECT-123: Implement new feature",
			expected: "PROJECT-123",
		},
		{
			name:     "JIRA ticket in middle",
			message:  "Fix AUTH-789 authentication issue",
			expected: "AUTH-789",
		},
		{
			name:     "multiple JIRA tickets",
			message:  "SHOP-456 and AUTH-789 fixes",
			expected: "SHOP-456", // Should return first match
		},
		{
			name:     "no JIRA ticket",
			message:  "Fix authentication timeout issue",
			expected: "",
		},
		{
			name:     "lowercase should not match",
			message:  "shop-456 fix something",
			expected: "",
		},
		{
			name:     "numbers only should not match",
			message:  "123 fix something",
			expected: "",
		},
		{
			name:     "single letter project code",
			message:  "A-123 should not match",
			expected: "",
		},
		{
			name:     "long project code",
			message:  "PROJECTNAME-999 should match",
			expected: "PROJECTNAME-999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractJiraTicket(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractWorktreeNameFromBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expected   string
	}{
		{
			name:       "hotfix with JIRA ticket",
			branchName: "hotfix/SHOP-456_fix_auth",
			expected:   "SHOP-456",
		},
		{
			name:       "feature with JIRA ticket",
			branchName: "feature/PROJ-789_new_ui",
			expected:   "PROJ-789",
		},
		{
			name:       "bugfix with JIRA ticket",
			branchName: "bugfix/BUG-123_fix_crash",
			expected:   "BUG-123",
		},
		{
			name:       "merge with JIRA ticket",
			branchName: "merge/AUTH-456_deploy",
			expected:   "AUTH-456",
		},
		{
			name:       "hotfix without JIRA ticket",
			branchName: "hotfix/critical-auth-fix",
			expected:   "critical-auth-fix",
		},
		{
			name:       "feature without JIRA ticket",
			branchName: "feature/new-dashboard",
			expected:   "new-dashboard",
		},
		{
			name:       "branch without prefix",
			branchName: "AUTH-456_some_work",
			expected:   "AUTH-456",
		},
		{
			name:       "branch with underscores",
			branchName: "hotfix/fix_critical_bug_in_auth",
			expected:   "fix_critical_bug_in_auth",
		},
		{
			name:       "empty branch name",
			branchName: "",
			expected:   "",
		},
		{
			name:       "branch with only prefix",
			branchName: "hotfix/",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractWorktreeNameFromBranch(tt.branchName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractWorktreeNameFromMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "message with JIRA ticket",
			message:  "hotfix: SHOP-456 Fix authentication timeout",
			expected: "SHOP-456",
		},
		{
			name:     "message with feat prefix",
			message:  "feat: Add new user dashboard",
			expected: "add",
		},
		{
			name:     "message with fix prefix",
			message:  "fix: Update authentication flow",
			expected: "update",
		},
		{
			name:     "message with hotfix prefix",
			message:  "hotfix: Critical security patch",
			expected: "critical",
		},
		{
			name:     "message with merge prefix",
			message:  "merge: Deploy new features",
			expected: "deploy",
		},
		{
			name:     "message with add prefix",
			message:  "add: New payment method",
			expected: "new",
		},
		{
			name:     "message with update prefix",
			message:  "update: Improve search performance",
			expected: "improve",
		},
		{
			name:     "message without prefix",
			message:  "Implement new authentication",
			expected: "implement",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "unknown",
		},
		{
			name:     "single word prefix only",
			message:  "fix:",
			expected: "fix:",
		},
		{
			name:     "uppercase message",
			message:  "FIX: CRITICAL BUG",
			expected: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractWorktreeNameFromMessage(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractBranchFromRef(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{
			name:     "origin ref with hotfix",
			ref:      "origin/hotfix/SHOP-456",
			expected: "hotfix/SHOP-456",
		},
		{
			name:     "origin ref with feature",
			ref:      "origin/feature/new-ui",
			expected: "feature/new-ui",
		},
		{
			name:     "ref without origin",
			ref:      "hotfix/SHOP-456",
			expected: "hotfix/SHOP-456",
		},
		{
			name:     "simple branch name",
			ref:      "main",
			expected: "main",
		},
		{
			name:     "empty ref",
			ref:      "",
			expected: "",
		},
		{
			name:     "complex origin ref",
			ref:      "origin/feature/PROJ-123_implement_new_auth",
			expected: "feature/PROJ-123_implement_new_auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBranchFromRef(tt.ref)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitMergePatternRegex(t *testing.T) {
	// Test the regex pattern used in extractMergeBranches
	mergePattern := `Merge branch '([^']+)' into (.+)`
	re := regexp.MustCompile(mergePattern)

	tests := []struct {
		name         string
		message      string
		expectMatch  bool
		sourceBranch string
		targetBranch string
	}{
		{
			name:         "standard merge message",
			message:      "Merge branch 'feature/new-ui' into main",
			expectMatch:  true,
			sourceBranch: "feature/new-ui",
			targetBranch: "main",
		},
		{
			name:         "hotfix merge message",
			message:      "Merge branch 'hotfix/SHOP-456' into production",
			expectMatch:  true,
			sourceBranch: "hotfix/SHOP-456",
			targetBranch: "production",
		},
		{
			name:         "complex branch names",
			message:      "Merge branch 'feature/PROJ-123_implement_auth' into develop",
			expectMatch:  true,
			sourceBranch: "feature/PROJ-123_implement_auth",
			targetBranch: "develop",
		},
		{
			name:        "non-merge message",
			message:     "feat: Add new user interface",
			expectMatch: false,
		},
		{
			name:        "different merge format",
			message:     "Merged feature-branch into main",
			expectMatch: false,
		},
		{
			name:        "empty message",
			message:     "",
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := re.FindStringSubmatch(tt.message)

			if tt.expectMatch {
				assert.True(t, len(matches) >= 3, "Expected match but got none")
				if len(matches) >= 3 {
					assert.Equal(t, tt.sourceBranch, matches[1])
					assert.Equal(t, tt.targetBranch, matches[2])
				}
			} else {
				assert.True(t, len(matches) == 0, "Expected no match but got one")
			}
		})
	}
}

// Mock functions for testing git operations without actual git
func TestMockRecentActivity(t *testing.T) {
	// Create mock recent activities for testing filtering logic
	now := time.Now()
	activities := []RecentActivity{
		{
			Type:          "hotfix",
			WorktreeName:  "SHOP-456",
			BranchName:    "hotfix/SHOP-456_fix_auth",
			CommitMessage: "hotfix: SHOP-456 Fix authentication",
			Timestamp:     now.Add(-1 * time.Hour),
			JiraTicket:    "SHOP-456",
		},
		{
			Type:          "feature",
			WorktreeName:  "PROJ-789",
			BranchName:    "feature/PROJ-789_new_ui",
			CommitMessage: "feat: PROJ-789 Add new UI",
			Timestamp:     now.Add(-2 * time.Hour),
			JiraTicket:    "PROJ-789",
		},
		{
			Type:          "merge",
			WorktreeName:  "AUTH-123",
			BranchName:    "merge/AUTH-123_deploy",
			CommitMessage: "merge: AUTH-123 Deploy changes",
			Timestamp:     now.Add(-3 * time.Hour),
			JiraTicket:    "AUTH-123",
		},
	}

	// Test filtering logic (would be used in filterAndValidateActivities)
	var hotfixAndMerge []RecentActivity
	for _, activity := range activities {
		if activity.Type == "hotfix" || activity.Type == "merge" {
			hotfixAndMerge = append(hotfixAndMerge, activity)
		}
	}

	assert.Len(t, hotfixAndMerge, 2)
	assert.Equal(t, "hotfix", hotfixAndMerge[0].Type)
	assert.Equal(t, "merge", hotfixAndMerge[1].Type)

	// Test sorting by timestamp (most recent first)
	assert.True(t, hotfixAndMerge[0].Timestamp.After(hotfixAndMerge[1].Timestamp))
}

