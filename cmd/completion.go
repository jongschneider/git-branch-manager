package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `Generate completion script for your shell.

To load completions:

Bash:
  source <(gbm completion bash)

  # To load completions for each session, execute once:
  # Find your system's completion directory and save the completion file:
  # Common locations: /etc/bash_completion.d/, /usr/local/etc/bash_completion.d/
  gbm completion bash > <completion_directory>/gbm

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  # Save to a directory in your $fpath, such as:
  gbm completion zsh > "${fpath[1]}/_gbm"

  # You will need to start a new shell for this setup to take effect.

fish:
  gbm completion fish | source

  # To load completions for each session, execute once:
  # Save to your fish completions directory (usually ~/.config/fish/completions/):
  gbm completion fish > ~/.config/fish/completions/gbm.fish

PowerShell:
  gbm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  gbm completion powershell > gbm.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return nil
	},
}

func init() {
}

// worktreeProvider interface defines the common interface for getting worktrees
// This is used by both switch and remove commands for completion
type worktreeProvider interface {
	GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
}

// getWorktreeCompletions creates tab-separated completion entries showing worktree names and branch info
// Returns entries in format "WORKTREE_NAME\t    branch_name" for shell completion
// The part before \t gets completed, the part after is just descriptive
func getWorktreeCompletions(provider worktreeProvider) []string {
	worktrees, err := provider.GetAllWorktrees()
	if err != nil {
		return nil
	}

	if len(worktrees) == 0 {
		return nil
	}

	// Find the maximum worktree name length for alignment
	maxNameLen := 0
	for name := range worktrees {
		maxNameLen = max(maxNameLen, len(name))
	}

	var completions []string
	for name, info := range worktrees {
		// Use tab separator: "WORKTREE_NAME\t    branch_name"
		// Everything before \t gets completed, everything after is just description
		padding := strings.Repeat(" ", maxNameLen-len(name)+4) // 4 spaces minimum
		completion := fmt.Sprintf("%s\t%s%s", name, padding, info.CurrentBranch)
		completions = append(completions, completion)
	}
	return completions
}

// getWorktreeCompletionsWithManager is a convenience function that creates a manager and gets completions
func getWorktreeCompletionsWithManager() []string {
	manager, err := createInitializedManager()
	if err != nil {
		if !errors.Is(err, ErrLoadGBMConfig) {
			return nil
		}
		PrintVerbose("%v", err)
	}
	return getWorktreeCompletions(manager)
}
