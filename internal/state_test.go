package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultState(t *testing.T) {
	state := DefaultState()

	assert.NotNil(t, state)
	assert.True(t, state.LastSync.IsZero())
	assert.True(t, state.LastMergebackCheck.IsZero())
	assert.Empty(t, state.TrackedVars)
	assert.Empty(t, state.AdHocWorktrees)
	assert.Empty(t, state.CurrentWorktree)
	assert.Empty(t, state.PreviousWorktree)
}

func TestStateLoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()

	// Test loading non-existent state (should return default)
	state, err := LoadState(tmpDir)
	require.NoError(t, err)
	assert.True(t, state.LastSync.IsZero())
	assert.True(t, state.LastMergebackCheck.IsZero())

	// Set some values
	now := time.Now()
	state.LastSync = now
	state.LastMergebackCheck = now.Add(-1 * time.Hour)
	state.TrackedVars = []string{"main", "dev"}
	state.CurrentWorktree = "main"

	// Save state
	err = state.Save(tmpDir)
	require.NoError(t, err)

	// Load state again
	loadedState, err := LoadState(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, state.LastSync.Unix(), loadedState.LastSync.Unix())
	assert.Equal(t, state.LastMergebackCheck.Unix(), loadedState.LastMergebackCheck.Unix())
	assert.Equal(t, state.TrackedVars, loadedState.TrackedVars)
	assert.Equal(t, state.CurrentWorktree, loadedState.CurrentWorktree)
}

func TestStateMergebackCheckTimestamp(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial state
	state := DefaultState()
	assert.True(t, state.LastMergebackCheck.IsZero())

	// Update mergeback check timestamp
	checkTime := time.Now()
	state.LastMergebackCheck = checkTime

	// Save state
	err := state.Save(tmpDir)
	require.NoError(t, err)

	// Load and verify
	loadedState, err := LoadState(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, checkTime.Unix(), loadedState.LastMergebackCheck.Unix())

	// Test time calculations
	timeSinceCheck := time.Since(loadedState.LastMergebackCheck)
	assert.True(t, timeSinceCheck >= 0)
	assert.True(t, timeSinceCheck < time.Second) // Should be very recent
}

func TestStateTomlFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create state with specific timestamp
	state := DefaultState()
	testTime := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)
	state.LastMergebackCheck = testTime
	state.LastSync = testTime.Add(1 * time.Hour)

	// Save state
	err := state.Save(tmpDir)
	require.NoError(t, err)

	// Read the file content directly to verify TOML format
	stateFile := filepath.Join(tmpDir, "state.toml")
	content, err := os.ReadFile(stateFile)
	require.NoError(t, err)

	// Verify TOML format contains our field
	assert.Contains(t, string(content), "last_mergeback_check")
	assert.Contains(t, string(content), "last_sync")
}
