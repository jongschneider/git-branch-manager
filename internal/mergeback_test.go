package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnvrcFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    []EnvVarMapping
		expectError bool
	}{
		{
			name: "standard three-tier configuration",
			content: `MAIN=main
PREVIEW=preview-v2.1
PROD=production-v2.0`,
			expected: []EnvVarMapping{
				{Name: "MAIN", Branch: "main", Order: 0},
				{Name: "PREVIEW", Branch: "preview-v2.1", Order: 1},
				{Name: "PROD", Branch: "production-v2.0", Order: 2},
			},
			expectError: false,
		},
		{
			name: "configuration with comments and empty lines",
			content: `# Environment configuration
MAIN=main

# Preview environment
PREVIEW=preview-v2.1

# Production environment
PROD=production-v2.0`,
			expected: []EnvVarMapping{
				{Name: "MAIN", Branch: "main", Order: 0},
				{Name: "PREVIEW", Branch: "preview-v2.1", Order: 1},
				{Name: "PROD", Branch: "production-v2.0", Order: 2},
			},
			expectError: false,
		},
		{
			name:    "single environment",
			content: `MAIN=main`,
			expected: []EnvVarMapping{
				{Name: "MAIN", Branch: "main", Order: 0},
			},
			expectError: false,
		},
		{
			name:        "empty file",
			content:     "",
			expected:    []EnvVarMapping{},
			expectError: false,
		},
		{
			name: "complex environment names",
			content: `MAIN=main
STAGING_2=staging-branch-v2
PREVIEW_PROD=preview-production-v3.0
PROD_FINAL=production-final-v1.0`,
			expected: []EnvVarMapping{
				{Name: "MAIN", Branch: "main", Order: 0},
				{Name: "STAGING_2", Branch: "staging-branch-v2", Order: 1},
				{Name: "PREVIEW_PROD", Branch: "preview-production-v3.0", Order: 2},
				{Name: "PROD_FINAL", Branch: "production-final-v1.0", Order: 3},
			},
			expectError: false,
		},
		{
			name: "mixed valid and invalid lines",
			content: `MAIN=main
invalid-line-without-equals
PREVIEW=preview-v2.1
# comment
PROD=production-v2.0
123INVALID=branch`,
			expected: []EnvVarMapping{
				{Name: "MAIN", Branch: "main", Order: 0},
				{Name: "PREVIEW", Branch: "preview-v2.1", Order: 1},
				{Name: "PROD", Branch: "production-v2.0", Order: 2},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".envrc")

			err := os.WriteFile(configPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			result, err := parseEnvrcFile(configPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseEnvrcFileErrors(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		_, err := parseEnvrcFile("/non/existent/file")
		assert.Error(t, err)
	})
}

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
	cwd, _ := os.Getwd()
	_, err := FindGitRoot(cwd)
	if err != nil {
		t.Skip("Not in a git repository, skipping integration test")
	}

	t.Run("missing .envrc file", func(t *testing.T) {
		result, err := CheckMergeBackStatus("/non/existent/.envrc")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty .envrc file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".envrc")

		err := os.WriteFile(configPath, []byte(""), 0644)
		require.NoError(t, err)

		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.MergeBacksNeeded)
	})

	t.Run("single environment", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".envrc")

		err := os.WriteFile(configPath, []byte("MAIN=main"), 0644)
		require.NoError(t, err)

		result, err := CheckMergeBackStatus(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.MergeBacksNeeded)
	})
}

func TestEnvVarMapping(t *testing.T) {
	t.Run("env var mapping structure", func(t *testing.T) {
		mapping := EnvVarMapping{
			Name:   "PROD",
			Branch: "production-v1.0",
			Order:  0,
		}

		assert.Equal(t, "PROD", mapping.Name)
		assert.Equal(t, "production-v1.0", mapping.Branch)
		assert.Equal(t, 0, mapping.Order)
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
