package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gbm/internal"
	"gbm/internal/testutils"

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

func TestListCommand_EmptyRepository(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(repo.GetLocalPath())

	assert.NoError(t, repo.WriteFile(internal.DefaultBranchConfigFilename, "# Empty gbm.branchconfig.yaml"))
	assert.NoError(t, repo.CommitChanges("Add empty gbm.branchconfig.yaml"))

	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err := cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()
	assert.Equal(t, "", strings.TrimSpace(outputStr))
}

func TestListCommand_WithGBMConfigWorktrees(t *testing.T) {
	repo := testutils.NewGBMConfigRepo(t, map[string]string{
		"main": "main",
		"dev":  "develop",
		"feat": "feature/auth",
	})

	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(repo.GetLocalPath())

	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err := cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()
	// Since no worktrees exist, list should show nothing
	assert.Equal(t, "", strings.TrimSpace(outputStr))
}

func TestListCommand_UntrackedWorktrees(t *testing.T) {
	// Create a repository with branches
	sourceRepo := testutils.NewMultiBranchRepo(t)

	// Clone the repository (no sync needed since NewMultiBranchRepo doesn't have gbm.branchconfig.yaml)
	setupClonedRepo(t, sourceRepo)

	// Create an additional worktree that's not in gbm.branchconfig.yaml (untracked)
	addCmd := newRootCommand()
	addCmd.SetArgs([]string{"add", "--new-branch", "UNTRACKED", "develop"})
	err := addCmd.Execute()
	require.NoError(t, err)

	// Now test the list command
	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err = cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()

	// Parse the table output
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	// Should have exactly 2 worktrees
	assert.Len(t, rows, 2)

	// Verify main worktree
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found, "main worktree should be present")
	assert.Equal(t, "main", mainWorktree.Branch)
	// With the improved clone command, main worktree is now properly tracked in config
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify UNTRACKED worktree
	untrackedWorktree, found := findWorktreeInRows(rows, "UNTRACKED")
	require.True(t, found, "UNTRACKED worktree should be present")
	assert.Equal(t, "develop", untrackedWorktree.Branch)
	assert.Contains(t, untrackedWorktree.SyncStatus, "UNTRACKED")
}

func TestListCommand_OrphanedWorktrees(t *testing.T) {
	// Create a repository with branches
	sourceRepo := testutils.NewGBMConfigRepo(t, map[string]string{
		"main": "main",
		"DEV":  "develop",
	})

	// Clone the repository and sync worktrees
	setupClonedRepoWithWorktrees(t, sourceRepo)

	// Remove DEV from gbm.branchconfig.yaml to make it orphaned
	err := os.WriteFile(internal.DefaultBranchConfigFilename, []byte("worktrees:\n  main:\n    branch: main\n"), 0o644)
	require.NoError(t, err)

	// Now test the list command
	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err = cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()

	// Parse the table output
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	// Should have exactly 2 worktrees
	assert.Len(t, rows, 2)

	// Verify MAIN worktree (should be in sync)
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found, "main worktree should be present")
	assert.Equal(t, "main", mainWorktree.Branch)
	assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

	// Verify DEV worktree (should be untracked/orphaned)
	devWorktree, found := findWorktreeInRows(rows, "DEV")
	require.True(t, found, "DEV worktree should be present")
	assert.Equal(t, "develop", devWorktree.Branch)
	assert.Contains(t, devWorktree.SyncStatus, "UNTRACKED")
}

func TestListCommand_GitStatus(t *testing.T) {
	// Create a repository with branches
	sourceRepo := testutils.NewMultiBranchRepo(t)

	// Clone the repository (no sync needed since NewMultiBranchRepo doesn't have gbm.branchconfig.yaml)
	repoPath := setupClonedRepo(t, sourceRepo)

	// Create a file in the main worktree to create git status
	mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
	err := os.WriteFile(filepath.Join(mainWorktreePath, "test.txt"), []byte("test content"), 0o644)
	require.NoError(t, err)

	// Now test the list command
	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err = cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()

	// Parse the table output
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	// Should have exactly 1 worktree (main)
	assert.Len(t, rows, 1)

	// Verify main worktree
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found, "main worktree should be present")
	assert.Equal(t, "main", mainWorktree.Branch)
	// Should have some git status indication (the exact symbol may vary)
	assert.NotEmpty(t, mainWorktree.GitStatus, "Git status should not be empty")
}

