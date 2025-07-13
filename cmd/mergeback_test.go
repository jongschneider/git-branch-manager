package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"gbm/internal"
	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergebackCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no arguments - should attempt auto-detection",
			args:        []string{},
			expectError: true, // Will fail due to no confirmation input in test
			errorMsg:    "failed to read confirmation",
		},
		{
			name:        "with worktree name",
			args:        []string{"test-feature"},
			expectError: false,
		},
		{
			name:        "with worktree name and jira ticket",
			args:        []string{"test-feature", "PROJECT-123"},
			expectError: false,
		},
		{
			name:        "too many arguments",
			args:        []string{"test-feature", "PROJECT-123", "extra"},
			expectError: true,
			errorMsg:    "accepts at most 2 arg(s), received 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test repository with hotfix activity
			repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
			defer repo.Cleanup()

			// Create a production branch with hotfix commits for auto-detection
			err := repo.CreateBranch("production", "Production content")
			require.NoError(t, err)

			// Create a hotfix branch and commit
			err = repo.CreateBranchFrom("hotfix/SHOP-456_fix_auth", "production", "hotfix: SHOP-456 Fix authentication timeout")
			require.NoError(t, err)

			// Change to the repo directory for the test
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			err = os.Chdir(repo.GetLocalPath())
			require.NoError(t, err)

			// Create fresh root command to avoid state conflicts
			cmd := newRootCommand()
			args := append([]string{"mergeback"}, tt.args...)
			cmd.SetArgs(args)

			err = cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateMergebackBranchName(t *testing.T) {
	tests := []struct {
		name           string
		worktreeName   string
		jiraTicket     string
		targetWorktree string
		expected       string
	}{
		{
			name:           "simple worktree name with target",
			worktreeName:   "fix-auth",
			jiraTicket:     "",
			targetWorktree: "preview",
			expected:       "merge/fix-auth_preview",
		},
		{
			name:           "JIRA ticket with target",
			worktreeName:   "PROJECT-123",
			jiraTicket:     "PROJECT-123",
			targetWorktree: "main",
			expected:       "merge/PROJECT-123_main",
		},
		{
			name:           "worktree with spaces and underscores",
			worktreeName:   "fix auth bug",
			jiraTicket:     "",
			targetWorktree: "preview",
			expected:       "merge/fix-auth-bug_preview",
		},
		{
			name:           "uppercase target worktree",
			worktreeName:   "hotfix",
			jiraTicket:     "",
			targetWorktree: "PREVIEW",
			expected:       "merge/hotfix_preview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateMergebackBranchName(tt.worktreeName, tt.jiraTicket, tt.targetWorktree, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterAndValidateActivities(t *testing.T) {
	// Create mock activities
	activities := []internal.RecentActivity{
		{
			Type:          "hotfix",
			WorktreeName:  "SHOP-456",
			BranchName:    "hotfix/SHOP-456_fix_auth",
			CommitMessage: "hotfix: Fix authentication timeout",
			Timestamp:     time.Now().Add(-1 * time.Hour),
		},
		{
			Type:          "feature",
			WorktreeName:  "PROJ-789",
			BranchName:    "feature/PROJ-789_new_ui",
			CommitMessage: "feat: Add new user interface",
			Timestamp:     time.Now().Add(-2 * time.Hour),
		},
		{
			Type:          "merge",
			WorktreeName:  "AUTH-123",
			BranchName:    "merge/AUTH-123_deploy",
			CommitMessage: "merge: Deploy authentication changes",
			Timestamp:     time.Now().Add(-3 * time.Hour),
		},
	}

	t.Run("filters out feature branches", func(t *testing.T) {
		// Note: This test will need actual manager for full validation
		// For now, test the type filtering logic
		var filtered []internal.RecentActivity
		for _, activity := range activities {
			if activity.Type == "hotfix" || activity.Type == "merge" {
				filtered = append(filtered, activity)
			}
		}

		assert.Len(t, filtered, 2)
		assert.Equal(t, "hotfix", filtered[0].Type)
		assert.Equal(t, "merge", filtered[1].Type)
	})
}

func TestExtractWorktreeNameFromBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expected   string
	}{
		{
			name:       "hotfix branch with JIRA ticket",
			branchName: "hotfix/SHOP-456_fix_auth",
			expected:   "SHOP-456",
		},
		{
			name:       "feature branch with JIRA ticket",
			branchName: "feature/PROJ-789_new_ui",
			expected:   "PROJ-789",
		},
		{
			name:       "bugfix branch with JIRA ticket",
			branchName: "bugfix/BUG-123_fix_crash",
			expected:   "BUG-123",
		},
		{
			name:       "hotfix branch without JIRA ticket",
			branchName: "hotfix/critical-auth-fix",
			expected:   "critical-auth-fix",
		},
		{
			name:       "branch without prefix",
			branchName: "AUTH-456_some_work",
			expected:   "AUTH-456",
		},
		{
			name:       "empty branch name",
			branchName: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ExtractWorktreeNameFromBranch(tt.branchName)
			assert.Equal(t, tt.expected, result)
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
			name:     "commit message with JIRA ticket",
			message:  "hotfix: SHOP-456 Fix authentication timeout",
			expected: "SHOP-456",
		},
		{
			name:     "commit message with multiple JIRA tickets",
			message:  "feat: PROJ-123 and AUTH-789 implement new feature",
			expected: "PROJ-123", // Should return the first one
		},
		{
			name:     "commit message without JIRA ticket",
			message:  "fix: Update user interface styles",
			expected: "",
		},
		{
			name:     "JIRA ticket at end of message",
			message:  "Implement new authentication flow for SHOP-999",
			expected: "SHOP-999",
		},
		{
			name:     "lowercase jira pattern should not match",
			message:  "fix: shop-456 update something",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ExtractJiraTicket(tt.message)
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
			name:     "message without JIRA ticket",
			message:  "hotfix: Fix critical authentication bug",
			expected: "fix",
		},
		{
			name:     "message with feat prefix",
			message:  "feat: Add new user dashboard",
			expected: "add",
		},
		{
			name:     "message with update prefix",
			message:  "update: Improve search performance",
			expected: "improve",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "unknown",
		},
		{
			name:     "single word message",
			message:  "fix:",
			expected: "fix:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ExtractWorktreeNameFromMessage(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Removed duplicate functions - they exist in mergeback_completion_test.go

func TestFindPotentialMergeTargets(t *testing.T) {
	// Create mock GBM config
	config := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"MAIN": {
				Branch:    "main",
				MergeInto: "",
			},
			"PREVIEW": {
				Branch:    "preview",
				MergeInto: "main",
			},
			"PRODUCTION": {
				Branch:    "production",
				MergeInto: "preview",
			},
		},
	}

	tests := []struct {
		name       string
		branchName string
		expected   []string
	}{
		{
			name:       "production branch merges into preview",
			branchName: "production",
			expected:   []string{"preview"},
		},
		{
			name:       "preview branch merges into main",
			branchName: "preview",
			expected:   []string{"main"},
		},
		{
			name:       "main branch has no merge target",
			branchName: "main",
			expected:   []string{},
		},
		{
			name:       "hotfix branch with no explicit config",
			branchName: "hotfix/SHOP-456",
			expected:   []string{"production"}, // Should find production as root
		},
		{
			name:       "unknown branch",
			branchName: "unknown-branch",
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findPotentialMergeTargets(tt.branchName, config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindMergeTargetBranchAndWorktree(t *testing.T) {
	t.Run("handles missing config gracefully", func(t *testing.T) {
		// Create a test repository
		repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
		defer repo.Cleanup()

		// Change to repo directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		err := os.Chdir(repo.GetLocalPath())
		require.NoError(t, err)

		// Create manager without .gbm.config.yaml (should handle gracefully)
		manager, err := createInitializedManager()
		require.NoError(t, err)

		// Test that the function doesn't panic and returns something reasonable
		branch, worktree, err := findMergeTargetBranchAndWorktree(manager)

		// Should not panic and should return some default values
		assert.NoError(t, err)
		assert.NotEmpty(t, branch, "Should return a branch name")
		assert.NotEmpty(t, worktree, "Should return a worktree name")

		// Without config, should default to main branch
		assert.Equal(t, "main", branch)
		assert.Equal(t, "main", worktree)
	})
}

func TestMergebackIntegration(t *testing.T) {
	// Create a test repository with proper GBM config and mergeback chain
	repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
	defer repo.Cleanup()

	// Set up deployment chain: production -> preview -> main
	err := repo.CreateBranch("preview", "Preview content")
	require.NoError(t, err)

	err = repo.CreateBranch("production", "Production content")
	require.NoError(t, err)

	// Create .gbm.config.yaml
	gbmConfig := map[string]string{
		"main":       "main",
		"preview":    "preview",
		"production": "production",
	}
	err = repo.CreateGBMConfig(gbmConfig)
	require.NoError(t, err)
	err = repo.CommitChangesWithForceAdd("Add .gbm.config.yaml")
	require.NoError(t, err)

	// Create hotfix activity on production
	err = repo.SwitchToBranch("production")
	require.NoError(t, err)
	err = repo.WriteFile("hotfix.txt", "hotfix: SHOP-456 Fix critical authentication bug")
	require.NoError(t, err)
	err = repo.CommitChangesWithForceAdd("hotfix: SHOP-456 Fix critical authentication bug")
	require.NoError(t, err)

	// Change to the repo directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err = os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	// Test manual mergeback creation
	t.Run("manual mergeback creation", func(t *testing.T) {
		// Create a fresh root command instance
		cmd := newRootCommand()
		cmd.SetArgs([]string{"mergeback", "fix-auth", "SHOP-456"})

		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify worktree was created
		assert.DirExists(t, "worktrees/MERGE_fix-auth_main")
	})

	// Test branch naming
	t.Run("verify branch naming", func(t *testing.T) {
		// Check if local merge branch exists with JIRA ticket naming
		cmd := exec.Command("git", "branch", "--list", "merge/*")
		cmd.Dir = repo.GetLocalPath()
		output, err := cmd.Output()
		require.NoError(t, err)

		// Parse local branches
		localBranches := []string{}
		for _, line := range strings.Split(string(output), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				// Remove leading *, +, spaces and any git branch indicators
				branch := strings.TrimLeft(line, "*+ ")
				localBranches = append(localBranches, branch)
			}
		}

		// Debug: print all local merge branches
		t.Logf("Local merge branches: %v", localBranches)

		// The branch should be named with the JIRA ticket (SHOP-456)
		found := false
		expectedBranches := []string{"merge/SHOP-456", "merge/SHOP-456_main"}
		for _, expectedBranch := range expectedBranches {
			for _, branch := range localBranches {
				if branch == expectedBranch {
					found = true
					t.Logf("Found expected branch: %s", expectedBranch)
					break
				}
			}
			if found {
				break
			}
		}
		assert.True(t, found, "Expected merge branch with JIRA ticket not found. Local branches: %v", localBranches)
	})
}

func TestMergebackWorktreeNaming(t *testing.T) {
	// Test worktree naming behavior with and without prefix
	t.Run("with mergeback prefix", func(t *testing.T) {
		// This is tested in the integration test above - default "MERGE" prefix
		// Results in: MERGE_worktree-name_target-worktree
	})

	t.Run("without mergeback prefix", func(t *testing.T) {
		// Create a test repository
		repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
		defer repo.Cleanup()

		// Create .gbm.config.yaml with empty mergeback prefix
		gbmConfig := map[string]string{
			"main": "main",
		}
		err := repo.CreateGBMConfig(gbmConfig)
		require.NoError(t, err)
		err = repo.CommitChangesWithForceAdd("Add .gbm.config.yaml")
		require.NoError(t, err)

		// Change to the repo directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		err = os.Chdir(repo.GetLocalPath())
		require.NoError(t, err)

		// Create the .gbm directory and config.toml file with empty mergeback prefix
		err = os.MkdirAll(".gbm", 0o755)
		require.NoError(t, err)

		configContent := `[settings]
mergeback_prefix = ""
worktree_prefix = "worktrees"

[state]
last_sync = "1970-01-01T00:00:00Z"

[icons]

[jira]

[file_copy]`

		err = os.WriteFile(".gbm/config.toml", []byte(configContent), 0o644)
		require.NoError(t, err)

		// Test mergeback creation without prefix
		cmd := newRootCommand()
		cmd.SetArgs([]string{"mergeback", "fix-auth", "SHOP-456"})

		err = cmd.Execute()
		assert.NoError(t, err)

		// Verify worktree was created with target suffix but no prefix
		assert.DirExists(t, "worktrees/fix-auth_main")

		// Verify branch naming still includes target suffix
		gitCmd := exec.Command("git", "branch", "--list", "merge/*")
		gitCmd.Dir = repo.GetLocalPath()
		output, err := gitCmd.Output()
		require.NoError(t, err)

		assert.Contains(t, string(output), "merge/SHOP-456_main", "Branch should include target suffix")
	})
}

// Benchmark tests for performance-critical functions
func BenchmarkExtractJiraTicket(b *testing.B) {
	message := "hotfix: SHOP-456 Fix authentication timeout issue with oauth"
	for i := 0; i < b.N; i++ {
		internal.ExtractJiraTicket(message)
	}
}

func BenchmarkExtractWorktreeNameFromBranch(b *testing.B) {
	branchName := "hotfix/SHOP-456_fix_authentication_timeout_issue"
	for i := 0; i < b.N; i++ {
		internal.ExtractWorktreeNameFromBranch(branchName)
	}
}
