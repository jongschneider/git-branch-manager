package cmd

import (
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize all worktrees with current gbm.branchconfig.yaml definitions",
		Long: `Synchronize all worktrees with current gbm.branchconfig.yaml definitions.

Fetches from remote first, then creates missing worktrees for new worktree configurations,
updates existing worktrees if branch references have changed, and optionally removes orphaned worktrees.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			syncDryRun, _ := cmd.Flags().GetBool("dry-run")
			syncForce, _ := cmd.Flags().GetBool("force")

			manager, err := createInitializedManagerStrict()
			if err != nil {
				return err
			}

			if syncDryRun {
				iconManager := internal.GetGlobalIconManager()
				PrintInfo("%s", internal.FormatStatusIcon(iconManager.DryRun(), "Dry run mode - showing what would be changed:"))
				status, err := manager.GetSyncStatus()
				if err != nil {
					return err
				}

				if status.InSync {
					PrintInfo("%s", internal.FormatSuccess("All worktrees are in sync"))
					return nil
				}

				if len(status.MissingWorktrees) > 0 {
					iconManager := internal.GetGlobalIconManager()
					PrintInfo("%s", internal.FormatStatusIcon(iconManager.Missing(), "Missing worktrees:"))
					for _, envVar := range status.MissingWorktrees {
						PrintInfo("  • %s", envVar)
					}
				}

				if len(status.BranchChanges) > 0 {
					iconManager := internal.GetGlobalIconManager()
					PrintInfo("%s", internal.FormatStatusIcon(iconManager.Changes(), "Branch changes needed:"))
					for envVar, change := range status.BranchChanges {
						PrintInfo("  • %s: %s → %s", envVar, change.OldBranch, change.NewBranch)
					}
				}

				if len(status.OrphanedWorktrees) > 0 {
					iconManager := internal.GetGlobalIconManager()
					PrintInfo("%s", internal.FormatStatusIcon(iconManager.Orphaned(), "Orphaned worktrees (use --force to remove):"))
					for _, envVar := range status.OrphanedWorktrees {
						PrintInfo("  • %s", envVar)
					}
				}

				return nil
			}

			PrintVerbose("Synchronizing worktrees (force=%v)", syncForce)

			// Create confirmation function for destructive operations
			// Only prompt when force is used (to confirm destructive actions)
			var confirmFunc internal.ConfirmationFunc
			if syncForce {
				confirmFunc = func(message string) bool {
					fmt.Print(message + " [y/N]: ")
					var response string
					fmt.Scanln(&response)
					return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
				}
			}

			if err := manager.SyncWithConfirmation(syncDryRun, syncForce, confirmFunc); err != nil {
				return err
			}

			PrintInfo("%s", internal.FormatSuccess("Successfully synchronized worktrees"))
			return nil
		},
	}

	cmd.Flags().Bool("dry-run", false, "show what would be changed without making changes")
	cmd.Flags().Bool("force", false, "skip confirmation prompts and remove orphaned worktrees")

	return cmd
}
