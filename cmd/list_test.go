package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"gbm/internal"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WorktreeRow represents a parsed row from the list command output
type WorktreeRow struct {
	Name       string
	Branch     string
	GitStatus  string
	SyncStatus string
}

// parseListOutput parses the table output from `gbm list` into structured data
func parseListOutput(output string) ([]WorktreeRow, error) {
	var rows []WorktreeRow
	lines := strings.Split(output, "\n")

	// Regex to match data rows (contains │ and isn't a separator line)
	dataRowRegex := regexp.MustCompile(`^│\s*([^│]*?)\s*│\s*([^│]*?)\s*│\s*([^│]*?)\s*│\s*([^│]*?)\s*│\s*$`)

	for _, line := range lines {
		// Skip empty lines, header separators, and header row
		if strings.TrimSpace(line) == "" ||
			strings.Contains(line, "┌") || strings.Contains(line, "├") || strings.Contains(line, "└") ||
			strings.Contains(line, "WORKTREE") {
			continue
		}

		matches := dataRowRegex.FindStringSubmatch(line)
		if len(matches) == 5 { // Full match + 4 groups
			rows = append(rows, WorktreeRow{
				Name:       strings.TrimSpace(matches[1]),
				Branch:     strings.TrimSpace(matches[2]),
				GitStatus:  strings.TrimSpace(matches[3]),
				SyncStatus: strings.TrimSpace(matches[4]),
			})
		}
	}

	return rows, nil
}

// findWorktreeInRows finds a worktree by name in the parsed rows
func findWorktreeInRows(rows []WorktreeRow, name string) (WorktreeRow, bool) {
	for _, row := range rows {
		if row.Name == name {
			return row, true
		}
	}
	return WorktreeRow{}, false
}

func TestHandleList_EmptyWorktrees(t *testing.T) {
	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            true,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{},
				BranchChanges:     map[string]internal.BranchChange{},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return map[string]*internal.WorktreeListInfo{}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	assert.Equal(t, "", strings.TrimSpace(outputStr))
}

func TestHandleList_WithTrackedWorktrees(t *testing.T) {
	worktrees := map[string]*internal.WorktreeListInfo{
		"main": {
			Path:           "/path/to/worktrees/main",
			ExpectedBranch: "main",
			CurrentBranch:  "main",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
		"dev": {
			Path:           "/path/to/worktrees/dev",
			ExpectedBranch: "develop",
			CurrentBranch:  "develop",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
	}

	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            true,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{},
				BranchChanges:     map[string]internal.BranchChange{},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return worktrees, nil
		},
		GetSortedWorktreeNamesFunc: func(wt map[string]*internal.WorktreeListInfo) []string {
			return []string{"dev", "main"}
		},
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
				"dev":  "develop",
			}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	assert.Len(t, rows, 2)

	// Verify main worktree
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found)
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify dev worktree
	devWorktree, found := findWorktreeInRows(rows, "dev")
	require.True(t, found)
	assert.Equal(t, "develop", devWorktree.Branch)
	assert.Contains(t, devWorktree.SyncStatus, "IN_SYNC")
}

func TestHandleList_UntrackedWorktrees(t *testing.T) {
	worktrees := map[string]*internal.WorktreeListInfo{
		"main": {
			Path:           "/path/to/worktrees/main",
			ExpectedBranch: "main",
			CurrentBranch:  "main",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
		"UNTRACKED": {
			Path:           "/path/to/worktrees/UNTRACKED",
			ExpectedBranch: "",
			CurrentBranch:  "develop",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
	}

	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            false,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{},
				BranchChanges:     map[string]internal.BranchChange{},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return worktrees, nil
		},
		GetSortedWorktreeNamesFunc: func(wt map[string]*internal.WorktreeListInfo) []string {
			return []string{"main", "UNTRACKED"}
		},
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
				// UNTRACKED is not in mapping
			}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	assert.Len(t, rows, 2)

	// Verify main worktree
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found)
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify UNTRACKED worktree
	untrackedWorktree, found := findWorktreeInRows(rows, "UNTRACKED")
	require.True(t, found)
	assert.Equal(t, "develop", untrackedWorktree.Branch)
	assert.Contains(t, untrackedWorktree.SyncStatus, "UNTRACKED")
}

