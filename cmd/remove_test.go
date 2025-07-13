package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// Helper function to verify worktree no longer exists
func verifyWorktreeRemoved(t *testing.T, repoPath, worktreeName string) {
	worktreePath := filepath.Join(repoPath, "worktrees", worktreeName)

	// Check directory doesn't exist
	_, err := os.Stat(worktreePath)
	assert.True(t, os.IsNotExist(err), "Worktree directory should not exist after removal")

	// Check git worktree list doesn't include it
	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to list worktrees")

	assert.NotContains(t, string(output), worktreePath, "Worktree should not appear in git worktree list")
}

// Helper function to verify worktree still exists
func verifyWorktreeExists(t *testing.T, repoPath, worktreeName string) {
	worktreePath := filepath.Join(repoPath, "worktrees", worktreeName)

	// Check directory exists
	_, err := os.Stat(worktreePath)
	assert.NoError(t, err, "Worktree directory should exist")

	// Check git worktree list includes it
	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to list worktrees")

	assert.Contains(t, string(output), worktreePath, "Worktree should appear in git worktree list")
}

// Helper function to create uncommitted changes in a worktree
func createUncommittedChanges(t *testing.T, worktreePath string) {
	filePath := filepath.Join(worktreePath, "uncommitted_changes.txt")
	err := os.WriteFile(filePath, []byte("These are uncommitted changes"), 0644)
	require.NoError(t, err, "Failed to create uncommitted changes")
}

// Helper function to check if worktree has uncommitted changes
func hasUncommittedChanges(t *testing.T, worktreePath string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to check git status")

	return len(strings.TrimSpace(string(output))) > 0
}

// Helper function to simulate user input for confirmation prompts
func simulateUserInput(input string, fn func() error) error {
	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// Write the input
	go func() {
		defer w.Close()
		w.Write([]byte(input + "\n"))
	}()

	// Execute the function
	err := fn()

	// Restore stdin
	os.Stdin = oldStdin
	r.Close()

	return err
}

// Helper function to capture command output
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestRemoveCommand_SuccessfulRemoval(t *testing.T) {

	// Create source repo with multiple branches and .envrc
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Verify FEAT worktree exists before removal
	verifyWorktreeExists(t, repoPath, "feat")

	// Remove FEAT worktree with user confirmation (simulate "y" input)
	err := simulateUserInput("y", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"remove", "feat"})
		return cmd.Execute()
	})

	require.NoError(t, err, "Remove command should succeed with user confirmation")

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "feat")
}

func TestRemoveCommand_NonexistentWorktree(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Try to remove nonexistent worktree
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "NONEXISTENT"})

	err := cmd.Execute()
	require.Error(t, err, "Remove should fail for nonexistent worktree")
	assert.Contains(t, err.Error(), "worktree 'NONEXISTENT' not found", "Error should mention worktree not found")
}

func TestRemoveCommand_NotInGitRepo(t *testing.T) {
	// Create empty temp directory (not a git repo)
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)


	// Try to remove in non-git directory
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "SOME_WORKTREE"})

	err := cmd.Execute()
	require.Error(t, err, "Remove should fail when not in a git repository")
	assert.Contains(t, err.Error(), "not in a git repository", "Error should mention not being in a git repository")
}

func TestRemoveCommand_UncommittedChangesWithoutForce(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Create uncommitted changes in main worktree
	mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
	createUncommittedChanges(t, mainWorktreePath)

	// Verify worktree has uncommitted changes
	assert.True(t, hasUncommittedChanges(t, mainWorktreePath), "main worktree should have uncommitted changes")

	// Try to remove without force flag
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "main"})

	err := cmd.Execute()
	require.Error(t, err, "Remove should fail with uncommitted changes when force not used")
	assert.Contains(t, err.Error(), "has uncommitted changes", "Error should mention uncommitted changes")
	assert.Contains(t, err.Error(), "Use --force to remove anyway", "Error should suggest using --force")

	// Verify worktree still exists
	verifyWorktreeExists(t, repoPath, "main")
}

func TestRemoveCommand_ForceWithUncommittedChanges(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Create uncommitted changes in PROD worktree
	prodWorktreePath := filepath.Join(repoPath, "worktrees", "prod")
	createUncommittedChanges(t, prodWorktreePath)

	// Verify worktree has uncommitted changes
	assert.True(t, hasUncommittedChanges(t, prodWorktreePath), "PROD worktree should have uncommitted changes")

	// Remove with force flag should succeed despite uncommitted changes
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "prod", "--force"})

	err := cmd.Execute()
	require.NoError(t, err, "Remove with --force should succeed even with uncommitted changes")

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "prod")
}

func TestRemoveCommand_ForceBypassesConfirmation(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Remove with force flag should bypass confirmation prompt
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "dev", "--force"})

	err := cmd.Execute()
	require.NoError(t, err, "Remove with --force should succeed without confirmation")

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "dev")
}

func TestRemoveCommand_UserAcceptsConfirmation(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Verify worktree exists before removal
	verifyWorktreeExists(t, repoPath, "feat")

	// Remove worktree with user accepting confirmation (simulate "y" input)
	err := simulateUserInput("y", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"remove", "feat"})
		return cmd.Execute()
	})

	require.NoError(t, err, "Remove should succeed when user accepts confirmation")

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "feat")
}

