package cmd

import (
	"errors"
	"fmt"
	"os"

	"gbm/internal"

	"github.com/spf13/cobra"
)

// worktreePusher interface for push operations to enable unit testing
//
//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreePusher.go . worktreePusher
type worktreePusher interface {
	PushAllWorktrees() error
	PushWorktree(worktreeName string) error
	IsInWorktree(currentPath string) (bool, string, error)
	GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
}

func newPushCommand() *cobra.Command {
	cmd := &cobra.Command{
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
			pushAll, _ := cmd.Flags().GetBool("all")

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

			if pushAll {
				return handlePushAll(manager)
			}

			if len(args) == 0 {
				return handlePushCurrent(manager, wd)
			}

			return handlePushNamed(manager, args[0])
		},
	}

	cmd.Flags().Bool("all", false, "Push all worktrees")

	// Add completion for worktree names
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getWorktreeNames(), cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func handlePushAll(pusher worktreePusher) error {
	PrintInfo("Pushing all worktrees...")
	return pusher.PushAllWorktrees()
}

func handlePushCurrent(pusher worktreePusher, currentPath string) error {
	// Check if we're in a worktree
	inWorktree, worktreeName, err := pusher.IsInWorktree(currentPath)
	if err != nil {
		return fmt.Errorf("failed to check if in worktree: %w", err)
	}

	if !inWorktree {
		return fmt.Errorf("not currently in a worktree. Use 'gbm push <worktree-name>' to push a specific worktree")
	}

	PrintInfo("Pushing current worktree '%s'...", worktreeName)
	return pusher.PushWorktree(worktreeName)
}

func handlePushNamed(pusher worktreePusher, worktreeName string) error {
	// Check if worktree exists
	worktrees, err := pusher.GetAllWorktrees()
	if err != nil {
		return fmt.Errorf("failed to get worktrees: %w", err)
	}

	if _, exists := worktrees[worktreeName]; !exists {
		return fmt.Errorf("worktree '%s' does not exist", worktreeName)
	}

	PrintInfo("Pushing worktree '%s'...", worktreeName)
	return pusher.PushWorktree(worktreeName)
}
