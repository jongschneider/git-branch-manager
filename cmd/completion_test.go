package cmd

import (
	"errors"
	"testing"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

// mockWorktreeProvider implements the worktreeProvider interface for testing
type mockWorktreeProvider struct {
	worktrees map[string]*internal.WorktreeListInfo
	err       error
}

func (m *mockWorktreeProvider) GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error) {
	return m.worktrees, m.err
}

func TestGetWorktreeCompletions(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func() *mockWorktreeProvider
		expectResult func(t *testing.T, completions []string)
	}{
		{
			name: "success - return formatted worktree completions",
			mockSetup: func() *mockWorktreeProvider {
				return &mockWorktreeProvider{
					worktrees: map[string]*internal.WorktreeListInfo{
						"main": {CurrentBranch: "main"},
						"dev":  {CurrentBranch: "develop"},
						"feat": {CurrentBranch: "feature-123"},
					},
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				// Expect tab-separated format with branch info and proper alignment
				expected := []string{
					"main\t    main",
					"dev\t     develop",
					"feat\t    feature-123",
				}
				assert.ElementsMatch(t, expected, completions)
			},
		},
		{
			name: "error - GetAllWorktrees fails",
			mockSetup: func() *mockWorktreeProvider {
				return &mockWorktreeProvider{
					err: errors.New("failed to get worktrees"),
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				assert.Nil(t, completions)
			},
		},
		{
			name: "success - empty worktrees",
			mockSetup: func() *mockWorktreeProvider {
				return &mockWorktreeProvider{
					worktrees: map[string]*internal.WorktreeListInfo{},
				}
			},
			expectResult: func(t *testing.T, completions []string) {
				assert.Empty(t, completions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			completions := getWorktreeCompletions(mock)

			// Validate results
			tt.expectResult(t, completions)
		})
	}
}