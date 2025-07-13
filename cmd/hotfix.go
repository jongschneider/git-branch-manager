package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var hotfixCmd = &cobra.Command{
	Use:     "hotfix <worktree-name> [jira-ticket]",
	Aliases: []string{"hf"},
	Short:   "Create a hotfix worktree from the production branch",
	Long: `Create a hotfix worktree based on the last branch in the mergeback chain.

The hotfix command automatically:
- Finds the production branch (bottom of mergeback chain) as the base
- Creates a worktree directory with configurable prefix (default: HOTFIX_)
- Creates a new branch with hotfix/ prefix
- Integrates with JIRA for branch naming if ticket provided

The worktree prefix can be configured in .gbm/config.toml under settings.hotfix_prefix.
Set to empty string to disable prefixing.

Examples:
  gbm hotfix critical-bug                  # Creates worktree HOTFIX_critical-bug with branch hotfix/critical-bug
  gbm hotfix PROJECT-123                   # Creates worktree HOTFIX_PROJECT-123 with branch hotfix/PROJECT-123_summary_from_jira
  gbm hf auth-fix PROJECT-456              # Creates HOTFIX_auth-fix worktree with JIRA integration`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		worktreeName := args[0]

		// Create manager
		manager, err := createInitializedManager()
		if err != nil {
			return err
		}

		// Find the production branch (last in mergeback chain)
		baseBranch, err := findProductionBranch(manager)
		if err != nil {
			return fmt.Errorf("failed to determine production branch: %w", err)
		}

		PrintInfo("Using production branch '%s' as base for hotfix", baseBranch)

		// Generate hotfix branch name
		var jiraTicket string
		if len(args) > 1 {
			jiraTicket = args[1]
		} else if internal.IsJiraKey(worktreeName) {
			jiraTicket = worktreeName
		}

		branchName, err := generateHotfixBranchName(worktreeName, jiraTicket, manager)
		if err != nil {
			return fmt.Errorf("failed to generate hotfix branch name: %w", err)
		}

		// Get hotfix prefix from config and build worktree name
		hotfixPrefix := manager.GetConfig().Settings.HotfixPrefix
		var hotfixWorktreeName string
		if hotfixPrefix != "" {
			hotfixWorktreeName = hotfixPrefix + "_" + worktreeName
		} else {
			hotfixWorktreeName = worktreeName
		}

		PrintInfo("Creating hotfix worktree '%s' on branch '%s'", hotfixWorktreeName, branchName)

		// Add the hotfix worktree
		if err := manager.AddWorktree(hotfixWorktreeName, branchName, true, baseBranch); err != nil {
			return fmt.Errorf("failed to add hotfix worktree: %w", err)
		}

		PrintInfo("Hotfix worktree '%s' added successfully", hotfixWorktreeName)
		PrintInfo("Remember to merge back through the deployment chain: %s → preview → main", baseBranch)

		return nil
	},
}

// findProductionBranch finds the branch at the bottom of the mergeback chain
// This is the branch that has no merge_into target configured
func findProductionBranch(manager *internal.Manager) (string, error) {
	// Get current working directory and find git root
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return "", fmt.Errorf("failed to find git root: %w", err)
	}

	// Look for .gbm.config.yaml file
	configPath := filepath.Join(repoRoot, ".gbm.config.yaml")
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		// If no config file, fall back to default branch
		PrintVerbose("No .gbm.config.yaml found, using default branch as production")
		return manager.GetGitManager().GetDefaultBranch()
	}

	// Find the branch with no merge_into target (bottom of chain)
	var productionWorktree string
	var productionBranch string

	for worktreeName, worktreeConfig := range config.Worktrees {
		if worktreeConfig.MergeInto == "" {
			// This is a root branch (no merge target)
			if productionWorktree != "" {
				// Multiple root branches found, need to determine which is production
				PrintVerbose("Multiple root branches found: %s and %s", productionWorktree, worktreeName)
				// Use heuristics: look for common production names
				if isProductionBranchName(worktreeConfig.Branch) {
					productionWorktree = worktreeName
					productionBranch = worktreeConfig.Branch
				}
			} else {
				productionWorktree = worktreeName
				productionBranch = worktreeConfig.Branch
			}
		}
	}

	if productionBranch == "" {
		return "", fmt.Errorf("no production branch found in mergeback configuration")
	}

	PrintVerbose("Found production branch: %s (worktree: %s)", productionBranch, productionWorktree)
	return productionBranch, nil
}

// isProductionBranchName returns true if the branch name suggests it's a production branch
func isProductionBranchName(branchName string) bool {
	prodNames := []string{"prod", "production", "master", "main", "release"}
	lowerBranch := strings.ToLower(branchName)

	for _, prodName := range prodNames {
		if strings.Contains(lowerBranch, prodName) {
			return true
		}
	}
	return false
}

// generateHotfixBranchName creates a hotfix branch name with proper formatting
func generateHotfixBranchName(worktreeName, jiraTicket string, manager *internal.Manager) (string, error) {
	var branchName string

	if jiraTicket != "" && internal.IsJiraKey(jiraTicket) {
		// Generate branch name from JIRA ticket
		if manager != nil {
			jiraBranchName, err := internal.GenerateBranchFromJira(jiraTicket, manager)
			if err != nil {
				PrintVerbose("Failed to generate branch name from JIRA issue %s: %v", jiraTicket, err)
				// Fallback to simple format
				branchName = fmt.Sprintf("hotfix/%s", strings.ToUpper(jiraTicket))
			} else {
				// Replace feature/ with hotfix/
				if strings.HasPrefix(jiraBranchName, "feature/") {
					branchName = "hotfix/" + jiraBranchName[8:] // Remove "feature/" prefix
				} else if strings.HasPrefix(jiraBranchName, "bugfix/") {
					branchName = "hotfix/" + jiraBranchName[7:] // Remove "bugfix/" prefix
				} else {
					branchName = "hotfix/" + jiraBranchName
				}
			}
		} else {
			// No manager available, use simple format
			branchName = fmt.Sprintf("hotfix/%s", strings.ToUpper(jiraTicket))
		}
	} else {
		// Generate from worktree name
		cleanName := strings.ReplaceAll(worktreeName, " ", "-")
		cleanName = strings.ReplaceAll(cleanName, "_", "-")
		cleanName = strings.ToLower(cleanName)
		branchName = "hotfix/" + cleanName
	}

	return branchName, nil
}

func init() {
	rootCmd.AddCommand(hotfixCmd)

	// Add JIRA key completions for both positional arguments
	hotfixCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 || len(args) == 1 {
			// Try to get config for JIRA completion
			manager, err := createInitializedManager()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// Complete JIRA keys with summaries for context
			jiraIssues, err := internal.GetJiraIssues(manager)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completions []string
			for _, issue := range jiraIssues {
				completion := fmt.Sprintf("%s\t%s", issue.Key, issue.Summary)
				completions = append(completions, completion)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}