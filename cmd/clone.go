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

var cloneCmd = &cobra.Command{
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

		PrintInfo("Setting up .envrc configuration...")
		if err := setupEnvrc(defaultBranch); err != nil {
			return fmt.Errorf("failed to setup .envrc: %w", err)
		}

		PrintInfo("Initializing worktree management...")
		if err := initializeWorktreeManagement(); err != nil {
			return fmt.Errorf("failed to initialize worktree management: %w", err)
		}

		PrintInfo("Repository cloned successfully!")
		return nil
	},
}

func runGitBareClone(repoUrl string) error {
	// Extract repository name from URL
	repo := extractRepoName(repoUrl)

	// Create directory for the repository
	if err := os.MkdirAll(repo, 0755); err != nil {
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
		os.Chdir("..")
		os.RemoveAll(repo)
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
	if err := os.MkdirAll("worktrees", 0755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// Create the main worktree
	cmd := exec.Command("git", "worktree", "add", "worktrees/MAIN", defaultBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create main worktree: %w", err)
	}

	return nil
}

func setupEnvrc(defaultBranch string) error {
	worktreeEnvrcPath := filepath.Join("worktrees", "MAIN", ".envrc")
	rootEnvrcPath := ".envrc"

	// Check if .envrc exists in the MAIN worktree
	if _, err := os.Stat(worktreeEnvrcPath); err == nil {
		PrintInfo("Found .envrc in MAIN worktree, copying to root...")
		if err := copyFile(worktreeEnvrcPath, rootEnvrcPath); err != nil {
			return fmt.Errorf("failed to copy .envrc from worktree: %w", err)
		}
	} else if os.IsNotExist(err) {
		PrintInfo("No .envrc found in MAIN worktree, creating new one...")
		if err := createDefaultEnvrc(rootEnvrcPath, defaultBranch); err != nil {
			return fmt.Errorf("failed to create default .envrc: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check .envrc in worktree: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func createDefaultEnvrc(path, defaultBranch string) error {
	content := fmt.Sprintf("# Git Branch Manager configuration\n# This file defines the mapping between environment variables and branches\n\nMAIN=%s\n", defaultBranch)

	return os.WriteFile(path, []byte(content), 0644)
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

	// Load .envrc configuration
	if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
		return fmt.Errorf("failed to load .envrc: %w", err)
	}

	// Initialize worktree management - create worktrees for each .envrc mapping
	if err := manager.Sync(false, false, false); err != nil {
		return fmt.Errorf("failed to sync worktrees: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
