package cmd

import (
	"gbm/internal"

	"github.com/spf13/cobra"
)

var (
	syncDryRun bool
	syncForce  bool
	syncFetch  bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize all worktrees with current .gbm.config.yaml definitions",
	Long: `Synchronize all worktrees with current .gbm.config.yaml definitions.

Creates missing worktrees for new worktree configurations, updates existing worktrees
if branch references have changed, and optionally removes orphaned worktrees.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		PrintVerbose("Synchronizing worktrees (force=%v, fetch=%v)", syncForce, syncFetch)
		if err := manager.Sync(syncDryRun, syncForce, syncFetch); err != nil {
			return err
		}

		PrintInfo("%s", internal.FormatSuccess("Successfully synchronized worktrees"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "show what would be changed without making changes")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "skip confirmation prompts and remove orphaned worktrees")
	syncCmd.Flags().BoolVar(&syncFetch, "fetch", false, "update remote tracking before sync")
}
