package cmd

import (
	"fmt"
	"os"

	"gbm/internal"
	"github.com/spf13/cobra"
)

var (
	initForce bool
	initFetch bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the current Git repository for branch management",
	Long: `Initialize the current Git repository for branch management.

Creates initial worktree structure based on .envrc file.
Validates that the repository is a valid Git repo and creates .gbm/ directory for metadata storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager, err := internal.NewManager(wd)
		if err != nil {
			return err
		}

		PrintVerbose("Loading .envrc configuration from: %s", GetConfigPath())
		if err := manager.LoadEnvMapping(GetConfigPath()); err != nil {
			return fmt.Errorf("failed to load .envrc: %w", err)
		}

		PrintVerbose("Initializing worktree management (force=%v, fetch=%v)", initForce, initFetch)
		if err := manager.Initialize(initForce, initFetch); err != nil {
			return err
		}

		PrintInfo("✅ Successfully initialized Git Branch Manager")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initForce, "force", false, "initialize even if worktrees already exist")
	initCmd.Flags().BoolVar(&initFetch, "fetch", false, "fetch remote branches during initialization")
}