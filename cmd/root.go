package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	configPath  string
	worktreeDir string
	verbose     bool
)

var rootCmd = &cobra.Command{
	Use:   "gbm",
	Short: "Git Branch Manager - Manage Git worktrees based on .envrc configuration",
	Long: `Git Branch Manager (gbm) is a CLI tool that manages Git repository branches
and worktrees based on environment variables defined in a .envrc file.

The tool synchronizes local worktrees with branch definitions and provides
notifications when configurations drift out of sync.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "specify custom .envrc path")
	rootCmd.PersistentFlags().StringVarP(&worktreeDir, "worktree-dir", "w", "", "override worktree directory location")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

func GetConfigPath() string {
	if configPath != "" {
		return configPath
	}
	return ".envrc"
}

func GetWorktreeDir() string {
	if worktreeDir != "" {
		return worktreeDir
	}
	return "worktrees"
}

func PrintInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func PrintVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

func IsVerbose() bool {
	return verbose
}
