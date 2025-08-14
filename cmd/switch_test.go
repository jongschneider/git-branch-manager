package cmd

import (
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

func TestHandleSwitchToWorktree_ExactMatch(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *worktreeSwitcherMock
		args      []string
		printPath bool
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "successful exact match",
			args: []string{"dev"},
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "dev", worktreeName)
						return "/path/to/dev", nil
					},
					SetCurrentWorktreeFunc: func(worktreeName string) error {
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
			name:      "successful exact match with print path",
			args:      []string{"main"},
			printPath: true,
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						assert.Equal(t, "main", worktreeName)
						return "/path/to/main", nil
					},
					SetCurrentWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "main", worktreeName)
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handleSwitchToWorktree(mock, tt.args[0], tt.printPath)
			tt.expectErr(t, err)

			// Verify mocks were called
			assert.Len(t, mock.GetWorktreePathCalls(), 1)
			assert.Len(t, mock.SetCurrentWorktreeCalls(), 1)
		})
	}
}

func TestHandleSwitchToWorktree_FuzzyMatch(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *worktreeSwitcherMock
		input     string
		expectErr func(t *testing.T, err error)
	}{
		{
			name:  "fuzzy match success",
			input: "fea",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetWorktreePathFunc: func(worktreeName string) (string, error) {
						if worktreeName == "fea" {
							return "", assert.AnError // Simulate exact match failing
						}
						if worktreeName == "feat" {
							return "/path/to/feat", nil
						}
						return "", assert.AnError
					},
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"feat": {Path: "/path/to/feat"},
							"dev":  {Path: "/path/to/dev"},
						}, nil
					},
					SetCurrentWorktreeFunc: func(worktreeName string) error {
						assert.Equal(t, "feat", worktreeName)
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handleSwitchToWorktree(mock, tt.input, false)
			tt.expectErr(t, err)

			// Verify GetWorktreePath was called twice (exact match + fuzzy match result)
			assert.Len(t, mock.GetWorktreePathCalls(), 2)
			// Verify GetAllWorktrees was called for fuzzy matching
			assert.Len(t, mock.GetAllWorktreesCalls(), 1)
			// Verify SetCurrentWorktree was called with the matched name
			assert.Len(t, mock.SetCurrentWorktreeCalls(), 1)
		})
	}
}

func TestHandleListWorktrees(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *worktreeSwitcherMock
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "list worktrees success",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"main": {Path: "/path/to/main", CurrentBranch: "main", ExpectedBranch: "main"},
							"dev":  {Path: "/path/to/dev", CurrentBranch: "dev", ExpectedBranch: "dev"},
						}, nil
					},
					GetSortedWorktreeNamesFunc: func(worktrees map[string]*internal.WorktreeListInfo) []string {
						return []string{"main", "dev"}
					},
					GetStatusIconFunc: func(gitStatus *internal.GitStatus) string {
						return "✓"
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "no worktrees found",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{}, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handleListWorktrees(mock)
			tt.expectErr(t, err)

			// Verify GetAllWorktrees was called
			assert.Len(t, mock.GetAllWorktreesCalls(), 1)

			// For the first test case, we expect additional method calls
			if tt.name == "list worktrees success" {
				assert.Len(t, mock.GetSortedWorktreeNamesCalls(), 1)
			}
		})
	}
}

func TestFindFuzzyMatch(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func() *worktreeSwitcherMock
		target       string
		expectedName string
	}{
		{
			name:   "exact case insensitive match",
			target: "DEV",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev":  {Path: "/path/to/dev"},
							"main": {Path: "/path/to/main"},
						}, nil
					},
				}
			},
			expectedName: "dev",
		},
		{
			name:   "substring match",
			target: "fea",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"feat": {Path: "/path/to/feat"},
							"dev":  {Path: "/path/to/dev"},
						}, nil
					},
				}
			},
			expectedName: "feat",
		},
		{
			name:   "prefix match preference",
			target: "mai",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"main":   {Path: "/path/to/main"},
							"manual": {Path: "/path/to/manual"},
						}, nil
					},
				}
			},
			expectedName: "main", // Only "main" starts with "mai"
		},
		{
			name:   "no match",
			target: "nonexistent",
			mockSetup: func() *worktreeSwitcherMock {
				return &worktreeSwitcherMock{
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return map[string]*internal.WorktreeListInfo{
							"dev":  {Path: "/path/to/dev"},
							"main": {Path: "/path/to/main"},
						}, nil
					},
				}
			},
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			result := findFuzzyMatch(mock, tt.target)
			assert.Equal(t, tt.expectedName, result)

			// Verify GetAllWorktrees was called
			assert.Len(t, mock.GetAllWorktreesCalls(), 1)
		})
	}
}
