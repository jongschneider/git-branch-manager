package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var (
	printPath bool
)

var switchCmd = &cobra.Command{
	Use:   "switch [WORKTREE_NAME]",
	Short: "Switch to a different worktree",
	Long: `Switch to a different worktree by environment variable name.

If no worktree name is provided, lists all available worktrees.
Use with shell integration for automatic directory switching:

  gbm-switch() {
      local target_dir=$(gbm switch --print-path "$1")
      if [ $? -eq 0 ] && [ -n "$target_dir" ]; then
          cd "$target_dir"
      else
          gbm switch "$@"
      fi
  }

Examples:
  gbm switch PROD      # Show path to PROD worktree
  gbm switch STAGING   # Show path to STAGING worktree
  gbm switch           # List all available worktrees`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintVerbose("Running switch command")

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		repoRoot, err := internal.FindGitRoot(wd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}
		PrintVerbose("Checking status from repository root: %s", repoRoot)

		manager, err := internal.NewManager(repoRoot)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return listWorktrees(manager)
		}

		worktreeName := strings.ToUpper(args[0])
		return switchToWorktree(manager, worktreeName)
	},
}

func switchToWorktree(manager *internal.Manager, worktreeName string) error {
	PrintVerbose("Switching to worktree: %s", worktreeName)

	// Try exact match first
	targetPath, err := manager.GetWorktreePath(worktreeName)
	if err != nil {
		// If exact match fails, try fuzzy matching
		PrintVerbose("Exact match failed, trying fuzzy matching")
		matchedName := findFuzzyMatch(manager, worktreeName)
		if matchedName != "" {
			PrintInfo("Fuzzy matched '%s' to '%s'", worktreeName, matchedName)
			targetPath, err = manager.GetWorktreePath(matchedName)
			if err != nil {
				return err
			}
			worktreeName = matchedName // Update for display
		} else {
			return err
		}
	}

	if printPath {
		fmt.Print(targetPath)
		return nil
	}

	// Check if shell integration is available by looking for gbm-switch function
	if os.Getenv("GBM_SHELL_INTEGRATION") != "" {
		// If shell integration is available, output cd command
		fmt.Printf("cd %s\n", targetPath)
		return nil
	}

	fmt.Printf("%s\n", internal.FormatInfo(fmt.Sprintf("Worktree %s is located at: %s", worktreeName, targetPath)))
	fmt.Println(internal.FormatSubtle("Use shell integration 'gbm-switch' function to automatically change directory"))
	fmt.Println(internal.FormatSubtle("Or run: cd " + targetPath))
	return nil
}

func findFuzzyMatch(manager *internal.Manager, target string) string {
	worktrees, err := manager.GetAllWorktrees()
	if err != nil {
		return ""
	}

	target = strings.ToUpper(target)

	// Look for case-insensitive substring matches
	var matches []string
	for name := range worktrees {
		upperName := strings.ToUpper(name)
		if strings.Contains(upperName, target) {
			matches = append(matches, name)
		}
	}

	// If we have exactly one match, return it
	if len(matches) == 1 {
		return matches[0]
	}

	// If we have multiple matches, prefer the one that starts with the target
	for _, match := range matches {
		upperMatch := strings.ToUpper(match)
		if strings.HasPrefix(upperMatch, target) {
			return match
		}
	}

	// If no prefix match, return the first match (if any)
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

func listWorktrees(manager *internal.Manager) error {
	PrintVerbose("Listing available worktrees")

	worktrees, err := manager.GetAllWorktrees()
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		fmt.Println(internal.FormatInfo("No worktrees found. Run 'gbm init' to create worktrees."))
		return nil
	}

	fmt.Println(internal.FormatSubHeader("Available worktrees:"))

	// Sort worktree names for consistent output
	var names []string
	for name := range worktrees {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		info := worktrees[name]
		status := "ready"
		if info.GitStatus != nil {
			status = manager.GetStatusIcon(info.GitStatus)
		}

		relPath, _ := filepath.Rel(".", info.Path)
		branchInfo := info.CurrentBranch
		if info.ExpectedBranch != info.CurrentBranch {
			branchInfo = fmt.Sprintf("%s -> %s", info.CurrentBranch, info.ExpectedBranch)
		}
		fmt.Printf("  %s %-12s %s (%s)\n", status, name, relPath, branchInfo)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().BoolVar(&printPath, "print-path", false, "Print the worktree path only (for shell integration)")

	// Add completion for worktree names
	switchCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getWorktreeNames(), cobra.ShellCompDirectiveNoFileComp
	}
}

func getWorktreeNames() []string {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	repoRoot, err := internal.FindGitRoot(wd)
	if err != nil {
		return nil
	}

	manager, err := internal.NewManager(repoRoot)
	if err != nil {
		return nil
	}

	worktrees, err := manager.GetAllWorktrees()
	if err != nil {
		return nil
	}

	var names []string
	for name := range worktrees {
		names = append(names, name)
	}
	return names
}
