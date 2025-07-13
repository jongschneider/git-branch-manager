package cmd

import (
	"fmt"

	"gbm/internal"

	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate .gbm.config.yaml syntax and branch references",
		Long: `Validate .gbm.config.yaml syntax and branch references.

Checks if referenced branches exist locally or remotely. Useful for CI/CD integration
and ensuring configuration correctness before syncing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := createInitializedManagerStrictWithCmd(cmd)
			if err != nil {
				return err
			}

		PrintVerbose("Validating branch references...")

		// Get the mapping to validate
		mapping, err := manager.GetWorktreeMapping()
		if err != nil {
			return err
		}

		// Create table for validation results
		table := internal.NewTable([]string{"WORKTREE", "BRANCH", "STATUS"})

		allValid := true
		for worktreeName, branchName := range mapping {
			exists, err := manager.BranchExists(branchName)
			if err != nil {
				table.AddRow([]string{worktreeName, branchName, internal.FormatError("ERROR")})
				allValid = false
				continue
			}

			if exists {
				table.AddRow([]string{worktreeName, branchName, internal.FormatSuccess("VALID")})
			} else {
				table.AddRow([]string{worktreeName, branchName, internal.FormatError("NOT FOUND")})
				allValid = false
			}
		}

		// Display validation header
		if allValid {
			PrintInfo("%s", internal.FormatSuccess(".gbm.config.yaml validation passed"))
		} else {
			PrintError("%s", internal.FormatError(".gbm.config.yaml validation failed"))
		}

		fmt.Println()
		table.Print()

		if !allValid {
			return fmt.Errorf("validation failed - one or more branches do not exist")
		}

		return nil
	},
	}

	return cmd
}

func init() {
}
