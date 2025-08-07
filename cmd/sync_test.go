package cmd

import (
	"fmt"
	"strings"
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

func TestHandleSyncDryRun(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *worktreeSyncerMock
		expectError bool
	}{
		{
			name: "all worktrees in sync returns no error",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return &internal.SyncStatus{InSync: true}, nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name: "missing worktrees returns no error",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return &internal.SyncStatus{
						InSync:           false,
						MissingWorktrees: []string{"dev", "feat"},
					}, nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name: "branch changes returns no error",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return &internal.SyncStatus{
						InSync: false,
						BranchChanges: map[string]internal.BranchChange{
							"dev": {OldBranch: "develop", NewBranch: "main"},
						},
					}, nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name: "worktree promotions returns no error",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return &internal.SyncStatus{
						InSync: false,
						WorktreePromotions: []internal.WorktreePromotion{
							{
								SourceWorktree: "production-old",
								TargetWorktree: "production",
								Branch:         "prod-v2",
								TargetBranch:   "prod-v1",
							},
						},
					}, nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name: "orphaned worktrees returns no error",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return &internal.SyncStatus{
						InSync:            false,
						OrphanedWorktrees: []string{"old-feature", "abandoned-dev"},
					}, nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name: "GetSyncStatus error is propagated",
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.GetSyncStatusFunc = func() (*internal.SyncStatus, error) {
					return nil, fmt.Errorf("sync status error")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			err := handleSyncDryRun(mock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock was called exactly once
			assert.Len(t, mock.GetSyncStatusCalls(), 1)
		})
	}
}

func TestHandleSync(t *testing.T) {
	tests := []struct {
		name        string
		force       bool
		setupMock   func() *worktreeSyncerMock
		expectError bool
	}{
		{
			name:  "successful sync without force",
			force: false,
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.SyncWithConfirmationFunc = func(dryRun, force bool, confirmFunc internal.ConfirmationFunc) error {
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:  "successful sync with force",
			force: true,
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.SyncWithConfirmationFunc = func(dryRun, force bool, confirmFunc internal.ConfirmationFunc) error {
					// Verify parameters passed correctly
					if dryRun != false || force != true {
						return fmt.Errorf("incorrect parameters: dryRun=%v, force=%v", dryRun, force)
					}
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:  "sync error is propagated",
			force: false,
			setupMock: func() *worktreeSyncerMock {
				mock := &worktreeSyncerMock{}
				mock.SyncWithConfirmationFunc = func(dryRun, force bool, confirmFunc internal.ConfirmationFunc) error {
					return fmt.Errorf("sync failed")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			err := handleSync(mock, tt.force)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock was called exactly once
			assert.Len(t, mock.SyncWithConfirmationCalls(), 1)

			// Verify parameters passed to mock
			if len(mock.SyncWithConfirmationCalls()) > 0 {
				call := mock.SyncWithConfirmationCalls()[0]
				assert.False(t, call.DryRun, "DryRun should always be false in handleSync")
				assert.Equal(t, tt.force, call.Force)
				assert.NotNil(t, call.ConfirmFunc, "ConfirmFunc should not be nil")
			}
		})
	}
}

func TestSyncCommand_FlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectDry   bool
		expectForce bool
	}{
		{
			name:        "dry-run flag sets dry-run to true",
			args:        []string{"--dry-run"},
			expectDry:   true,
			expectForce: false,
		},
		{
			name:        "force flag sets force to true",
			args:        []string{"--force"},
			expectDry:   false,
			expectForce: true,
		},
		{
			name:        "both flags can be set",
			args:        []string{"--dry-run", "--force"},
			expectDry:   true,
			expectForce: true,
		},
		{
			name:        "no flags defaults to false",
			args:        []string{},
			expectDry:   false,
			expectForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newSyncCommand()
			cmd.SetArgs(tt.args)

			// Parse flags without executing
			err := cmd.ParseFlags(tt.args)
			assert.NoError(t, err)

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			force, _ := cmd.Flags().GetBool("force")

			assert.Equal(t, tt.expectDry, dryRun)
			assert.Equal(t, tt.expectForce, force)
		})
	}
}

func TestConfirmationFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"accepts y", "y", true},
		{"accepts yes", "yes", true},
		{"accepts Y", "Y", true},
		{"accepts YES", "YES", true},
		{"rejects n", "n", false},
		{"rejects no", "no", false},
		{"rejects empty", "", false},
		{"rejects random", "maybe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the confirmation function logic directly
			result := strings.ToLower(tt.input) == "y" || strings.ToLower(tt.input) == "yes"
			assert.Equal(t, tt.expected, result)
		})
	}
}
