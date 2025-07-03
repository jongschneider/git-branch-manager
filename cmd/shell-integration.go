package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shellIntegrationCmd = &cobra.Command{
	Use:   "shell-integration",
	Short: "Generate shell integration code for automatic checking",
	Long: `Generate shell integration code that can be evaluated in your shell to enable automatic checking.

This command outputs shell code that can be added to your .bashrc, .zshrc, or other shell configuration
to automatically check for worktree drift and display status in your prompt.

Example usage:
  eval "$(gbm shell-integration)"

Or add to your shell configuration:
  echo 'eval "$(gbm shell-integration)"' >> ~/.bashrc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintVerbose("Generating shell integration code")
		shellCode := `
# Git Branch Manager (gbm) shell integration
# Automatically check for worktree drift and display status

__gbm_prompt() {
    if [ -f ".envrc" ] && [ -d ".git" ]; then
        local status=$(gbm check --format=prompt 2>/dev/null)
        if [ $? -eq 0 ] && [ -n "$status" ]; then
            echo "$status"
        fi
    fi
}

# Add gbm status to PS1 for bash/zsh
if [ -n "$BASH_VERSION" ] || [ -n "$ZSH_VERSION" ]; then
    if [[ "$PS1" != *"__gbm_prompt"* ]]; then
        PS1='$(__gbm_prompt)'$PS1
    fi
fi

# Function to manually check gbm status
gbm-status() {
    gbm check --format=text
}

# Function to quickly sync worktrees
gbm-sync() {
    gbm sync "$@"
}
`
		PrintVerbose("Outputting shell integration script")
		fmt.Print(shellCode)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(shellIntegrationCmd)
}