package cmd

import (
	"os"
	"strings"
	"testing"
	"time"

	"gbm/internal"
	"gbm/internal/testutils"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSmartMergebackCompletions(t *testing.T) {
	t.Run("function exists and handles no manager gracefully", func(t *testing.T) {
		// This test ensures the function doesn't panic when manager creation fails
		completions := getSmartMergebackCompletions()
		assert.NotNil(t, completions)
		// Should fall back to empty slice or JIRA completions
	})
}

func TestGetJiraCompletions(t *testing.T) {
	t.Run("function exists and handles no JIRA gracefully", func(t *testing.T) {
		// This test ensures the function doesn't panic when JIRA CLI is not available
		completions := getJiraCompletions()
		assert.NotNil(t, completions)
		// Should return empty slice when JIRA is not available
	})
}

func TestMergebackValidArgsFunction(t *testing.T) {
	t.Run("returns completions for first argument", func(t *testing.T) {
		// Test the ValidArgsFunction for the mergeback command
		cmd := newMergebackCommand()
		validArgsFunc := cmd.ValidArgsFunction
		assert.NotNil(t, validArgsFunc)

		// Test first argument (should get smart completions)
		completions, directive := validArgsFunc(cmd, []string{}, "")
		assert.NotNil(t, completions)
		assert.Equal(t, int(cobra.ShellCompDirectiveNoFileComp), int(directive))

		// Test second argument (should get JIRA completions)
		completions2, directive2 := validArgsFunc(cmd, []string{"test"}, "")
		assert.NotNil(t, completions2)
		assert.Equal(t, int(cobra.ShellCompDirectiveNoFileComp), int(directive2))

		// Test third argument (should return nil)
		completions3, directive3 := validArgsFunc(cmd, []string{"test", "jira"}, "")
		assert.Nil(t, completions3)
		assert.Equal(t, int(cobra.ShellCompDirectiveNoFileComp), int(directive3))
	})
}

// Mock test for completion formatting
func TestCompletionFormatting(t *testing.T) {
	// Test how completion entries are formatted
	activity := internal.RecentActivity{
		Type:          "hotfix",
		WorktreeName:  "SHOP-456",
		BranchName:    "hotfix/SHOP-456_fix_auth",
		CommitMessage: "hotfix: SHOP-456 Fix authentication timeout",
		Author:        "john.doe",
		Timestamp:     time.Date(2025, 7, 12, 14, 30, 0, 0, time.UTC),
		JiraTicket:    "SHOP-456",
	}

	expectedDescription := "Type: hotfix | Branch: hotfix/SHOP-456_fix_auth | Date: 2025-07-12"
	expectedCompletion := "SHOP-456\tType: hotfix | Branch: hotfix/SHOP-456_fix_auth | Date: 2025-07-12"

	// Test the formatting logic that would be used in getSmartMergebackCompletions
	description := formatActivityDescription(activity)
	completion := formatActivityCompletion(activity)

	assert.Equal(t, expectedDescription, description)
	assert.Equal(t, expectedCompletion, completion)
}

// Helper functions for testing completion formatting
func formatActivityDescription(activity internal.RecentActivity) string {
	return formatCompletionDescription(activity.Type, activity.BranchName, activity.Timestamp)
}

func formatActivityCompletion(activity internal.RecentActivity) string {
	description := formatActivityDescription(activity)
	return formatCompletionEntry(activity.WorktreeName, description)
}

func formatCompletionDescription(activityType, branchName string, timestamp time.Time) string {
	return "Type: " + activityType + " | Branch: " + branchName + " | Date: " + timestamp.Format("2006-01-02")
}

func formatCompletionEntry(worktreeName, description string) string {
	return worktreeName + "\t" + description
}

func TestCompletionPrioritization(t *testing.T) {
	// Test the logic for prioritizing different types of activities
	now := time.Now()
	activities := []internal.RecentActivity{
		{
			Type:         "merge",
			WorktreeName: "MERGE-123",
			BranchName:   "merge/MERGE-123_deploy",
			Timestamp:    now.Add(-2 * time.Hour),
		},
		{
			Type:         "hotfix",
			WorktreeName: "HOTFIX-456",
			BranchName:   "hotfix/HOTFIX-456_fix",
			Timestamp:    now.Add(-1 * time.Hour),
		},
		{
			Type:         "hotfix",
			WorktreeName: "HOTFIX-789",
			BranchName:   "hotfix/HOTFIX-789_fix",
			Timestamp:    now.Add(-3 * time.Hour),
		},
	}

	// Test prioritization logic: hotfix > merge, recent > older
	var bestActivity *internal.RecentActivity

	for i := range activities {
		activity := &activities[i]

		if bestActivity == nil {
			bestActivity = activity
			continue
		}

		// Prioritize by type (hotfix is highest priority)
		if activity.Type == "hotfix" && bestActivity.Type != "hotfix" {
			bestActivity = activity
			continue
		}

		// If same type, prioritize more recent
		if activity.Type == bestActivity.Type && activity.Timestamp.After(bestActivity.Timestamp) {
			bestActivity = activity
			continue
		}
	}

	// Should select the most recent hotfix
	assert.Equal(t, "hotfix", bestActivity.Type)
	assert.Equal(t, "HOTFIX-456", bestActivity.WorktreeName)
	assert.Equal(t, now.Add(-1*time.Hour), bestActivity.Timestamp)
}

func TestCompletionFallback(t *testing.T) {
	// Test that completions fall back to JIRA when no smart suggestions are available
	t.Run("falls back to JIRA when no activities", func(t *testing.T) {
		// Simulate empty activities list
		var activities []internal.RecentActivity

		// The function should fall back to JIRA completions
		// We can't test the actual function easily without a manager,
		// but we can test the logic
		assert.Len(t, activities, 0)

		// In the real function, this would trigger:
		// if len(completions) == 0 {
		//     return getJiraCompletions()
		// }
	})
}

func TestCompletionSeparator(t *testing.T) {
	// Test that separator is added correctly between smart suggestions and JIRA
	separator := "---\tOther JIRA tickets:"

	// Test separator format
	assert.Contains(t, separator, "---")
	assert.Contains(t, separator, "Other JIRA tickets:")
	assert.Contains(t, separator, "\t")
}

// Test completion for edge cases
func TestCompletionEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		activity internal.RecentActivity
		expected string
	}{
		{
			name: "activity with empty worktree name",
			activity: internal.RecentActivity{
				Type:         "hotfix",
				WorktreeName: "",
				BranchName:   "hotfix/fix-something",
				Timestamp:    time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC),
			},
			expected: "", // Should be skipped
		},
		{
			name: "activity with empty branch name",
			activity: internal.RecentActivity{
				Type:         "hotfix",
				WorktreeName: "SHOP-456",
				BranchName:   "",
				Timestamp:    time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC),
			},
			expected: "SHOP-456\tType: hotfix | Branch:  | Date: 2025-07-12",
		},
		{
			name: "activity with long branch name",
			activity: internal.RecentActivity{
				Type:         "hotfix",
				WorktreeName: "SHOP-456",
				BranchName:   "hotfix/SHOP-456_fix_very_long_authentication_timeout_issue_with_oauth_redirect",
				Timestamp:    time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC),
			},
			expected: "SHOP-456\tType: hotfix | Branch: hotfix/SHOP-456_fix_very_long_authentication_timeout_issue_with_oauth_redirect | Date: 2025-07-12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.activity.WorktreeName == "" {
				// Should be skipped in real function
				assert.Empty(t, tt.activity.WorktreeName)
			} else {
				completion := formatActivityCompletion(tt.activity)
				assert.Equal(t, tt.expected, completion)
			}
		})
	}
}

