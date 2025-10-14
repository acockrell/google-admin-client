package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// cacheStatusCmd represents the cache status command
var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache statistics",
	Long: `Display statistics about the current cache state including:
  - Cache location
  - Total size
  - Number of entries
  - Age of oldest and newest entries
  - Default TTL setting
  - Cache enabled/disabled status

Examples:
  gac cache status
  gac cache status --format json`,
	Run: cacheStatusRunFunc,
}

func init() {
	cacheCmd.AddCommand(cacheStatusCmd)
}

type cacheStatusOutput struct {
	Enabled      bool   `json:"enabled"`
	Location     string `json:"location"`
	TotalSize    string `json:"total_size"`
	TotalSizeRaw int64  `json:"total_size_bytes"`
	EntryCount   int    `json:"entry_count"`
	OldestEntry  string `json:"oldest_entry,omitempty"`
	NewestEntry  string `json:"newest_entry,omitempty"`
	DefaultTTL   string `json:"default_ttl"`
}

func cacheStatusRunFunc(cmd *cobra.Command, args []string) {
	stats, err := getCacheStats()
	if err != nil {
		exitWithError(fmt.Sprintf("Failed to get cache statistics: %s", err))
	}

	// Prepare output data
	output := cacheStatusOutput{
		Enabled:      stats.CacheEnabled,
		Location:     stats.CacheLocation,
		TotalSize:    formatBytes(stats.TotalSize),
		TotalSizeRaw: stats.TotalSize,
		EntryCount:   stats.EntryCount,
		DefaultTTL:   stats.DefaultTTL.String(),
	}

	if !stats.OldestEntry.IsZero() {
		output.OldestEntry = formatTimeAgo(stats.OldestEntry)
	}

	if !stats.NewestEntry.IsZero() {
		output.NewestEntry = formatTimeAgo(stats.NewestEntry)
	}

	// Use unified formatter for structured output
	if outputFormat == OutputFormatJSON || outputFormat == OutputFormatYAML {
		if err := FormatOutput(output, nil); err != nil {
			exitWithError(fmt.Sprintf("Failed to format output: %s", err))
		}
		return
	}

	// Plain text output (default)
	fmt.Println("Cache Status")
	fmt.Println("============")
	fmt.Printf("Enabled:       %v\n", stats.CacheEnabled)
	fmt.Printf("Location:      %s\n", stats.CacheLocation)
	fmt.Printf("Total Size:    %s\n", formatBytes(stats.TotalSize))
	fmt.Printf("Entry Count:   %d\n", stats.EntryCount)

	if stats.EntryCount > 0 {
		fmt.Printf("Oldest Entry:  %s\n", formatTimeAgo(stats.OldestEntry))
		fmt.Printf("Newest Entry:  %s\n", formatTimeAgo(stats.NewestEntry))
	}

	fmt.Printf("Default TTL:   %s\n", stats.DefaultTTL)

	if !stats.CacheEnabled {
		fmt.Println("\nNote: Caching is currently disabled")
		fmt.Println("Enable with: cache.enabled=true in config or remove --no-cache flag")
	}
}

// formatTimeAgo formats a time as a relative duration from now
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}
