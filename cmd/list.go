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

Shows environment variable mappings and indicates sync status for each entry.
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
		table := internal.NewTable([]string{"WORKTREE", "BRANCH", "GIT STATUS", "SYNC STATUS", "PATH"})

		// Get sorted worktree names (.envrc first, then ad hoc by creation time desc)
		sortedNames := manager.GetSortedWorktreeNames(worktrees)

		for _, worktreeName := range sortedNames {
			info := worktrees[worktreeName]
			var syncStatus string

			// Check for missing worktrees
			if slices.Contains(status.MissingWorktrees, worktreeName) {
				syncStatus = "MISSING"
			}

			// Check for branch changes
			if change, exists := status.BranchChanges[worktreeName]; exists {
				syncStatus = fmt.Sprintf("OUT_OF_SYNC (%s → %s)", change.OldBranch, change.NewBranch)
			}

			// Check for orphaned worktrees
			if slices.Contains(status.OrphanedWorktrees, worktreeName) {
				syncStatus = "UNTRACKED"
			}

			// Check if this is an untracked worktree (not in .envrc)
			if syncStatus == "" {
				envMapping, err := manager.GetEnvMapping()
				if err == nil {
					if _, exists := envMapping[worktreeName]; !exists {
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

		table.Print()

		if !status.InSync {
			fmt.Println()
			PrintInfo("%s", internal.FormatInfo("Run 'gbm sync' to synchronize changes"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
