package internal

import (
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeBackDetection_BasicThreeTierScenario(t *testing.T) {
	mainBranch := uuid.NewString()
	previewBranch := uuid.NewString()
	prodBranch := uuid.NewString()

	// Create a repository with synchronized branch hierarchy first
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch(mainBranch),
		testutils.WithUser("Test User", "test@example.com"),
	)

	// Create .envrc on main branch first
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    mainBranch,
		"PREVIEW": previewBranch,
		"PROD":    prodBranch,
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create synchronized branches from main
	err = repo.CreateSynchronizedBranch(previewBranch)
	require.NoError(t, err)

	err = repo.SwitchToBranch(mainBranch)
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch(prodBranch)
	require.NoError(t, err)

	// Now create a hotfix commit on PROD that doesn't exist in PREVIEW
	err = repo.WriteFile("hotfix.txt", "Critical security fix")
	require.NoError(t, err)

	err = repo.CommitChanges("Fix critical security vulnerability")
	require.NoError(t, err)

	// Push the prod changes
	err = repo.PushBranch(prodBranch)
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch(mainBranch)
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect that PROD has commits that need to be merged to PREVIEW
		assert.Len(t, status.MergeBacksNeeded, 1)
		assert.True(t, status.HasUserCommits)

		mergeBack := status.MergeBacksNeeded[0]
		assert.Equal(t, "PROD", mergeBack.FromBranch)
		assert.Equal(t, "PREVIEW", mergeBack.ToBranch)
		assert.Equal(t, 1, mergeBack.TotalCount)
		assert.Equal(t, 1, mergeBack.UserCount)
		assert.Len(t, mergeBack.UserCommits, 1)
		assert.Contains(t, mergeBack.UserCommits[0].Message, "Fix critical security vulnerability")

		alertMsg := FormatMergeBackAlert(status)
		assert.Contains(t, alertMsg, "⚠️  Merge-back required in tracked branches:\n\nPROD → PREVIEW: 1 commits need merge-back (1 by you)")
		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackDetection_MultipleCommits(t *testing.T) {
	mainBranch := uuid.NewString()
	previewBranch := uuid.NewString()
	prodBranch := uuid.NewString()
	// Create repository to test multiple commits from same user
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch(mainBranch),
		testutils.WithUser("Alice", "alice@example.com"),
	)

	// Create .envrc on main first
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    mainBranch,
		"PREVIEW": previewBranch,
		"PROD":    prodBranch,
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create synchronized branches
	err = repo.CreateSynchronizedBranch(previewBranch)
	require.NoError(t, err)

	err = repo.SwitchToBranch(mainBranch)
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch(prodBranch)
	require.NoError(t, err)

	// First commit
	err = repo.WriteFile("fix1.txt", "First bug fix")
	require.NoError(t, err)

	err = repo.CommitChanges("Fix database connection issue")
	require.NoError(t, err)

	// Second commit
	err = repo.WriteFile("fix2.txt", "Second bug fix")
	require.NoError(t, err)

	err = repo.CommitChanges("Fix memory leak in auth module")
	require.NoError(t, err)

	// Third commit
	err = repo.WriteFile("fix3.txt", "Third bug fix")
	require.NoError(t, err)

	err = repo.CommitChanges("Fix race condition in cache")
	require.NoError(t, err)

	// Push the prod changes
	err = repo.PushBranch(prodBranch)
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch(mainBranch)
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect PROD -> PREVIEW merge-back needed
		assert.Len(t, status.MergeBacksNeeded, 1)
		assert.True(t, status.HasUserCommits)

		mergeBack := status.MergeBacksNeeded[0]
		assert.Equal(t, "PROD", mergeBack.FromBranch)
		assert.Equal(t, "PREVIEW", mergeBack.ToBranch)
		assert.Equal(t, 3, mergeBack.TotalCount) // All three commits
		assert.Equal(t, 3, mergeBack.UserCount)  // All commits by Alice
		assert.Len(t, mergeBack.UserCommits, 3)

		// Check commit messages
		messages := make([]string, len(mergeBack.UserCommits))
		for i, commit := range mergeBack.UserCommits {
			messages[i] = commit.Message
		}
		assert.Contains(t, messages, "Fix database connection issue")
		assert.Contains(t, messages, "Fix memory leak in auth module")
		assert.Contains(t, messages, "Fix race condition in cache")

		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackDetection_CascadingMergebacks(t *testing.T) {
	// Test scenario where PROD has commits, and PREVIEW also has different commits
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Developer", "dev@example.com"),
	)

	// Create .envrc on main first
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    "main",
		"PREVIEW": "preview",
		"PROD":    "prod",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create synchronized branches after .envrc is committed
	err = repo.CreateSynchronizedBranch("preview")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("prod")
	require.NoError(t, err)

	// Add commits to PROD
	err = repo.SwitchToBranch("prod")
	require.NoError(t, err)

	err = repo.WriteFile("prod_hotfix.txt", "Production hotfix")
	require.NoError(t, err)

	err = repo.CommitChanges("Fix critical production bug")
	require.NoError(t, err)

	err = repo.PushBranch("prod")
	require.NoError(t, err)

	// Add commits to PREVIEW
	err = repo.SwitchToBranch("preview")
	require.NoError(t, err)

	err = repo.WriteFile("preview_feature.txt", "Preview feature")
	require.NoError(t, err)

	err = repo.CommitChanges("Add new preview feature")
	require.NoError(t, err)

	err = repo.PushBranch("preview")
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect both PROD -> PREVIEW and PREVIEW -> MAIN
		assert.Len(t, status.MergeBacksNeeded, 2)
		assert.True(t, status.HasUserCommits)

		// Find the merge-backs by source branch
		var prodToPreview, previewToMain *MergeBackInfo
		for i := range status.MergeBacksNeeded {
			if status.MergeBacksNeeded[i].FromBranch == "PROD" {
				prodToPreview = &status.MergeBacksNeeded[i]
			} else if status.MergeBacksNeeded[i].FromBranch == "PREVIEW" {
				previewToMain = &status.MergeBacksNeeded[i]
			}
		}

		require.NotNil(t, prodToPreview)
		assert.Equal(t, "PREVIEW", prodToPreview.ToBranch)
		assert.Equal(t, 1, prodToPreview.TotalCount)
		assert.Equal(t, 1, prodToPreview.UserCount)

		require.NotNil(t, previewToMain)
		assert.Equal(t, "MAIN", previewToMain.ToBranch)
		assert.Equal(t, 1, previewToMain.TotalCount)
		assert.Equal(t, 1, previewToMain.UserCount)

		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackDetection_NoMergeBacksNeeded(t *testing.T) {
	// Test scenario where all branches are in sync
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Developer", "dev@example.com"),
	)

	// Create .envrc on main first
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    "main",
		"PREVIEW": "preview",
		"PROD":    "prod",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create synchronized branches with no additional commits
	err = repo.CreateSynchronizedBranch("preview")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("prod")
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect no merge-backs needed
		assert.Len(t, status.MergeBacksNeeded, 0)
		assert.False(t, status.HasUserCommits)

		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackDetection_NonExistentBranches(t *testing.T) {
	// Test scenario with .envrc referencing branches that don't exist
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Developer", "dev@example.com"),
	)

	// Create .envrc with non-existent branches
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    "main",
		"PREVIEW": "nonexistent-preview",
		"PROD":    "nonexistent-prod",
		"STAGING": "also-nonexistent",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect no merge-backs due to missing branches
		assert.Len(t, status.MergeBacksNeeded, 0)
		assert.False(t, status.HasUserCommits)

		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackDetection_DynamicHierarchy(t *testing.T) {
	// Test with a more complex hierarchy
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("DevOps", "devops@example.com"),
	)

	// Create .envrc first on main
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    "main",
		"STAGING": "staging",
		"PREVIEW": "preview",
		"PROD":    "prod",
		"HOTFIX":  "hotfix",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create a five-tier hierarchy with synchronized branches
	err = repo.CreateSynchronizedBranch("staging")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("preview")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("prod")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("hotfix")
	require.NoError(t, err)

	// Add a commit to HOTFIX (bottom of hierarchy)
	err = repo.SwitchToBranch("hotfix")
	require.NoError(t, err)

	err = repo.WriteFile("emergency.txt", "Emergency patch")
	require.NoError(t, err)

	err = repo.CommitChanges("Emergency security patch")
	require.NoError(t, err)

	err = repo.PushBranch("hotfix")
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	// Test merge-back detection
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Should detect only HOTFIX -> PROD merge-back (immediate parent)
		assert.Len(t, status.MergeBacksNeeded, 1)
		assert.True(t, status.HasUserCommits)

		mergeBack := status.MergeBacksNeeded[0]
		assert.Equal(t, "HOTFIX", mergeBack.FromBranch)
		assert.Equal(t, "PROD", mergeBack.ToBranch)
		assert.Equal(t, 1, mergeBack.TotalCount)
		assert.Equal(t, 1, mergeBack.UserCount)

		return nil
	})
	require.NoError(t, err)
}

