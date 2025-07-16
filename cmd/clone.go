package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <repository-url>",
		Short: "Clone a repository as a bare repo and create the main worktree",
		Long: `Clone a repository as a bare repository and create the main worktree
using the HEAD branch. This sets up the repository structure for
worktree-based development.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoUrl := args[0]

			PrintInfo("Cloning repository using git-bare-clone.sh...")
			if err := runGitBareClone(repoUrl); err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}

			PrintInfo("Discovering default branch...")
			defaultBranch, err := getDefaultBranch()
			if err != nil {
				return fmt.Errorf("failed to discover default branch: %w", err)
			}
			PrintInfo("Default branch: %s", defaultBranch)

			PrintInfo("Creating main worktree...")
			if err := createMainWorktree(defaultBranch); err != nil {
				return fmt.Errorf("failed to create main worktree: %w", err)
			}

			PrintInfo("Setting up gbm.branchconfig.yaml configuration...")
			if err := setupGBMConfig(defaultBranch); err != nil {
				return fmt.Errorf("failed to setup gbm.branchconfig.yaml: %w", err)
			}

			PrintInfo("Initializing worktree management...")
			if err := initializeWorktreeManagement(); err != nil {
				return fmt.Errorf("failed to initialize worktree management: %w", err)
			}

			PrintInfo("Repository cloned successfully!")
			return nil
		},
	}

	return cmd
}

func runGitBareClone(repoUrl string) error {
	// Extract repository name from URL
	repo := extractRepoName(repoUrl)

	// Create directory for the repository
	if err := os.MkdirAll(repo, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", repo, err)
	}

	// Change to the repository directory
	if err := os.Chdir(repo); err != nil {
		return fmt.Errorf("failed to change to directory %s: %w", repo, err)
	}

	PrintInfo("Cloning bare repository to .git...")
	// Clone bare repository to .git
	cmd := exec.Command("git", "clone", "--bare", repoUrl, ".git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Clean up the directory if cloning fails
		_ = os.Chdir("..")
		_ = os.RemoveAll(repo)
		return fmt.Errorf("failed to clone bare repository: %w", err)
	}

	PrintInfo("Adjusting origin fetch locations...")
	// Change to .git directory and configure remote
	if err := os.Chdir(".git"); err != nil {
		return fmt.Errorf("failed to change to .git directory: %w", err)
	}

	// Set remote origin fetch configuration
	cmd = exec.Command("git", "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure remote origin fetch: %w", err)
	}

	PrintInfo("Fetching all branches from remote...")
	// Fetch all branches from remote
	cmd = exec.Command("git", "fetch", "origin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch from origin: %w", err)
	}

	// Change back to repository root
	if err := os.Chdir(".."); err != nil {
		return fmt.Errorf("failed to change back to repository root: %w", err)
	}

	return nil
}

func extractRepoName(repoUrl string) string {
	// Remove .git suffix if present
	url := strings.TrimSuffix(repoUrl, ".git")

	// Extract the last part of the URL (repository name)
	parts := strings.Split(url, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}

	return "repository"
}

func getDefaultBranch() (string, error) {
	// First, try to set the remote HEAD reference
	cmd := exec.Command("git", "remote", "set-head", "origin", "-a")
	if err := cmd.Run(); err != nil {
		// If that fails, try to get the remote HEAD manually
		cmd = exec.Command("git", "ls-remote", "--symref", "origin", "HEAD")
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get default branch: %w", err)
		}

		// Parse the output to extract branch name
		// Output format: ref: refs/heads/main	HEAD
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ref: refs/heads/") {
				parts := strings.Split(line, "\t")
				if len(parts) > 0 {
					refPath := parts[0]
					branchName := strings.TrimPrefix(refPath, "ref: refs/heads/")
					return branchName, nil
				}
			}
		}
		return "", fmt.Errorf("could not determine default branch from remote")
	}

	// Now try to get the symbolic ref
	cmd = exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default branch: %w", err)
	}

	// Parse the output to extract branch name
	// Output format: refs/remotes/origin/main
	refPath := strings.TrimSpace(string(output))
	parts := strings.Split(refPath, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("unexpected symbolic-ref output format: %s", refPath)
	}

	return parts[len(parts)-1], nil
}

func createMainWorktree(defaultBranch string) error {
	// Create worktrees directory
	if err := os.MkdirAll("worktrees", 0o755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// Create the main worktree using the default branch name as the worktree name
	worktreeName := defaultBranch
	cmd := exec.Command("git", "worktree", "add", filepath.Join("worktrees", worktreeName), defaultBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create main worktree: %w", err)
	}

	return nil
}

func setupGBMConfig(defaultBranch string) error {
	worktreeConfigPath := filepath.Join("worktrees", defaultBranch, internal.DefaultBranchConfigFilename)
	branchConfigPath := internal.DefaultBranchConfigFilename

	// Check if gbm.branchconfig.yaml exists in the default branch worktree
	if _, err := os.Stat(worktreeConfigPath); err == nil {
		PrintInfo("Found gbm.branchconfig.yaml in %s worktree, copying to root...", defaultBranch)
		if err := copyFile(worktreeConfigPath, branchConfigPath); err != nil {
			return fmt.Errorf("failed to copy gbm.branchconfig.yaml from worktree: %w", err)
		}
	} else if os.IsNotExist(err) {
		// Check if gbm.branchconfig.yaml already exists in root (from repository contents)
		if _, err := os.Stat(branchConfigPath); err == nil {
			PrintInfo("Found gbm.branchconfig.yaml in root, keeping existing configuration...")
			// Don't overwrite existing gbm.branchconfig.yaml from repository
			return nil
		} else if os.IsNotExist(err) {
			PrintInfo("No gbm.branchconfig.yaml found in %s worktree, creating new one...", defaultBranch)
			if err := createDefaultGBMConfig(branchConfigPath, defaultBranch); err != nil {
				return fmt.Errorf("failed to create default gbm.branchconfig.yaml: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check gbm.branchconfig.yaml in root: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check gbm.branchconfig.yaml in worktree: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func createDefaultGBMConfig(path, defaultBranch string) error {
	content := fmt.Sprintf(`# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  # Primary worktree - no merge_into (root of merge chain)
  %s:
    branch: %s
    description: "Main production branch"
`, defaultBranch, defaultBranch)

	return os.WriteFile(path, []byte(content), 0o644)
}

func initializeWorktreeManagement() error {
	// Get current working directory (repository root)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create manager
	manager, err := internal.NewManager(wd)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	// Save default config and state to create .gbm directory
	if err := manager.SaveConfig(); err != nil {
		return fmt.Errorf("failed to initialize .gbm/config.toml: %w", err)
	}

	if err := manager.SaveState(); err != nil {
		return fmt.Errorf("failed to initialize .gbm/state.toml: %w", err)
	}

	// Load gbm.branchconfig.yaml configuration
	branchConfigPath := internal.DefaultBranchConfigFilename

	if err := manager.LoadGBMConfig(branchConfigPath); err != nil {
		return fmt.Errorf("failed to load %s: %w", internal.DefaultBranchConfigFilename, err)
	}

	// Initialize worktree management - create worktrees for each branch config mapping
	// Use a more permissive sync that doesn't fail on invalid branches during clone
	if err := manager.Sync(false, false); err != nil {
		// For clone operations, we want to be more permissive and not fail
		// if there are invalid branch references in the branch config file
		PrintInfo("Warning: some branch references in %s may be invalid: %v", internal.DefaultBranchConfigFilename, err)
		PrintInfo("You can run 'gbm sync' later to resolve any issues")
	}

	return nil
}

func init() {
}
