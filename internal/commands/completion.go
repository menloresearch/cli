package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:    "completion [shell]",
	Short:  "Generate shell completion scripts",
	Hidden: true,
	Long: `Generate shell completion scripts for menlo.

To load completions:

Bash:

  $ source <(menlo completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ menlo completion bash > /etc/bash_completion.d/menlo
  # macOS:
  $ menlo completion bash > /usr/local/etc/bash_completion.d/menlo

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ menlo completion zsh > "${fpath[1]}/_menlo"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ menlo completion fish | source

  # To load completions for each session, execute once:
  $ menlo completion fish > ~/.config/fish/completions/menlo.fish

PowerShell:

  PS> menlo completion powershell | Out-String | Invoke-Expression

  # To load completions for every session, add the output of the previous
  # command to your powershell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}