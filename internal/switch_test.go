package internal

import (
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_SetCurrentWorktree_GetPreviousWorktree(t *testing.T) {
	// Setup repository
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)

	// Initially no previous worktree
	previous := manager.GetPreviousWorktree()
	assert.Empty(t, previous)

	// Set current worktree to dev
	err = manager.SetCurrentWorktree("dev")
	assert.NoError(t, err)
	
	// Still no previous since this is the first
	previous = manager.GetPreviousWorktree()
	assert.Empty(t, previous)

	// Set current worktree to feat
	err = manager.SetCurrentWorktree("feat")
	assert.NoError(t, err)

	// Now previous should be dev
	previous = manager.GetPreviousWorktree()
	assert.Equal(t, "dev", previous)

	// Set current worktree to main
	err = manager.SetCurrentWorktree("main")
	assert.NoError(t, err)

	// Now previous should be feat
	previous = manager.GetPreviousWorktree()
	assert.Equal(t, "feat", previous)
}

func TestManager_GetSortedWorktreeNames(t *testing.T) {
	// Setup repository
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)

	// Create test worktree info map
	worktrees := map[string]*WorktreeListInfo{
		"feat": {
			Path:           "/path/to/feat",
			CurrentBranch:  "feat",
			ExpectedBranch: "feat",
		},
		"main": {
			Path:           "/path/to/main",
			CurrentBranch:  "main", 
			ExpectedBranch: "main",
		},
		"dev": {
			Path:           "/path/to/dev",
			CurrentBranch:  "dev",
			ExpectedBranch: "dev",
		},
	}

	sortedNames := manager.GetSortedWorktreeNames(worktrees)
	
	// Should return all names
	assert.Len(t, sortedNames, 3)
	assert.Contains(t, sortedNames, "main")
	assert.Contains(t, sortedNames, "dev") 
	assert.Contains(t, sortedNames, "feat")
}

func TestManager_GetStatusIcon(t *testing.T) {
	// Setup repository
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)

	tests := []struct {
		name       string
		gitStatus  *GitStatus
		expectIcon string
	}{
		{
			name:       "nil status",
			gitStatus:  nil,
			expectIcon: "", // Should handle nil gracefully
		},
		{
			name: "clean status",
			gitStatus: &GitStatus{
				IsDirty:   false,
				Untracked: 0,
				Modified:  0,
				Staged:    0,
			},
			expectIcon: "✅", // Clean should return checkmark
		},
		{
			name: "modified files",
			gitStatus: &GitStatus{
				IsDirty:   true,
				Untracked: 0,
				Modified:  1,
				Staged:    0,
			},
			expectIcon: "✚", // Changes should return plus sign
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := manager.GetStatusIcon(tt.gitStatus)
			if tt.expectIcon != "" {
				assert.Equal(t, tt.expectIcon, icon)
			}
			// For empty expectIcon, just verify it doesn't panic
		})
	}
}