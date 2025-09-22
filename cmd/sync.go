package cmd

import (
	"fmt"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeSyncer.go . worktreeSyncer

// worktreeSyncer interface abstracts the Manager operations needed for sync operations
type worktreeSyncer interface {
	GetSyncStatus() (*internal.SyncStatus, error)
	SyncWithConfirmation(dryRun, force bool, removeOrphans bool, confirmFunc internal.ConfirmationFunc) error
}

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize all worktrees with current gbm.branchconfig.yaml definitions",
		Long: `Synchronize all worktrees with current gbm.branchconfig.yaml definitions.

Fetches from remote first, then creates missing worktrees for new worktree configurations,
updates existing worktrees if branch references have changed. Use --remove-orphans to also
remove untracked worktrees not defined in the configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			syncDryRun, _ := cmd.Flags().GetBool("dry-run")
			syncForce, _ := cmd.Flags().GetBool("force")
			removeOrphans, _ := cmd.Flags().GetBool("remove-orphans")

			manager, err := createInitializedManager()
			if err != nil {
				return err
			}

			if syncDryRun {
				return handleSyncDryRun(manager, removeOrphans)
			}

			return handleSync(manager, syncForce, removeOrphans)
		},
	}

	cmd.Flags().Bool("dry-run", false, "show what would be changed without making changes")
	cmd.Flags().Bool("force", false, "skip confirmation prompts for sync operations")
	cmd.Flags().Bool("remove-orphans", false, "remove untracked worktrees not in gbm.branchconfig.yaml")

	return cmd
}

func handleSyncDryRun(syncer worktreeSyncer, removeOrphans bool) error {
	iconManager := internal.GetGlobalIconManager()
	PrintInfo("%s", internal.FormatStatusIcon(iconManager.DryRun(), "Dry run mode - showing what would be changed:"))
	status, err := syncer.GetSyncStatus()
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

	if len(status.WorktreePromotions) > 0 {
		iconManager := internal.GetGlobalIconManager()
		PrintInfo("%s", internal.FormatStatusIcon(iconManager.Changes(), "Worktree promotions (destructive):"))
		for _, promotion := range status.WorktreePromotions {
			PrintInfo("  • %s (%s) will be promoted to %s", promotion.SourceWorktree, promotion.Branch, promotion.TargetWorktree)
			PrintInfo("    1. Worktree %s (%s) will be removed", promotion.TargetWorktree, promotion.TargetBranch)
			PrintInfo("    2. Worktree %s (%s) will be moved to %s", promotion.SourceWorktree, promotion.Branch, promotion.TargetWorktree)
		}
	}

	if removeOrphans && len(status.OrphanedWorktrees) > 0 {
		iconManager := internal.GetGlobalIconManager()
		PrintInfo("%s", internal.FormatStatusIcon(iconManager.Orphaned(), "Orphaned worktrees (will be removed):"))
		for _, envVar := range status.OrphanedWorktrees {
			PrintInfo("  • %s", envVar)
		}
	}

	return nil
}

func handleSync(syncer worktreeSyncer, force bool, removeOrphans bool) error {
	PrintVerbose("Synchronizing worktrees (force=%v)", force)

	// Create confirmation function for destructive operations
	// Always provide confirmation for promotions; only for orphaned worktrees when force is used
	confirmFunc := func(message string) bool {
		fmt.Print(message + " [y/N]: ")
		var response string
		_, _ = fmt.Scanln(&response)
		return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
	}

	if err := syncer.SyncWithConfirmation(false, force, removeOrphans, confirmFunc); err != nil {
		return err
	}

	PrintInfo("%s", internal.FormatSuccess("Successfully synchronized worktrees"))
	return nil
}
