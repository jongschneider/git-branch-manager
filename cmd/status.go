package cmd

import (
	"fmt"
	"os"
	"slices"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current sync status between .envrc and actual worktrees",
	Long: `Show current sync status between .envrc and actual worktrees.

Displays which branches are out of sync, lists missing worktrees, and shows orphaned worktrees.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}
		PrintVerbose("Checking status from repository root: %s", repoRoot)

		manager, err := internal.NewManager(repoRoot)
		if err != nil {
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		PrintVerbose("Retrieving sync status from manager")
		status, err := manager.GetSyncStatus()
		if err != nil {
			return err
		}

		PrintVerbose("Fetching worktree list for status display")
		// Get worktree information for table
		worktrees, err := manager.GetWorktreeList()
		if err != nil {
			return fmt.Errorf("failed to get worktree list: %w", err)
		}

		PrintVerbose("Building status table with %d worktrees", len(worktrees))
		table := internal.NewTable([]string{"ENV VAR", "BRANCH", "GIT STATUS", "SYNC STATUS"})

		// Add rows for each worktree
		for envVar, info := range worktrees {
			var syncStatus string

			// Check for missing worktrees
			if slices.Contains(status.MissingWorktrees, envVar) {
				syncStatus = "MISSING"
			}

			// Check for branch changes
			if change, exists := status.BranchChanges[envVar]; exists {
				syncStatus = fmt.Sprintf("OUT_OF_SYNC (%s → %s)", change.OldBranch, change.NewBranch)
			}

			// Check for orphaned worktrees
			if slices.Contains(status.OrphanedWorktrees, envVar) {
				syncStatus = "ORPHANED"
			}

			// Default to in sync if no issues
			if syncStatus == "" {
				syncStatus = internal.FormatStatusIcon("✅", "IN_SYNC")
			}

			// Get git status icon
			gitStatusIcon := internal.FormatGitStatus(info.GitStatus)

			table.AddRow([]string{envVar, info.ExpectedBranch, gitStatusIcon, syncStatus})
		}

		table.Print()

		if !status.InSync {
			fmt.Println()
			PrintInfo("%s", internal.FormatStatusIcon("💡", "Run 'gbm sync' to synchronize changes"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
