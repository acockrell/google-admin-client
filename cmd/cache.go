package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	// Cache configuration flags
	noCacheFlag  bool
	cacheTTLFlag string
)

// CacheEntry represents a cached data entry with metadata
type CacheEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	TTL       int64       `json:"ttl"` // TTL in seconds
	Data      interface{} `json:"data"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	Directory     string
	TotalSize     int64
	EntryCount    int
	OldestEntry   time.Time
	NewestEntry   time.Time
	DefaultTTL    time.Duration
	CacheEnabled  bool
	CacheLocation string
}

// getCacheDir returns the cache directory path
func getCacheDir() (string, error) {
	// Check if cache directory is configured
	cacheDir := viper.GetString("cache.directory")

	if cacheDir == "" {
		// Default to ~/.cache/gac
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %w", err)
		}
		cacheDir = filepath.Join(usr.HomeDir, ".cache", "gac")
	}

	// Expand ~ if present
	if strings.HasPrefix(cacheDir, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %w", err)
		}
		cacheDir = filepath.Join(usr.HomeDir, cacheDir[2:])
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// isCacheEnabled returns whether caching is enabled
func isCacheEnabled() bool {
	// Check if --no-cache flag is set
	if noCacheFlag {
		return false
	}

	// Check config file setting (default to true)
	return viper.GetBool("cache.enabled")
}

// getCacheTTL returns the cache TTL duration
func getCacheTTL() time.Duration {
	// Check if --cache-ttl flag is set
	if cacheTTLFlag != "" {
		duration, err := time.ParseDuration(cacheTTLFlag)
		if err != nil {
			Logger.Warn().Err(err).Str("ttl", cacheTTLFlag).Msg("Invalid cache TTL, using default")
		} else {
			return duration
		}
	}

	// Check config file setting
	configTTL := viper.GetString("cache.ttl")
	if configTTL != "" {
		duration, err := time.ParseDuration(configTTL)
		if err != nil {
			Logger.Warn().Err(err).Str("ttl", configTTL).Msg("Invalid cache TTL in config, using default")
		} else {
			return duration
		}
	}

	// Default to 15 minutes
	return 15 * time.Minute
}

// getCacheKey generates a cache key from resource type, domain, and filters
func getCacheKey(resourceType, domain string, filters map[string]string) string {
	// Start with resource type and domain
	keyParts := []string{resourceType, domain}

	// Add sorted filter keys and values for consistency
	if len(filters) > 0 {
		// Create a deterministic string from filters
		filterStr := ""
		for k, v := range filters {
			filterStr += fmt.Sprintf("%s=%s,", k, v)
		}
		// Hash the filter string to keep filename manageable
		hash := sha256.Sum256([]byte(filterStr))
		filterHash := fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes of hash
		keyParts = append(keyParts, filterHash)
	} else {
		keyParts = append(keyParts, "default")
	}

	// Join with hyphens and add .json extension
	return strings.Join(keyParts, "-") + ".json"
}

// readFromCache reads data from cache if it exists and is not expired
func readFromCache(key string, ttl time.Duration) (interface{}, error) {
	if !isCacheEnabled() {
		return nil, errors.New("cache disabled")
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	cachePath := filepath.Join(cacheDir, key)

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, errors.New("cache miss")
	}

	// #nosec G304 - Path is constructed from validated cache directory
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	// Check if cache is expired
	age := time.Since(entry.Timestamp)
	effectiveTTL := time.Duration(entry.TTL) * time.Second

	// Use provided TTL if it's different from stored TTL
	if ttl > 0 && ttl != effectiveTTL {
		effectiveTTL = ttl
	}

	if age > effectiveTTL {
		Logger.Debug().
			Str("key", key).
			Dur("age", age).
			Dur("ttl", effectiveTTL).
			Msg("Cache expired")
		return nil, errors.New("cache expired")
	}

	Logger.Debug().
		Str("key", key).
		Dur("age", age).
		Dur("ttl", effectiveTTL).
		Msg("Cache hit")

	return entry.Data, nil
}

// writeToCache writes data to cache
func writeToCache(key string, data interface{}, ttl time.Duration) error {
	if !isCacheEnabled() {
		return nil // Silently skip if cache disabled
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	entry := CacheEntry{
		Timestamp: time.Now(),
		TTL:       int64(ttl.Seconds()),
		Data:      data,
	}

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	cachePath := filepath.Join(cacheDir, key)

	// #nosec G306 - Cache files should be user-readable only (0600)
	if err := os.WriteFile(cachePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	Logger.Debug().
		Str("key", key).
		Dur("ttl", ttl).
		Msg("Cache written")

	return nil
}

// clearCache clears cache entries for a specific resource type or all caches
func clearCache(resourceType string) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	clearedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// If resourceType is specified, only delete matching files
		if resourceType != "" && resourceType != "all" {
			if !strings.HasPrefix(entry.Name(), resourceType+"-") {
				continue
			}
		}

		cachePath := filepath.Join(cacheDir, entry.Name())
		if err := os.Remove(cachePath); err != nil {
			Logger.Warn().Err(err).Str("file", entry.Name()).Msg("Failed to remove cache file")
			continue
		}
		clearedCount++
	}

	Logger.Info().
		Str("resource_type", resourceType).
		Int("cleared_count", clearedCount).
		Msg("Cache cleared")

	return nil
}

// getCacheStats returns cache statistics
func getCacheStats() (*CacheStats, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	stats := &CacheStats{
		Directory:     cacheDir,
		TotalSize:     0,
		EntryCount:    0,
		DefaultTTL:    getCacheTTL(),
		CacheEnabled:  isCacheEnabled(),
		CacheLocation: cacheDir,
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		stats.TotalSize += info.Size()
		stats.EntryCount++

		modTime := info.ModTime()
		if stats.OldestEntry.IsZero() || modTime.Before(stats.OldestEntry) {
			stats.OldestEntry = modTime
		}
		if stats.NewestEntry.IsZero() || modTime.After(stats.NewestEntry) {
			stats.NewestEntry = modTime
		}
	}

	return stats, nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
