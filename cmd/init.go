package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_repositoryInitializer.go . repositoryInitializer

// repositoryInitializer interface abstracts the Manager operations needed for initializing repositories
type repositoryInitializer interface {
	AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error
	SaveConfig() error
	SaveState() error
	GetRepoPath() string
	GetConfig() *internal.Config
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [directory] [--branch=<branch-name>]",
		Short: "Initialize a new git repository with gbm structure",
		Long: `Initialize a new git repository as a bare repo with gbm worktree management.

This command creates:
- A bare git repository (.git directory)
- Worktree directory structure
- Main worktree for the default branch
- gbm.branchconfig.yaml configuration file
- .gbm directory with configuration files
- Initial commit with gbm setup

Examples:
  gbm init                    # Initialize in current directory with default branch
  gbm init my-project         # Initialize in 'my-project' directory  
  gbm init --branch=develop   # Initialize with 'develop' as default branch
  gbm init my-project --branch=main  # Initialize 'my-project' with 'main' branch`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchFlag, _ := cmd.Flags().GetString("branch")

			// Resolve target directory
			targetDir, err := resolveTargetDirectory(args)
			if err != nil {
				return err
			}

			// Handle the initialization
			return handleInit(targetDir, branchFlag)
		},
	}

	cmd.Flags().String("branch", "", "Override default branch name (defaults to git's init.defaultBranch setting)")
	return cmd
}

// resolveTargetDirectory determines the target directory for initialization
func resolveTargetDirectory(args []string) (string, error) {
	if len(args) == 0 {
		// Use current directory
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return wd, nil
	}

	// Use specified directory (convert to absolute path)
	targetDir := args[0]
	if !filepath.IsAbs(targetDir) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		targetDir = filepath.Join(wd, targetDir)
	}

	return targetDir, nil
}

func handleInit(targetDir, branchFlag string) error {
	PrintInfo("Initializing gbm repository in %s", targetDir)

	// Phase 1: Pre-repository validation and setup
	PrintInfo("Validating directory...")
	if err := validateInitDirectory(targetDir); err != nil {
		return err
	}

	PrintInfo("Creating directory structure...")
	if err := createInitDirectory(targetDir); err != nil {
		return err
	}

	PrintInfo("Initializing bare git repository...")
	if err := initializeBareRepository(targetDir); err != nil {
		return err
	}

	// Determine branch name
	branchName := branchFlag
	if branchName == "" {
		PrintInfo("Detecting default branch name...")
		var err error
		branchName, err = getNativeDefaultBranch()
		if err != nil {
			return fmt.Errorf("failed to determine default branch name: %w", err)
		}
	}
	PrintInfo("Using branch: %s", branchName)

	// Phase 2: Manager-based setup (after repository exists)
	PrintInfo("Creating manager...")
	manager, err := internal.NewManager(targetDir)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	PrintInfo("Setting up worktree structure...")
	if err := setupWorktreeStructure(manager, branchName); err != nil {
		return fmt.Errorf("failed to setup worktree structure: %w", err)
	}

	PrintInfo("Creating gbm configuration...")
	if err := createGBMConfig(manager, branchName); err != nil {
		return fmt.Errorf("failed to create gbm configuration: %w", err)
	}

	PrintInfo("Initializing gbm state...")
	if err := initializeGBMState(manager); err != nil {
		return fmt.Errorf("failed to initialize gbm state: %w", err)
	}

	PrintInfo("Creating initial commit...")
	if err := createInitialCommit(manager, branchName); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	PrintInfo("Repository initialized successfully!")
	PrintInfo("Main worktree available at: %s/worktrees/%s", targetDir, branchName)
	return nil
}

// Phase 1: Pre-repository helper functions
func validateInitDirectory(path string) error {
	// Check if we're currently in a git repository
	if isGitRepository(".") {
		return fmt.Errorf("current directory is already in a git repository")
	}

	// Check if path exists
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil // Directory doesn't exist - will be created
	}
	if err != nil {
		return fmt.Errorf("failed to check directory: %w", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	// Directory exists - check if it already contains a git repository
	if isGitRepository(path) {
		return fmt.Errorf("directory already contains a git repository: %s", path)
	}

	return nil
}

func createInitDirectory(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func initializeBareRepository(path string) error {
	gitDir := filepath.Join(path, ".git")

	if err := internal.ExecGitCommandSilent(path, "init", "--bare", gitDir); err != nil {
		return fmt.Errorf("failed to initialize bare repository: %w", err)
	}

	if err := internal.ExecGitCommandSilent(gitDir, "config", "core.bare", "false"); err != nil {
		return fmt.Errorf("failed to configure repository: %w", err)
	}

	return nil
}

func getNativeDefaultBranch() (string, error) {
	output, err := internal.ExecGitCommand("", "config", "--get", "init.defaultBranch")
	if err == nil && len(output) > 0 {
		branchName := strings.TrimSpace(string(output))
		if branchName != "" {
			return branchName, nil
		}
	}
	return "main", nil
}

func isGitRepository(path string) bool {
	gitPath := filepath.Join(path, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return true
	}

	_, err := internal.ExecGitCommand(path, "rev-parse", "--git-dir")
	return err == nil
}

// Phase 2: Manager-based helper functions
func setupWorktreeStructure(initializer repositoryInitializer, branchName string) error {
	if err := initializer.AddWorktree(branchName, branchName, true, ""); err != nil {
		return fmt.Errorf("failed to create main worktree: %w", err)
	}
	return nil
}

func createGBMConfig(initializer repositoryInitializer, branchName string) error {
	configPath := filepath.Join(initializer.GetRepoPath(), internal.DefaultBranchConfigFilename)

	content := fmt.Sprintf(`# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  # Primary worktree - no merge_into (root of merge chain)
  %s:
    branch: %s
    description: "Main production branch"
`, branchName, branchName)

	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to create %s: %w", internal.DefaultBranchConfigFilename, err)
	}

	return nil
}

func initializeGBMState(initializer repositoryInitializer) error {
	if err := initializer.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save initial config: %w", err)
	}

	if err := initializer.SaveState(); err != nil {
		return fmt.Errorf("failed to save initial state: %w", err)
	}

	return nil
}

func createInitialCommit(initializer repositoryInitializer, branchName string) error {
	worktreePath := filepath.Join(initializer.GetRepoPath(), initializer.GetConfig().Settings.WorktreePrefix, branchName)

	// Copy gbm.branchconfig.yaml from repository root to worktree
	configFile := internal.DefaultBranchConfigFilename
	sourceConfig := filepath.Join(initializer.GetRepoPath(), configFile)
	targetConfig := filepath.Join(worktreePath, configFile)

	// Read the config from the repository root
	configContent, err := os.ReadFile(sourceConfig)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configFile, err)
	}

	// Write it to the worktree
	if err := os.WriteFile(targetConfig, configContent, 0o644); err != nil {
		return fmt.Errorf("failed to copy %s to worktree: %w", configFile, err)
	}

	// Add gbm.branchconfig.yaml to the commit
	if err := internal.ExecGitCommandSilent(worktreePath, "add", configFile); err != nil {
		return fmt.Errorf("failed to add %s to git: %w", configFile, err)
	}

	// Create initial commit
	commitMessage := "Initial commit with gbm configuration"
	if err := internal.ExecGitCommandSilent(worktreePath, "commit", "-m", commitMessage); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}
