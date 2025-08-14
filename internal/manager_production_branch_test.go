package internal

import (
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_FindProductionBranch_Simple(t *testing.T) {
	t.Run("finds start of deployment chain with real Git branches", func(t *testing.T) {
		// Create test repository
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		// Create real Git branches for deployment chain
		err := repo.CreateBranch("development", "Development environment content")
		require.NoError(t, err)

		err = repo.CreateBranch("staging", "Staging environment content")
		require.NoError(t, err)

		// main branch already exists as default

		// Create deployment chain config: main -> staging -> production (production is deepest leaf for hotfixes)
		worktrees := map[string]testutils.WorktreeConfig{
			"main": {
				Branch:      "main",
				MergeInto:   "",
				Description: "Default branch (root)",
			},
			"staging": {
				Branch:      "staging",
				MergeInto:   "main",
				Description: "Staging environment",
			},
			"production": {
				Branch:      "development",
				MergeInto:   "staging",
				Description: "Production environment (hotfix base)",
			},
		}

		err = repo.CreateGBMConfig(worktrees)
		require.NoError(t, err)

		// Test the full Manager.FindProductionBranch() integration
		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		result, err := manager.FindProductionBranch()
		require.NoError(t, err)

		// Should find "development" as the deepest leaf (production branch for hotfixes)
		assert.Equal(t, "development", result)

		// Verify the branches actually exist in Git
		branches := []string{"development", "staging", "main"}
		for _, branch := range branches {
			exists, err := manager.BranchExists(branch)
			require.NoError(t, err)
			assert.True(t, exists, "Branch %s should exist", branch)
		}
	})

	t.Run("production branch with development alongside", func(t *testing.T) {
		// Test realistic scenario: main -> staging -> production (hotfix base) + development (feature base)
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		// Create real Git branches
		err := repo.CreateBranch("staging", "Staging environment content")
		require.NoError(t, err)

		err = repo.CreateBranch("production", "Production environment content")
		require.NoError(t, err)

		err = repo.CreateBranch("development", "Development environment content")
		require.NoError(t, err)

		// Create deployment structure: main -> staging -> production (depth 2, hotfix base)
		//                               main -> development (depth 1, feature base)
		worktrees := map[string]testutils.WorktreeConfig{
			"main": {
				Branch:      "main",
				MergeInto:   "",
				Description: "Default branch (root)",
			},
			"staging": {
				Branch:      "staging",
				MergeInto:   "main",
				Description: "Staging environment",
			},
			"production": {
				Branch:      "production",
				MergeInto:   "staging",
				Description: "Production environment (hotfix base)",
			},
			"development": {
				Branch:      "development",
				MergeInto:   "main",
				Description: "Development environment (feature base)",
			},
		}

		err = repo.CreateGBMConfig(worktrees)
		require.NoError(t, err)

		// Test tree-based algorithm finds the deepest leaf for hotfixes
		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		result, err := manager.FindProductionBranch()
		require.NoError(t, err)

		// Should find "production" as the deepest leaf node (depth 2, ideal for hotfixes)
		assert.Equal(t, "production", result)

		// Verify all branches exist in Git
		branches := []string{"main", "staging", "production", "development"}
		for _, branch := range branches {
			exists, err := manager.BranchExists(branch)
			require.NoError(t, err)
			assert.True(t, exists, "Branch %s should exist", branch)
		}
	})

	t.Run("multiple production chains - deepest leaf wins", func(t *testing.T) {
		// Test: main -> mobile-prod (depth 1), main -> staging -> web-prod (depth 2)
		// Hotfix should prefer web-prod (deeper, more tested)
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		// Create real Git branches for multiple production environments
		err := repo.CreateBranch("mobile-production", "Mobile production content")
		require.NoError(t, err)

		err = repo.CreateBranch("staging", "Staging environment content")
		require.NoError(t, err)

		err = repo.CreateBranch("web-production", "Web production content")
		require.NoError(t, err)

		// Create structure with different depth production branches
		worktrees := map[string]testutils.WorktreeConfig{
			"main": {
				Branch:      "main",
				MergeInto:   "",
				Description: "Main branch (root)",
			},
			"mobile-prod": {
				Branch:      "mobile-production",
				MergeInto:   "main",
				Description: "Mobile production (depth 1)",
			},
			"staging": {
				Branch:      "staging",
				MergeInto:   "main",
				Description: "Staging environment",
			},
			"web-prod": {
				Branch:      "web-production",
				MergeInto:   "staging",
				Description: "Web production (depth 2)",
			},
		}

		err = repo.CreateGBMConfig(worktrees)
		require.NoError(t, err)

		// Test tree-based algorithm finds the deepest production branch
		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		result, err := manager.FindProductionBranch()
		require.NoError(t, err)

		// Should find "web-production" as the deepest leaf (depth 2, most tested path)
		assert.Equal(t, "web-production", result)

		// Verify all branches exist in Git
		branches := []string{"main", "mobile-production", "staging", "web-production"}
		for _, branch := range branches {
			exists, err := manager.BranchExists(branch)
			require.NoError(t, err)
			assert.True(t, exists, "Branch %s should exist", branch)
		}
	})

	t.Run("falls back to default branch when no config", func(t *testing.T) {
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		result, err := manager.FindProductionBranch()
		require.NoError(t, err)
		assert.Equal(t, "main", result)
	})

	t.Run("falls back to root branch when no deployment chains", func(t *testing.T) {
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		// Create config with only root branches
		worktrees := map[string]testutils.WorktreeConfig{
			"main": {Branch: "main", MergeInto: ""},
			"dev":  {Branch: "development", MergeInto: ""},
		}

		err := repo.CreateGBMConfig(worktrees)
		require.NoError(t, err)

		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		result, err := manager.FindProductionBranch()
		require.NoError(t, err)

		// Should return one of the root branches (main or development)
		assert.Contains(t, []string{"main", "development"}, result)
	})

	t.Run("returns error when empty config", func(t *testing.T) {
		repo := testutils.NewGitTestRepo(t,
			testutils.WithDefaultBranch("main"),
			testutils.WithUser("Test User", "test@example.com"),
		)
		t.Cleanup(func() { repo.Cleanup() })

		// Create empty config
		err := repo.CreateGBMConfig(map[string]testutils.WorktreeConfig{})
		require.NoError(t, err)

		manager, err := NewManager(repo.GetLocalPath())
		require.NoError(t, err)

		_, err = manager.FindProductionBranch()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no production branch found")
	})
}
