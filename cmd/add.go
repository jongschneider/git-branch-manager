package cmd

import (
	"errors"
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <worktree-name> [branch-name] [base-branch]",
		Short: "Add a new worktree",
		Long: `Add a new worktree with various options:
- Create on existing branch: gbm add INGSVC-5544 existing-branch-name
- Create on new branch: gbm add INGSVC-5544 feature/new-branch -b
- Create on new branch with base: gbm add INGSVC-5544 feature/new-branch main -b
- Interactive mode: gbm add INGSVC-5544 --interactive
- Tab completion: Shows JIRA keys with summaries, suggests branch names when needed

The third argument specifies which branch/commit to use as the starting point for new branches.
If not specified for new branches, the repository's default branch (main/master) is used.
This matches the behavior of 'git worktree add'.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			newBranch, _ := cmd.Flags().GetBool("new-branch")
			interactive, _ := cmd.Flags().GetBool("interactive")

			worktreeName := args[0]

			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return err
				}

				PrintVerbose("%v", err)
			}

			var branchName string
			var baseBranch string

			if interactive {
				// Handle interactive mode
				result, err := handleInteractiveWithParams(manager)
				if err != nil {
					return fmt.Errorf("failed to handle interactive mode: %w", err)
				}
				branchName = result.branchName
				newBranch = result.newBranch
			} else {
				// Handle direct specification
				if len(args) > 1 {
					branchName = args[1]
					// Check for third argument (base branch)
					if len(args) > 2 {
						baseBranch = args[2]
					}
				} else if newBranch {
					// Generate branch name from worktree name
					branchName = generateBranchName(worktreeName, manager)
				} else if internal.IsJiraKey(worktreeName) {
					// Auto-suggest branch name for JIRA keys
					suggestedBranch := generateBranchName(worktreeName, manager)
					return fmt.Errorf("branch name required. Suggested: %s\n\nTry: gbm add %s %s -b", suggestedBranch, worktreeName, suggestedBranch)
				} else {
					return fmt.Errorf("branch name required when not creating new branch (use -b to create new branch)")
				}
			}

			PrintInfo("Adding worktree '%s' on branch '%s'", worktreeName, branchName)

			// Determine base branch for new branches
			var resolvedBaseBranch string
			if newBranch {
				if baseBranch != "" {
					// Validate that the base branch exists
					exists, err := manager.GetGitManager().BranchExists(baseBranch)
					if err != nil {
						return fmt.Errorf("failed to check if base branch exists: %w", err)
					}
					if !exists {
						return fmt.Errorf("base branch '%s' does not exist", baseBranch)
					}
					resolvedBaseBranch = baseBranch
				} else {
					// Use default branch
					defaultBranch, err := manager.GetGitManager().GetDefaultBranch()
					if err != nil {
						return fmt.Errorf("failed to get default branch: %w", err)
					}
					resolvedBaseBranch = defaultBranch
					PrintInfo("Using default base branch: %s", resolvedBaseBranch)
				}
			}

			// Add the worktree
			if err := manager.AddWorktree(worktreeName, branchName, newBranch, resolvedBaseBranch); err != nil {
				return fmt.Errorf("failed to add worktree: %w", err)
			}

			PrintInfo("Worktree '%s' added successfully", worktreeName)

			return nil
		},
	}

	cmd.Flags().BoolP("new-branch", "b", false, "Create a new branch for the worktree")
	cmd.Flags().BoolP("interactive", "i", false, "Interactive mode to select branch")

	// Add JIRA key completions for the first positional argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// Try to get config for JIRA completion, but don't fail if it's not available
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return nil, cobra.ShellCompDirectiveNoFileComp

				}

				PrintVerbose("%v", err)
			}

			// Complete JIRA keys with summaries for context
			jiraIssues, err := internal.GetJiraIssues(manager)
			if err != nil {
				// If JIRA CLI is not available, fall back to no completions
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completions []string
			for _, issue := range jiraIssues {
				// Format: "KEY\tSummary" - clean completion of just the key with summary context
				completion := fmt.Sprintf("%s\t%s", issue.Key, issue.Summary)
				completions = append(completions, completion)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		} else if len(args) == 1 {
			// Complete branch name based on the JIRA key
			worktreeName := args[0]
			if internal.IsJiraKey(worktreeName) {
				// Try to get config for JIRA completion
				manager, err := createInitializedManager()
				if err != nil {
					// Fallback to default branch name generation
					branchName := fmt.Sprintf("feature/%s", strings.ToLower(worktreeName))
					return []string{branchName}, cobra.ShellCompDirectiveNoFileComp
				}

				branchName, err := internal.GenerateBranchFromJira(worktreeName, manager)
				if err != nil {
					// Fallback to default branch name generation
					branchName = fmt.Sprintf("feature/%s", strings.ToLower(worktreeName))
				}
				return []string{branchName}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

type interactiveResult struct {
	branchName string
	newBranch  bool
}

func handleInteractiveWithParams(manager *internal.Manager) (*interactiveResult, error) {
	// Get available branches
	branches, err := manager.GetRemoteBranches()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	fmt.Println(internal.FormatSubHeader("Available branches:"))
	for i, branch := range branches {
		fmt.Printf("  %s\n", internal.FormatInfo(fmt.Sprintf("%d. %s", i+1, branch)))
	}
	fmt.Printf("  %s\n", internal.FormatInfo(fmt.Sprintf("%d. Create new branch", len(branches)+1)))

	var choice int
	fmt.Print(internal.FormatPrompt("Select a branch: "))
	if _, err := fmt.Scanln(&choice); err != nil {
		return nil, fmt.Errorf("failed to read choice: %w", err)
	}

	if choice < 1 || choice > len(branches)+1 {
		return nil, fmt.Errorf("invalid choice")
	}

	if choice == len(branches)+1 {
		// Create new branch
		fmt.Print(internal.FormatPrompt("Enter new branch name: "))
		var branchName string
		if _, err := fmt.Scanln(&branchName); err != nil {
			return nil, fmt.Errorf("failed to read branch name: %w", err)
		}

		return &interactiveResult{
			branchName: branchName,
			newBranch:  true,
		}, nil
	} else {
		// Use existing branch
		return &interactiveResult{
			branchName: branches[choice-1],
			newBranch:  false,
		}, nil
	}
}

func generateBranchName(worktreeName string, manager *internal.Manager) string {
	// Check if this is a JIRA key first
	if internal.IsJiraKey(worktreeName) {
		branchName, err := internal.GenerateBranchFromJira(worktreeName, manager)
		if err != nil {
			PrintVerbose("Failed to generate branch name from JIRA issue %s: %v", worktreeName, err)
			PrintInfo("Could not fetch JIRA issue details. Using default branch name format.")
			// Generate a basic branch name from the JIRA key
			return fmt.Sprintf("feature/%s", strings.ToLower(worktreeName))
		} else {
			PrintInfo("Auto-generated branch name from JIRA issue: %s", branchName)
			return branchName
		}
	}

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
