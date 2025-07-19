package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newMergebackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mergeback [worktree-name] [jira-ticket]",
		Aliases: []string{"mb"},
		Short:   "Create a mergeback worktree to merge changes up the deployment chain",
		Long: `Create a mergeback worktree to merge changes up the deployment chain.

The mergeback command automatically:
- Detects which branch in the mergeback chain needs the merge based on configuration
- Creates a worktree directory with configurable prefix (default: MERGE_<worktree>_<base>)
- Creates a new branch with merge/ prefix
- Integrates with JIRA for branch naming if ticket provided

The worktree prefix can be configured in .gbm/config.toml under settings.mergeback_prefix.
Set to empty string to disable prefixing (worktrees will still include target suffix for namespace separation).

Examples:
  gbm mergeback                            # Auto-detects recent hotfix/merge and creates appropriate mergeback
  gbm mergeback <TAB>                      # Shows smart suggestions from recent git activity (press Tab)
  gbm mergeback fix-auth                   # Creates worktree MERGE_fix-auth_preview with branch merge/fix-auth_preview
  gbm mergeback PROJECT-123                # Creates worktree MERGE_PROJECT-123_main with branch merge/PROJECT-123_summary_main
  gbm mb deploy-hotfix PROJECT-456         # Creates MERGE_deploy-hotfix_<base> worktree with JIRA integration

Tab Completion:
  Press TAB to see intelligent suggestions based on recent hotfix/merge activity,
  or press ENTER for automatic detection with confirmation prompt.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var worktreeName string
			var jiraTicket string

			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				return err
			}

			// Handle auto-detection if no worktree name provided
			if len(args) == 0 {
				worktreeName, jiraTicket, err = autoDetectMergebackTarget(manager)
				if err != nil {
					return fmt.Errorf("failed to auto-detect mergeback target: %w", err)
				}
				PrintVerbose("Auto-detected mergeback target: %s (JIRA: %s)", worktreeName, jiraTicket)
			} else {
				worktreeName = args[0]
				// Check for explicit JIRA ticket in second argument
				if len(args) > 1 {
					jiraTicket = args[1]
				} else if internal.IsJiraKey(worktreeName) {
					jiraTicket = worktreeName
				}
			}

			// Find the target branch for merging
			baseBranch, baseWorktreeName, err := findMergeTargetBranchAndWorktree(manager)
			if err != nil {
				return fmt.Errorf("failed to determine merge target branch: %w", err)
			}

			PrintInfo("Using branch '%s' (worktree '%s') as base for mergeback", baseBranch, baseWorktreeName)

			// Generate mergeback branch name (jiraTicket already set above)

			branchName, err := generateMergebackBranchName(worktreeName, jiraTicket, baseWorktreeName, manager)
			if err != nil {
				return fmt.Errorf("failed to generate mergeback branch name: %w", err)
			}

			// Get mergeback prefix from config and build worktree name
			mergebackPrefix := manager.GetConfig().Settings.MergebackPrefix
			var mergebackWorktreeName string
			if mergebackPrefix != "" {
				mergebackWorktreeName = mergebackPrefix + "_" + worktreeName + "_" + baseWorktreeName
			} else {
				mergebackWorktreeName = worktreeName + "_" + baseWorktreeName
			}

			PrintInfo("Creating mergeback worktree '%s' on branch '%s'", mergebackWorktreeName, branchName)

			// Add the mergeback worktree
			if err := manager.AddWorktree(mergebackWorktreeName, branchName, true, baseBranch); err != nil {
				return fmt.Errorf("failed to add mergeback worktree: %w", err)
			}

			PrintInfo("Mergeback worktree '%s' added successfully", mergebackWorktreeName)
			PrintInfo("Ready to merge changes into '%s'", baseBranch)

			return nil
		},
	}

	// Add smart auto-detection results as tab completion for first argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// First argument: provide smart detection results + JIRA completions
			return getSmartMergebackCompletions(), cobra.ShellCompDirectiveNoFileComp
		} else if len(args) == 1 {
			// Second argument: JIRA ticket completions
			return getJiraCompletions(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

// findMergeTargetBranchAndWorktree finds the branch and worktree that should be used as the base for mergeback
// Uses tree structure and git log to find branches that need merging
func findMergeTargetBranchAndWorktree(manager *internal.Manager) (string, string, error) {
	// Get current working directory and find git root
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return "", "", fmt.Errorf("failed to find git root: %w", err)
	}

	// Load config
	configPath := filepath.Join(repoRoot, internal.DefaultBranchConfigFilename)
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		PrintVerbose("No gbm.branchconfig.yaml found, using default branch as merge target")
		defaultBranch, err := manager.GetGitManager().GetDefaultBranch()
		if err != nil {
			return "", "", err
		}
		return defaultBranch, "main", nil
	}

	// Get deepest leaf nodes (production branches) from all trees
	deepestLeaves := config.Tree.GetAllDeepestLeafNodes()
	if len(deepestLeaves) == 0 {
		return "", "", fmt.Errorf("no leaf nodes found in branch configuration")
	}

	// Check each deepest leaf node (production branch) to see if it needs mergeback
	gitManager := manager.GetGitManager()
	for _, leaf := range deepestLeaves {
		// Check if this leaf has commits that need to be merged into its parent
		if leaf.Parent != nil {
			hasCommits, err := hasCommitsBetweenBranches(gitManager, leaf.Parent.Config.Branch, leaf.Config.Branch)
			if err != nil {
				PrintVerbose("Error checking commits between %s and %s: %v", leaf.Parent.Config.Branch, leaf.Config.Branch, err)
				continue
			}

			if hasCommits {
				PrintVerbose("Found mergeback needed: %s -> %s (worktree '%s')", leaf.Config.Branch, leaf.Parent.Config.Branch, leaf.Parent.Name)
				return leaf.Parent.Config.Branch, leaf.Parent.Name, nil
			}
		}
	}

	// If no deepest leaves need mergeback, check their parents recursively
	return findNextMergeTargetInChain(config.Tree, deepestLeaves, gitManager)
}

// hasCommitsBetweenBranches checks if source has commits that target doesn't have
func hasCommitsBetweenBranches(gitManager *internal.GitManager, targetBranch, sourceBranch string) (bool, error) {
	// First try with origin/ prefix
	cmd := exec.Command("git", "log", fmt.Sprintf("origin/%s..origin/%s", targetBranch, sourceBranch), "--oneline")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to local branches if origin branches don't exist
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..%s", targetBranch, sourceBranch), "--oneline")
		output, err = cmd.Output()
		if err != nil {
			return false, fmt.Errorf("failed to check commits between branches: %w", err)
		}
	}

	// If there's output, there are commits to merge
	return strings.TrimSpace(string(output)) != "", nil
}

// findNextMergeTargetInChain recursively checks parent branches for merge opportunities
func findNextMergeTargetInChain(tree *internal.WorktreeManager, leaves []*internal.WorktreeNode, gitManager *internal.GitManager) (string, string, error) {
	// Check parent branches of the leaves
	checkedBranches := make(map[string]bool)

	for _, leaf := range leaves {
		if leaf.Parent != nil && !checkedBranches[leaf.Parent.Config.Branch] {
			checkedBranches[leaf.Parent.Config.Branch] = true

			// Check if parent needs mergeback to its parent
			if leaf.Parent.Parent != nil {
				hasCommits, err := hasCommitsBetweenBranches(gitManager, leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch)
				if err != nil {
					PrintVerbose("Error checking commits between %s and %s: %v", leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch, err)
					continue
				}

				if hasCommits {
					PrintVerbose("Found mergeback needed: %s -> %s (worktree '%s')", leaf.Parent.Config.Branch, leaf.Parent.Parent.Config.Branch, leaf.Parent.Parent.Name)
					return leaf.Parent.Parent.Config.Branch, leaf.Parent.Parent.Name, nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("no mergeback targets found")
}

// generateMergebackBranchName creates a mergeback branch name with proper formatting
// Now includes target branch suffix to prevent conflicts: merge/PROJECT-123_fix_preview
var generateMergebackBranchName = func(worktreeName, jiraTicket, targetWorktree string, manager *internal.Manager) (string, error) {
	generator := createBranchNameGenerator("merge")
	return generator(worktreeName, jiraTicket, strings.ToLower(targetWorktree), manager)
}

// autoDetectMergebackTarget analyzes recent git history to suggest a mergeback target
func autoDetectMergebackTarget(manager *internal.Manager) (string, string, error) {
	// Get recent mergeable activity from git history (only hotfix and merge types)
	activities, err := manager.GetGitManager().GetRecentMergeableActivity(7) // Last 7 days
	if err != nil {
		return "", "", fmt.Errorf("failed to analyze git history: %w", err)
	}

	// Filter to only hotfix and merge branches, and check if they're ahead
	filteredActivities, err := filterAndValidateActivities(activities, manager)
	if err != nil {
		return "", "", fmt.Errorf("failed to filter activities: %w", err)
	}

	if len(filteredActivities) == 0 {
		return "", "", fmt.Errorf("no recent hotfix or merge activity found that needs mergeback in the last 7 days")
	}

	// Find the most relevant recent activity
	var bestActivity *internal.RecentActivity

	// Prioritize: hotfix > merge, and more recent over older
	for i := range filteredActivities {
		activity := &filteredActivities[i]

		if bestActivity == nil {
			bestActivity = activity
			continue
		}

		// Prioritize by type (hotfix is highest priority)
		if activity.Type == "hotfix" && bestActivity.Type != "hotfix" {
			bestActivity = activity
			continue
		}

		// If same type, prioritize more recent
		if activity.Type == bestActivity.Type && activity.Timestamp.After(bestActivity.Timestamp) {
			bestActivity = activity
			continue
		}
	}

	if bestActivity == nil || bestActivity.WorktreeName == "" {
		return "", "", fmt.Errorf("could not determine worktree name from recent activity")
	}

	PrintInfo("Found recent %s activity: %s (%s)", bestActivity.Type, bestActivity.WorktreeName, bestActivity.CommitMessage)

	// Show user what was found and ask for confirmation
	fmt.Printf("\n%s\n", internal.FormatSubHeader("Recent Activity Detected:"))
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Type"), bestActivity.Type)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Worktree"), bestActivity.WorktreeName)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Branch"), bestActivity.BranchName)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Commit"), bestActivity.CommitMessage)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Author"), bestActivity.Author)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Date"), bestActivity.Timestamp.Format("2006-01-02 15:04"))

	// Ask for confirmation
	fmt.Printf("\n%s ", internal.FormatPrompt("Use this for mergeback? (y/n):"))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return "", "", fmt.Errorf("failed to read confirmation: %w", err)
	}

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		return "", "", fmt.Errorf("mergeback cancelled by user")
	}

	return bestActivity.WorktreeName, bestActivity.JiraTicket, nil
}

// filterAndValidateActivities filters activities to only hotfix/merge types and validates they're ahead
func filterAndValidateActivities(activities []internal.RecentActivity, manager *internal.Manager) ([]internal.RecentActivity, error) {
	var validActivities []internal.RecentActivity

	for _, activity := range activities {
		// Only consider hotfix and merge types
		if activity.Type != "hotfix" && activity.Type != "merge" {
			PrintVerbose("Skipping %s activity: %s (not hotfix or merge)", activity.Type, activity.WorktreeName)
			continue
		}

		// Check if the branch is ahead of target branches
		if isActivityRelevantForMergeback(activity, manager) {
			validActivities = append(validActivities, activity)
			PrintVerbose("Including %s activity: %s (branch is ahead)", activity.Type, activity.WorktreeName)
		} else {
			PrintVerbose("Skipping %s activity: %s (branch not ahead of targets)", activity.Type, activity.WorktreeName)
		}
	}

	return validActivities, nil
}

// isActivityRelevantForMergeback checks if the activity's branch is ahead of potential merge targets
func isActivityRelevantForMergeback(activity internal.RecentActivity, manager *internal.Manager) bool {
	if activity.BranchName == "" {
		// If we don't have branch info, assume it's relevant
		PrintVerbose("No branch info for %s, assuming relevant", activity.WorktreeName)
		return true
	}

	// Get current working directory and find git root
	wd, err := os.Getwd()
	if err != nil {
		PrintVerbose("Could not get working directory: %v", err)
		return true // Assume relevant if we can't check
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		PrintVerbose("Could not find git root: %v", err)
		return true // Assume relevant if we can't check
	}

	// Look for gbm.branchconfig.yaml file to get merge targets
	configPath := filepath.Join(repoRoot, internal.DefaultBranchConfigFilename)
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		PrintVerbose("No gbm.branchconfig.yaml found, assuming activity is relevant")
		return true // If no config, assume relevant
	}

	// Find potential merge targets (branches that this branch should merge into)
	potentialTargets := findPotentialMergeTargets(activity.BranchName, config)
	if len(potentialTargets) == 0 {
		PrintVerbose("No merge targets found for %s, assuming relevant", activity.BranchName)
		return true
	}

	// Check if the source branch is ahead of any of the potential targets
	for _, targetBranch := range potentialTargets {
		isAhead, err := isBranchAheadOf(activity.BranchName, targetBranch, manager)
		if err != nil {
			PrintVerbose("Could not check if %s is ahead of %s: %v", activity.BranchName, targetBranch, err)
			continue
		}

		if isAhead {
			PrintVerbose("Branch %s is ahead of %s, activity is relevant", activity.BranchName, targetBranch)
			return true
		}
	}

	PrintVerbose("Branch %s is not ahead of any merge targets", activity.BranchName)
	return false
}

// findPotentialMergeTargets finds branches that the given branch should merge into
func findPotentialMergeTargets(branchName string, config *internal.GBMConfig) []string {
	targets := []string{} // Initialize as empty slice, not nil

	// Look through the mergeback configuration to find what this branch merges into
	for _, worktreeConfig := range config.Worktrees {
		if worktreeConfig.Branch == branchName && worktreeConfig.MergeInto != "" {
			targets = append(targets, worktreeConfig.MergeInto)
		}
	}

	// If no explicit merge target found, check common merge chains
	// For hotfix branches, they typically merge into production first, then up the chain
	if len(targets) == 0 && strings.Contains(strings.ToLower(branchName), "hotfix") {
		// Find the production/deployment branch (one that has MergeInto set, not the root)
		// Production is typically the branch that merges into preview/staging
		for _, worktreeConfig := range config.Worktrees {
			if worktreeConfig.MergeInto != "" && strings.Contains(strings.ToLower(worktreeConfig.Branch), "prod") {
				targets = append(targets, worktreeConfig.Branch)
				break
			}
		}
		// If no production found, try to find any branch that merges into something (deployment branch)
		if len(targets) == 0 {
			for _, worktreeConfig := range config.Worktrees {
				if worktreeConfig.MergeInto != "" {
					targets = append(targets, worktreeConfig.Branch)
					break
				}
			}
		}
	}

	return targets
}

// isBranchAheadOf checks if sourceBranch has commits that targetBranch doesn't have
func isBranchAheadOf(sourceBranch, targetBranch string, _ *internal.Manager) (bool, error) {
	// Get the repo root from working directory
	wd, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return false, fmt.Errorf("failed to find git root: %w", err)
	}

	// Use git rev-list to count commits that are in source but not in target
	output, err := internal.ExecGitCommand(repoRoot, "rev-list", "--count", fmt.Sprintf("%s..%s", targetBranch, sourceBranch))
	if err != nil {
		return false, fmt.Errorf("failed to check branch relationship: %w", err)
	}

	countStr := strings.TrimSpace(string(output))
	var count int
	if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil {
		return false, fmt.Errorf("failed to parse commit count: %w", err)
	}

	return count > 0, nil
}

// getSmartMergebackCompletions provides intelligent tab completion based on recent activity
func getSmartMergebackCompletions() []string {
	var completions []string

	// Try to get smart detection results
	manager, err := createInitializedManager() // Legacy call in helper function
	if err != nil {
		// Fallback to JIRA completions if manager fails
		return getJiraCompletions()
	}

	// Get recent mergeable activity (same logic as auto-detection)
	activities, err := manager.GetGitManager().GetRecentMergeableActivity(7)
	if err != nil {
		// Fallback to JIRA completions if git analysis fails
		return getJiraCompletions()
	}

	// Filter and validate activities
	filteredActivities, err := filterAndValidateActivities(activities, manager)
	if err != nil {
		// Fallback to JIRA completions if filtering fails
		return getJiraCompletions()
	}

	// Convert filtered activities to completions in priority order
	for _, activity := range filteredActivities {
		if activity.WorktreeName == "" {
			continue
		}

		// Format: "WORKTREE_NAME\tType: hotfix | Branch: hotfix/PROJECT-123 | Date: 2025-07-12"
		description := fmt.Sprintf("Type: %s | Branch: %s | Date: %s",
			activity.Type,
			activity.BranchName,
			activity.Timestamp.Format("2006-01-02"))

		completion := fmt.Sprintf("%s\t%s", activity.WorktreeName, description)
		completions = append(completions, completion)
	}

	// If no smart suggestions, fall back to JIRA
	if len(completions) == 0 {
		return getJiraCompletions()
	}

	// Add a separator and JIRA completions as additional options
	jiraCompletions := getJiraCompletions()
	if len(jiraCompletions) > 0 {
		// Add separator
		completions = append(completions, "---\tOther JIRA tickets:")
		// Add JIRA completions
		completions = append(completions, jiraCompletions...)
	}

	return completions
}

// getJiraCompletions provides JIRA ticket completions as fallback
func getJiraCompletions() []string {
	var completions []string

	// Try to get config for JIRA completion
	manager, err := createInitializedManager() // Legacy call in helper function
	if err != nil {
		return completions
	}

	// Complete JIRA keys with summaries for context
	jiraIssues, err := internal.GetJiraIssues(manager)
	if err != nil {
		return completions
	}

	for _, issue := range jiraIssues {
		completion := fmt.Sprintf("%s\t%s", issue.Key, issue.Summary)
		completions = append(completions, completion)
	}

	return completions
}
