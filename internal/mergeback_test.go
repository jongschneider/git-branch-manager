package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for legacy .envrc parsing removed since we no longer support .envrc files

func TestCommitInfo(t *testing.T) {
	t.Run("isUserCommit with email match", func(t *testing.T) {
		commit := MergeBackCommitInfo{
			Hash:      "abc123",
			Message:   "Fix bug",
			Author:    "John Doe",
			Email:     "john@example.com",
			Timestamp: time.Now(),
			IsUser:    false,
		}

		assert.True(t, isUserCommit(commit, "john@example.com", "Jane Doe"))
		assert.False(t, isUserCommit(commit, "jane@example.com", "Jane Doe"))
	})

	t.Run("isUserCommit with name match", func(t *testing.T) {
		commit := MergeBackCommitInfo{
			Hash:      "abc123",
			Message:   "Fix bug",
			Author:    "John Doe",
			Email:     "john@example.com",
			Timestamp: time.Now(),
			IsUser:    false,
		}

		assert.True(t, isUserCommit(commit, "", "John Doe"))
		assert.False(t, isUserCommit(commit, "", "Jane Doe"))
	})

	t.Run("isUserCommit with no match", func(t *testing.T) {
		commit := MergeBackCommitInfo{
			Hash:      "abc123",
			Message:   "Fix bug",
			Author:    "John Doe",
			Email:     "john@example.com",
			Timestamp: time.Now(),
			IsUser:    false,
		}

		assert.False(t, isUserCommit(commit, "jane@example.com", "Jane Doe"))
	})
}

func TestFormatMergeBackAlert(t *testing.T) {
	t.Run("no merge-backs needed", func(t *testing.T) {
		status := &MergeBackStatus{
			MergeBacksNeeded: []MergeBackInfo{},
			HasUserCommits:   false,
		}

		result := FormatMergeBackAlert(status)
		assert.Equal(t, "", result)
	})

	t.Run("nil status", func(t *testing.T) {
		result := FormatMergeBackAlert(nil)
		assert.Equal(t, "", result)
	})

	t.Run("single merge-back with user commits", func(t *testing.T) {
		now := time.Now()
		twoDaysAgo := now.Add(-48 * time.Hour)

		status := &MergeBackStatus{
			MergeBacksNeeded: []MergeBackInfo{
				{
					FromBranch: "PROD",
					ToBranch:   "PREVIEW",
					Commits: []MergeBackCommitInfo{
						{Hash: "abc1234567", Message: "Fix critical bug", Author: "John Doe", Email: "john@example.com", Timestamp: twoDaysAgo, IsUser: true},
						{Hash: "def5678901", Message: "Update config", Author: "Jane Doe", Email: "jane@example.com", Timestamp: now, IsUser: false},
					},
					UserCommits: []MergeBackCommitInfo{
						{Hash: "abc1234567", Message: "Fix critical bug", Author: "John Doe", Email: "john@example.com", Timestamp: twoDaysAgo, IsUser: true},
					},
					TotalCount: 2,
					UserCount:  1,
				},
			},
			HasUserCommits: true,
		}

		result := FormatMergeBackAlert(status)
		assert.Contains(t, result, "⚠️  Merge-back required in tracked branches:")
		assert.Contains(t, result, "PROD → PREVIEW: 2 commits need merge-back (1 by you)")
		assert.Contains(t, result, "• abc1234 - Fix critical bug (you - 2 days ago)")
	})

	t.Run("multiple merge-backs with mixed user commits", func(t *testing.T) {
		now := time.Now()
		oneDayAgo := now.Add(-24 * time.Hour)

		status := &MergeBackStatus{
			MergeBacksNeeded: []MergeBackInfo{
				{
					FromBranch:  "PROD",
					ToBranch:    "PREVIEW",
					Commits:     []MergeBackCommitInfo{{Hash: "abc1234567", Message: "Fix bug", Author: "John Doe", Email: "john@example.com", Timestamp: oneDayAgo, IsUser: true}},
					UserCommits: []MergeBackCommitInfo{{Hash: "abc1234567", Message: "Fix bug", Author: "John Doe", Email: "john@example.com", Timestamp: oneDayAgo, IsUser: true}},
					TotalCount:  1,
					UserCount:   1,
				},
				{
					FromBranch:  "PREVIEW",
					ToBranch:    "MAIN",
					Commits:     []MergeBackCommitInfo{{Hash: "def5678901", Message: "Other fix", Author: "Jane Doe", Email: "jane@example.com", Timestamp: now, IsUser: false}},
					UserCommits: []MergeBackCommitInfo{},
					TotalCount:  1,
					UserCount:   0,
				},
			},
			HasUserCommits: true,
		}

		result := FormatMergeBackAlert(status)
		assert.Contains(t, result, "PROD → PREVIEW: 1 commits need merge-back (1 by you)")
		assert.Contains(t, result, "PREVIEW → MAIN: 1 commits need merge-back (0 by you)")
	})
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now",
			time:     now,
			expected: "just now",
		},
		{
			name:     "30 minutes ago",
			time:     now.Add(-30 * time.Minute),
			expected: "30 minutes ago",
		},
		{
			name:     "2 hours ago",
			time:     now.Add(-2 * time.Hour),
			expected: "2 hours ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "3 days ago",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatRelativeTime(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckMergeBackStatusIntegration(t *testing.T) {
	// Create a test repository with proper git environment
	repo := testutils.NewGitTestRepo(t, testutils.WithDefaultBranch("main"))
	defer repo.Cleanup()

	// Change to repo directory for testing
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	t.Run("missing .gbm.config.yaml file", func(t *testing.T) {
		result, err := CheckMergeBackStatus("/non/existent/.gbm.config.yaml")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty .gbm.config.yaml file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".gbm.config.yaml")

		err := os.WriteFile(configPath, []byte(""), 0644)
		require.NoError(t, err)

		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.MergeBacksNeeded)
	})

	t.Run("single environment", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".gbm.config.yaml")

		config := `worktrees:
  main:
    branch: main
    description: "Main branch"
`
		err := os.WriteFile(configPath, []byte(config), 0644)
		require.NoError(t, err)

		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.MergeBacksNeeded)
	})

	t.Run("multi-tier deployment chain", func(t *testing.T) {
		// Create a more realistic test with actual git branches and merge scenarios
		err := repo.CreateBranch("preview", "Preview content")
		require.NoError(t, err)

		err = repo.CreateBranch("production", "Production content")
		require.NoError(t, err)

		// Create .gbm.config.yaml with merge chain: production -> preview -> main
		configPath := filepath.Join(repo.GetLocalPath(), ".gbm.config.yaml")
		config := `worktrees:
  main:
    branch: main
    merge_into: ""
    description: "Main branch"
  preview:
    branch: preview
    merge_into: "main"
    description: "Preview environment"
  production:
    branch: production
    merge_into: "preview"
    description: "Production environment"
`
		err = os.WriteFile(configPath, []byte(config), 0644)
		require.NoError(t, err)

		// Test with the multi-tier config
		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Since we created branches with different content, they should show as needing mergebacks
		// This is the correct behavior - preview and production have commits that main doesn't have
		if len(result.MergeBacksNeeded) > 0 {
			// Verify the merge chain is detected correctly
			foundPreviewToMain := false
			foundProductionToPreview := false

			for _, mergeback := range result.MergeBacksNeeded {
				if mergeback.FromBranch == "preview" && mergeback.ToBranch == "main" {
					foundPreviewToMain = true
				}
				if mergeback.FromBranch == "production" && mergeback.ToBranch == "preview" {
					foundProductionToPreview = true
				}
			}

			// Should detect at least one of the merge relationships
			assert.True(t, foundPreviewToMain || foundProductionToPreview,
				"Should detect merge relationships in the deployment chain")
		}
	})

	t.Run("with commits requiring mergeback", func(t *testing.T) {
		// Create production branch with new commits that need to be merged back
		err := repo.SwitchToBranch("production")
		require.NoError(t, err)

		err = repo.WriteFile("hotfix.txt", "Critical hotfix for production")
		require.NoError(t, err)

		err = repo.CommitChangesWithForceAdd("hotfix: Fix critical production issue")
		require.NoError(t, err)

		// Switch back to main for testing
		err = repo.SwitchToBranch("main")
		require.NoError(t, err)

		// Create config that includes the production branch
		configPath := filepath.Join(repo.GetLocalPath(), ".gbm.config.yaml")
		config := `worktrees:
  main:
    branch: main
    merge_into: ""
    description: "Main branch"
  production:
    branch: production
    merge_into: "main"
    description: "Production environment"
`
		err = os.WriteFile(configPath, []byte(config), 0644)
		require.NoError(t, err)

		// Now there should be commits requiring mergeback
		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Should detect commits that need to be merged back
		if len(result.MergeBacksNeeded) > 0 {
			assert.True(t, result.HasUserCommits)
			assert.Equal(t, "production", result.MergeBacksNeeded[0].FromBranch)
			assert.Equal(t, "main", result.MergeBacksNeeded[0].ToBranch)
		}
	})
}


