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
		// Get all worktrees (including those created with gbm add)
		worktrees, err := manager.GetAllWorktrees()
		if err != nil {
			return fmt.Errorf("failed to get worktree list: %w", err)
		}

		PrintVerbose("Found %d worktrees to display", len(worktrees))
		if len(worktrees) == 0 {
			return nil
		}

		PrintVerbose("Building worktree list table")
		table := internal.NewTable([]string{"WORKTREE", "BRANCH", "STATUS", "PATH"})

		iconManager := internal.GetGlobalIconManager()

		// Get sorted worktree names (.envrc first, then ad hoc by creation time desc)
		sortedNames := manager.GetSortedWorktreeNames(worktrees)

		for _, worktreeName := range sortedNames {
			info := worktrees[worktreeName]
			statusIcon := iconManager.Success()
			statusText := "OK"

			// Check if this worktree has issues
			if slices.Contains(status.MissingWorktrees, worktreeName) {
				statusIcon = iconManager.Error()
				statusText = "MISSING"
			}

			if change, exists := status.BranchChanges[worktreeName]; exists {
				statusIcon = iconManager.Warning()
				statusText = "OUT_OF_SYNC"
				info.CurrentBranch = change.OldBranch
			}

			if slices.Contains(status.OrphanedWorktrees, worktreeName) {
				statusIcon = iconManager.Orphaned()
				statusText = "UNTRACKED"
			}

			// For worktrees not in .envrc, mark as "UNTRACKED"
			if info.ExpectedBranch == info.CurrentBranch && info.ExpectedBranch != "" {
				// Check if this worktree is actually tracked in .envrc
				envMapping, err := manager.GetEnvMapping()
				if err == nil {
					if _, exists := envMapping[worktreeName]; !exists {
						statusIcon = iconManager.Info()
						statusText = "UNTRACKED"
					}
				}
			}

			branchDisplay := info.CurrentBranch
			if info.ExpectedBranch != "" && info.ExpectedBranch != info.CurrentBranch {
				branchDisplay = fmt.Sprintf("%s (expected: %s)", info.CurrentBranch, info.ExpectedBranch)
			}

			table.AddRow([]string{worktreeName, branchDisplay, internal.FormatStatusIcon(statusIcon, statusText), info.Path})
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
