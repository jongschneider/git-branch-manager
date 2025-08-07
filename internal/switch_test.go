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

	// Create gbm.branchconfig.yaml to test tracked vs ad hoc worktree sorting
	err := repo.CreateGBMConfig(map[string]testutils.WorktreeConfig{
		"main": {
			Branch:      "main",
			Description: "Main branch",
		},
		"dev": {
			Branch:      "develop",
			MergeInto:   "main",
			Description: "Development branch",
		},
		"preview": {
			Branch:      "preview",
			MergeInto:   "main", 
			Description: "Preview branch",
		},
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml configuration")
	require.NoError(t, err)

	// Create Manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)
	
	// Load GBM config to enable tracked/ad hoc worktree distinction 
	err = manager.LoadGBMConfig("")
	require.NoError(t, err)

	// Create test worktree info map with mix of tracked and ad hoc worktrees
	worktrees := map[string]*WorktreeListInfo{
		"feat": {
			Path:           "/fake/path/to/feat", // ad hoc worktree (not in config)
			CurrentBranch:  "feat",
			ExpectedBranch: "feat",
		},
		"main": {
			Path:           "/fake/path/to/main", // tracked worktree (in config)
			CurrentBranch:  "main", 
			ExpectedBranch: "main",
		},
		"dev": {
			Path:           "/fake/path/to/dev", // tracked worktree (in config)
			CurrentBranch:  "develop",
			ExpectedBranch: "develop",
		},
		"preview": {
			Path:           "/fake/path/to/preview", // tracked worktree (in config)
			CurrentBranch:  "preview",
			ExpectedBranch: "preview",
		},
		"bugfix": {
			Path:           "/fake/path/to/bugfix", // ad hoc worktree (not in config)
			CurrentBranch:  "bugfix",
			ExpectedBranch: "bugfix",
		},
	}

	sortedNames := manager.GetSortedWorktreeNames(worktrees)
	
	// Should return all names in sorted order
	assert.Len(t, sortedNames, 5)
	
	// Tracked worktrees (from gbm.branchconfig.yaml) should come first, sorted alphabetically,
	// followed by ad hoc worktrees sorted alphabetically (since fake paths cause os.Stat to fail)
	expectedOrder := []string{"dev", "main", "preview", "bugfix", "feat"}
	assert.Equal(t, expectedOrder, sortedNames)
	
	// Also verify it returns a consistent order on multiple calls
	sortedNames2 := manager.GetSortedWorktreeNames(worktrees)
	assert.Equal(t, sortedNames, sortedNames2)
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