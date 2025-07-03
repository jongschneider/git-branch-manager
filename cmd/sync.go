package cmd

import (
	"fmt"
	"os"

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
	Short: "Synchronize all worktrees with current .envrc definitions",
	Long: `Synchronize all worktrees with current .envrc definitions.

Creates missing worktrees for new environment variables, updates existing worktrees
if branch references have changed, and optionally removes orphaned worktrees.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager, err := internal.NewManager(wd)
		if err != nil {
			return err
		}


		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		if syncDryRun {
			PrintInfo("🔍 Dry run mode - showing what would be changed:")
			status, err := manager.GetSyncStatus()
			if err != nil {
				return err
			}

			if status.InSync {
				PrintInfo("✅ All worktrees are in sync")
				return nil
			}

			if len(status.MissingWorktrees) > 0 {
				PrintInfo("📁 Missing worktrees:")
				for _, envVar := range status.MissingWorktrees {
					PrintInfo("  • %s", envVar)
				}
			}

			if len(status.BranchChanges) > 0 {
				PrintInfo("🔄 Branch changes needed:")
				for envVar, change := range status.BranchChanges {
					PrintInfo("  • %s: %s → %s", envVar, change.OldBranch, change.NewBranch)
				}
			}

			if len(status.OrphanedWorktrees) > 0 {
				PrintInfo("🗑️  Orphaned worktrees (use --force to remove):")
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

		PrintInfo("✅ Successfully synchronized worktrees")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "show what would be changed without making changes")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "skip confirmation prompts and remove orphaned worktrees")
	syncCmd.Flags().BoolVar(&syncFetch, "fetch", false, "update remote tracking before sync")
}