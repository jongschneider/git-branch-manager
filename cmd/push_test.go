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

// Helper function to make local changes in a worktree
func makeLocalChanges(t *testing.T, worktreePath, filename, content string) {
	filePath := filepath.Join(worktreePath, filename)
	err := os.WriteFile(filePath, []byte(content), 0o644)
	require.NoError(t, err, "Failed to write file %s", filePath)

	// Stage and commit the changes
	gitCmd := exec.Command("git", "add", filename)
	gitCmd.Dir = worktreePath
	err = gitCmd.Run()
	require.NoError(t, err, "Failed to stage file %s", filename)

	gitCmd = exec.Command("git", "commit", "-m", "Local change to "+filename)
	gitCmd.Dir = worktreePath
	err = gitCmd.Run()
	require.NoError(t, err, "Failed to commit changes to %s", filename)
}

// Helper function to verify a commit exists in the remote repository
func verifyRemoteHasCommit(t *testing.T, repo *testutils.GitTestRepo, branch, commitMessage string) {
	err := repo.InLocalRepo(func() error {
		if err := repo.SwitchToBranch(branch); err != nil {
			return err
		}

		// Fetch latest changes from remote first
		gitCmd := exec.Command("git", "pull", "origin", branch)
		gitCmd.Dir = repo.GetLocalPath()
		if err := gitCmd.Run(); err != nil {
			return err
		}

		gitCmd = exec.Command("git", "log", "--oneline", "-n", "10")
		gitCmd.Dir = repo.GetLocalPath()
		output, err := gitCmd.Output()
		if err != nil {
			return err
		}

		if !strings.Contains(string(output), commitMessage) {
			t.Errorf("Commit message '%s' not found in remote branch %s", commitMessage, branch)
		}

		return nil
	})
	require.NoError(t, err, "Failed to verify remote commit")
}

// Helper function to get remote commit hash
func getRemoteCommitHash(t *testing.T, repo *testutils.GitTestRepo, branch string) string {
	var hash string
	err := repo.InLocalRepo(func() error {
		if err := repo.SwitchToBranch(branch); err != nil {
			return err
		}

		// Fetch latest changes from remote first
		gitCmd := exec.Command("git", "pull", "origin", branch)
		gitCmd.Dir = repo.GetLocalPath()
		if err := gitCmd.Run(); err != nil {
			return err
		}

		gitCmd = exec.Command("git", "rev-parse", "HEAD")
		gitCmd.Dir = repo.GetLocalPath()
		output, err := gitCmd.Output()
		if err != nil {
			return err
		}

		hash = strings.TrimSpace(string(output))
		return nil
	})
	require.NoError(t, err, "Failed to get remote commit hash")
	return hash
}

// Helper function to check if upstream is configured
func checkUpstreamExists(t *testing.T, worktreePath string) bool {
	gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	gitCmd.Dir = worktreePath
	err := gitCmd.Run()
	return err == nil
}

func TestPushCommand_CurrentWorktree(t *testing.T) {
	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into the dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial remote commit hash
	initialHash := getRemoteCommitHash(t, sourceRepo, "develop")

	// Make local changes
	makeLocalChanges(t, devWorktreePath, "push_test.txt", "Local changes to push")

	// Push current worktree (should push DEV since we're in it)
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.NoError(t, err, "Push command should succeed")

	// Verify the changes were pushed to remote
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to push_test.txt")

	// Verify remote commit hash changed
	newHash := getRemoteCommitHash(t, sourceRepo, "develop")
	assert.NotEqual(t, initialHash, newHash, "Remote commit hash should change after push")
}

func TestPushCommand_NamedWorktree(t *testing.T) {
	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root (not in a worktree)
	os.Chdir(repoPath)

	// Get initial remote commit hash for feature/auth branch
	initialHash := getRemoteCommitHash(t, sourceRepo, "feature/auth")

	// Make local changes in feat worktree
	featWorktreePath := filepath.Join(repoPath, "worktrees", "feat")
	makeLocalChanges(t, featWorktreePath, "named_push.txt", "Changes pushed by name")

	// Push specific worktree by name
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push", "feat"})

	err := cmd.Execute()
	require.NoError(t, err, "Push command should succeed")

	// Verify the changes were pushed to remote
	verifyRemoteHasCommit(t, sourceRepo, "feature/auth", "Local change to named_push.txt")

	// Verify remote commit hash changed
	newHash := getRemoteCommitHash(t, sourceRepo, "feature/auth")
	assert.NotEqual(t, initialHash, newHash, "Remote commit hash should change after push")
}

func TestPushCommand_AllWorktrees(t *testing.T) {
	// Create source repo with multiple branches and .gbm.config.yaml
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Get initial remote commit hashes for all branches
	initialMainHash := getRemoteCommitHash(t, sourceRepo, "main")
	initialDevHash := getRemoteCommitHash(t, sourceRepo, "develop")
	initialFeatHash := getRemoteCommitHash(t, sourceRepo, "feature/auth")

	// Make local changes in multiple worktrees
	mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	featWorktreePath := filepath.Join(repoPath, "worktrees", "feat")

	makeLocalChanges(t, mainWorktreePath, "main_push.txt", "Main branch changes")
	makeLocalChanges(t, devWorktreePath, "dev_push.txt", "Development changes")
	makeLocalChanges(t, featWorktreePath, "feat_push.txt", "Feature changes")

	// Push all worktrees
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push", "--all"})

	err := cmd.Execute()
	require.NoError(t, err, "Push all command should succeed")

	// Verify all changes were pushed to remote
	verifyRemoteHasCommit(t, sourceRepo, "main", "Local change to main_push.txt")
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to dev_push.txt")
	verifyRemoteHasCommit(t, sourceRepo, "feature/auth", "Local change to feat_push.txt")

	// Verify all remote commit hashes changed
	newMainHash := getRemoteCommitHash(t, sourceRepo, "main")
	newDevHash := getRemoteCommitHash(t, sourceRepo, "develop")
	newFeatHash := getRemoteCommitHash(t, sourceRepo, "feature/auth")

	assert.NotEqual(t, initialMainHash, newMainHash, "main remote commit hash should change")
	assert.NotEqual(t, initialDevHash, newDevHash, "DEV remote commit hash should change")
	assert.NotEqual(t, initialFeatHash, newFeatHash, "FEAT remote commit hash should change")
}

