package cmd

import (
	"fmt"
	"os"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var (
	pushAll bool
)

var pushCmd = &cobra.Command{
	Use:   "push [worktree-name]",
	Short: "Push worktree changes to remote",
	Long: `Push changes from a worktree to the remote repository.

Usage:
  gbm push                    # Push current worktree (if in a worktree)
  gbm push <worktree-name>    # Push specific worktree
  gbm push --all              # Push all worktrees

The command will automatically set upstream (-u) if not already set.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoPath, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("not in a git repository: %w", err)
		}

		manager, err := internal.NewManager(repoPath)
		if err != nil {
			return fmt.Errorf("failed to create manager: %w", err)
		}

		if pushAll {
			return handlePushAll(manager)
		}

		if len(args) == 0 {
			return handlePushCurrent(manager, wd)
		}

		return handlePushNamed(manager, args[0])
	},
}

func handlePushAll(manager *internal.Manager) error {
	PrintInfo("Pushing all worktrees...")
	return manager.PushAllWorktrees()
}

func handlePushCurrent(manager *internal.Manager, currentPath string) error {
	// Check if we're in a worktree
	inWorktree, worktreeName, err := manager.IsInWorktree(currentPath)
	if err != nil {
		return fmt.Errorf("failed to check if in worktree: %w", err)
	}

	if !inWorktree {
		return fmt.Errorf("not currently in a worktree. Use 'gbm push <worktree-name>' to push a specific worktree")
	}

	PrintInfo("Pushing current worktree '%s'...", worktreeName)
	return manager.PushWorktree(worktreeName)
}

func handlePushNamed(manager *internal.Manager, worktreeName string) error {
	// Check if worktree exists
	worktrees, err := manager.GetAllWorktrees()
	if err != nil {
		return fmt.Errorf("failed to get worktrees: %w", err)
	}

	if _, exists := worktrees[worktreeName]; !exists {
		return fmt.Errorf("worktree '%s' does not exist", worktreeName)
	}

	PrintInfo("Pushing worktree '%s'...", worktreeName)
	return manager.PushWorktree(worktreeName)
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVar(&pushAll, "all", false, "Push all worktrees")

	// Add completion for worktree names
	pushCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getWorktreeNames(), cobra.ShellCompDirectiveNoFileComp
	}
}