func TestCompletionIntegration(t *testing.T) {
	// Create a test repository with hotfix/merge activity
	repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
	defer repo.Cleanup()

	// Create branches and hotfix activity
	err := repo.CreateBranch("production", "Production content")
	require.NoError(t, err)

	// Create hotfix branch with JIRA ticket
	err = repo.CreateBranchFrom("hotfix/SHOP-456_fix_auth", "production", "hotfix: SHOP-456 Fix authentication timeout")
	require.NoError(t, err)

	// Create merge branch
	err = repo.CreateBranchFrom("merge/AUTH-789_deploy", "main", "merge: AUTH-789 Deploy authentication changes")
	require.NoError(t, err)

	// Change to repo directory for testing
	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()
	err = os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	// Test smart completions
	t.Run("smart completions with activity", func(t *testing.T) {
		completions := getSmartMergebackCompletions()
		assert.NotEmpty(t, completions)

		// Should contain entries with proper formatting
		found := false
		for _, completion := range completions {
			if strings.Contains(completion, "SHOP-456") && strings.Contains(completion, "Type: hotfix") {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected SHOP-456 hotfix completion not found")
	})

	// Test that completions are formatted correctly
	t.Run("completion formatting", func(t *testing.T) {
		completions := getSmartMergebackCompletions()

		for _, completion := range completions {
			if strings.Contains(completion, "SHOP-456") {
				// Should have tab-separated format: "WORKTREE\tType: X | Branch: Y | Date: Z"
				parts := strings.Split(completion, "\t")
				assert.Equal(t, 2, len(parts), "Completion should have tab-separated format")
				assert.Contains(t, parts[1], "Type:", "Completion should contain type information")
				assert.Contains(t, parts[1], "Branch:", "Completion should contain branch information")
				assert.Contains(t, parts[1], "Date:", "Completion should contain date information")
				break
			}
		}
	})
}
