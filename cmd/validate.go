package cmd

import (
	"fmt"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate .envrc syntax and branch references",
	Long: `Validate .envrc syntax and branch references.

Checks if referenced branches exist locally or remotely. Useful for CI/CD integration
and ensuring configuration correctness before syncing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := createInitializedManagerStrict()
		if err != nil {
			return err
		}

		PrintVerbose("Validating branch references...")

		// Get the mapping to validate
		mapping, err := manager.GetEnvMapping()
		if err != nil {
			return err
		}

		// Create table for validation results
		table := internal.NewTable([]string{"ENV VARIABLE", "BRANCH", "STATUS"})

		allValid := true
		for envVar, branchName := range mapping {
			exists, err := manager.BranchExists(branchName)
			if err != nil {
				table.AddRow([]string{envVar, branchName, internal.FormatError("ERROR")})
				allValid = false
				continue
			}

			if exists {
				table.AddRow([]string{envVar, branchName, internal.FormatSuccess("VALID")})
			} else {
				table.AddRow([]string{envVar, branchName, internal.FormatError("NOT FOUND")})
				allValid = false
			}
		}

		// Display validation header
		if allValid {
			PrintInfo("%s", internal.FormatSuccess(".envrc validation passed"))
		} else {
			PrintError("%s", internal.FormatError(".envrc validation failed"))
		}

		fmt.Println()
		table.Print()

		if !allValid {
			return fmt.Errorf("validation failed - one or more branches do not exist")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
