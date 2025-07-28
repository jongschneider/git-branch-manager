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
// These tests use the worktreePuller interface with mocks to test business logic
// without requiring real git operations. They run in milliseconds.

func TestHandlePullAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *worktreePullerMock
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "success - pull all worktrees",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					PullAllWorktreesFunc: func() error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - pull all fails with git error",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					PullAllWorktreesFunc: func() error {
						return errors.New("git pull failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "git pull failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handlePullAll(mock)
			tt.expectErr(t, err)
		})
	}
}

func TestHandlePullCurrent(t *testing.T) {
	tests := []struct {
		name        string
		currentPath string
		mockSetup   func() *worktreePullerMock
		expectErr   func(t *testing.T, err error)
	}{
		{
			name:        "success - pull current worktree",
			currentPath: "/test/worktrees/dev",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return true, "dev", nil
					},
					PullWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "dev", worktreeName)
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
			currentPath: "/test/repo",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
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
			currentPath: "/test/invalid",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return false, "", errors.New("failed to determine worktree status")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to check if in worktree")
			},
		},
		{
			name:        "error - PullWorktree fails",
			currentPath: "/test/worktrees/dev",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					IsInWorktreeFunc: func(currentPath string) (bool, string, error) {
						return true, "dev", nil
					},
					PullWorktreeFunc: func(worktreeName string) error {
						return errors.New("git pull failed for worktree")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "git pull failed for worktree")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handlePullCurrent(mock, tt.currentPath)
			tt.expectErr(t, err)
		})
	}
}

func TestHandlePullNamed(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		mockSetup    func() *worktreePullerMock
		expectErr    func(t *testing.T, err error)
	}{
		{
			name:         "success - pull named worktree",
			worktreeName: "dev",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {
								Path:           "/test/worktrees/dev",
								ExpectedBranch: "develop",
								CurrentBranch:  "develop",
							},
						}, nil
					},
					PullWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "dev", worktreeName)
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
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {
								Path:           "/test/worktrees/dev",
								ExpectedBranch: "develop",
								CurrentBranch:  "develop",
							},
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
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return nil, errors.New("failed to enumerate worktrees")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get worktrees")
			},
		},
		{
			name:         "error - PullWorktree fails",
			worktreeName: "dev",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev": {
								Path:           "/test/worktrees/dev",
								ExpectedBranch: "develop",
								CurrentBranch:  "develop",
							},
						}, nil
					},
					PullWorktreeFunc: func(worktreeName string) error {
						return errors.New("git pull failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "git pull failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handlePullNamed(mock, tt.worktreeName)
			tt.expectErr(t, err)
		})
	}
}

// ============================================================================
// NOTE: Integration tests have been moved to internal/pull_test.go
// ============================================================================
// The integration tests that use real git repositories have been moved to the
// internal package to test Manager methods directly. This follows the pattern
// established in cmd/add.go and internal/git_add_test.go. Only fast unit tests
// using mocks remain in this file.
