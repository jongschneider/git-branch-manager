package cmd

import (
	"fmt"
	"os"

	"gbm/internal"

	"slices"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed worktrees and their status",
	Long: `List all managed worktrees and their status.

Shows environment variable mappings and indicates sync status for each entry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		PrintVerbose("Listing worktrees from working directory: %s", wd)

		manager, err := internal.NewManager(wd)
		if err != nil {
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		PrintVerbose("Retrieving sync status for list operation")
		status, err := manager.GetSyncStatus()
		if err != nil {
			return err
		}

		PrintVerbose("Fetching detailed worktree information")
		// Get worktree information
		worktrees, err := manager.GetWorktreeList()
		if err != nil {
			return fmt.Errorf("failed to get worktree list: %w", err)
		}

		PrintVerbose("Found %d worktrees to display", len(worktrees))
		if len(worktrees) == 0 {
			return nil
		}

		PrintVerbose("Building worktree list table")
		table := internal.NewTable([]string{"ENV VAR", "BRANCH", "STATUS", "PATH"})

		for envVar, info := range worktrees {
			statusIcon := "✅"
			statusText := "OK"

			// Check if this worktree has issues
			if slices.Contains(status.MissingWorktrees, envVar) {
				statusIcon = "❌"
				statusText = "MISSING"
			}

			if change, exists := status.BranchChanges[envVar]; exists {
				statusIcon = "⚠️"
				statusText = "OUT_OF_SYNC"
				info.CurrentBranch = change.OldBranch
			}

			if slices.Contains(status.OrphanedWorktrees, envVar) {
				statusIcon = "🗑️"
				statusText = "ORPHANED"
			}

			table.AddRow([]string{envVar, info.ExpectedBranch, statusIcon + " " + statusText, info.Path})
		}

		table.Print()

		fmt.Println()
		if !status.InSync {
			PrintInfo("💡 Run 'gbm status' for detailed information or 'gbm sync' to fix issues")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
