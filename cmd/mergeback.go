package cmd

import (
	"errors"
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
		Use:     "mergeback [worktree-name]",
		Aliases: []string{"mb"},
		Short:   "Create a mergeback worktree to merge changes up the deployment chain",
		Long: `Create a mergeback worktree to merge changes up the deployment chain.

The mergeback command automatically:
- Detects which branch in the mergeback chain needs the merge based on configuration
- Creates a worktree directory with configurable prefix (default: MERGE_<worktree>_<base>)
- Creates a new branch with merge/ prefix
- Offers to perform the merge automatically with user confirmation

The worktree prefix can be configured in .gbm/config.toml under settings.mergeback_prefix.
Set to empty string to disable prefixing (worktrees will still include target suffix for namespace separation).

After creating the worktree, gbm will show you which commits will be merged and ask if you 
want to perform the merge automatically. If conflicts occur, gbm will let you know and you 
can resolve them manually in the mergeback worktree.

Examples:
  gbm mergeback                            # Auto-detects recent merge activity and creates appropriate mergeback
  gbm mergeback <TAB>                      # Shows smart suggestions from recent git activity (press Tab)
  gbm mergeback fix-auth                   # Creates worktree MERGE_fix-auth_preview with branch merge/fix-auth_preview
  gbm mb deploy-hotfix                     # Creates MERGE_deploy-hotfix_<base> worktree

Tab Completion:
  Press TAB to see intelligent suggestions based on recent merge activity,
  or press ENTER for automatic detection with confirmation prompt.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return err
				}

				PrintVerbose("%v", err)
			}

			// Find the source and target branches for merging
			sourceBranch, baseBranch, baseWorktreeName, sourceWorktreeName, err := findMergeTargetBranchAndWorktree(manager)
			if err != nil {
				return fmt.Errorf("failed to determine merge target branch: %w", err)
			}

			PrintInfo("Mergeback needed: '%s' → '%s'", sourceWorktreeName, baseWorktreeName)
			PrintVerbose("Will merge from '%s' into '%s'", sourceBranch, baseBranch)

			// Use source worktree name for naming (e.g., "production" for production → preview)
			// User can override by passing worktree name as argument
			var worktreeName string
			if len(args) == 0 {
				worktreeName = sourceWorktreeName
			} else {
				worktreeName = args[0]
			}

			// Generate mergeback branch name
			branchName := fmt.Sprintf("merge/%s_%s", worktreeName, strings.ToLower(baseWorktreeName))

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

			// Offer to perform the merge automatically
			if err := offerMergeExecution(manager, mergebackWorktreeName, worktreeName, sourceBranch, baseBranch); err != nil {
				return fmt.Errorf("merge execution failed: %w", err)
			}

			return nil
		},
	}

	// Add smart auto-detection results as tab completion for first argument
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// First argument: provide smart detection results
			return getSmartMergebackCompletions(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

// findMergeTargetBranchAndWorktree finds the source branch with changes and target branch/worktree for mergeback
// Uses tree structure and git log to find branches that need merging
// Returns: sourceBranch, targetBranch, targetWorktreeName, sourceWorktreeName, error
func findMergeTargetBranchAndWorktree(manager *internal.Manager) (string, string, string, string, error) {
	// Get current working directory and find git root
	wd, err := os.Getwd()
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to find git root: %w", err)
	}

	// Load config
	configPath := filepath.Join(repoRoot, internal.DefaultBranchConfigFilename)
	config, err := internal.ParseGBMConfig(configPath)
	if err != nil {
		PrintVerbose("No gbm.branchconfig.yaml found, using default branch as merge target")
		defaultBranch, err := manager.GetGitManager().GetDefaultBranch()
		if err != nil {
			return "", "", "", "", err
		}
		// Without config, we can't determine source branch reliably
		return "", defaultBranch, "main", "", fmt.Errorf("cannot determine source branch without gbm.branchconfig.yaml")
	}

	// Get deepest leaf nodes (production branches) from all trees
	deepestLeaves := config.Tree.GetAllDeepestLeafNodes()
	if len(deepestLeaves) == 0 {
		return "", "", "", "", fmt.Errorf("no leaf nodes found in branch configuration")
	}

	// Check each deepest leaf node (production branch) to see if it needs mergeback
	for _, leaf := range deepestLeaves {
		// Check if this leaf has commits that need to be merged into its parent
		if leaf.Parent != nil {
			hasCommits, err := hasCommitsBetweenBranches(leaf.Parent.Config.Branch, leaf.Config.Branch)
			if err != nil {
				PrintVerbose("Error checking commits between %s and %s: %v", leaf.Parent.Config.Branch, leaf.Config.Branch, err)
				continue
			}

			if hasCommits {
				PrintVerbose("Found mergeback needed: %s -> %s (worktree '%s')", leaf.Config.Branch, leaf.Parent.Config.Branch, leaf.Parent.Name)
				// Return source branch (leaf with changes), target branch, target worktree name, source worktree name
				// Use origin/ prefix to ensure we merge from remote state
				return "origin/" + leaf.Config.Branch, leaf.Parent.Config.Branch, leaf.Parent.Name, leaf.Name, nil
			}
		}
	}

	// If no deepest leaves need mergeback, check their parents recursively
	return findNextMergeTargetInChain(deepestLeaves)
}

// hasCommitsBetweenBranches checks if source has commits that target doesn't have
func hasCommitsBetweenBranches(targetBranch, sourceBranch string) (bool, error) {
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
func findNextMergeTargetInChain(leaves []*internal.WorktreeNode) (string, string, string, string, error) {
	// Check parent branches of the leaves
	checkedBranches := make(map[string]bool)

	for _, leaf := range leaves {
		if leaf.Parent != nil && !checkedBranches[leaf.Parent.Config.Branch] {
			checkedBranches[leaf.Parent.Config.Branch] = true

			// Check if parent needs mergeback to its parent
			if leaf.Parent.Parent != nil {
				hasCommits, err := hasCommitsBetweenBranches(leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch)
				if err != nil {
					PrintVerbose("Error checking commits between %s and %s: %v", leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch, err)
					continue
				}

				if hasCommits {
					PrintVerbose("Found mergeback needed: %s -> %s (worktree '%s')", leaf.Parent.Config.Branch, leaf.Parent.Parent.Config.Branch, leaf.Parent.Parent.Name)
					// Return source branch, target branch, target worktree name, source worktree name
					// Use origin/ prefix to ensure we merge from remote state
					return "origin/" + leaf.Parent.Config.Branch, leaf.Parent.Parent.Config.Branch, leaf.Parent.Parent.Name, leaf.Parent.Name, nil
				}
			}
		}
	}

	return "", "", "", "", fmt.Errorf("no mergeback targets found")
}

// autoDetectMergebackTarget analyzes recent git history to suggest a mergeback target
func autoDetectMergebackTarget(manager *internal.Manager) (string, error) {
	// Get recent mergeable activity from git history (only hotfix and merge types)
	activities, err := manager.GetGitManager().GetRecentMergeableActivity(7) // Last 7 days
	if err != nil {
		return "", fmt.Errorf("failed to analyze git history: %w", err)
	}

	// Filter to only hotfix and merge branches, and check if they're ahead
	filteredActivities, err := filterAndValidateActivities(activities, manager)
	if err != nil {
		return "", fmt.Errorf("failed to filter activities: %w", err)
	}

	if len(filteredActivities) == 0 {
		return "", fmt.Errorf("no recent hotfix or merge activity found that needs mergeback in the last 7 days")
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
		return "", fmt.Errorf("could not determine worktree name from recent activity")
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
		return "", fmt.Errorf("failed to read confirmation: %w", err)
	}

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		return "", fmt.Errorf("mergeback cancelled by user")
	}

	return bestActivity.WorktreeName, nil
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
	completions := make([]string, 0)

	// Try to get smart detection results
	manager, err := createInitializedManager()
	if err != nil {
		if !errors.Is(err, ErrLoadGBMConfig) {
			return completions
		}

		PrintVerbose("%v", err)
		return completions
	}

	// Get recent mergeable activity (same logic as auto-detection)
	activities, err := manager.GetGitManager().GetRecentMergeableActivity(7)
	if err != nil {
		return completions
	}

	// Filter and validate activities
	filteredActivities, err := filterAndValidateActivities(activities, manager)
	if err != nil {
		return completions
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

	return completions
}

// offerMergeExecution prompts user to perform the merge and executes it if confirmed
func offerMergeExecution(manager *internal.Manager, mergebackWorktreeName, sourceName, sourceBranch, targetBranch string) error {
	// Get git root
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return fmt.Errorf("failed to find git root: %w", err)
	}

	// Get worktree path
	worktreePath := filepath.Join(repoRoot, internal.DefaultWorktreeDirname, mergebackWorktreeName)

	// Get commits that will be merged
	mergeBranch := fmt.Sprintf("merge/%s_%s", sourceName, strings.ToLower(targetBranch))
	commits, err := getCommitsToMerge(repoRoot, targetBranch, sourceBranch)
	if err != nil {
		PrintVerbose("Could not get commits to merge: %v", err)
		commits = []string{"(unable to determine commits)"}
	}

	// Display merge information
	fmt.Printf("\n%s\n", internal.FormatSubHeader("Merge Information:"))
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Source"), sourceName)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Source Branch"), sourceBranch)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Target Branch"), targetBranch)
	fmt.Printf("  %s: %s\n", internal.FormatInfo("Merge Branch"), mergeBranch)
	fmt.Printf("  %s: %d commits\n", internal.FormatInfo("Commits to Merge"), len(commits))

	if len(commits) > 0 && commits[0] != "(unable to determine commits)" {
		fmt.Printf("\n%s\n", internal.FormatSubHeader("Recent Commits:"))
		// Show up to 5 most recent commits
		maxCommits := len(commits)
		if maxCommits > 5 {
			maxCommits = 5
		}
		for i := 0; i < maxCommits; i++ {
			fmt.Printf("  • %s\n", commits[i])
		}
		if len(commits) > 5 {
			fmt.Printf("  ... and %d more commits\n", len(commits)-5)
		}
	}

	// Ask for confirmation
	fmt.Printf("\n%s ", internal.FormatPrompt("Perform the merge automatically? (y/n):"))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		PrintInfo("Merge cancelled. You can perform the merge manually in the worktree '%s'", mergebackWorktreeName)
		return nil
	}

	// Perform the merge
	PrintInfo("Performing merge of '%s' into '%s'...", sourceBranch, targetBranch)

	// Execute the merge in the worktree
	if err := performMerge(worktreePath, sourceBranch, targetBranch); err != nil {
		if isMergeConflict(err) {
			PrintInfo("Merge conflicts detected. Please resolve conflicts manually in worktree '%s'", mergebackWorktreeName)
			PrintInfo("After resolving conflicts, use: git add . && git commit")
			return nil
		}
		return fmt.Errorf("merge failed: %w", err)
	}

	PrintInfo("Merge completed successfully!")
	PrintInfo("Review the merge in worktree '%s' before pushing", mergebackWorktreeName)

	return nil
}

// getCommitsToMerge gets the list of commits that will be merged
func getCommitsToMerge(repoRoot, targetBranch, sourceBranch string) ([]string, error) {
	// Verify the source branch exists
	if _, err := internal.ExecGitCommand(repoRoot, "rev-parse", "--verify", sourceBranch); err != nil {
		// Try with origin/ prefix if not found
		originBranch := "origin/" + sourceBranch
		if _, err := internal.ExecGitCommand(repoRoot, "rev-parse", "--verify", originBranch); err != nil {
			return nil, fmt.Errorf("could not find source branch %s or %s", sourceBranch, originBranch)
		}
		sourceBranch = originBranch
	}

	// Get commits that are in source but not in target
	var targetRef string
	if strings.Contains(targetBranch, "/") {
		targetRef = targetBranch
	} else {
		targetRef = "origin/" + targetBranch
	}

	output, err := internal.ExecGitCommand(repoRoot, "log", "--oneline", fmt.Sprintf("%s..%s", targetRef, sourceBranch))
	if err != nil {
		return nil, fmt.Errorf("failed to get commit list: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	return lines, nil
}

// performMerge executes the actual merge operation
// sourceBranch is the branch being merged FROM (e.g., "production")
// targetBranch is the branch being merged INTO (e.g., "preview")
// The worktree should already be on a merge branch created from targetBranch
func performMerge(worktreePath, sourceBranch, targetBranch string) error {
	// Verify we can access the source branch
	if _, err := internal.ExecGitCommand(worktreePath, "rev-parse", "--verify", sourceBranch); err != nil {
		// Try with origin/ prefix
		originBranch := "origin/" + sourceBranch
		if _, err := internal.ExecGitCommand(worktreePath, "rev-parse", "--verify", originBranch); err != nil {
			return fmt.Errorf("could not find source branch %s or %s", sourceBranch, originBranch)
		}
		sourceBranch = originBranch
	}

	// Perform the merge
	output, err := internal.ExecGitCommandCombined(worktreePath, "merge", "--no-ff", "-m", fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch), sourceBranch)
	if err != nil {
		// Include output in error for better debugging
		return fmt.Errorf("git merge failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// isMergeConflict checks if the error indicates a merge conflict
func isMergeConflict(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "conflict") ||
		strings.Contains(errorStr, "merge conflict") ||
		strings.Contains(errorStr, "automatic merge failed")
}
