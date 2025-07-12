package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFilesToWorktree_AdHocOnly(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "gbm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with test configuration (no gbmConfig means all worktrees are ad-hoc)
	manager := &Manager{
		repoPath: tmpDir,
		config: &Config{
			Settings: ConfigSettings{
				WorktreePrefix: DefaultWorktreeDirname,
			},
			FileCopy: ConfigFileCopy{
				Rules: []FileCopyRule{
					{
						SourceWorktree: "master",
						Files:          []string{".env", "config/"},
					},
				},
			},
		},
		gbmConfig: nil, // No config means all worktrees are ad-hoc
	}

	// Create source worktree directory and files
	sourceWorktreePath := filepath.Join(tmpDir, "worktrees", "master")
	require.NoError(t, os.MkdirAll(sourceWorktreePath, 0755))

	// Create test files in source worktree
	envContent := "DATABASE_URL=postgres://localhost/test\nAPI_KEY=secret123"
	require.NoError(t, os.WriteFile(filepath.Join(sourceWorktreePath, ".env"), []byte(envContent), 0644))

	configDir := filepath.Join(sourceWorktreePath, "config")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	configContent := `{"env": "development", "debug": true}`
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "development.json"), []byte(configContent), 0644))

	// Create target worktree directory
	targetWorktreePath := filepath.Join(tmpDir, "worktrees", "feature-branch")
	require.NoError(t, os.MkdirAll(targetWorktreePath, 0755))

	// Test file copying
	err = manager.copyFilesToWorktree("feature-branch")
	require.NoError(t, err)

	// Verify .env file was copied
	copiedEnvPath := filepath.Join(targetWorktreePath, ".env")
	assert.FileExists(t, copiedEnvPath)

	copiedEnvContent, err := os.ReadFile(copiedEnvPath)
	require.NoError(t, err)
	assert.Equal(t, envContent, string(copiedEnvContent))

	// Verify config directory was copied
	copiedConfigDir := filepath.Join(targetWorktreePath, "config")
	assert.DirExists(t, copiedConfigDir)

	copiedConfigFile := filepath.Join(copiedConfigDir, "development.json")
	assert.FileExists(t, copiedConfigFile)

	copiedConfigContent, err := os.ReadFile(copiedConfigFile)
	require.NoError(t, err)
	assert.Equal(t, configContent, string(copiedConfigContent))
}

func TestCopyFilesToWorktree_NoRules(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	manager := &Manager{
		repoPath: tmpDir,
		config: &Config{
			Settings: ConfigSettings{
				WorktreePrefix: DefaultWorktreeDirname,
			},
			FileCopy: ConfigFileCopy{
				Rules: []FileCopyRule{},
			},
		},
	}

	// Should not fail when no rules are configured
	err = manager.copyFilesToWorktree("test-worktree")
	assert.NoError(t, err)
}

func TestCopyFilesToWorktree_SourceNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	manager := &Manager{
		repoPath: tmpDir,
		config: &Config{
			Settings: ConfigSettings{
				WorktreePrefix: DefaultWorktreeDirname,
			},
			FileCopy: ConfigFileCopy{
				Rules: []FileCopyRule{
					{
						SourceWorktree: "nonexistent",
						Files:          []string{".env"},
					},
				},
			},
		},
	}

	// Should not fail when source worktree doesn't exist
	err = manager.copyFilesToWorktree("test-worktree")
	assert.NoError(t, err)
}

