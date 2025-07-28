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
	LastSync           time.Time         `toml:"last_sync"`
	TrackedVars        []string          `toml:"tracked_vars"`
	AdHocWorktrees     []string          `toml:"ad_hoc_worktrees"`
	CurrentWorktree    string            `toml:"current_worktree"`
	PreviousWorktree   string            `toml:"previous_worktree"`
	LastMergebackCheck time.Time         `toml:"last_mergeback_check"`
	WorktreeBaseBranch map[string]string `toml:"worktree_base_branch"`
}

// DefaultState returns a new State with default values
func DefaultState() *State {
	return &State{
		LastSync:           time.Time{},
		TrackedVars:        []string{},
		AdHocWorktrees:     []string{},
		CurrentWorktree:    "",
		PreviousWorktree:   "",
		LastMergebackCheck: time.Time{},
		WorktreeBaseBranch: make(map[string]string),
	}
}

// LoadState loads the state from .gbm/state.toml
func LoadState(gbmDir string) (*State, error) {
	statePath := filepath.Join(gbmDir, DefaultStateFilename)

	// If state.toml exists, load it directly
	if _, err := os.Stat(statePath); err == nil {
		var state State
		if _, err := toml.DecodeFile(statePath, &state); err != nil {
			return nil, fmt.Errorf("failed to decode state file: %w", err)
		}
		// Initialize WorktreeBaseBranch map if it doesn't exist (for backward compatibility)
		if state.WorktreeBaseBranch == nil {
			state.WorktreeBaseBranch = make(map[string]string)
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

	statePath := filepath.Join(gbmDir, DefaultStateFilename)
	file, err := os.Create(statePath)
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer func() { _ = file.Close() }()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf("failed to encode state: %w", err)
	}

	return nil
}

// SetWorktreeBaseBranch stores the base branch for a worktree
func (s *State) SetWorktreeBaseBranch(worktreeName, baseBranch string) {
	if s.WorktreeBaseBranch == nil {
		s.WorktreeBaseBranch = make(map[string]string)
	}
	s.WorktreeBaseBranch[worktreeName] = baseBranch
}

// GetWorktreeBaseBranch retrieves the base branch for a worktree
func (s *State) GetWorktreeBaseBranch(worktreeName string) (string, bool) {
	if s.WorktreeBaseBranch == nil {
		return "", false
	}
	baseBranch, exists := s.WorktreeBaseBranch[worktreeName]
	return baseBranch, exists
}

// RemoveWorktreeBaseBranch removes base branch information for a worktree
func (s *State) RemoveWorktreeBaseBranch(worktreeName string) {
	if s.WorktreeBaseBranch != nil {
		delete(s.WorktreeBaseBranch, worktreeName)
	}
}