func TestMergeBackAlertFormatting_RealScenario(t *testing.T) {
	// Test the alert formatting with real commit data
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Engineer", "engineer@company.com"),
	)

	// Create .envrc on main first
	err := repo.CreateEnvrc(map[string]string{
		"MAIN":    "main",
		"PREVIEW": "preview",
		"PROD":    "prod",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add .envrc configuration")
	require.NoError(t, err)

	// Create synchronized branches
	err = repo.CreateSynchronizedBranch("preview")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("prod")
	require.NoError(t, err)

	err = repo.WriteFile("security_patch.txt", "Security vulnerability fix")
	require.NoError(t, err)

	err = repo.CommitChanges("CVE-2024-1234: Fix SQL injection in user auth")
	require.NoError(t, err)

	// Push the prod changes
	err = repo.PushBranch("prod")
	require.NoError(t, err)

	// Switch back to main for testing
	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	// Test merge-back detection and alert formatting
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), ".envrc")
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		alert := FormatMergeBackAlert(status)
		assert.NotEmpty(t, alert)
		assert.Contains(t, alert, "⚠️  Merge-back required in tracked branches:")
		assert.Contains(t, alert, "PROD → PREVIEW: 1 commits need merge-back (1 by you)")
		assert.Contains(t, alert, "CVE-2024-1234: Fix SQL injection in user auth")
		assert.Contains(t, alert, "(you -")

		return nil
	})
	require.NoError(t, err)
}
