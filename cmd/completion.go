package cmd

import (
	"os"

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