func TestAddWorktree_TrackedWorktreeNoFileCopy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create git repo structure
	repoPath := tmpDir
	require.NoError(t, os.MkdirAll(filepath.Join(repoPath, ".git"), 0755))

	// Create manager with gbmConfig that tracks a worktree
	manager := &Manager{
		repoPath: repoPath,
		config: &Config{
			Settings: ConfigSettings{
				WorktreePrefix: DefaultWorktreeDirname,
			},
			FileCopy: ConfigFileCopy{
				Rules: []FileCopyRule{
					{
						SourceWorktree: "master",
						Files:          []string{".env"},
					},
				},
			},
		},
		gbmConfig: &GBMConfig{
			Worktrees: map[string]WorktreeConfig{
				"tracked-worktree": {
					Branch: "feature/tracked",
				},
			},
		},
	}

	// Create source worktree with .env file
	sourceWorktreePath := filepath.Join(repoPath, "worktrees", "master")
	require.NoError(t, os.MkdirAll(sourceWorktreePath, 0755))
	envContent := "DATABASE_URL=postgres://localhost/test"
	require.NoError(t, os.WriteFile(filepath.Join(sourceWorktreePath, ".env"), []byte(envContent), 0644))

	// Create target worktree directory (simulating git worktree add)
	targetWorktreePath := filepath.Join(repoPath, "worktrees", "tracked-worktree")
	require.NoError(t, os.MkdirAll(targetWorktreePath, 0755))

	// Simulate what AddWorktree does for file copying logic
	isAdHoc := true
	if manager.gbmConfig != nil {
		if _, exists := manager.gbmConfig.Worktrees["tracked-worktree"]; exists {
			isAdHoc = false
		}
	}

	// Only copy files for ad-hoc worktrees
	if isAdHoc {
		err := manager.copyFilesToWorktree("tracked-worktree")
		require.NoError(t, err)
	}

	// Verify .env file was NOT copied (since this is a tracked worktree)
	copiedEnvPath := filepath.Join(targetWorktreePath, ".env")
	assert.NoFileExists(t, copiedEnvPath)
}

func TestAddWorktree_AdHocWorktreeFileCopy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create git repo structure
	repoPath := tmpDir
	require.NoError(t, os.MkdirAll(filepath.Join(repoPath, ".git"), 0755))

	// Create manager with gbmConfig
	manager := &Manager{
		repoPath: repoPath,
		config: &Config{
			Settings: ConfigSettings{
				WorktreePrefix: DefaultWorktreeDirname,
			},
			FileCopy: ConfigFileCopy{
				Rules: []FileCopyRule{
					{
						SourceWorktree: "master",
						Files:          []string{".env"},
					},
				},
			},
		},
		gbmConfig: &GBMConfig{
			Worktrees: map[string]WorktreeConfig{
				"tracked-worktree": {
					Branch: "feature/tracked",
				},
			},
		},
	}

	// Create source worktree with .env file
	sourceWorktreePath := filepath.Join(repoPath, "worktrees", "master")
	require.NoError(t, os.MkdirAll(sourceWorktreePath, 0755))
	envContent := "DATABASE_URL=postgres://localhost/test"
	require.NoError(t, os.WriteFile(filepath.Join(sourceWorktreePath, ".env"), []byte(envContent), 0644))

	// Create target worktree directory (simulating git worktree add)
	targetWorktreePath := filepath.Join(repoPath, "worktrees", "adhoc-worktree")
	require.NoError(t, os.MkdirAll(targetWorktreePath, 0755))

	// Simulate what AddWorktree does for file copying logic
	isAdHoc := true
	if manager.gbmConfig != nil {
		if _, exists := manager.gbmConfig.Worktrees["adhoc-worktree"]; exists {
			isAdHoc = false
		}
	}

	// Only copy files for ad-hoc worktrees
	if isAdHoc {
		err := manager.copyFilesToWorktree("adhoc-worktree")
		require.NoError(t, err)
	}

	// Verify .env file WAS copied (since this is an ad-hoc worktree)
	copiedEnvPath := filepath.Join(targetWorktreePath, ".env")
	assert.FileExists(t, copiedEnvPath)

	copiedEnvContent, err := os.ReadFile(copiedEnvPath)
	require.NoError(t, err)
	assert.Equal(t, envContent, string(copiedEnvContent))
}