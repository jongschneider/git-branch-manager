package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	configPath string
	worktreeDir string
	debug      bool
	logFile    *os.File
)

var rootCmd = &cobra.Command{
	Use:   "gbm",
	Short: "Git Branch Manager - Manage Git worktrees based on .envrc configuration",
	Long: `Git Branch Manager (gbm) is a CLI tool that manages Git repository branches
and worktrees based on environment variables defined in a .envrc file.

The tool synchronizes local worktrees with branch definitions and provides
notifications when configurations drift out of sync.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		InitializeLogging()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "specify custom .envrc path")
	rootCmd.PersistentFlags().StringVarP(&worktreeDir, "worktree-dir", "w", "", "override worktree directory location")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging to ./gbm.log")
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

func InitializeLogging() {
	if debug {
		var err error
		logFile, err = tea.LogToFile("gbm.log", "gbm")
		if err != nil {
			PrintError("Failed to initialize log file: %v", err)
		}
	}
}

func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		fmt.Fprintf(logFile, "%s [INFO] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintVerbose(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if logFile != nil {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		fmt.Fprintf(logFile, "%s [DEBUG] %s:%d %s\n", timestamp, file, line, msg)
	}
}

func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
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