func TestRemoveCommand_UserAcceptsConfirmationWithYes(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Verify worktree exists before removal
	verifyWorktreeExists(t, repoPath, "dev")

	// Remove worktree with user accepting confirmation (simulate "yes" input)
	err := simulateUserInput("yes", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"remove", "dev"})
		return cmd.Execute()
	})

	require.NoError(t, err, "Remove should succeed when user types 'yes'")

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "dev")
}

func TestRemoveCommand_UserDeclinesConfirmation(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Verify worktree exists before attempted removal
	verifyWorktreeExists(t, repoPath, "main")

	// Remove worktree with user declining confirmation (simulate "n" input)
	err := simulateUserInput("n", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"remove", "main"})
		return cmd.Execute()
	})

	require.NoError(t, err, "Remove should complete without error when user declines")

	// Verify worktree still exists
	verifyWorktreeExists(t, repoPath, "main")
}

func TestRemoveCommand_UserDeclinesWithEmptyInput(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Verify worktree exists before attempted removal
	verifyWorktreeExists(t, repoPath, "prod")

	// Remove worktree with user providing empty input (just hitting enter)
	err := simulateUserInput("", func() error {
		cmd := newRootCommand()
		cmd.SetArgs([]string{"remove", "prod"})
		return cmd.Execute()
	})

	require.NoError(t, err, "Remove should complete without error when user provides empty input")

	// Verify worktree still exists
	verifyWorktreeExists(t, repoPath, "prod")
}

func TestRemoveCommand_RemovalFromWorktreeDirectory(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into the FEAT worktree directory
	featWorktreePath := filepath.Join(repoPath, "worktrees", "feat")
	os.Chdir(featWorktreePath)

	// Verify we're in the worktree directory (resolve any symlinks for comparison)
	currentDir, _ := os.Getwd()
	currentDir, _ = filepath.EvalSymlinks(currentDir)
	featWorktreePath, _ = filepath.EvalSymlinks(featWorktreePath)
	assert.Equal(t, featWorktreePath, currentDir, "Should be in FEAT worktree directory")

	// Remove the worktree we're currently in (with force to avoid confirmation)
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "feat", "--force"})

	err := cmd.Execute()
	require.NoError(t, err, "Remove should succeed even when executed from within the worktree")

	// Change back to repo root to verify removal
	os.Chdir(repoPath)

	// Verify worktree was removed
	verifyWorktreeRemoved(t, repoPath, "feat")
}

func TestRemoveCommand_UpdatesWorktreeList(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// First, verify all expected worktrees exist
	verifyWorktreeExists(t, repoPath, "main")
	verifyWorktreeExists(t, repoPath, "dev")
	verifyWorktreeExists(t, repoPath, "feat")
	verifyWorktreeExists(t, repoPath, "prod")

	// Get initial worktree count
	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = repoPath
	initialOutput, err := cmd.Output()
	require.NoError(t, err, "Failed to list worktrees initially")
	initialWorktrees := strings.Split(strings.TrimSpace(string(initialOutput)), "\n")

	// Remove one worktree
	removeCmd := newRootCommand()
	removeCmd.SetArgs([]string{"remove", "dev", "--force"})

	err = removeCmd.Execute()
	require.NoError(t, err, "Remove command should succeed")

	// Get updated worktree count
	cmd = exec.Command("git", "worktree", "list")
	cmd.Dir = repoPath
	updatedOutput, err := cmd.Output()
	require.NoError(t, err, "Failed to list worktrees after removal")
	updatedWorktrees := strings.Split(strings.TrimSpace(string(updatedOutput)), "\n")

	// Verify worktree count decreased by 1
	assert.Equal(t, len(initialWorktrees)-1, len(updatedWorktrees), "Worktree count should decrease by 1")

	// Verify DEV worktree no longer appears in the list
	for _, line := range updatedWorktrees {
		assert.NotContains(t, line, "worktrees/dev", "dev worktree should not appear in git worktree list")
	}

	// Verify other worktrees still exist
	verifyWorktreeExists(t, repoPath, "main")
	verifyWorktreeExists(t, repoPath, "feat")
	verifyWorktreeExists(t, repoPath, "prod")
}

func TestRemoveCommand_CleanupFilesystem(t *testing.T) {

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Add some files to the main worktree
	mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
	testFile := filepath.Join(mainWorktreePath, "test_file.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err, "Failed to create test file in worktree")

	// Verify file exists
	_, err = os.Stat(testFile)
	require.NoError(t, err, "Test file should exist before removal")

	// Remove worktree
	cmd := newRootCommand()
	cmd.SetArgs([]string{"remove", "main", "--force"})

	err = cmd.Execute()
	require.NoError(t, err, "Remove command should succeed")

	// Verify entire worktree directory is gone
	_, err = os.Stat(mainWorktreePath)
	assert.True(t, os.IsNotExist(err), "Worktree directory should be completely removed")

	// Verify test file is also gone
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err), "Files within worktree should be removed")
}
