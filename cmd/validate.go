package cmd

import (
	"fmt"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeValidator.go . worktreeValidator

// worktreeValidator abstracts validation dependencies for unit testing
type worktreeValidator interface {
	GetWorktreeMapping() (map[string]string, error)
	BranchExists(branch string) (bool, error)
}

func newValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate gbm.branchconfig.yaml syntax and branch references",
		Long: `Validate gbm.branchconfig.yaml syntax and branch references.

Checks if referenced branches exist locally or remotely. Useful for CI/CD integration
and ensuring configuration correctness before syncing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := createInitializedManager()
			if err != nil {
				return err
			}

			return handleValidate(manager)
		},
	}

	return cmd
}

// handleValidate performs validation using the provided validator dependency.
// It mirrors the original command logic to keep output and behavior identical.
func handleValidate(validator worktreeValidator) error {
	PrintVerbose("Validating branch references...")

	// Get the mapping to validate
	mapping, err := validator.GetWorktreeMapping()
	if err != nil {
		return err
	}

	// Create table for validation results
	table := internal.NewTable([]string{"WORKTREE", "BRANCH", "STATUS"})

	allValid := true
	for worktreeName, branchName := range mapping {
		exists, err := validator.BranchExists(branchName)
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
		PrintInfo("%s", internal.FormatSuccess("gbm.branchconfig.yaml validation passed"))
	} else {
		PrintError("%s", internal.FormatError("gbm.branchconfig.yaml validation failed"))
	}

	fmt.Println()
	table.Print()

	if !allValid {
		return fmt.Errorf("validation failed - one or more branches do not exist")
	}

	return nil
}
