package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <worktree-name>",
		Short: "Remove a worktree",
		Long: `Remove a worktree and clean up its directory.

This command removes the specified worktree and its associated directory.
If the worktree contains uncommitted changes, use --force to remove anyway.

Examples:
  gbm remove FEATURE-123
  gbm remove FEATURE-123 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			worktreeName := args[0]

			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return err
				}

				PrintVerbose("%v", err)
			}

			// Check if worktree exists
			worktreePath, err := manager.GetWorktreePath(worktreeName)
			if err != nil {
				return fmt.Errorf("worktree '%s' not found: %w", worktreeName, err)
			}

			// Check if worktree has uncommitted changes (unless force is used)
			if !force {
				gitStatus, err := manager.GetWorktreeStatus(worktreePath)
				if err != nil {
					return fmt.Errorf("failed to check worktree status: %w", err)
				}

				if gitStatus.HasChanges() {
					return fmt.Errorf("worktree '%s' has uncommitted changes. Use --force to remove anyway", worktreeName)
				}
			}

			// Confirm removal (unless force is used)
			if !force {
				fmt.Printf("Are you sure you want to remove worktree '%s'? [y/N]: ", worktreeName)
				var response string
				_, _ = fmt.Scanln(&response)
				if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
					PrintInfo("Removal cancelled")
					return nil
				}
			}

			// Remove the worktree
			if err := manager.RemoveWorktree(worktreeName); err != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}

			PrintInfo("Worktree '%s' removed successfully", worktreeName)
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force removal even if worktree has uncommitted changes")

	// Add completion for worktree names
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// Create manager
			manager, err := createInitializedManager()
			if err != nil {
				if !errors.Is(err, ErrLoadGBMConfig) {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}

				PrintVerbose("%v", err)
			}

			// Get all worktrees
			worktrees, err := manager.GetAllWorktrees()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completions []string
			for worktreeName := range worktrees {
				completions = append(completions, worktreeName)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}
