package cmd

import (
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var (
	newBranch   bool
	baseBranch  string
	interactive bool
)

var addCmd = &cobra.Command{
	Use:   "add <worktree-name> [branch-name]",
	Short: "Add a new worktree",
	Long: `Add a new worktree with various options:
- Create on existing branch
- Create on new branch (--new-branch)
- Interactive mode (--interactive)`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		worktreeName := args[0]

		// Find git repository root
		repoPath, err := internal.FindGitRoot(".")
		if err != nil {
			return fmt.Errorf("not in a git repository: %w", err)
		}

		// Create manager
		manager, err := internal.NewManager(repoPath)
		if err != nil {
			return fmt.Errorf("failed to create manager: %w", err)
		}

		var branchName string

		if interactive {
			// Handle interactive mode
			branchName, err = handleInteractive(manager)
			if err != nil {
				return fmt.Errorf("failed to handle interactive mode: %w", err)
			}
		} else {
			// Handle direct specification
			if len(args) > 1 {
				branchName = args[1]
			} else if newBranch {
				// Generate branch name from worktree name
				branchName = generateBranchName(worktreeName)
			} else {
				return fmt.Errorf("branch name required when not creating new branch")
			}
		}

		PrintInfo("Adding worktree '%s' on branch '%s'", worktreeName, branchName)

		// Add the worktree
		if err := manager.AddWorktree(worktreeName, branchName, newBranch); err != nil {
			return fmt.Errorf("failed to add worktree: %w", err)
		}

		PrintInfo("Worktree '%s' added successfully", worktreeName)

		return nil
	},
}

func handleInteractive(manager *internal.Manager) (string, error) {
	// Get available branches
	branches, err := manager.GetRemoteBranches()
	if err != nil {
		return "", fmt.Errorf("failed to get remote branches: %w", err)
	}

	fmt.Println(internal.FormatSubHeader("Available branches:"))
	for i, branch := range branches {
		fmt.Printf("  %s\n", internal.FormatInfo(fmt.Sprintf("%d. %s", i+1, branch)))
	}
	fmt.Printf("  %s\n", internal.FormatInfo(fmt.Sprintf("%d. Create new branch", len(branches)+1)))

	var choice int
	fmt.Print(internal.FormatPrompt("Select a branch: "))
	if _, err := fmt.Scanln(&choice); err != nil {
		return "", fmt.Errorf("failed to read choice: %w", err)
	}

	if choice < 1 || choice > len(branches)+1 {
		return "", fmt.Errorf("invalid choice")
	}

	if choice == len(branches)+1 {
		// Create new branch
		newBranch = true
		fmt.Print(internal.FormatPrompt("Enter new branch name: "))
		var branchName string
		if _, err := fmt.Scanln(&branchName); err != nil {
			return "", fmt.Errorf("failed to read branch name: %w", err)
		}
		return branchName, nil
	} else {
		// Use existing branch
		newBranch = false
		return branches[choice-1], nil
	}
}

func generateBranchName(worktreeName string) string {
	// Convert worktree name to a valid branch name
	branchName := strings.ReplaceAll(worktreeName, " ", "-")
	branchName = strings.ReplaceAll(branchName, "_", "-")
	branchName = strings.ToLower(branchName)

	// Add feature/ prefix if not already present
	if !strings.HasPrefix(branchName, "feature/") && !strings.HasPrefix(branchName, "bugfix/") && !strings.HasPrefix(branchName, "hotfix/") {
		branchName = "feature/" + branchName
	}

	return branchName
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&newBranch, "new-branch", "b", false, "Create a new branch for the worktree")
	addCmd.Flags().StringVar(&baseBranch, "base", "", "Base branch for new branch (default: current branch)")
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to select branch")
}
