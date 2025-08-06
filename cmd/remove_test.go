package cmd

import (
	"errors"
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// UNIT TESTS (Using mocks - these are fast and don't require real git operations)
// ============================================================================
// These tests use the worktreeRemover interface with mocks to test business logic
// without requiring real git operations. They run in milliseconds.

func TestHandleRemoveWithConfirmation(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		force        bool
		confirmFunc  confirmationFunc
		mockSetup    func() *worktreeRemoverMock
		assertMocks  func(t *testing.T, mock *worktreeRemoverMock)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name:         "success - remove with force bypasses checks",
			worktreeName: "test-worktree",
			force:        true,
			confirmFunc:  nil, // Not used when force=true
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "test-worktree", worktreeName)
						return "/path/to/worktree", nil
					},
					RemoveWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "test-worktree", worktreeName)
						return nil
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 0) // Should be skipped with force
				assert.Len(t, mock.RemoveWorktreeCalls(), 1)
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "success - remove without uncommitted changes with confirmation",
			worktreeName: "clean-worktree",
			force:        false,
			confirmFunc:  func(worktreeName string) bool { return true }, // User confirms
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "clean-worktree", worktreeName)
						return "/path/to/clean-worktree", nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						assert.Equal(t, "/path/to/clean-worktree", worktreePath)
						return &internal.GitStatus{}, nil // Clean status
					},
					RemoveWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "clean-worktree", worktreeName)
						return nil
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 1)
				assert.Len(t, mock.RemoveWorktreeCalls(), 1)
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "success - user cancels removal",
			worktreeName: "cancel-worktree",
			force:        false,
			confirmFunc:  func(worktreeName string) bool { return false }, // User cancels
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "cancel-worktree", worktreeName)
						return "/path/to/cancel-worktree", nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						assert.Equal(t, "/path/to/cancel-worktree", worktreePath)
						return &internal.GitStatus{}, nil // Clean status
					},
					// RemoveWorktree should not be called when user cancels
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 1)
				assert.Len(t, mock.RemoveWorktreeCalls(), 0) // Should not be called when user cancels
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err) // Cancellation is not an error
			},
		},
		{
			name:         "error - worktree not found",
			worktreeName: "nonexistent",
			force:        false,
			confirmFunc:  nil, // Not used since error occurs before confirmation
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "nonexistent", worktreeName)
						return "", errors.New("worktree not found")
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 0) // Should not be called due to early error
				assert.Len(t, mock.RemoveWorktreeCalls(), 0)    // Should not be called due to error
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "worktree 'nonexistent' not found")
			},
		},
		{
			name:         "error - uncommitted changes without force",
			worktreeName: "dirty-worktree",
			force:        false,
			confirmFunc:  nil, // Not used since error occurs before confirmation
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "dirty-worktree", worktreeName)
						return "/path/to/dirty-worktree", nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						assert.Equal(t, "/path/to/dirty-worktree", worktreePath)
						// Return status with changes
						status := &internal.GitStatus{
							IsDirty:   true,
							Untracked: 1,
						}
						return status, nil
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 1)
				assert.Len(t, mock.RemoveWorktreeCalls(), 0) // Should not be called due to error
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "has uncommitted changes")
				assert.Contains(t, err.Error(), "Use --force to remove anyway")
			},
		},
		{
			name:         "error - status check fails",
			worktreeName: "status-error",
			force:        false,
			confirmFunc:  nil, // Not used since error occurs before confirmation
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "status-error", worktreeName)
						return "/path/to/status-error", nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						assert.Equal(t, "/path/to/status-error", worktreePath)
						return nil, errors.New("failed to get git status")
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 1)
				assert.Len(t, mock.RemoveWorktreeCalls(), 0) // Should not be called due to error
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to check worktree status")
			},
		},
		{
			name:         "error - removal fails",
			worktreeName: "remove-fails",
			force:        true,
			confirmFunc:  nil, // Not used when force=true
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "remove-fails", worktreeName)
						return "/path/to/remove-fails", nil
					},
					RemoveWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "remove-fails", worktreeName)
						return errors.New("git worktree remove failed")
					},
				}
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetWorktreePathCalls(), 1)
				assert.Len(t, mock.GetWorktreeStatusCalls(), 0) // Should be skipped with force
				assert.Len(t, mock.RemoveWorktreeCalls(), 1)    // Should be called even though it fails
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to remove worktree")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handleRemoveWithConfirmation(mock, tt.worktreeName, tt.force, tt.confirmFunc)

			// Assert mock calls
			tt.assertMocks(t, mock)

			// Assert error conditions
			tt.assertErr(t, err)
		})
	}
}

func TestGetWorktreeCompletions(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func() *worktreeRemoverMock
		expectResult func(t *testing.T, completions []string)
		assertMocks  func(t *testing.T, mock *worktreeRemoverMock)
		assertErr    func(t *testing.T, err error)
	}{
		{
			name: "success - return worktree names",
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"main": {},
							"dev":  {},
							"feat": {},
						}, nil
					},
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				assert.ElementsMatch(t, []string{"main", "dev", "feat"}, completions)
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetAllWorktreesCalls(), 1)
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - GetAllWorktrees fails",
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return nil, errors.New("failed to get worktrees")
					},
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				assert.Nil(t, completions)
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetAllWorktreesCalls(), 1)
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get worktrees")
			},
		},
		{
			name: "success - empty worktrees",
			mockSetup: func() *worktreeRemoverMock {
				return &worktreeRemoverMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{}, nil
					},
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				assert.Empty(t, completions)
			},
			assertMocks: func(t *testing.T, mock *worktreeRemoverMock) {
				assert.Len(t, mock.GetAllWorktreesCalls(), 1)
			},
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			completions, err := getWorktreeCompletions(mock)

			// Validate results
			tt.expectResult(t, completions)

			// Assert mock calls
			tt.assertMocks(t, mock)

			// Assert error conditions
			tt.assertErr(t, err)
		})
	}
}
