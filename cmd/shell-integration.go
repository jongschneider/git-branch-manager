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

# Set environment variable to indicate shell integration is active
export GBM_SHELL_INTEGRATION=1

# Enable completion
if [ -n "$ZSH_VERSION" ]; then
    # For zsh, enable completion
    autoload -U compinit && compinit
    source <(gbm completion zsh)
elif [ -n "$BASH_VERSION" ]; then
    # For bash, enable completion
    source <(gbm completion bash)
fi

__gbm_prompt() {
    if [ -f ".envrc" ] && [ -d ".git" ]; then
        local gbm_status=$(gbm check --format=prompt 2>/dev/null)
        if [ $? -eq 0 ] && [ -n "$gbm_status" ]; then
            echo "$gbm_status"
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

# Function to switch between worktrees
gbm-switch() {
    if [ $# -eq 0 ]; then
        gbm switch
        return
    fi

    local target_dir=$(gbm switch --print-path "$1" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$target_dir" ]; then
        cd "$target_dir"
        echo "Switched to worktree: $1"
    else
        gbm switch "$@"
    fi
}

# Override gbm function to handle switch command specially
gbm() {
    if [ "$1" = "switch" ] && [ $# -gt 1 ]; then
        # For switch command with arguments, execute the cd command
        local cmd_output=$(command gbm "$@" 2>/dev/null)
        if [ $? -eq 0 ] && [[ "$cmd_output" =~ ^cd ]]; then
            eval "$cmd_output"
            echo "Switched to worktree: $2"
        else
            command gbm "$@"
        fi
    else
        # For all other commands, just pass through
        command gbm "$@"
    fi
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
