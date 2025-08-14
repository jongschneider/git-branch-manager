package internal

import (
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetWorktreeMapping_Integration(t *testing.T) {
	// Create a repository with branches
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Create branches
	err := repo.CreateBranch("develop", "main")
	require.NoError(t, err)

	// Create basic gbm.branchconfig.yaml
	gbmContent := `worktrees:
  main:
    branch: main
  dev:
    branch: develop
`
	err = repo.WriteFile(DefaultBranchConfigFilename, gbmContent)
	require.NoError(t, err)

	// Create manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)

	err = manager.LoadGBMConfig("")
	require.NoError(t, err)

	// Test GetWorktreeMapping
	mapping, err := manager.GetWorktreeMapping()
	require.NoError(t, err)

	expected := map[string]string{
		"main": "main",
		"dev":  "develop",
	}
	assert.Equal(t, expected, mapping)
}

func TestManager_GetSyncStatus_Integration(t *testing.T) {
	// Create a repository with branches
	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)
	t.Cleanup(func() {
		if repo != nil {
			repo.Cleanup()
		}
	})

	// Create basic gbm.branchconfig.yaml
	gbmContent := `worktrees:
  main:
    branch: main
`
	err := repo.WriteFile(DefaultBranchConfigFilename, gbmContent)
	require.NoError(t, err)

	// Create manager
	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)

	err = manager.LoadGBMConfig("")
	require.NoError(t, err)

	// Test GetSyncStatus - should show missing worktrees since no worktrees exist yet
	status, err := manager.GetSyncStatus()
	require.NoError(t, err)

	assert.False(t, status.InSync)
	assert.Contains(t, status.MissingWorktrees, "main")
	assert.Empty(t, status.OrphanedWorktrees)
	assert.Empty(t, status.BranchChanges)
}
