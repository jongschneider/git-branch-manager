package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"gbm/internal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var logFile *os.File

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gbm",
		Short: "Git Branch Manager - Manage Git worktrees based on gbm.branchconfig.yaml",
		Long: `Git Branch Manager (gbm) is a CLI tool that manages Git repository branches
and worktrees based on configuration defined in gbm.branchconfig.yaml.

The tool synchronizes local worktrees with branch definitions and provides
notifications when configurations drift out of sync.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			InitializeLogging(cmd)
			checkAndDisplayMergeBackAlerts()
		},
	}

	// Add persistent flags
	rootCmd.PersistentFlags().String("worktree-dir", "", "override worktree directory location")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging to ./gbm.log")

	// Create manager for commands that need it
	manager, err := createInitializedManager()
	if err != nil {
		// For commands that can work without manager, we still add them
		// but pass nil manager and they handle gracefully
		if !errors.Is(err, ErrLoadGBMConfig) {
			// Hard error - something is seriously wrong
			PrintError("Failed to initialize manager: %v", err)
		} else {
			PrintVerbose("Manager initialization failed: %v", err)
		}
	}

	// Add all subcommands
	rootCmd.AddCommand(newAddCommand(manager))
	rootCmd.AddCommand(newPushCommand())
	rootCmd.AddCommand(newCloneCommand())
	rootCmd.AddCommand(newInitCommand())
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(newHotfixCommand())
	rootCmd.AddCommand(newInfoCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newMergebackCommand())
	rootCmd.AddCommand(newPullCommand())
	rootCmd.AddCommand(newRemoveCommand())
	rootCmd.AddCommand(shellIntegrationCmd)
	rootCmd.AddCommand(newSwitchCommand())
	rootCmd.AddCommand(newSyncCommand())
	rootCmd.AddCommand(newValidateCommand())

	return rootCmd
}

func Execute() error {
	return newRootCommand().Execute()
}

func isDebugEnabled(cmd *cobra.Command) bool {
	debug, _ := cmd.Flags().GetBool("debug")
	return debug
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

func PrintInfo(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", internal.FormatInfo(msg))
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		_, _ = fmt.Fprintf(logFile, "%s [INFO] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintVerbose(format string, args ...any) {
	// For backwards compatibility, assume debug mode from global logFile state
	msg := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprintf(os.Stderr, "%s\n", internal.FormatVerbose(msg))
	}
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		_, _ = fmt.Fprintf(logFile, "%s [DEBUG] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", internal.FormatError("ERROR: "+msg))
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		_, _ = fmt.Fprintf(logFile, "%s [ERROR] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func CloseLogFile() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

// createInitializedManager creates a new manager and requires branch config to exist.
// It returns an error if branch config cannot be loaded.
func createInitializedManager() (*internal.Manager, error) {
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

	configPath := internal.DefaultBranchConfigFilename

	PrintVerbose("Loading %s configuration from: %s", internal.DefaultBranchConfigFilename, configPath)
	err = manager.LoadGBMConfig(configPath)
	if err != nil {
		err = fmt.Errorf("failed to load %s: %w: %w", internal.DefaultBranchConfigFilename, err, ErrLoadGBMConfig)
	}

	return manager, err
}

var ErrLoadGBMConfig = fmt.Errorf("failed to load gbm.branchconfig.yaml")

func checkAndDisplayMergeBackAlerts() {
	// Check if merge-back alerts should be shown
	if !shouldShowMergeBackAlerts() {
		return
	}

	configPath := internal.DefaultBranchConfigFilename

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
		fmt.Fprintln(os.Stderr, alert)

		// Update the LastMergebackCheck timestamp since we showed an alert
		updateLastMergebackCheck()
	}
}

// updateLastMergebackCheck updates the LastMergebackCheck timestamp in state
func updateLastMergebackCheck() {
	wd, err := os.Getwd()
	if err != nil {
		PrintVerbose("Failed to get working directory when updating mergeback timestamp: %v", err)
		return
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		PrintVerbose("Not in a git repository when updating mergeback timestamp: %v", err)
		return
	}

	gbmDir := internal.GetGBMDir(repoRoot)
	state, err := internal.LoadState(gbmDir)
	if err != nil {
		PrintVerbose("Failed to load state when updating mergeback timestamp: %v", err)
		return
	}

	state.LastMergebackCheck = time.Now()
	if err := state.Save(gbmDir); err != nil {
		PrintVerbose("Failed to save state after updating mergeback timestamp: %v", err)
	}
}

// updateLastMergebackCheckWithError is a version that returns errors for testing
func updateLastMergebackCheckWithError() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	gbmDir := internal.GetGBMDir(repoRoot)
	state, err := internal.LoadState(gbmDir)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	state.LastMergebackCheck = time.Now()
	if err := state.Save(gbmDir); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// shouldShowMergeBackAlerts checks configuration and timestamp to determine
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

	// If merge back alerts are disabled, don't show them (early short-circuit)
	if !config.Settings.MergeBackAlerts {
		return false
	}

	// Load state to check timestamp
	state, err := internal.LoadState(gbmDir)
	if err != nil {
		PrintVerbose("Failed to load state: %v", err)
		return true // Default to showing alerts if we can't load state
	}

	// Check if enough time has passed since last check
	timeSinceLastCheck := time.Since(state.LastMergebackCheck)

	// Quick check if we need to determine user commits first
	configPath := internal.DefaultBranchConfigFilename
	status, err := internal.CheckMergeBackStatus(configPath)
	if err != nil {
		PrintVerbose("Failed to check merge-back status for timestamp logic: %v", err)
		// Use the normal interval if we can't determine user commits
		return timeSinceLastCheck >= config.Settings.MergeBackCheckInterval
	}

	// If no mergebacks needed, don't show alerts regardless of timing
	if status == nil || len(status.MergeBacksNeeded) == 0 {
		return false
	}

	// Use appropriate interval based on whether user has commits
	var interval time.Duration
	if status.HasUserCommits {
		interval = config.Settings.MergeBackUserCommitInterval
	} else {
		interval = config.Settings.MergeBackCheckInterval
	}

	return timeSinceLastCheck >= interval
}
