package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/require"
)

// setupClonedRepo clones a source repository to a temp directory and changes to that directory.
// The original working directory is automatically restored when the test completes via t.Cleanup.
// Returns the cloned repository path.
func setupClonedRepo(t *testing.T, sourceRepo *testutils.GitTestRepo) string {
	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()

	// Register cleanup to restore original directory when test completes
	t.Cleanup(func() { os.Chdir(originalDir) })

	os.Chdir(targetDir)

	// Clone the repository
	cloneCmd := newRootCommand()
	cloneCmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})
	err := cloneCmd.Execute()
	require.NoError(t, err, "Failed to clone repository")

	// Navigate to cloned repo
	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)
	os.Chdir(repoPath)

	return repoPath
}

// setupClonedRepoWithWorktrees clones a source repository to a temp directory and syncs worktrees.
// The original working directory is automatically restored when the test completes via t.Cleanup.
// Returns the cloned repository path.
func setupClonedRepoWithWorktrees(t *testing.T, sourceRepo *testutils.GitTestRepo) string {
	repoPath := setupClonedRepo(t, sourceRepo)

	// Sync worktrees
	syncCmd := newRootCommand()
	syncCmd.SetArgs([]string{"sync"})
	err := syncCmd.Execute()
	require.NoError(t, err, "Failed to sync worktrees")

	return repoPath
}
