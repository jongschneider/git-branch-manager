package cmd

import (
	"errors"
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

func TestHandlePushAll(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *worktreePusherMock
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "success - push all worktrees",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					PushAllWorktreesFunc: func() error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - push all worktrees fails",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					PushAllWorktreesFunc: func() error {
						return errors.New("push failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "push failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()

			err := handlePushAll(mock)
			tt.expectErr(t, err)

			// Verify the mock was called
			assert.Equal(t, 1, len(mock.PushAllWorktreesCalls()))
		})
	}
}

func TestHandlePushCurrent(t *testing.T) {
	tests := []struct {
		name        string
		currentPath string
		setupMock   func() *worktreePusherMock
		expectErr   func(t *testing.T, err error)
	}{
		{
			name:        "success - push current worktree",
			currentPath: "/some/path",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return true, "dev", nil
					},
					PushWorktreeFunc: func(worktreeName string) error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:        "error - not in worktree",
			currentPath: "/some/path",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return false, "", nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not currently in a worktree")
			},
		},
		{
			name:        "error - IsInWorktree fails",
			currentPath: "/some/path",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return false, "", errors.New("failed to check worktree")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to check if in worktree")
			},
		},
		{
			name:        "error - push worktree fails",
			currentPath: "/some/path",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return true, "dev", nil
					},
					PushWorktreeFunc: func(worktreeName string) error {
						return errors.New("push failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "push failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()

			err := handlePushCurrent(mock, tt.currentPath)
			tt.expectErr(t, err)
		})
	}
}

func TestHandlePushNamed(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		setupMock    func() *worktreePusherMock
		expectErr    func(t *testing.T, err error)
	}{
		{
			name:         "success - push named worktree",
			worktreeName: "dev",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {Path: "/path/to/dev"},
						}, nil
					},
					PushWorktreeFunc: func(worktreeName string) error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "error - worktree does not exist",
			worktreeName: "nonexistent",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {Path: "/path/to/dev"},
						}, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "worktree 'nonexistent' does not exist")
			},
		},
		{
			name:         "error - GetAllWorktrees fails",
			worktreeName: "dev",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return nil, errors.New("failed to get worktrees")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get worktrees")
			},
		},
		{
			name:         "error - push worktree fails",
			worktreeName: "dev",
			setupMock: func() *worktreePusherMock {
				return &worktreePusherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {Path: "/path/to/dev"},
						}, nil
					},
					PushWorktreeFunc: func(worktreeName string) error {
						return errors.New("push failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "push failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()

			err := handlePushNamed(mock, tt.worktreeName)
			tt.expectErr(t, err)
		})
	}
}
