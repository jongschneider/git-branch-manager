package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to make remote changes and push them
func makeRemoteChanges(t *testing.T, repo *testutils.GitTestRepo, branch, filename, content string) {
	err := repo.InLocalRepo(func() error {
		if err := repo.SwitchToBranch(branch); err != nil {
			return err
		}

		if err := repo.WriteFile(filename, content); err != nil {
			return err
		}

		if err := repo.CommitChanges("Remote change to " + filename); err != nil {
			return err
		}

		return repo.PushBranch(branch)
	})
	require.NoError(t, err, "Failed to make remote changes")
}

// Helper function to verify file content in a worktree
func verifyWorktreeContent(t *testing.T, worktreePath, filename, expectedContent string) {
	filePath := filepath.Join(worktreePath, filename)
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read file %s", filePath)
	assert.Equal(t, expectedContent, string(content), "File content mismatch in %s", filePath)
}

// Helper function to get current commit hash in a directory
func getCurrentCommitHash(t *testing.T, dir string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to get commit hash")
	return strings.TrimSpace(string(output))
}


func TestPullCommand_CurrentWorktree(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into the dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial commit hash
	initialHash := getCurrentCommitHash(t, devWorktreePath)

	// Make remote changes to the develop branch
	makeRemoteChanges(t, sourceRepo, "develop", "new_feature.txt", "New feature content")

	// Pull current worktree (should pull dev since we're in it)
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err := cmd.Execute()
	require.NoError(t, err, "Pull command should succeed")

	// Verify the changes were pulled
	verifyWorktreeContent(t, devWorktreePath, "new_feature.txt", "New feature content")

	// Verify commit hash changed
	newHash := getCurrentCommitHash(t, devWorktreePath)
	assert.NotEqual(t, initialHash, newHash, "Commit hash should change after pull")
}

func TestPullCommand_NamedWorktree(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root (not in a worktree)
	os.Chdir(repoPath)

	// Get initial commit hash of feat worktree
	featWorktreePath := filepath.Join(repoPath, "worktrees", "feat")
	initialHash := getCurrentCommitHash(t, featWorktreePath)

	// Make remote changes to the feature/auth branch
	makeRemoteChanges(t, sourceRepo, "feature/auth", "auth_changes.txt", "Authentication improvements")

	// Pull specific worktree by name
	cmd := rootCmd
	cmd.SetArgs([]string{"pull", "feat"})

	err := cmd.Execute()
	require.NoError(t, err, "Pull command should succeed")

	// Verify the changes were pulled to feat worktree
	verifyWorktreeContent(t, featWorktreePath, "auth_changes.txt", "Authentication improvements")

	// Verify commit hash changed
	newHash := getCurrentCommitHash(t, featWorktreePath)
	assert.NotEqual(t, initialHash, newHash, "Commit hash should change after pull")
}

func TestPullCommand_AllWorktrees(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Get initial commit hashes for all worktrees
	mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	featWorktreePath := filepath.Join(repoPath, "worktrees", "feat")

	initialMainHash := getCurrentCommitHash(t, mainWorktreePath)
	initialDevHash := getCurrentCommitHash(t, devWorktreePath)
	initialFeatHash := getCurrentCommitHash(t, featWorktreePath)

	// Make remote changes to multiple branches
	makeRemoteChanges(t, sourceRepo, "main", "main_update.txt", "Main branch update")
	makeRemoteChanges(t, sourceRepo, "develop", "dev_update.txt", "Development update")
	makeRemoteChanges(t, sourceRepo, "feature/auth", "feat_update.txt", "Feature update")

	// Pull all worktrees
	cmd := rootCmd
	cmd.SetArgs([]string{"pull", "--all"})

	err := cmd.Execute()
	require.NoError(t, err, "Pull all command should succeed")

	assert.True(t, pullAll, "pullAll flag should be set to true")

	// Verify all worktrees were updated
	verifyWorktreeContent(t, mainWorktreePath, "main_update.txt", "Main branch update")
	verifyWorktreeContent(t, devWorktreePath, "dev_update.txt", "Development update")
	verifyWorktreeContent(t, featWorktreePath, "feat_update.txt", "Feature update")

	// Verify all commit hashes changed
	newMainHash := getCurrentCommitHash(t, mainWorktreePath)
	newDevHash := getCurrentCommitHash(t, devWorktreePath)
	newFeatHash := getCurrentCommitHash(t, featWorktreePath)

	assert.NotEqual(t, initialMainHash, newMainHash, "main commit hash should change")
	assert.NotEqual(t, initialDevHash, newDevHash, "dev commit hash should change")
	assert.NotEqual(t, initialFeatHash, newFeatHash, "feat commit hash should change")
}

