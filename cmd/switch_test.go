package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwitchCommand_BasicWorktreeSwitching(t *testing.T) {

	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	_ = os.Chdir(repoPath)

	tests := []struct {
		name         string
		worktreeName string
	}{
		{
			name:         "switch to main worktree",
			worktreeName: "main",
		},
		{
			name:         "switch to dev worktree",
			worktreeName: "dev",
		},
		{
			name:         "switch to feat worktree",
			worktreeName: "feat",
		},
		{
			name:         "switch to prod worktree",
			worktreeName: "prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs([]string{"switch", tt.worktreeName})

			output, err := captureOutput(func() error {
				return cmd.Execute()
			})

			require.NoError(t, err, "Switch command should succeed for %s", tt.worktreeName)
			assert.Contains(t, output, "worktrees/"+tt.worktreeName, "Output should contain correct worktree path")
		})
	}
}

func TestSwitchCommand_PrintPathFlag(t *testing.T) {
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	tests := []struct {
		name         string
		worktreeName string
		expectedPath string
	}{
		{
			name:         "print path for main",
			worktreeName: "main",
			expectedPath: filepath.Join(repoPath, "worktrees", "main"),
		},
		{
			name:         "print path for dev",
			worktreeName: "dev",
			expectedPath: filepath.Join(repoPath, "worktrees", "dev"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs([]string{"switch", "--print-path", tt.worktreeName})

			output, err := captureOutput(func() error {
				return cmd.Execute()
			})

			require.NoError(t, err, "Switch with --print-path should succeed")
			// Use Contains instead of Equal to handle symlink path differences
			assert.Contains(t, output, tt.expectedPath, "Should contain the correct worktree path")
		})
	}
}

func TestSwitchCommand_FuzzyMatching(t *testing.T) {
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	tests := []struct {
		name          string
		input         string
		expectedMatch string
		expectedPath  string
		expectError   bool
		errorContains string
	}{
		{
			name:          "case insensitive match - dev",
			input:         "dev",
			expectedMatch: "dev",
			expectedPath:  filepath.Join(repoPath, "worktrees", "dev"),
		},
		{
			name:          "case insensitive match - main",
			input:         "main",
			expectedMatch: "main",
			expectedPath:  filepath.Join(repoPath, "worktrees", "main"),
		},
		{
			name:          "substring match - fea",
			input:         "fea",
			expectedMatch: "feat",
			expectedPath:  filepath.Join(repoPath, "worktrees", "feat"),
		},
		{
			name:          "prefix match preference - ma",
			input:         "ma",
			expectedMatch: "main",
			expectedPath:  filepath.Join(repoPath, "worktrees", "main"),
		},
		{
			name:          "nonexistent worktree",
			input:         "NONEXISTENT",
			expectError:   true,
			errorContains: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			cmd.SetArgs([]string{"switch", tt.input})

			output, err := captureOutput(func() error {
				return cmd.Execute()
			})

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err, "Fuzzy match should succeed for %s", tt.input)
				assert.Contains(t, output, "worktrees/"+tt.expectedMatch, "Should contain correct worktree path")
				if tt.expectedMatch != tt.input {
					assert.Contains(t, output, tt.expectedMatch, "Should mention the matched worktree name")
				}
			}
		})
	}
}

func TestSwitchCommand_ListWorktrees(t *testing.T) {
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"switch"})

	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	require.NoError(t, err, "List worktrees should succeed")

	// Check that all expected worktrees are listed
	expectedWorktrees := []string{"main", "dev", "feat", "prod"}
	for _, worktree := range expectedWorktrees {
		assert.Contains(t, output, worktree, "Should list worktree %s", worktree)
	}

	// Check for header
	assert.Contains(t, output, "Available worktrees", "Should show header")
}

func TestSwitchCommand_PreviousWorktree(t *testing.T) {
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	// First, switch to dev to establish a previous worktree
	cmd := newRootCommand()
	cmd.SetArgs([]string{"switch", "dev"})
	err := cmd.Execute()
	require.NoError(t, err, "Initial switch to dev should succeed")

	// Then switch to main
	cmd = newRootCommand()
	cmd.SetArgs([]string{"switch", "main"})
	err = cmd.Execute()
	require.NoError(t, err, "Switch to main should succeed")

	// Now switch back to previous (should be dev) using --print-path to get the actual path
	cmd = newRootCommand()
	cmd.SetArgs([]string{"switch", "--print-path", "-"})

	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	require.NoError(t, err, "Switch to previous worktree should succeed")
	assert.Contains(t, output, "worktrees/dev", "Should return path to dev worktree")
}

func TestSwitchCommand_NoPreviousWorktree(t *testing.T) {
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	// Try to switch to previous without any previous worktree
	cmd := newRootCommand()
	cmd.SetArgs([]string{"switch", "-"})

	err := cmd.Execute()
	require.Error(t, err, "Switch to previous should fail when no previous worktree")
	assert.Contains(t, err.Error(), "no previous worktree available", "Should mention no previous worktree")
}

func TestSwitchCommand_ShellIntegration(t *testing.T) {

	sourceRepo := testutils.NewStandardGBMConfigRepo(t)
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	_ = os.Chdir(repoPath)

	// Set shell integration environment variable
	t.Setenv("GBM_SHELL_INTEGRATION", "1")

	cmd := newRootCommand()
	cmd.SetArgs([]string{"switch", "FEAT"})

	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	require.NoError(t, err, "Switch with shell integration should succeed")

	// Extract the path from the cd command
	assert.Contains(t, output, "cd ", "Should output cd command")
	assert.Contains(t, output, "worktrees/FEAT", "Should contain correct worktree path")

	// Parse the cd command to get the target directory and verify it exists
	lines := strings.TrimSpace(output)
	require.True(t, strings.HasPrefix(lines, "cd "), "Output should start with 'cd ' command")

	targetPath := strings.TrimPrefix(lines, "cd ")

	// Verify the target path exists and is a directory
	info, err := os.Stat(targetPath)
	require.NoError(t, err, "Target directory should exist")
	assert.True(t, info.IsDir(), "Target should be a directory")

	// Verify it contains a .git file (worktree marker)
	gitFile := filepath.Join(targetPath, ".git")
	_, err = os.Stat(gitFile)
	assert.NoError(t, err, "Worktree should have .git file")
}

func TestSwitchCommand_ErrorConditions(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) string
		args          []string
		errorContains string
	}{
		{
			name: "not in git repository",
			setupFunc: func(t *testing.T) string {
				tempDir := t.TempDir()
				return tempDir
			},
			args:          []string{"switch", "dev"},
			errorContains: "not in a git repository",
		},
		{
			name: "worktree does not exist",
			setupFunc: func(t *testing.T) string {
				// Create a basic repo with worktrees
				sourceRepo := testutils.NewStandardGBMConfigRepo(t)
				repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)
				return repoPath
			},
			args:          []string{"switch", "NONEXISTENT"},
			errorContains: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc(t)
			originalDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(originalDir) }()

			_ = os.Chdir(repoPath)

			cmd := newRootCommand()
			cmd.SetArgs(tt.args)

			output, err := captureOutput(func() error {
				return cmd.Execute()
			})

			// For some cases, we expect error in the error return
			// For others (like no worktrees), it's in the output
			if err != nil {
				assert.Contains(t, err.Error(), tt.errorContains, "Error should contain expected message")
			} else {
				assert.Contains(t, output, tt.errorContains, "Output should contain expected message")
			}
		})
	}
}
