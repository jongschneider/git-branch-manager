package cmd

import (
	"errors"
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeRemover.go . worktreeRemover

// worktreeRemover interface abstracts the Manager operations needed for removing worktrees
type worktreeRemover interface {
	GetWorktreePath(worktreeName string) (string, error)
	GetWorktreeStatus(worktreePath string) (*internal.GitStatus, error)
	RemoveWorktree(worktreeName string) error
	GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
}

// confirmationFunc is a function type for confirming actions
type confirmationFunc func(worktreeName string) bool

// defaultConfirmation is the default confirmation function that prompts the user
func defaultConfirmation(worktreeName string) bool {
	fmt.Printf("Are you sure you want to remove worktree '%s'? [y/N]: ", worktreeName)
	var response string
	_, _ = fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

// handleRemove handles the removal of a worktree with the specified options
func handleRemove(remover worktreeRemover, worktreeName string, force bool) error {
	return handleRemoveWithConfirmation(remover, worktreeName, force, defaultConfirmation)
}

// handleRemoveWithConfirmation handles the removal with a custom confirmation function
func handleRemoveWithConfirmation(remover worktreeRemover, worktreeName string, force bool, confirm confirmationFunc) error {
	// Check if worktree exists
	worktreePath, err := remover.GetWorktreePath(worktreeName)
	if err != nil {
		return fmt.Errorf("worktree '%s' not found: %w", worktreeName, err)
	}

	// Check if worktree has uncommitted changes (unless force is used)
	if !force {
		gitStatus, err := remover.GetWorktreeStatus(worktreePath)
		if err != nil {
			return fmt.Errorf("failed to check worktree status: %w", err)
		}

		if gitStatus.HasChanges() {
			return fmt.Errorf("worktree '%s' has uncommitted changes. Use --force to remove anyway", worktreeName)
		}
	}

	// Confirm removal (unless force is used)
	if !force {
		if !confirm(worktreeName) {
			PrintInfo("Removal cancelled")
			return nil
		}
	}

	// Remove the worktree
	if err := remover.RemoveWorktree(worktreeName); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	PrintInfo("Worktree '%s' removed successfully", worktreeName)
	return nil
}


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

			return handleRemove(manager, worktreeName, force)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force removal even if worktree has uncommitted changes")

	// Add completion for worktree names
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getWorktreeCompletionsWithManager(), cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}
