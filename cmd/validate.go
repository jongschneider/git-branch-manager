package cmd

import (
	"fmt"
	"os"

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
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}

		manager, err := internal.NewManager(repoRoot)
		if err != nil {
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
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
				table.AddRow([]string{envVar, branchName, internal.FormatStatusIcon("❌", "ERROR")})
				allValid = false
				continue
			}

			if exists {
				table.AddRow([]string{envVar, branchName, internal.FormatStatusIcon("✅", "VALID")})
			} else {
				table.AddRow([]string{envVar, branchName, internal.FormatStatusIcon("❌", "NOT FOUND")})
				allValid = false
			}
		}

		// Display validation header
		if allValid {
			PrintInfo("%s", internal.FormatStatusIcon("✅", ".envrc validation passed"))
		} else {
			PrintError("%s", internal.FormatStatusIcon("❌", ".envrc validation failed"))
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

