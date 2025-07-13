package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// ConfigState represents the old state structure for migration purposes
type ConfigState struct {
	LastSync         time.Time `toml:"last_sync"`
	TrackedVars      []string  `toml:"tracked_vars"`
	AdHocWorktrees   []string  `toml:"ad_hoc_worktrees"`
	CurrentWorktree  string    `toml:"current_worktree"`
	PreviousWorktree string    `toml:"previous_worktree"`
}

// State represents the runtime state data that is frequently modified
// This will be stored in a separate .gbm/state.toml file
type State struct {
	LastSync         time.Time `toml:"last_sync"`
	TrackedVars      []string  `toml:"tracked_vars"`
	AdHocWorktrees   []string  `toml:"ad_hoc_worktrees"`
	CurrentWorktree  string    `toml:"current_worktree"`
	PreviousWorktree string    `toml:"previous_worktree"`
}

// DefaultState returns a new State with default values
func DefaultState() *State {
	return &State{
		LastSync:         time.Time{},
		TrackedVars:      []string{},
		AdHocWorktrees:   []string{},
		CurrentWorktree:  "",
		PreviousWorktree: "",
	}
}

// LoadState loads the state from .gbm/state.toml
func LoadState(gbmDir string) (*State, error) {
	statePath := filepath.Join(gbmDir, "state.toml")

	// If state.toml exists, load it directly
	if _, err := os.Stat(statePath); err == nil {
		var state State
		if _, err := toml.DecodeFile(statePath, &state); err != nil {
			return nil, fmt.Errorf("failed to decode state file: %w", err)
		}
		return &state, nil
	}

	// Neither file exists, return default state
	return DefaultState(), nil
}

// Save saves the state to .gbm/state.toml
func (s *State) Save(gbmDir string) error {
	if err := os.MkdirAll(gbmDir, 0o755); err != nil {
		return fmt.Errorf("failed to create .gbm directory: %w", err)
	}

	statePath := filepath.Join(gbmDir, "state.toml")
	file, err := os.Create(statePath)
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf("failed to encode state: %w", err)
	}

	return nil
}
