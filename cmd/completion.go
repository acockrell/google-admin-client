package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for gac.

The completion command generates shell completion scripts for bash, zsh, and fish shells.
These scripts enable tab completion for gac commands, flags, and arguments.

Installation Instructions:

Bash:
  # Linux:
  gac completion bash > /etc/bash_completion.d/gac

  # macOS (Homebrew):
  gac completion bash > $(brew --prefix)/etc/bash_completion.d/gac

  # Manual installation:
  gac completion bash > ~/.gac-completion.bash
  echo 'source ~/.gac-completion.bash' >> ~/.bashrc

Zsh:
  # Add to .zshrc:
  gac completion zsh > "${fpath[1]}/_gac"

  # Or with oh-my-zsh:
  gac completion zsh > ~/.oh-my-zsh/completions/_gac

  # Then reload:
  autoload -U compinit && compinit

Fish:
  gac completion fish > ~/.config/fish/completions/gac.fish

After installing the completion script, restart your shell or source your shell
configuration file.

Examples:
  # Generate bash completion
  gac completion bash

  # Generate zsh completion
  gac completion zsh

  # Generate fish completion
  gac completion fish

  # Install bash completion (Linux)
  sudo gac completion bash > /etc/bash_completion.d/gac

  # Install zsh completion
  gac completion zsh > ~/.oh-my-zsh/completions/_gac && exec zsh
`,
	Args:                  cobra.NoArgs,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	DisableFlagsInUseLine: true,
}

var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	Long: `Generate bash completion script for gac.

To install bash completion:

Linux:
  gac completion bash | sudo tee /etc/bash_completion.d/gac > /dev/null

macOS (with Homebrew):
  gac completion bash > $(brew --prefix)/etc/bash_completion.d/gac

Manual installation:
  gac completion bash > ~/.gac-completion.bash
  echo 'source ~/.gac-completion.bash' >> ~/.bashrc

After installation, restart your shell or run:
  source ~/.bashrc
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.GenBashCompletionV2(os.Stdout, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating bash completion: %v\n", err)
			os.Exit(1)
		}
	},
}

var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	Long: `Generate zsh completion script for gac.

To install zsh completion:

  # Standard installation
  gac completion zsh > "${fpath[1]}/_gac"

  # With oh-my-zsh
  gac completion zsh > ~/.oh-my-zsh/completions/_gac

After installation, restart your shell or run:
  autoload -U compinit && compinit
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating zsh completion: %v\n", err)
			os.Exit(1)
		}
	},
}

var fishCompletionCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completion script",
	Long: `Generate fish completion script for gac.

To install fish completion:

  gac completion fish > ~/.config/fish/completions/gac.fish

After installation, fish will automatically load the completions.
No need to restart the shell.
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.GenFishCompletion(os.Stdout, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating fish completion: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(bashCompletionCmd)
	completionCmd.AddCommand(zshCompletionCmd)
	completionCmd.AddCommand(fishCompletionCmd)
}