func TestHandleList_OrphanedWorktrees(t *testing.T) {
	worktrees := map[string]*internal.WorktreeListInfo{
		"main": {
			Path:           "/path/to/worktrees/main",
			ExpectedBranch: "main",
			CurrentBranch:  "main",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
		"DEV": {
			Path:           "/path/to/worktrees/DEV",
			ExpectedBranch: "develop",
			CurrentBranch:  "develop",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
	}

	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            false,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{"DEV"}, // DEV is orphaned
				BranchChanges:     map[string]internal.BranchChange{},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return worktrees, nil
		},
		GetSortedWorktreeNamesFunc: func(wt map[string]*internal.WorktreeListInfo) []string {
			return []string{"main", "DEV"}
		},
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
				// DEV not in mapping anymore (orphaned)
			}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	assert.Len(t, rows, 2)

	// Verify main worktree (should be in sync)
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found)
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify DEV worktree (should be untracked/orphaned)
	devWorktree, found := findWorktreeInRows(rows, "DEV")
	require.True(t, found)
	assert.Equal(t, "develop", devWorktree.Branch)
	assert.Contains(t, devWorktree.SyncStatus, "UNTRACKED")
}

func TestHandleList_BranchChanges(t *testing.T) {
	worktrees := map[string]*internal.WorktreeListInfo{
		"main": {
			Path:           "/path/to/worktrees/main",
			ExpectedBranch: "main",
			CurrentBranch:  "main",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
		"DEV": {
			Path:           "/path/to/worktrees/DEV",
			ExpectedBranch: "develop",
			CurrentBranch:  "feature/auth",
			GitStatus:      &internal.GitStatus{IsDirty: false},
		},
	}

	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            false,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{},
				BranchChanges: map[string]internal.BranchChange{
					"DEV": {
						OldBranch: "develop",
						NewBranch: "feature/auth",
					},
				},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return worktrees, nil
		},
		GetSortedWorktreeNamesFunc: func(wt map[string]*internal.WorktreeListInfo) []string {
			return []string{"main", "DEV"}
		},
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
				"DEV":  "develop",
			}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	assert.Len(t, rows, 2)

	// Verify main worktree
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found)
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify DEV worktree shows branch change
	devWorktree, found := findWorktreeInRows(rows, "DEV")
	require.True(t, found)
	assert.Equal(t, "feature/auth (expected: develop)", devWorktree.Branch)
	assert.Contains(t, devWorktree.SyncStatus, "OUT_OF_SYNC (develop → feature/auth)")
}

func TestHandleList_GitStatusDisplay(t *testing.T) {
	worktrees := map[string]*internal.WorktreeListInfo{
		"main": {
			Path:           "/path/to/worktrees/main",
			ExpectedBranch: "main",
			CurrentBranch:  "main",
			GitStatus:      &internal.GitStatus{IsDirty: true, Modified: 1},
		},
	}

	mock := &worktreeListerMock{
		GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
			return &internal.SyncStatus{
				InSync:            true,
				MissingWorktrees:  []string{},
				OrphanedWorktrees: []string{},
				BranchChanges:     map[string]internal.BranchChange{},
			}, nil
		},
		GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
			return worktrees, nil
		},
		GetSortedWorktreeNamesFunc: func(wt map[string]*internal.WorktreeListInfo) []string {
			return []string{"main"}
		},
		GetWorktreeMappingFunc: func() (map[string]string, error) {
			return map[string]string{
				"main": "main",
			}, nil
		},
	}

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err := handleList(mock, cmd)
	require.NoError(t, err)

	outputStr := output.String()
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	assert.Len(t, rows, 1)

	// Verify main worktree shows git status
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found)
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.NotEmpty(t, mainWorktree.GitStatus, "Git status should not be empty")
}

func TestHandleList_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func() *worktreeListerMock
		expectError bool
		errorMsg    string
	}{
		{
			name: "GetSyncStatus error",
			mockSetup: func() *worktreeListerMock {
				return &worktreeListerMock{
					GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
						return nil, fmt.Errorf("sync status error")
					},
				}
			},
			expectError: true,
			errorMsg:    "sync status error",
		},
		{
			name: "GetAllWorktrees error",
			mockSetup: func() *worktreeListerMock {
				return &worktreeListerMock{
					GetSyncStatusFunc: func() (*internal.SyncStatus, error) {
						return &internal.SyncStatus{}, nil
					},
					GetAllWorktreesFunc: func() (map[string]*internal.WorktreeListInfo, error) {
						return nil, fmt.Errorf("worktree list error")
					},
				}
			},
			expectError: true,
			errorMsg:    "failed to get worktree list: worktree list error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			cmd := &cobra.Command{}
			var output bytes.Buffer
			cmd.SetOut(&output)

			err := handleList(mock, cmd)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
