package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"gbm/internal"

	"slices"

	"github.com/spf13/cobra"
)

var (
	checkFormat   string
	checkExitCode bool
)

type CheckOutput struct {
	InSync    bool     `json:"in_sync"`
	Status    string   `json:"status"`
	Issues    []string `json:"issues,omitempty"`
	Indicator string   `json:"indicator"`
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Quick check for drift between .envrc and worktrees",
	Long: `Quick check for drift between .envrc and worktrees.

Can be used for shell integration or automated checking. Returns non-zero exit code if out of sync.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}
		PrintVerbose("Performing check from repository root: %s", repoRoot)

		manager, err := internal.NewManager(repoRoot)
		if err != nil {
			if checkExitCode {
				os.Exit(1)
			}
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			if checkExitCode {
				os.Exit(1)
			}
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		PrintVerbose("Retrieving sync status for check operation")
		status, err := manager.GetSyncStatus()
		if err != nil {
			if checkExitCode {
				os.Exit(1)
			}
			return err
		}

		output := CheckOutput{
			InSync: status.InSync,
			Issues: []string{},
		}

		if status.InSync {
			output.Status = "in_sync"
			output.Indicator = "✅"
		} else {
			output.Status = "out_of_sync"
			output.Indicator = "⚠️"

			for _, envVar := range status.MissingWorktrees {
				output.Issues = append(output.Issues, fmt.Sprintf("Missing worktree: %s", envVar))
			}

			for envVar, change := range status.BranchChanges {
				output.Issues = append(output.Issues, fmt.Sprintf("Branch change: %s (%s → %s)", envVar, change.OldBranch, change.NewBranch))
			}

			for _, envVar := range status.OrphanedWorktrees {
				output.Issues = append(output.Issues, fmt.Sprintf("Orphaned worktree: %s", envVar))
			}
		}

		PrintVerbose("Formatting output as: %s", checkFormat)
		switch checkFormat {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(output); err != nil {
				return fmt.Errorf("failed to encode JSON: %w", err)
			}
		case "prompt":
			if !status.InSync {
				fmt.Printf("[gbm:%s] ", output.Indicator)
			}
		case "table":
			// Get worktree information for table
			worktrees, err := manager.GetWorktreeList()
			if err != nil {
				return fmt.Errorf("failed to get worktree list: %w", err)
			}
			PrintVerbose("Generating table format for %d worktrees", len(worktrees))

			table := internal.NewTable([]string{"ENV VAR", "BRANCH", "STATUS", "ISSUES"})

			for envVar, info := range worktrees {
				statusText := "OK"
				issues := ""

				// Check for issues
				if slices.Contains(status.MissingWorktrees, envVar) {
					statusText = "MISSING"
					issues = "Worktree missing"
				}

				if change, exists := status.BranchChanges[envVar]; exists {
					statusText = "OUT_OF_SYNC"
					issues = fmt.Sprintf("%s → %s", change.OldBranch, change.NewBranch)
				}

				if slices.Contains(status.OrphanedWorktrees, envVar) {
					statusText = "ORPHANED"
					issues = "Variable removed"
				}

				table.AddRow([]string{envVar, info.ExpectedBranch, statusText, issues})
			}

			table.Print()
		case "text":
			fallthrough
		default:
			if status.InSync {
				PrintInfo("%s All worktrees in sync", output.Indicator)
			} else {
				PrintInfo("%s %d issue(s) detected", output.Indicator, len(output.Issues))
				for _, issue := range output.Issues {
					PrintInfo("  • %s", issue)
				}
			}
		}

		if checkExitCode {
			if status.InSync {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVar(&checkFormat, "format", "text", "output format (prompt|json|text|table)")
	checkCmd.Flags().BoolVar(&checkExitCode, "exit-code", false, "return status code only")
}