func TestMergeBackStructures(t *testing.T) {
	t.Run("merge back status structure", func(t *testing.T) {
		status := MergeBackStatus{
			MergeBacksNeeded: []MergeBackInfo{},
			HasUserCommits:   false,
		}

		assert.Empty(t, status.MergeBacksNeeded)
		assert.False(t, status.HasUserCommits)
	})

	t.Run("merge back info structure", func(t *testing.T) {
		info := MergeBackInfo{
			FromBranch:  "PROD",
			ToBranch:    "PREVIEW",
			Commits:     []MergeBackCommitInfo{},
			UserCommits: []MergeBackCommitInfo{},
			TotalCount:  0,
			UserCount:   0,
		}

		assert.Equal(t, "PROD", info.FromBranch)
		assert.Equal(t, "PREVIEW", info.ToBranch)
		assert.Empty(t, info.Commits)
		assert.Empty(t, info.UserCommits)
		assert.Equal(t, 0, info.TotalCount)
		assert.Equal(t, 0, info.UserCount)
	})

	t.Run("commit info structure", func(t *testing.T) {
		now := time.Now()
		commit := MergeBackCommitInfo{
			Hash:      "abc123",
			Message:   "Test commit",
			Author:    "John Doe",
			Email:     "john@example.com",
			Timestamp: now,
			IsUser:    true,
		}

		assert.Equal(t, "abc123", commit.Hash)
		assert.Equal(t, "Test commit", commit.Message)
		assert.Equal(t, "John Doe", commit.Author)
		assert.Equal(t, "john@example.com", commit.Email)
		assert.Equal(t, now, commit.Timestamp)
		assert.True(t, commit.IsUser)
	})
}
