package cmd

import (
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newAddCommand(manager *internal.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <worktree-name> [branch-name] [base-branch]",
		Short: "Add a new worktree",
		Long: `Add a new worktree with various options:
- Create on existing branch: gbm add INGSVC-5544 existing-branch-name
- Create on new branch: gbm add INGSVC-5544 feature/new-branch -b
- Create on new branch with base: gbm add INGSVC-5544 feature/new-branch main -b
- Tab completion: Shows JIRA keys with summaries, suggests branch names when needed

The third argument specifies which branch/commit to use as the starting point for new branches.
If not specified for new branches, the repository's default branch (main/master) is used.
This matches the behavior of 'git worktree add'.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			newBranch, _ := cmd.Flags().GetBool("new-branch")

			worktreeName := args[0]

			// Use the manager passed in during command creation
			if manager == nil {
				return fmt.Errorf("manager not available - ensure you're in a git repository with gbm.branchconfig.yaml")
			}

			var branchName string
			var baseBranch string

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

	// Add JIRA key completions for the first positional argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// Use the manager passed in during command creation
			if manager == nil {
				PrintVerbose("Manager not available for completion")
				return nil, cobra.ShellCompDirectiveNoFileComp
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
				// Use the manager passed in during command creation
				if manager == nil {
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
