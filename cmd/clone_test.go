package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal"
	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseGBMConfig reads and unmarshals a gbm.branchconfig.yaml file
func parseGBMConfig(t *testing.T, path string) *internal.GBMConfig {
	t.Helper()

	config, err := internal.ParseGBMConfig(path)
	require.NoError(t, err)

	return config
}

func TestCloneCommand_Basic(t *testing.T) {
	sourceRepo := testutils.NewMultiBranchRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(targetDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)

	assert.DirExists(t, repoPath)
	assert.DirExists(t, filepath.Join(repoPath, ".git"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
	assert.FileExists(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	
	// Verify .gbm directory and files are created
	assert.DirExists(t, filepath.Join(repoPath, ".gbm"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "config.toml"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "state.toml"))

	config := parseGBMConfig(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	expected := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				Description: "Main production branch",
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestCloneCommand_WithExistingGBMConfig(t *testing.T) {
	sourceRepo := testutils.NewGBMConfigRepo(t, map[string]string{
		"main": "main",
		"dev":  "develop",
		"feat": "feature/auth",
	})

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(targetDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)

	assert.FileExists(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	
	// Verify .gbm directory and files are created
	assert.DirExists(t, filepath.Join(repoPath, ".gbm"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "config.toml"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "state.toml"))

	config := parseGBMConfig(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	expected := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				Description: "Main branch",
			},
			"dev": {
				Branch:      "develop",
				MergeInto:   "main",
				Description: "Dev branch",
			},
			"feat": {
				Branch:      "feature/auth",
				MergeInto:   "dev",
				Description: "Feat branch",
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestCloneCommand_WithoutGBMConfig(t *testing.T) {
	sourceRepo := testutils.NewBasicRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(targetDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)

	assert.FileExists(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	
	// Verify .gbm directory and files are created
	assert.DirExists(t, filepath.Join(repoPath, ".gbm"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "config.toml"))
	assert.FileExists(t, filepath.Join(repoPath, ".gbm", "state.toml"))

	config := parseGBMConfig(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
	expected := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				Description: "Main production branch",
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestCloneCommand_DifferentDefaultBranches(t *testing.T) {
	tests := []struct {
		name          string
		defaultBranch string
	}{
		{"master branch", "master"},
		{"develop branch", "develop"},
		{"custom branch", "custom-main"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepo := testutils.NewGitTestRepo(t,
				testutils.WithDefaultBranch(tt.defaultBranch),
				testutils.WithUser("Test User", "test@example.com"),
			)

			targetDir := t.TempDir()
			originalDir, _ := os.Getwd()
			t.Cleanup(func() { _ = os.Chdir(originalDir) })

			_ = os.Chdir(targetDir)

			cmd := newRootCommand()
			cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

			err := cmd.Execute()
			require.NoError(t, err)

			repoName := sourceRepo.GetRepoName()
			repoPath := filepath.Join(targetDir, repoName)

			assert.DirExists(t, filepath.Join(repoPath, "worktrees", tt.defaultBranch))

			config := parseGBMConfig(t, filepath.Join(repoPath, internal.DefaultBranchConfigFilename))
			expected := &internal.GBMConfig{
				Worktrees: map[string]internal.WorktreeConfig{
					tt.defaultBranch: {
						Branch:      tt.defaultBranch,
						Description: "Main production branch",
					},
				},
			}
			assert.Equal(t, expected, config)
		})
	}
}

func TestCloneCommand_DirectoryStructure(t *testing.T) {
	sourceRepo := testutils.NewMultiBranchRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(targetDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := sourceRepo.GetRepoName()
	repoPath := filepath.Join(targetDir, repoName)

	expectedDirs := []string{
		".git",
		".gbm",
		"worktrees",
		"worktrees/main",
	}

	for _, dir := range expectedDirs {
		assert.DirExists(t, filepath.Join(repoPath, dir), "Expected directory %s to exist", dir)
	}

	expectedFiles := []string{
		internal.DefaultBranchConfigFilename,
		".gbm/config.toml",
		".gbm/state.toml",
		"worktrees/main/README.md",
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, filepath.Join(repoPath, file), "Expected file %s to exist", file)
	}
}

func TestCloneCommand_InvalidRepository(t *testing.T) {
	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	_ = os.Chdir(targetDir)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"clone", "/nonexistent/repo"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to clone repository")
}


func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		expected string
	}{
		{"GitHub HTTPS", "https://github.com/user/repo.git", "repo"},
		{"GitHub SSH", "git@github.com:user/repo.git", "repo"},
		{"Without .git", "https://github.com/user/repo", "repo"},
		{"Local path", "/path/to/repo", "repo"},
		{"Local path with .git", "/path/to/repo.git", "repo"},
		{"Empty string", "", "repository"},
		{"Single path", "repo", "repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRepoName(tt.repoURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateDefaultGBMConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, internal.DefaultBranchConfigFilename)

	err := createDefaultGBMConfig(configPath, "main")
	require.NoError(t, err)

	assert.FileExists(t, configPath)

	config := parseGBMConfig(t, configPath)
	expected := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				Description: "Main production branch",
			},
		},
	}
	assert.Equal(t, expected, config)

	// Also verify the file contains expected comments by reading raw content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configContent := string(content)
	assert.Contains(t, configContent, "# Git Branch Manager Configuration")
}
