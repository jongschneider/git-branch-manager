package cmd

import (
	"errors"
	"fmt"
	"os"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [worktree-name]",
		Short: "Pull worktree changes from remote",
		Long: `Pull changes from the remote repository to a worktree.

Usage:
  gbm pull                    # Pull current worktree (if in a worktree)
  gbm pull <worktree-name>    # Pull specific worktree
  gbm pull --all              # Pull all worktrees`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pullAll, _ := cmd.Flags().GetBool("all")

			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return err
				}

				PrintVerbose("%v", err)
			}

			if pullAll {
				return handlePullAll(manager)
			}

			if len(args) == 0 {
				return handlePullCurrent(manager, wd)
			}

			return handlePullNamed(manager, args[0])
		},
	}

	cmd.Flags().Bool("all", false, "Pull all worktrees")

	// Add completion for worktree names
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getWorktreeNames(), cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func handlePullAll(manager *internal.Manager) error {
	PrintInfo("Pulling all worktrees...")
	return manager.PullAllWorktrees()
}

func handlePullCurrent(manager *internal.Manager, currentPath string) error {
	// Check if we're in a worktree
	inWorktree, worktreeName, err := manager.IsInWorktree(currentPath)
	if err != nil {
		return fmt.Errorf("failed to check if in worktree: %w", err)
	}

	if !inWorktree {
		return fmt.Errorf("not currently in a worktree. Use 'gbm pull <worktree-name>' to pull a specific worktree")
	}

	PrintInfo("Pulling current worktree '%s'...", worktreeName)
	return manager.PullWorktree(worktreeName)
}

func handlePullNamed(manager *internal.Manager, worktreeName string) error {
	// Check if worktree exists
	worktrees, err := manager.GetAllWorktrees()
	if err != nil {
		return fmt.Errorf("failed to get worktrees: %w", err)
	}

	if _, exists := worktrees[worktreeName]; !exists {
		return fmt.Errorf("worktree '%s' does not exist", worktreeName)
	}

	PrintInfo("Pulling worktree '%s'...", worktreeName)
	return manager.PullWorktree(worktreeName)
}

