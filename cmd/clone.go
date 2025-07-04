package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

		// Extract repository name from URL
		repoName := extractRepoName(repoUrl)
		bareRepoPath := repoName + ".git"

		PrintInfo("Cloning repository as bare...")
		if err := cloneBareRepo(repoUrl, bareRepoPath); err != nil {
			return fmt.Errorf("failed to clone bare repository: %w", err)
		}

		PrintInfo("Creating main worktree...")
		mainWorktreePath := filepath.Join(bareRepoPath, "worktrees", "MAIN")
		if err := createMainWorktree(bareRepoPath, mainWorktreePath); err != nil {
			return fmt.Errorf("failed to create main worktree: %w", err)
		}

		PrintInfo("Checking for .envrc file...")
		envrcPath := filepath.Join(mainWorktreePath, ".envrc")
		if _, err := os.Stat(envrcPath); err == nil {
			PrintInfo("Found .envrc file in main worktree")
			// Copy .envrc to repository root for reference
			repoEnvrcPath := filepath.Join(bareRepoPath, ".envrc")
			if err := copyFile(envrcPath, repoEnvrcPath); err != nil {
				PrintError("Failed to copy .envrc to repository root: %v", err)
			} else {
				PrintInfo("Copied .envrc to repository root")
			}
		} else {
			PrintInfo("No .envrc file found in main worktree")
			PrintInfo("Consider creating a .envrc file to define environment variables for your worktrees")
			PrintInfo("You can generate one based on the initial worktree structure")
		}

		PrintInfo("Repository cloned successfully!")
		PrintInfo("Bare repository: %s", bareRepoPath)
		PrintInfo("Main worktree: %s", mainWorktreePath)

		return nil
	},
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

func cloneBareRepo(repoUrl, bareRepoPath string) error {
	cmd := exec.Command("git", "clone", "--bare", repoUrl, bareRepoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createMainWorktree(bareRepoPath, mainWorktreePath string) error {
	// Create worktrees directory
	worktreesDir := filepath.Dir(mainWorktreePath)
	if err := os.MkdirAll(worktreesDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// Get the default branch (HEAD)
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = bareRepoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback to main/master detection
		cmd = exec.Command("git", "branch", "-r")
		cmd.Dir = bareRepoPath
		output, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get remote branches: %w", err)
		}

		branches := strings.Split(string(output), "\n")
		var defaultBranch string
		for _, branch := range branches {
			branch = strings.TrimSpace(branch)
			if strings.Contains(branch, "origin/main") {
				defaultBranch = "main"
				break
			} else if strings.Contains(branch, "origin/master") {
				defaultBranch = "master"
				break
			}
		}

		if defaultBranch == "" {
			return fmt.Errorf("could not determine default branch")
		}

		// Create worktree with the detected default branch
		cmd = exec.Command("git", "worktree", "add", mainWorktreePath, defaultBranch)
	} else {
		// Extract branch name from refs/remotes/origin/HEAD -> refs/remotes/origin/main
		refPath := strings.TrimSpace(string(output))
		branchName := strings.TrimPrefix(refPath, "refs/remotes/origin/")

		// Create worktree with the HEAD branch
		cmd = exec.Command("git", "worktree", "add", mainWorktreePath, branchName)
	}

	cmd.Dir = bareRepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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

	_, err = srcFile.WriteTo(dstFile)
	return err
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
