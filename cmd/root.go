package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"gbm/internal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	logFile *os.File
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gbm",
		Short: "Git Branch Manager - Manage Git worktrees based on .envrc configuration",
		Long: `Git Branch Manager (gbm) is a CLI tool that manages Git repository branches
and worktrees based on environment variables defined in a .envrc file.

The tool synchronizes local worktrees with branch definitions and provides
notifications when configurations drift out of sync.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			InitializeLogging(cmd)
			checkAndDisplayMergeBackAlerts(cmd)
		},
	}

	// Add persistent flags
	rootCmd.PersistentFlags().String("config", "", "specify custom .gbm.config.yaml path")
	rootCmd.PersistentFlags().String("worktree-dir", "", "override worktree directory location")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging to ./gbm.log")

	// Add all subcommands
	rootCmd.AddCommand(newPushCommand())
	rootCmd.AddCommand(newAddCommand())
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(hotfixCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(mergebackCmd)
	rootCmd.AddCommand(newPullCommand())
	rootCmd.AddCommand(newRemoveCommand())
	rootCmd.AddCommand(shellIntegrationCmd)
	rootCmd.AddCommand(newSwitchCommand())
	rootCmd.AddCommand(newSyncCommand())
	rootCmd.AddCommand(validateCmd)

	return rootCmd
}

func Execute() error {
	return newRootCommand().Execute()
}

func getConfigPath(cmd *cobra.Command) string {
	configPath, _ := cmd.Flags().GetString("config")
	if configPath != "" {
		return configPath
	}
	return ".gbm.config.yaml"
}

func getWorktreeDir(cmd *cobra.Command) string {
	worktreeDir, _ := cmd.Flags().GetString("worktree-dir")
	if worktreeDir != "" {
		return worktreeDir
	}
	return "worktrees"
}

func isDebugEnabled(cmd *cobra.Command) bool {
	debug, _ := cmd.Flags().GetBool("debug")
	return debug
}

// Legacy functions for backwards compatibility
func GetConfigPath() string {
	// This is kept for backwards compatibility, but should be phased out
	return ".gbm.config.yaml"
}

func GetWorktreeDir() string {
	// This is kept for backwards compatibility, but should be phased out
	return "worktrees"
}

func InitializeLogging(cmd *cobra.Command) {
	if isDebugEnabled(cmd) {
		var err error
		logFile, err = tea.LogToFile("gbm.log", "gbm")
		if err != nil {
			PrintError("Failed to initialize log file: %v", err)
		}
	}
}

func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", internal.FormatInfo(msg))
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		fmt.Fprintf(logFile, "%s [INFO] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintVerbose(format string, args ...interface{}) {
	// For backwards compatibility, assume debug mode from global logFile state
	msg := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprintf(os.Stderr, "%s\n", internal.FormatVerbose(msg))
	}
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		fmt.Fprintf(logFile, "%s [DEBUG] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", internal.FormatError("ERROR: "+msg))
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		fmt.Fprintf(logFile, "%s [ERROR] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

// createInitializedManager creates a new manager with git root discovery and gbm config loaded.
// It gracefully handles missing .gbm.config.yaml files by logging a verbose message.
func createInitializedManager() (*internal.Manager, error) {
	// Legacy function - uses default config path
	return createInitializedManagerWithCmd(nil)
}

func createInitializedManagerWithCmd(cmd *cobra.Command) (*internal.Manager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	repoPath, err := internal.FindGitRoot(wd)
	if err != nil {
		return nil, fmt.Errorf("not in a git repository: %w", err)
	}

	manager, err := internal.NewManager(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	configPath := GetConfigPath() // Default fallback
	if cmd != nil {
		configPath = getConfigPath(cmd)
	}

	if err := manager.LoadGBMConfig(configPath); err != nil {
		PrintVerbose("No .gbm.config.yaml found or failed to load: %v", err)
	}

	return manager, nil
}

// createInitializedManagerStrict creates a new manager and requires .gbm.config.yaml to exist.
// It returns an error if .gbm.config.yaml cannot be loaded.
func createInitializedManagerStrict() (*internal.Manager, error) {
	// Legacy function - uses default config path
	return createInitializedManagerStrictWithCmd(nil)
}

func createInitializedManagerStrictWithCmd(cmd *cobra.Command) (*internal.Manager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	repoPath, err := internal.FindGitRoot(wd)
	if err != nil {
		return nil, fmt.Errorf("failed to find git repository root: %w", err)
	}

	manager, err := internal.NewManager(repoPath)
	if err != nil {
		return nil, err
	}

	configPath := GetConfigPath() // Default fallback
	if cmd != nil {
		configPath = getConfigPath(cmd)
	}

	PrintVerbose("Loading .gbm.config.yaml configuration from: %s", configPath)
	if err := manager.LoadGBMConfig(configPath); err != nil {
		return nil, fmt.Errorf("failed to load .gbm.config.yaml: %w", err)
	}

	return manager, nil
}

// createInitializedGitManager creates a new git manager with git root discovery.
// Used by commands that need direct git operations without .gbm.config.yaml dependency.
func createInitializedGitManager() (*internal.GitManager, error) {
	// Legacy function - uses default worktree dir
	return createInitializedGitManagerWithCmd(nil)
}

func createInitializedGitManagerWithCmd(cmd *cobra.Command) (*internal.GitManager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	gitRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return nil, fmt.Errorf("not in a git repository: %w", err)
	}

	worktreeDir := internal.DefaultWorktreeDirname // Default fallback
	if cmd != nil {
		worktreeDir = getWorktreeDir(cmd)
	}

	gitManager, err := internal.NewGitManager(gitRoot, worktreeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize git manager: %w", err)
	}

	return gitManager, nil
}

func checkAndDisplayMergeBackAlerts(cmd *cobra.Command) {
	// Check if merge-back alerts should be shown
	if !shouldShowMergeBackAlerts() {
		return
	}

	configPath := GetConfigPath() // Default fallback
	if cmd != nil {
		configPath = getConfigPath(cmd)
	}

	status, err := internal.CheckMergeBackStatus(configPath)
	if err != nil {
		PrintVerbose("Failed to check merge-back status: %v", err)
		return
	}

	if status == nil {
		return
	}

	alert := internal.FormatMergeBackAlert(status)
	if alert != "" {
		fmt.Fprint(os.Stderr, alert)
	}
}

// shouldShowMergeBackAlerts checks configuration to determine
// if merge-back alerts should be displayed
func shouldShowMergeBackAlerts() bool {
	// Check configuration file
	wd, err := os.Getwd()
	if err != nil {
		PrintVerbose("Failed to get working directory: %v", err)
		return false // Default to disabled
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		PrintVerbose("Not in a git repository: %v", err)
		return false // Default to disabled
	}

	gbmDir := internal.GetGBMDir(repoRoot)
	config, err := internal.LoadConfig(gbmDir)
	if err != nil {
		PrintVerbose("Failed to load config: %v", err)
		return false // Default to disabled
	}

	return config.Settings.MergeBackAlerts
}
