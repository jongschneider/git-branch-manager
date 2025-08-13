package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeAdder.go . worktreeAdder

// worktreeAdder interface abstracts the Manager operations needed for adding worktrees
type worktreeAdder interface {
	AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
	GetDefaultBranch() (string, error)
	BranchExists(branch string) (bool, error)
	GetJiraIssues() ([]internal.JiraIssue, error)
	GenerateBranchFromJira(jiraKey string) (string, error)
}

// WorktreeArgs represents the resolved arguments for creating a worktree
type WorktreeArgs struct {
	WorktreeName       string
	BranchName         string
	NewBranch          bool
	ResolvedBaseBranch string
}

// ArgsResolver handles the complex logic of resolving command arguments
type ArgsResolver struct {
	manager worktreeAdder
}

// ResolveArgs processes command arguments and flags to determine worktree parameters
func (r *ArgsResolver) ResolveArgs(cmdArgs []string, newBranchFlag bool) (*WorktreeArgs, error) {
	if len(cmdArgs) == 0 {
		return nil, fmt.Errorf("worktree name is required")
	}

	args := &WorktreeArgs{
		WorktreeName: cmdArgs[0],
		NewBranch:    newBranchFlag,
	}

	// Resolve branch name
	branchName, err := r.resolveBranchName(cmdArgs, newBranchFlag, args.WorktreeName)
	if err != nil {
		return nil, err
	}

	args.BranchName = branchName

	// Resolve base branch
	var baseBranch string
	if len(cmdArgs) > 2 {
		baseBranch = cmdArgs[2]
	}

	resolvedBaseBranch, err := r.resolveBaseBranch(newBranchFlag, baseBranch)
	if err != nil {
		return nil, err
	}
	args.ResolvedBaseBranch = resolvedBaseBranch

	return args, nil
}

// resolveBranchName determines the branch name based on arguments and flags
func (r *ArgsResolver) resolveBranchName(cmdArgs []string, newBranchFlag bool, worktreeName string) (string, error) {
	// Handle direct specification
	if len(cmdArgs) > 1 {
		return cmdArgs[1], nil
	}

	if newBranchFlag {
		// Generate branch name from worktree name
		return generateBranchName(worktreeName, r.manager), nil
	}

	if internal.IsJiraKey(worktreeName) {
		// Auto-suggest branch name for JIRA keys
		suggestedBranch := generateBranchName(worktreeName, r.manager)
		return "", fmt.Errorf("branch name required. Suggested: %s\n\nTry: gbm add %s %s -b", suggestedBranch, worktreeName, suggestedBranch)
	}

	return "", fmt.Errorf("branch name required when not creating new branch (use -b to create new branch)")
}

// resolveBaseBranch determines the base branch for new branch creation
func (r *ArgsResolver) resolveBaseBranch(newBranchFlag bool, baseBranch string) (string, error) {
	if !newBranchFlag {
		return "", nil
	}

	if baseBranch == "" {
		// Use default branch
		PrintInfo("Using default base branch")
		return r.manager.GetDefaultBranch()
	}

	// Validate that the base branch exists
	exists, err := r.manager.BranchExists(baseBranch)
	if err != nil {
		return "", fmt.Errorf("failed to check if base branch exists: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("base branch '%s' does not exist", baseBranch)
	}

	return baseBranch, nil
}

func newAddCommand(manager worktreeAdder) *cobra.Command {
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
			if manager == nil {
				return fmt.Errorf("manager not available - ensure you're in a git repository with gbm.branchconfig.yaml")
			}

			newBranch, _ := cmd.Flags().GetBool("new-branch")

			resolver := &ArgsResolver{manager: manager}
			worktreeArgs, err := resolver.ResolveArgs(args, newBranch)
			if err != nil {
				return err
			}

			PrintInfo("Adding worktree '%s' on branch '%s'", worktreeArgs.WorktreeName, worktreeArgs.BranchName)

			if err := manager.AddWorktree(
				worktreeArgs.WorktreeName,
				worktreeArgs.BranchName,
				worktreeArgs.NewBranch,
				worktreeArgs.ResolvedBaseBranch,
			); err != nil {
				return fmt.Errorf("failed to add worktree: %w", err)
			}

			PrintInfo("Worktree '%s' added successfully", worktreeArgs.WorktreeName)

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
			jiraIssues, err := manager.GetJiraIssues()
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

				branchName, err := manager.GenerateBranchFromJira(worktreeName)
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

func generateBranchName(worktreeName string, manager worktreeAdder) string {
	// Check if this is a JIRA key first
	if internal.IsJiraKey(worktreeName) {
		branchName, err := manager.GenerateBranchFromJira(worktreeName)
		if err != nil {
			PrintVerbose("Failed to generate branch name from JIRA issue %s: %v", worktreeName, err)
			PrintInfo("Could not fetch JIRA issue details. Using default branch name format.")
			// Generate a basic branch name from the JIRA key
			return fmt.Sprintf("feature/%s", strings.ToLower(worktreeName))
		}
		PrintInfo("Auto-generated branch name from JIRA issue: %s", branchName)
		return branchName
	}

	// Convert worktree name to a valid branch name
	branchName := strings.ReplaceAll(worktreeName, " ", "-")
	branchName = strings.ReplaceAll(branchName, "_", "-")
	branchName = strings.ToLower(branchName)

	// Add feature/ prefix if not already present
	// Regex matches: "feature/foo", "bugfix/bar", "hotfix/baz", "docs/update"
	// Regex does NOT match: "foo", "bar-baz", "/invalid", ""
	hasPrefixPattern := regexp.MustCompile(`^[^/]+/`).MatchString(branchName)
	if !hasPrefixPattern {
		branchName = "feature/" + branchName
	}

	return branchName
}