func TestListCommand_ExpectedBranchDisplay(t *testing.T) {
	// Create a repository with branches
	sourceRepo := testutils.NewGBMConfigRepo(t, map[string]string{
		"main": "main",
		"DEV":  "develop",
	})

	// Clone the repository and sync worktrees
	repoPath := setupClonedRepoWithWorktrees(t, sourceRepo)

	// Change DEV worktree to a different branch to test expected branch display
	devWorktreePath := filepath.Join(repoPath, "worktrees", "DEV")
	_ = os.Chdir(devWorktreePath)

	// Switch to feature/auth branch instead of develop to create branch mismatch
	gitOutput, err := exec.Command("git", "checkout", "feature/auth").CombinedOutput()
	if err != nil {
		t.Logf("Git checkout failed: %s", string(gitOutput))
	}
	require.NoError(t, err)

	// Change back to repo root
	_ = os.Chdir(repoPath)

	// Now test the list command
	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err = cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()

	// Parse the table output
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	// Should have exactly 2 worktrees
	assert.Len(t, rows, 2)

	// Verify main worktree (should show just "main")
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found, "main worktree should be present")
	assert.Equal(t, "main", mainWorktree.Branch)

	// Verify DEV worktree (should show "feature/auth (expected: develop)")
	devWorktree, found := findWorktreeInRows(rows, "DEV")
	require.True(t, found, "DEV worktree should be present")
	assert.Equal(t, "feature/auth (expected: develop)", devWorktree.Branch)
}

func TestListCommand_SortedOutput(t *testing.T) {
	// Create a repository with branches
	sourceRepo := testutils.NewGBMConfigRepo(t, map[string]string{
		"main": "main",
		"dev":  "develop",
		"feat": "feature/auth",
	})

	// Clone the repository and sync worktrees
	setupClonedRepoWithWorktrees(t, sourceRepo)

	// Create an additional ad-hoc worktree to test sorting
	addCmd := newRootCommand()
	addCmd.SetArgs([]string{"add", "--new-branch", "adhoc", "production/v2.0"})
	err := addCmd.Execute()
	require.NoError(t, err)

	// Now test the list command
	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	var output bytes.Buffer
	cmd.SetOut(&output)

	err = cmd.Execute()
	require.NoError(t, err)

	outputStr := output.String()

	// Parse the table output
	rows, err := parseListOutput(outputStr)
	require.NoError(t, err)

	// Should have exactly 4 worktrees
	assert.Len(t, rows, 4)

	// Verify all expected worktrees are present
	mainWorktree, found := findWorktreeInRows(rows, "main")
	require.True(t, found, "main worktree should be present")
	assert.Equal(t, "main", mainWorktree.Branch)

	devWorktree, found := findWorktreeInRows(rows, "dev")
	require.True(t, found, "dev worktree should be present")
	assert.Equal(t, "develop", devWorktree.Branch)

	featWorktree, found := findWorktreeInRows(rows, "feat")
	require.True(t, found, "feat worktree should be present")
	assert.Equal(t, "feature/auth", featWorktree.Branch)

	adhocWorktree, found := findWorktreeInRows(rows, "adhoc")
	require.True(t, found, "adhoc worktree should be present")
	assert.Equal(t, "production/v2.0", adhocWorktree.Branch)

	// Verify sorting: gbm.branchconfig.yaml worktrees (main, dev, feat) should come before ad-hoc (adhoc)
	// The order in the rows slice should reflect the display order
	worktreeNames := make([]string, len(rows))
	for i, row := range rows {
		worktreeNames[i] = row.Name
	}

	// Find positions
	mainPos := -1
	devPos := -1
	featPos := -1
	adhocPos := -1
	for i, name := range worktreeNames {
		switch name {
		case "main":
			mainPos = i
		case "dev":
			devPos = i
		case "feat":
			featPos = i
		case "adhoc":
			adhocPos = i
		}
	}

	// gbm.branchconfig.yaml worktrees should come before ad-hoc worktrees
	assert.True(t, mainPos < adhocPos, "main should come before adhoc")
	assert.True(t, devPos < adhocPos, "dev should come before adhoc")
	assert.True(t, featPos < adhocPos, "feat should come before adhoc")

	// Verify the exact order of worktrees as they appear in the output
	// Tracked worktrees should be sorted alphabetically, followed by ad-hoc worktrees
	expectedOrder := []string{"dev", "feat", "main", "adhoc"}
	actualOrder := worktreeNames
	assert.Equal(t, expectedOrder, actualOrder, "Worktrees should be in the expected sorted order")
}

func TestListCommand_NoGitRepository(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(tempDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find git repository root")
}

func TestListCommand_NoGBMConfigFile(t *testing.T) {
	repo := testutils.NewBasicRepo(t)

	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(repo.GetLocalPath())

	gbmConfigPath := filepath.Join(repo.GetLocalPath(), internal.DefaultBranchConfigFilename)
	if _, err := os.Stat(gbmConfigPath); err == nil {
		_ = os.Remove(gbmConfigPath)
	}

	cmd := newRootCommand()
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load gbm.branchconfig.yaml")
}
