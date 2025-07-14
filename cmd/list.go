package cmd

import (
	"fmt"
	"slices"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all managed worktrees and their status",
		Long: `List all managed worktrees and their status.

Shows environment variable mappings and indicates sync status for each entry.
Displays which branches are out of sync, lists missing worktrees, and shows orphaned worktrees.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := createInitializedManagerStrict()
			if err != nil {
				return err
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
			table := internal.NewTable([]string{"WORKTREE", "BRANCH", "GIT STATUS", "SYNC STATUS", "PATH"})

			// Get sorted worktree names (tracked first, then ad hoc by creation time desc)
			sortedNames := manager.GetSortedWorktreeNames(worktrees)

			for _, worktreeName := range sortedNames {
				info := worktrees[worktreeName]
				var syncStatus string

				// Check for branch changes
				if change, exists := status.BranchChanges[worktreeName]; exists {
					syncStatus = fmt.Sprintf("OUT_OF_SYNC (%s → %s)", change.OldBranch, change.NewBranch)
				}

				// Check for orphaned worktrees
				if slices.Contains(status.OrphanedWorktrees, worktreeName) {
					syncStatus = "UNTRACKED"
				}

				// Check if this is an untracked worktree (not in gbm.branchconfig.yaml)
				if syncStatus == "" {
					worktreeMapping, err := manager.GetWorktreeMapping()
					if err == nil {
						if _, exists := worktreeMapping[worktreeName]; !exists {
							syncStatus = internal.FormatInfo("UNTRACKED")
						} else {
							syncStatus = internal.FormatSuccess("IN_SYNC")
						}
					} else {
						syncStatus = internal.FormatSuccess("IN_SYNC")
					}
				}

				// Get git status icon
				gitStatusIcon := internal.FormatGitStatus(info.GitStatus)

				branchDisplay := info.CurrentBranch
				if info.ExpectedBranch != "" && info.ExpectedBranch != info.CurrentBranch {
					branchDisplay = fmt.Sprintf("%s (expected: %s)", info.CurrentBranch, info.ExpectedBranch)
				}

				table.AddRow([]string{worktreeName, branchDisplay, gitStatusIcon, syncStatus, info.Path})
			}

			fmt.Fprint(cmd.OutOrStdout(), table.String())
			fmt.Fprintln(cmd.OutOrStdout())

			// Only show sync hint if there are actual sync issues with existing worktrees
			hasExistingSyncIssues := len(status.BranchChanges) > 0 || len(status.OrphanedWorktrees) > 0
			if hasExistingSyncIssues {
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprint(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))
			}

			return nil
		},
	}

	return cmd
}

func init() {
}
