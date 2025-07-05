package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "repository"
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
