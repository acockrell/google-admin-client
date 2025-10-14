package cmd

import (
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command group
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache for user and group listings",
	Long: `Manage the local cache used to speed up user and group listing operations.

The cache stores API responses locally to reduce API calls and improve performance.
Cache entries expire based on the configured TTL (Time-To-Live).

Available subcommands:
  status  - Show cache statistics
  clear   - Clear cache entries

Examples:
  gac cache status
  gac cache clear users
  gac cache clear all`,
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
