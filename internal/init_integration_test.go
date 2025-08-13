package internal

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for init functionality using real git operations
// These tests verify that initialized repositories work correctly with Manager operations

func TestInitIntegration_BasicWorkflow(t *testing.T) {
	// Create a temporary directory for our init test
	tempDir := t.TempDir()
	initDir := filepath.Join(tempDir, "test-init-repo")

	tests := []struct {
		name         string
		branchName   string
		setup        func(t *testing.T) string
		expectErr    func(t *testing.T, err error)
		verifyResult func(t *testing.T, repoPath string)
	}{
		{
			name:       "init with default main branch",
			branchName: "main",
			setup: func(t *testing.T) string {
				return initDir + "-main"
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			verifyResult: func(t *testing.T, repoPath string) {
				// Verify basic repository structure
				assert.DirExists(t, filepath.Join(repoPath, ".git"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees"))
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "main"))
				assert.DirExists(t, filepath.Join(repoPath, ".gbm"))
				assert.FileExists(t, filepath.Join(repoPath, "gbm.branchconfig.yaml"))

				// Verify Manager can be created and works
				manager, err := NewManager(repoPath)
				require.NoError(t, err)

				// Verify branch exists
				exists, err := manager.BranchExists("main")
				require.NoError(t, err)
				assert.True(t, exists)

				// Verify initial commit exists and main worktree is set up
				gitManager, err := NewGitManager(repoPath, "worktrees")
				require.NoError(t, err)
				worktrees, err := gitManager.GetWorktrees()
				require.NoError(t, err)

				// Find the main worktree (there might be additional worktrees)
				var mainWorktree *WorktreeInfo
				for _, wt := range worktrees {
					if wt.Branch == "main" {
						mainWorktree = wt
						break
					}
				}
				require.NotNil(t, mainWorktree, "Main worktree should exist")
				assert.Equal(t, "main", mainWorktree.Branch)

				// Verify gbm.branchconfig.yaml was committed
				mainWorktreePath := filepath.Join(repoPath, "worktrees", "main")
				assert.FileExists(t, filepath.Join(mainWorktreePath, "gbm.branchconfig.yaml"))
			},
		},
		{
			name:       "init with custom develop branch",
			branchName: "develop",
			setup: func(t *testing.T) string {
				return initDir + "-develop"
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			verifyResult: func(t *testing.T, repoPath string) {
				// Verify develop branch structure
				assert.DirExists(t, filepath.Join(repoPath, "worktrees", "develop"))

				// Verify Manager can work with custom branch
				manager, err := NewManager(repoPath)
				require.NoError(t, err)

				// Verify custom branch exists
				exists, err := manager.BranchExists("develop")
				require.NoError(t, err)
				assert.True(t, exists)

				// Verify config reflects custom branch
				err = manager.LoadGBMConfig("")
				require.NoError(t, err)
				config := manager.GetGBMConfig()
				require.NotNil(t, config)
				assert.Contains(t, config.Worktrees, "develop")
				assert.Equal(t, "develop", config.Worktrees["develop"].Branch)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setup(t)

			// Simulate the init process that cmd/init.go does
			err := simulateInitProcess(t, repoPath, tt.branchName)
			tt.expectErr(t, err)

			if err == nil {
				tt.verifyResult(t, repoPath)
			}
		})
	}
}

func TestInitIntegration_ErrorCases(t *testing.T) {
	// Test error cases that would occur in real usage
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "fail when directory already has git repo",
			setup: func(t *testing.T) string {
				// Create a directory with an existing git repo
				existingRepoPath := filepath.Join(tempDir, "existing-repo")
				repo := testutils.NewBasicRepo(t)

				// Copy the repo to our test location
				err := copyDirectory(repo.GetLocalPath(), existingRepoPath)
				require.NoError(t, err)

				return existingRepoPath
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// Should fail because directory already contains git repo
			},
		},
		{
			name: "fail when path is a file",
			setup: func(t *testing.T) string {
				filePath := filepath.Join(tempDir, "test-file")
				err := os.WriteFile(filePath, []byte("test"), 0o644)
				require.NoError(t, err)
				return filePath
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				// Should fail because path is a file, not a directory
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetPath := tt.setup(t)
			err := simulateInitProcess(t, targetPath, "main")
			tt.expectErr(t, err)
		})
	}
}

// simulateInitProcess replicates what cmd/init.go does to initialize a repository
func simulateInitProcess(t *testing.T, targetDir, branchName string) error {
	t.Helper()

	// Phase 1: Directory validation and creation
	if err := simulateValidateInitDirectory(targetDir); err != nil {
		return err
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	// Phase 2: Initialize bare git repository
	gitDir := filepath.Join(targetDir, ".git")
	if err := ExecGitCommandSilent(targetDir, "init", "--bare", gitDir); err != nil {
		return err
	}

	if err := ExecGitCommandSilent(gitDir, "config", "core.bare", "false"); err != nil {
		return err
	}

	// Phase 3: Create Manager and use it to set up gbm structure
	manager, err := NewManager(targetDir)
	if err != nil {
		return err
	}

	// Create main worktree
	if err := manager.AddWorktree(branchName, branchName, true, ""); err != nil {
		return err
	}

	// Create gbm.branchconfig.yaml
	configPath := filepath.Join(targetDir, DefaultBranchConfigFilename)
	content := `# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  # Primary worktree - no merge_into (root of merge chain)
  ` + branchName + `:
    branch: ` + branchName + `
    description: "Main production branch"
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		return err
	}

	// Save manager state
	if err := manager.SaveConfig(); err != nil {
		return err
	}

	if err := manager.SaveState(); err != nil {
		return err
	}

	// Create initial commit in the worktree
	worktreePath := filepath.Join(targetDir, "worktrees", branchName)

	// Copy gbm.branchconfig.yaml to worktree
	sourceConfig := configPath
	targetConfig := filepath.Join(worktreePath, DefaultBranchConfigFilename)
	configContent, err := os.ReadFile(sourceConfig)
	if err != nil {
		return err
	}
	if err := os.WriteFile(targetConfig, configContent, 0o644); err != nil {
		return err
	}

	// Commit the config file
	if err := ExecGitCommandSilent(worktreePath, "add", DefaultBranchConfigFilename); err != nil {
		return err
	}

	if err := ExecGitCommandSilent(worktreePath, "commit", "-m", "Initial commit with gbm configuration"); err != nil {
		return err
	}

	return nil
}

// simulateValidateInitDirectory replicates the validation logic from cmd/init.go
func simulateValidateInitDirectory(path string) error {
	// Check if path exists
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil // Directory doesn't exist - will be created
	}
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("path exists but is not a directory")
	}

	// Check if it already contains a git repository
	gitPath := filepath.Join(path, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return errors.New("directory already contains a git repository")
	}

	return nil
}

// copyDirectory is a helper to copy directory contents for testing
func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() { _ = dstFile.Close() }()

		_, err = srcFile.WriteTo(dstFile)
		return err
	})
}
