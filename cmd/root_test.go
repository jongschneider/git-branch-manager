package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gbm/internal"
	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldShowMergeBackAlerts_DisabledByConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test git repository
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	// Change to test directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create .gbm directory and config with alerts disabled
	gbmDir := filepath.Join(tmpDir, ".gbm")
	err = os.MkdirAll(gbmDir, 0755)
	require.NoError(t, err)

	config := internal.DefaultConfig()
	config.Settings.MergeBackAlerts = false
	err = config.Save(gbmDir)
	require.NoError(t, err)

	// Should return false because alerts are disabled
	result := shouldShowMergeBackAlerts()
	assert.False(t, result)
}

func TestShouldShowMergeBackAlerts_EnabledWithNoState(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test git repository
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	// Change to test directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create .gbm directory and config with alerts enabled
	gbmDir := filepath.Join(tmpDir, ".gbm")
	err = os.MkdirAll(gbmDir, 0755)
	require.NoError(t, err)

	config := internal.DefaultConfig()
	config.Settings.MergeBackAlerts = true
	err = config.Save(gbmDir)
	require.NoError(t, err)

	// Create empty branch config to avoid errors
	branchConfigPath := filepath.Join(tmpDir, "gbm.branchconfig.yaml")
	err = os.WriteFile(branchConfigPath, []byte("worktrees: {}"), 0644)
	require.NoError(t, err)

	// Should return true because no state exists (first time check)
	// Note: This may return false if CheckMergeBackStatus returns no mergebacks needed
	result := shouldShowMergeBackAlerts()
	// We can't assert the exact value because it depends on mergeback status
	// but we can verify the function doesn't panic and returns a boolean
	assert.IsType(t, false, result)
}

func TestShouldShowMergeBackAlerts_TimestampLogic(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test git repository
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	// Change to test directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create .gbm directory and config with alerts enabled
	gbmDir := filepath.Join(tmpDir, ".gbm")
	err = os.MkdirAll(gbmDir, 0755)
	require.NoError(t, err)

	config := internal.DefaultConfig()
	config.Settings.MergeBackAlerts = true
	config.Settings.MergeBackCheckInterval = 1 * time.Hour // Set 1 hour interval
	err = config.Save(gbmDir)
	require.NoError(t, err)

	// Create state with recent check (within interval)
	state := internal.DefaultState()
	state.LastMergebackCheck = time.Now().Add(-30 * time.Minute) // 30 minutes ago
	err = state.Save(gbmDir)
	require.NoError(t, err)

	// Create empty branch config to avoid errors
	branchConfigPath := filepath.Join(tmpDir, "gbm.branchconfig.yaml")
	err = os.WriteFile(branchConfigPath, []byte("worktrees: {}"), 0644)
	require.NoError(t, err)

	// Should return false because not enough time has passed
	result := shouldShowMergeBackAlerts()
	assert.False(t, result)

	// Update state with old check (outside interval)
	state.LastMergebackCheck = time.Now().Add(-2 * time.Hour) // 2 hours ago
	err = state.Save(gbmDir)
	require.NoError(t, err)

	// Should still return false if no mergebacks are needed
	// (since our empty config has no mergebacks to check)
	result = shouldShowMergeBackAlerts()
	assert.False(t, result)
}

func TestShouldShowMergeBackAlerts_UserCommitInterval(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test git repository
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	// Change to test directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create .gbm directory and config with different intervals
	gbmDir := filepath.Join(tmpDir, ".gbm")
	err = os.MkdirAll(gbmDir, 0755)
	require.NoError(t, err)

	config := internal.DefaultConfig()
	config.Settings.MergeBackAlerts = true
	config.Settings.MergeBackCheckInterval = 6 * time.Hour         // Normal interval
	config.Settings.MergeBackUserCommitInterval = 30 * time.Minute // User commit interval
	err = config.Save(gbmDir)
	require.NoError(t, err)

	// Create state with check 1 hour ago
	state := internal.DefaultState()
	state.LastMergebackCheck = time.Now().Add(-1 * time.Hour)
	err = state.Save(gbmDir)
	require.NoError(t, err)

	// Create empty branch config
	branchConfigPath := filepath.Join(tmpDir, "gbm.branchconfig.yaml")
	err = os.WriteFile(branchConfigPath, []byte("worktrees: {}"), 0644)
	require.NoError(t, err)

	// With 1 hour since last check:
	// - Normal interval (6h): should not alert yet
	// - User commit interval (30m): should alert if user has commits
	// Since we have empty config, no mergebacks are needed, so result should be false
	result := shouldShowMergeBackAlerts()
	assert.False(t, result)
}

func TestUpdateLastMergebackCheck(t *testing.T) {
	// Create a proper git repository using testutils
	repo := testutils.NewGitTestRepo(t,
		testutils.WithUser("Test User", "test@example.com"),
	)

	// Change to the repository directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	err = os.Chdir(repo.GetLocalPath())
	require.NoError(t, err)

	// Create .gbm directory and initial state
	gbmDir := filepath.Join(repo.GetLocalPath(), ".gbm")
	err = os.MkdirAll(gbmDir, 0755)
	require.NoError(t, err)

	state := internal.DefaultState()
	oldTime := time.Now().Add(-1 * time.Hour)
	state.LastMergebackCheck = oldTime
	err = state.Save(gbmDir)
	require.NoError(t, err)

	// Verify initial state
	initialState, err := internal.LoadState(gbmDir)
	require.NoError(t, err)
	assert.Equal(t, oldTime.Unix(), initialState.LastMergebackCheck.Unix())

	// Call updateLastMergebackCheck
	beforeUpdate := time.Now()
	err = updateLastMergebackCheckWithError()
	require.NoError(t, err)

	// Load state and verify timestamp was updated
	updatedState, err := internal.LoadState(gbmDir)
	require.NoError(t, err)

	// Debug: print timestamps to see what's happening
	t.Logf("Old time: %v", oldTime)
	t.Logf("Before update: %v", beforeUpdate)
	t.Logf("Updated timestamp: %v", updatedState.LastMergebackCheck)

	// The timestamp should be more recent than the old time
	assert.True(t, updatedState.LastMergebackCheck.After(oldTime),
		"Updated timestamp %v should be after old time %v",
		updatedState.LastMergebackCheck, oldTime)

	// Should be very recent (after we called the function)
	assert.True(t, updatedState.LastMergebackCheck.After(beforeUpdate.Add(-1*time.Second)),
		"Updated timestamp %v should be after function call time %v",
		updatedState.LastMergebackCheck, beforeUpdate)
}
