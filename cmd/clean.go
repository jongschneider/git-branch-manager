package cmd

import (
	"fmt"
	"os"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var (
	cleanForce bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove orphaned worktrees (interactive by default)",
	Long: `Remove orphaned worktrees that are no longer defined in .envrc.

By default, this command is interactive and will prompt for confirmation before removing
each orphaned worktree. Use --force to skip confirmations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}
		PrintVerbose("Starting cleanup from repository root: %s", repoRoot)

		manager, err := internal.NewManager(repoRoot)
		if err != nil {
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		PrintVerbose("Retrieving sync status to identify orphaned worktrees")
		status, err := manager.GetSyncStatus()
		if err != nil {
			return err
		}

		PrintVerbose("Found %d orphaned worktrees to process", len(status.OrphanedWorktrees))
		if len(status.OrphanedWorktrees) == 0 {
			PrintInfo("%s", internal.FormatSuccess("No orphaned worktrees found"))
			return nil
		}

		iconManager := internal.GetGlobalIconManager()
		PrintInfo("%s", internal.FormatStatusIcon(iconManager.Orphaned(), fmt.Sprintf("Found %d orphaned worktree(s):", len(status.OrphanedWorktrees))))
		for _, envVar := range status.OrphanedWorktrees {
			PrintInfo("  • %s", envVar)
		}
		fmt.Println()

		PrintVerbose("Cleanup mode: force=%v", cleanForce)
		if cleanForce {
			PrintInfo("%s", internal.FormatStatusIcon("🔥", "Force mode enabled - removing all orphaned worktrees..."))
		} else {
			PrintInfo("%s", internal.FormatInfo("Interactive mode - you will be prompted for each worktree"))
		}

		PrintVerbose("Initiating cleanup of orphaned worktrees")
		if err := manager.CleanOrphaned(cleanForce); err != nil {
			return err
		}

		PrintInfo("%s", internal.FormatSuccess("Orphaned worktree cleanup completed"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&cleanForce, "force", false, "force removal without confirmation prompts")
}
