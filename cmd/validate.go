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
		if err := manager.ValidateEnvrc(); err != nil {
			PrintError("%v", err)
			return fmt.Errorf("validation failed")
		}

		PrintInfo("✅ .envrc validation passed")
		PrintInfo("📋 Configuration summary:")

		// Get the mapping to show what was validated
		mapping, err := manager.GetEnvMapping()
		if err != nil {
			return err
		}

		for envVar, branch := range mapping {
			PrintInfo("  • %s → %s", envVar, branch)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

