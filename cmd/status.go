package cmd

import (
	"fmt"
	"os"

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

		manager, err := internal.NewManager(wd)
		if err != nil {
			return err
		}

		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		status, err := manager.GetSyncStatus()
		if err != nil {
			return err
		}

		// Get worktree information for table
		worktrees, err := manager.GetWorktreeList()
		if err != nil {
			return fmt.Errorf("failed to get worktree list: %w", err)
		}

		table := internal.NewTable([]string{"ENV VAR", "BRANCH", "GIT STATUS", "SYNC STATUS"})

		// Add rows for each worktree
		for envVar, info := range worktrees {
			var syncStatus string

			// Check for missing worktrees
			for _, missing := range status.MissingWorktrees {
				if missing == envVar {
					syncStatus = "MISSING"
					break
				}
			}

			// Check for branch changes
			if change, exists := status.BranchChanges[envVar]; exists {
				syncStatus = fmt.Sprintf("OUT_OF_SYNC (%s → %s)", change.OldBranch, change.NewBranch)
			}

			// Check for orphaned worktrees
			for _, orphaned := range status.OrphanedWorktrees {
				if orphaned == envVar {
					syncStatus = "ORPHANED"
					break
				}
			}

			// Default to in sync if no issues
			if syncStatus == "" {
				syncStatus = "✅ IN_SYNC"
			}

			// Get git status icon
			gitStatusIcon := manager.GetStatusIcon(info.GitStatus)

			table.AddRow([]string{envVar, info.ExpectedBranch, gitStatusIcon, syncStatus})
		}

		table.Print()

		if !status.InSync {
			fmt.Println()
			PrintInfo("💡 Run 'gbm sync' to synchronize changes")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

