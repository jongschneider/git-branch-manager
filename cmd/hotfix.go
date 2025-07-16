package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newHotfixCommand() *cobra.Command {
	cmd := &cobra.Command{
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
			var branchName string
			
			if len(args) > 1 {
				// If second argument is provided, check if it's already a branch name or a JIRA ticket
				secondArg := args[1]
				if strings.HasPrefix(secondArg, "hotfix/") {
					// It's already a branch name (from autocompletion), use it directly
					branchName = secondArg
				} else if internal.IsJiraKey(secondArg) {
					// It's a JIRA ticket, generate branch name from it
					branchName, err = generateHotfixBranchName(worktreeName, secondArg, manager)
					if err != nil {
						return fmt.Errorf("failed to generate hotfix branch name: %w", err)
					}
				} else {
					// It's some other string, treat as JIRA ticket anyway
					branchName, err = generateHotfixBranchName(worktreeName, secondArg, manager)
					if err != nil {
						return fmt.Errorf("failed to generate hotfix branch name: %w", err)
					}
				}
			} else if internal.IsJiraKey(worktreeName) {
				// First argument is a JIRA ticket, use it
				branchName, err = generateHotfixBranchName(worktreeName, worktreeName, manager)
				if err != nil {
					return fmt.Errorf("failed to generate hotfix branch name: %w", err)
				}
			} else {
				// No JIRA ticket provided, generate simple branch name
				branchName, err = generateHotfixBranchName(worktreeName, "", manager)
				if err != nil {
					return fmt.Errorf("failed to generate hotfix branch name: %w", err)
				}
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
			// Build the deployment chain message
		deploymentChain := buildDeploymentChain(baseBranch, manager)
		PrintInfo("Remember to merge back through the deployment chain: %s", deploymentChain)

			return nil
		},
	}

	// Add JIRA key completions for first argument, JIRA summary for second argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// First argument: JIRA keys with summaries for context
			manager, err := createInitializedManager() // Legacy call in completion function
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
		} else if len(args) == 1 {
			// Second argument: Generate hotfix branch name based on the JIRA key
			worktreeName := args[0]
			if internal.IsJiraKey(worktreeName) {
				// Try to get config for JIRA completion
				manager, err := createInitializedManager()
				if err != nil {
					// Fallback to default branch name generation
					branchName := fmt.Sprintf("hotfix/%s", strings.ToUpper(worktreeName))
					return []string{branchName}, cobra.ShellCompDirectiveNoFileComp
				}

				// Generate hotfix branch name from JIRA
				// In autocompletion, worktreeName is the JIRA key, so use it as both worktree and JIRA ticket
				branchName, err := generateHotfixBranchName(worktreeName, worktreeName, manager)
				if err != nil {
					// Fallback to default branch name generation
					branchName = fmt.Sprintf("hotfix/%s", strings.ToUpper(worktreeName))
				}
				return []string{branchName}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
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

	// Look for gbm.branchconfig.yaml file
	configPath := filepath.Join(repoRoot, internal.DefaultBranchConfigFilename)
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		// If no config file, fall back to default branch
		PrintVerbose("No gbm.branchconfig.yaml found, using default branch as production")
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
				// Replace any prefix with hotfix/
				parts := strings.Split(jiraBranchName, "/")
				if len(parts) > 1 {
					parts[0] = "hotfix"
					branchName = strings.Join(parts, "/")
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

// buildDeploymentChain builds the complete deployment chain from base branch to final target
func buildDeploymentChain(baseBranch string, manager *internal.Manager) string {
	config := loadGBMConfig()
	if config == nil {
		return baseBranch
	}

	chain := buildMergeChain(baseBranch, config)
	if len(chain) <= 1 {
		return baseBranch
	}

	return strings.Join(chain, " → ")
}

// loadGBMConfig loads the gbm.branchconfig.yaml file, returning nil if not found
func loadGBMConfig() *internal.GBMConfig {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return nil
	}

	configPath := filepath.Join(repoRoot, internal.DefaultBranchConfigFilename)
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		return nil
	}

	return config
}

// buildMergeChain traverses the merge configuration to build the complete chain
func buildMergeChain(baseBranch string, config *internal.GBMConfig) []string {
	chain := []string{baseBranch}
	currentBranch := baseBranch

	const maxIterations = 10 // Prevent infinite loops
	for i := 0; i < maxIterations; i++ {
		nextBranch := findBranchThatMergesInto(currentBranch, config)
		if nextBranch == "" {
			break
		}
		
		chain = append(chain, nextBranch)
		currentBranch = nextBranch
	}

	return chain
}

// findBranchThatMergesInto finds the branch that merges into the given target branch
func findBranchThatMergesInto(targetBranch string, config *internal.GBMConfig) string {
	for _, worktreeConfig := range config.Worktrees {
		if worktreeConfig.MergeInto == targetBranch {
			return worktreeConfig.Branch
		}
	}
	return ""
}
