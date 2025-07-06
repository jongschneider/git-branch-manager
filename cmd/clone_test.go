package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneCommand_Basic(t *testing.T) {
	sourceRepo := testutils.NewMultiBranchRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)

	assert.DirExists(t, repoPath)
	assert.DirExists(t, filepath.Join(repoPath, ".git"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees"))
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
	assert.FileExists(t, filepath.Join(repoPath, ".envrc"))

	content, err := os.ReadFile(filepath.Join(repoPath, ".envrc"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "MAIN=main")
}

func TestCloneCommand_WithExistingEnvrc(t *testing.T) {
	sourceRepo := testutils.NewEnvrcRepo(t, map[string]string{
		"MAIN": "main",
		"DEV":  "develop",
		"FEAT": "feature/auth",
	})

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)

	assert.FileExists(t, filepath.Join(repoPath, ".envrc"))

	content, err := os.ReadFile(filepath.Join(repoPath, ".envrc"))
	require.NoError(t, err)
	envrcContent := string(content)
	assert.Contains(t, envrcContent, "MAIN=main")
	assert.Contains(t, envrcContent, "DEV=develop")
	assert.Contains(t, envrcContent, "FEAT=feature/auth")
}

func TestCloneCommand_WithoutEnvrc(t *testing.T) {
	sourceRepo := testutils.NewBasicRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)

	assert.FileExists(t, filepath.Join(repoPath, ".envrc"))

	content, err := os.ReadFile(filepath.Join(repoPath, ".envrc"))
	require.NoError(t, err)
	envrcContent := string(content)
	assert.Contains(t, envrcContent, "MAIN=main")
	assert.Contains(t, envrcContent, "# Git Branch Manager configuration")
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
			defer os.Chdir(originalDir)

			os.Chdir(targetDir)

			cmd := rootCmd
			cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

			err := cmd.Execute()
			require.NoError(t, err)

			repoName := extractRepoName(sourceRepo.GetRemotePath())
			repoPath := filepath.Join(targetDir, repoName)

			assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))

			content, err := os.ReadFile(filepath.Join(repoPath, ".envrc"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "MAIN="+tt.defaultBranch)
		})
	}
}

func TestCloneCommand_DirectoryStructure(t *testing.T) {
	sourceRepo := testutils.NewMultiBranchRepo(t)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)

	expectedDirs := []string{
		".git",
		"worktrees",
		"worktrees/MAIN",
	}

	for _, dir := range expectedDirs {
		assert.DirExists(t, filepath.Join(repoPath, dir), "Expected directory %s to exist", dir)
	}

	expectedFiles := []string{
		".envrc",
		"worktrees/MAIN/README.md",
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, filepath.Join(repoPath, file), "Expected file %s to exist", file)
	}
}

func TestCloneCommand_InvalidRepository(t *testing.T) {
	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", "/nonexistent/repo"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to clone repository")
}

func TestCloneCommand_EmptyRepository(t *testing.T) {
	sourceRepo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Test User", "test@example.com"),
	)

	targetDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(targetDir)

	cmd := rootCmd
	cmd.SetArgs([]string{"clone", sourceRepo.GetRemotePath()})

	err := cmd.Execute()
	require.NoError(t, err)

	repoName := extractRepoName(sourceRepo.GetRemotePath())
	repoPath := filepath.Join(targetDir, repoName)

	assert.DirExists(t, repoPath)
	assert.DirExists(t, filepath.Join(repoPath, "worktrees", "MAIN"))
	assert.FileExists(t, filepath.Join(repoPath, ".envrc"))
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

func TestCreateDefaultEnvrc(t *testing.T) {
	tempDir := t.TempDir()
	envrcPath := filepath.Join(tempDir, ".envrc")

	err := createDefaultEnvrc(envrcPath, "main")
	require.NoError(t, err)

	assert.FileExists(t, envrcPath)

	content, err := os.ReadFile(envrcPath)
	require.NoError(t, err)

	envrcContent := string(content)
	assert.Contains(t, envrcContent, "MAIN=main")
	assert.Contains(t, envrcContent, "# Git Branch Manager configuration")
	assert.Contains(t, envrcContent, "# This file defines the mapping between environment variables and branches")
}
