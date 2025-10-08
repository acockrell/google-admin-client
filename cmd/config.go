package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long: `Manage and validate gac configuration.

The config command provides utilities for working with gac configuration,
including validation of configuration files, credentials, and OAuth2 tokens.

Examples:
  # Validate current configuration
  gac config validate

  # Show current configuration
  gac config show
`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