func TestPullCommand_NotInWorktree(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root (not in a worktree)
	os.Chdir(repoPath)

	// Reset global flag state
	pullAll = false

	// Try to pull without specifying worktree name
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err := cmd.Execute()
	require.Error(t, err, "Pull should fail when not in a worktree")
	assert.Contains(t, err.Error(), "failed to check if in worktree", "Error should mention worktree check failure")
}

func TestPullCommand_NonexistentWorktree(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Reset global flag state
	pullAll = false

	// Try to pull nonexistent worktree
	cmd := rootCmd
	cmd.SetArgs([]string{"pull", "NONEXISTENT"})

	err := cmd.Execute()
	require.Error(t, err, "Pull should fail for nonexistent worktree")
	assert.Contains(t, err.Error(), "worktree 'NONEXISTENT' does not exist", "Error should mention worktree doesn't exist")
}

func TestPullCommand_NotInGitRepo(t *testing.T) {
	// Create empty temp directory (not a git repo)
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Reset global flag state
	pullAll = false

	// Try to pull in non-git directory
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err := cmd.Execute()
	require.Error(t, err, "Pull should fail when not in a git repository")
	assert.Contains(t, err.Error(), "not in a git repository", "Error should mention not being in a git repository")
}

func TestPullCommand_FastForward(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial commit hash
	initialHash := getCurrentCommitHash(t, devWorktreePath)

	// Make remote changes (clean fast-forward scenario)
	makeRemoteChanges(t, sourceRepo, "develop", "fast_forward.txt", "Fast forward content")

	// Pull changes
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err := cmd.Execute()
	require.NoError(t, err, "Fast-forward pull should succeed")

	// Verify changes were pulled
	verifyWorktreeContent(t, devWorktreePath, "fast_forward.txt", "Fast forward content")

	// Verify commit hash changed
	newHash := getCurrentCommitHash(t, devWorktreePath)
	assert.NotEqual(t, initialHash, newHash, "Commit hash should change after fast-forward")
}

func TestPullCommand_UpToDate(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial commit hash
	initialHash := getCurrentCommitHash(t, devWorktreePath)

	// Pull without any remote changes (should be up to date)
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err := cmd.Execute()
	require.NoError(t, err, "Pull should succeed even when up to date")

	// Verify commit hash unchanged
	newHash := getCurrentCommitHash(t, devWorktreePath)
	assert.Equal(t, initialHash, newHash, "Commit hash should remain same when up to date")
}

func TestPullCommand_WithLocalChanges(t *testing.T) {
	// Reset global flag state
	pullAll = false

	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Make local uncommitted changes
	localFilePath := filepath.Join(devWorktreePath, "local_changes.txt")
	err := os.WriteFile(localFilePath, []byte("Local uncommitted changes"), 0644)
	require.NoError(t, err, "Failed to create local file")

	// Make remote changes
	makeRemoteChanges(t, sourceRepo, "develop", "remote_changes.txt", "Remote changes")

	// Pull should succeed (merge scenario)
	cmd := rootCmd
	cmd.SetArgs([]string{"pull"})

	err = cmd.Execute()
	require.NoError(t, err, "Pull should handle local uncommitted changes")

	// Verify both local and remote changes exist
	verifyWorktreeContent(t, devWorktreePath, "remote_changes.txt", "Remote changes")

	// Local changes should still exist
	localContent, err := os.ReadFile(localFilePath)
	require.NoError(t, err, "Local file should still exist")
	assert.Equal(t, "Local uncommitted changes", string(localContent), "Local changes should be preserved")
}