func TestPushCommand_WithExistingUpstream(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Verify upstream exists (should be set during sync)
	assert.True(t, checkUpstreamExists(t, devWorktreePath), "Upstream should be configured")

	// Make local changes
	makeLocalChanges(t, devWorktreePath, "upstream_test.txt", "Changes with existing upstream")

	// Push changes
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.NoError(t, err, "Push with existing upstream should succeed")

	// Verify changes were pushed
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to upstream_test.txt")
}

func TestPushCommand_WithoutUpstream(t *testing.T) {
	// Create source repo and manually create a new branch without upstream
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Create a new branch in one of the worktrees without upstream
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Create and switch to new branch
	gitCmd := exec.Command("git", "checkout", "-b", "new-feature")
	gitCmd.Dir = devWorktreePath
	err := gitCmd.Run()
	require.NoError(t, err, "Failed to create new branch")

	// Verify no upstream exists for new branch
	assert.False(t, checkUpstreamExists(t, devWorktreePath), "New branch should not have upstream")

	// Make local changes
	makeLocalChanges(t, devWorktreePath, "no_upstream.txt", "Changes on new branch")

	// Push changes (should set upstream)
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err = cmd.Execute()
	require.NoError(t, err, "Push without upstream should succeed")

	// Verify upstream was set
	assert.True(t, checkUpstreamExists(t, devWorktreePath), "Upstream should be set after push")
}

func TestPushCommand_NotInWorktree(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root (not in a worktree)
	os.Chdir(repoPath)

	// Try to push without specifying worktree name
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.Error(t, err, "Push should fail when not in a worktree")
	assert.Contains(t, err.Error(), "failed to check if in worktree", "Error should mention worktree check failure")
}

func TestPushCommand_NonexistentWorktree(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Stay in repo root
	os.Chdir(repoPath)

	// Try to push nonexistent worktree
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push", "NONEXISTENT"})

	err := cmd.Execute()
	require.Error(t, err, "Push should fail for nonexistent worktree")
	assert.Contains(t, err.Error(), "worktree 'NONEXISTENT' does not exist", "Error should mention worktree doesn't exist")
}

func TestPushCommand_NotInGitRepo(t *testing.T) {
	// Create empty temp directory (not a git repo)
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Try to push in non-git directory
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.Error(t, err, "Push should fail when not in a git repository")
	assert.Contains(t, err.Error(), "not in a git repository", "Error should mention not being in a git repository")
}

func TestPushCommand_WithLocalCommits(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial remote commit hash
	initialHash := getRemoteCommitHash(t, sourceRepo, "develop")

	// Make multiple local commits
	makeLocalChanges(t, devWorktreePath, "commit1.txt", "First commit")
	makeLocalChanges(t, devWorktreePath, "commit2.txt", "Second commit")
	makeLocalChanges(t, devWorktreePath, "commit3.txt", "Third commit")

	// Push all local commits
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.NoError(t, err, "Push with multiple commits should succeed")

	// Verify all commits were pushed
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to commit1.txt")
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to commit2.txt")
	verifyRemoteHasCommit(t, sourceRepo, "develop", "Local change to commit3.txt")

	// Verify remote commit hash changed
	newHash := getRemoteCommitHash(t, sourceRepo, "develop")
	assert.NotEqual(t, initialHash, newHash, "Remote commit hash should change after pushing multiple commits")
}

func TestPushCommand_UpToDate(t *testing.T) {
	// Create source repo with worktrees
	sourceRepo := testutils.NewStandardGBMConfigRepo(t)

	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Navigate into dev worktree
	devWorktreePath := filepath.Join(repoPath, "worktrees", "dev")
	os.Chdir(devWorktreePath)

	// Get initial remote commit hash
	initialHash := getRemoteCommitHash(t, sourceRepo, "develop")

	// Push without any local changes (should be up to date)
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push"})

	err := cmd.Execute()
	require.NoError(t, err, "Push should succeed even when up to date")

	// Verify remote commit hash unchanged
	newHash := getRemoteCommitHash(t, sourceRepo, "develop")
	assert.Equal(t, initialHash, newHash, "Remote commit hash should remain same when up to date")
}

func TestPushCommand_EmptyWorktreeList(t *testing.T) {
	// Create basic repo without worktrees
	sourceRepo := testutils.NewBasicRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	// Clone the repository but don't sync worktrees
	cloneCmd := newRootCommand()
	cloneCmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})
	err := cloneCmd.Execute()
	require.NoError(t, err, "Failed to clone repository")

	// Navigate to cloned repo
	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)
	os.Chdir(repoPath)

	// Try to push all worktrees when none exist
	cmd := newRootCommand()
	cmd.SetArgs([]string{"push", "--all"})

	err = cmd.Execute()
	require.NoError(t, err, "Push all should succeed even with no worktrees")
}
