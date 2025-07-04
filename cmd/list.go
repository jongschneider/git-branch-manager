package cmd

import (
	"fmt"
	"os"
	"slices"

	"gbm/internal"

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

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}
		PrintVerbose("Listing worktrees from repository root: %s", repoRoot)

		manager, err := internal.NewManager(repoRoot)
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

		iconManager := internal.GetGlobalIconManager()
		for envVar, info := range worktrees {
			statusIcon := iconManager.Success()
			statusText := "OK"

			// Check if this worktree has issues
			if slices.Contains(status.MissingWorktrees, envVar) {
				statusIcon = iconManager.Error()
				statusText = "MISSING"
			}

			if change, exists := status.BranchChanges[envVar]; exists {
				statusIcon = iconManager.Warning()
				statusText = "OUT_OF_SYNC"
				info.CurrentBranch = change.OldBranch
			}

			if slices.Contains(status.OrphanedWorktrees, envVar) {
				statusIcon = iconManager.Orphaned()
				statusText = "ORPHANED"
			}

			table.AddRow([]string{envVar, info.ExpectedBranch, internal.FormatStatusIcon(statusIcon, statusText), info.Path})
		}

		table.Print()

		fmt.Println()
		if !status.InSync {
			PrintInfo("%s", internal.FormatInfo("Run 'gbm status' for detailed information or 'gbm sync' to fix issues"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
