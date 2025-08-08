package internal

import (
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestManager_ValidateConfig_AllBranchesValid(t *testing.T) {
	repo := testutils.NewMultiBranchRepo(t)

	worktrees := map[string]testutils.WorktreeConfig{
		"main": {Branch: "main", Description: "Main"},
		"dev":  {Branch: "develop", MergeInto: "main", Description: "Dev"},
		"feat": {Branch: "feature/auth", MergeInto: "dev", Description: "Feat"},
		"prod": {Branch: "production/v1.0", MergeInto: "feat", Description: "Prod"},
	}
	require.NoError(t, repo.CreateGBMConfig(worktrees))
	require.NoError(t, repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml"))

	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)
	require.NoError(t, manager.LoadGBMConfig(""))

	// ValidateConfig should succeed
	require.NoError(t, manager.ValidateConfig())

	// GetWorktreeMapping should reflect the YAML
	mapping, err := manager.GetWorktreeMapping()
	require.NoError(t, err)
	require.Equal(t, "main", mapping["main"])
	require.Equal(t, "develop", mapping["dev"])
	require.Equal(t, "feature/auth", mapping["feat"])
	require.Equal(t, "production/v1.0", mapping["prod"])
}

func TestManager_ValidateConfig_MissingBranches(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	// Create config with missing branches
	worktrees := map[string]testutils.WorktreeConfig{
		"main":    {Branch: "main", Description: "Main"},
		"missing": {Branch: "does-not-exist", MergeInto: "main", Description: "Missing"},
	}
	require.NoError(t, repo.CreateGBMConfig(worktrees))
	require.NoError(t, repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml"))

	manager, err := NewManager(repo.GetLocalPath())
	require.NoError(t, err)
	require.NoError(t, manager.LoadGBMConfig(""))

	// ValidateConfig should fail due to missing branch
	err = manager.ValidateConfig()
	require.Error(t, err)
}
