package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	clearAllFlag bool
)

// cacheClearCmd represents the cache clear command
var cacheClearCmd = &cobra.Command{
	Use:   "clear [users|groups|all]",
	Short: "Clear cache entries",
	Long: `Clear cache entries for users, groups, or all cached data.

This is useful when you want to force a refresh from the API,
or if cached data has become stale.

Examples:
  gac cache clear users        # Clear user listing cache
  gac cache clear groups       # Clear group listing cache
  gac cache clear all          # Clear all caches
  gac cache clear --all        # Clear all caches (alternative)`,
	Run:       cacheClearRunFunc,
	ValidArgs: []string{"users", "groups", "all"},
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	cacheClearCmd.Flags().BoolVar(&clearAllFlag, "all", false, "clear all cache entries")
}

func cacheClearRunFunc(cmd *cobra.Command, args []string) {
	// Determine what to clear
	resourceType := "all"
	if len(args) > 0 {
		resourceType = args[0]
	} else if clearAllFlag {
		resourceType = "all"
	}

	// Validate resource type
	validTypes := map[string]bool{
		"users":  true,
		"groups": true,
		"all":    true,
	}

	if !validTypes[resourceType] {
		exitWithError(fmt.Sprintf("Invalid resource type: %s. Valid types: users, groups, all", resourceType))
	}

	// Clear the cache
	if err := clearCache(resourceType); err != nil {
		exitWithError(fmt.Sprintf("Failed to clear cache: %s", err))
	}

	// Success message
	switch resourceType {
	case "all":
		Logger.Info().Msg("All cache entries cleared successfully")
	default:
		Logger.Info().Str("resource_type", resourceType).Msg("Cache cleared successfully")
	}
}
