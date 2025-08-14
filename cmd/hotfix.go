package cmd

import (
	"errors"
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_hotfixCreator.go . hotfixCreator

// hotfixCreator interface abstracts the Manager operations needed for hotfix creation
type hotfixCreator interface {
	AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error
	GetConfig() *internal.Config
	GetGBMConfig() *internal.GBMConfig
	FindProductionBranch() (string, error)
	GetJiraIssues() ([]internal.JiraIssue, error)
	GenerateBranchFromJira(jiraKey string) (string, error)
	GetDefaultBranch() (string, error)
}

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
			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return err
				}

				PrintVerbose("%v", err)
			}

			return handleHotfix(manager, args)
		},
	}

	// Add JIRA key completions for first argument, JIRA summary for second argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// First argument: JIRA keys with summaries for context
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}

				PrintVerbose("%v", err)
			}

			// Complete JIRA keys with summaries for context
			jiraIssues, err := manager.GetJiraIssues()
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
				branchName, err := generateHotfixBranchNameWithCreator(worktreeName, worktreeName, manager)
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

// handleHotfix handles the hotfix creation logic using the hotfixCreator interface
func handleHotfix(creator hotfixCreator, args []string) error {
	worktreeName := args[0]

	// Find the production branch (last in mergeback chain)
	baseBranch, err := creator.FindProductionBranch()
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
			branchName, err = generateHotfixBranchNameWithCreator(worktreeName, secondArg, creator)
			if err != nil {
				return fmt.Errorf("failed to generate hotfix branch name: %w", err)
			}
		} else {
			// It's some other string, treat as JIRA ticket anyway
			branchName, err = generateHotfixBranchNameWithCreator(worktreeName, secondArg, creator)
			if err != nil {
				return fmt.Errorf("failed to generate hotfix branch name: %w", err)
			}
		}
	} else if internal.IsJiraKey(worktreeName) {
		// First argument is a JIRA ticket, use it
		branchName, err = generateHotfixBranchNameWithCreator(worktreeName, worktreeName, creator)
		if err != nil {
			return fmt.Errorf("failed to generate hotfix branch name: %w", err)
		}
	} else {
		// No JIRA ticket provided, generate simple branch name
		branchName, err = generateHotfixBranchNameWithCreator(worktreeName, "", creator)
		if err != nil {
			return fmt.Errorf("failed to generate hotfix branch name: %w", err)
		}
	}

	// Get hotfix prefix from config and build worktree name
	hotfixPrefix := creator.GetConfig().Settings.HotfixPrefix
	var hotfixWorktreeName string
	if hotfixPrefix != "" {
		hotfixWorktreeName = hotfixPrefix + "_" + worktreeName
	} else {
		hotfixWorktreeName = worktreeName
	}

	PrintInfo("Creating hotfix worktree '%s' on branch '%s'", hotfixWorktreeName, branchName)

	// Add the hotfix worktree
	if err := creator.AddWorktree(hotfixWorktreeName, branchName, true, baseBranch); err != nil {
		return fmt.Errorf("failed to add hotfix worktree: %w", err)
	}

	PrintInfo("Hotfix worktree '%s' added successfully", hotfixWorktreeName)
	// Build the deployment chain message
	deploymentChain := buildDeploymentChain(baseBranch, creator.GetGBMConfig())
	PrintInfo("Remember to merge back through the deployment chain: %s", deploymentChain)

	return nil
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

// buildDeploymentChain builds the complete deployment chain from base branch to final target
func buildDeploymentChain(baseBranch string, config *internal.GBMConfig) string {
	if config == nil {
		return baseBranch
	}

	chain := buildMergeChain(baseBranch, config)
	if len(chain) <= 1 {
		return baseBranch
	}

	return strings.Join(chain, " → ")
}

// buildMergeChain traverses the merge configuration to build the complete chain
func buildMergeChain(baseBranch string, config *internal.GBMConfig) []string {
	chain := []string{baseBranch}
	currentBranch := baseBranch

	const maxIterations = 10 // Prevent infinite loops
	for range maxIterations {
		nextBranch := findMergeIntoTarget(currentBranch, config)
		if nextBranch == "" {
			break
		}

		chain = append(chain, nextBranch)
		currentBranch = nextBranch
	}

	return chain
}

// findMergeIntoTarget finds where the given branch merges into
func findMergeIntoTarget(sourceBranch string, config *internal.GBMConfig) string {
	for _, worktreeConfig := range config.Worktrees {
		if worktreeConfig.Branch == sourceBranch {
			return worktreeConfig.MergeInto
		}
	}
	return ""
}

// generateHotfixBranchName creates a hotfix branch name with proper formatting
var generateHotfixBranchName = func(worktreeName, jiraTicket string, manager *internal.Manager) (string, error) {
	generator := createBranchNameGenerator("hotfix")
	return generator(worktreeName, jiraTicket, "", manager)
}

// generateHotfixBranchNameWithCreator creates a hotfix branch name using the hotfixCreator interface
func generateHotfixBranchNameWithCreator(worktreeName, jiraTicket string, creator hotfixCreator) (string, error) {
	generator := createBranchNameGeneratorWithCreator("hotfix", creator)
	return generator(worktreeName, jiraTicket, "")
}
