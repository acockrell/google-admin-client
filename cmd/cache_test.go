package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetCacheDir(t *testing.T) {
	// Save original config
	originalCacheDir := viper.GetString("cache.directory")
	defer viper.Set("cache.directory", originalCacheDir)

	tests := []struct {
		name      string
		cacheDir  string
		wantErr   bool
		shouldSet bool
	}{
		{
			name:      "use default cache directory",
			cacheDir:  "",
			wantErr:   false,
			shouldSet: false,
		},
		{
			name:      "use custom cache directory",
			cacheDir:  os.TempDir() + "/test-cache",
			wantErr:   false,
			shouldSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSet {
				viper.Set("cache.directory", tt.cacheDir)
			} else {
				viper.Set("cache.directory", "")
			}

			dir, err := getCacheDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("getCacheDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory exists
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Errorf("getCacheDir() directory does not exist: %s", dir)
				}

				// Clean up test directory
				if tt.shouldSet {
					if err := os.RemoveAll(tt.cacheDir); err != nil {
						t.Logf("Failed to clean up test directory: %v", err)
					}
				}
			}
		})
	}
}

func TestIsCacheEnabled(t *testing.T) {
	// Save original values
	originalNoCache := noCacheFlag
	originalEnabled := viper.GetBool("cache.enabled")
	defer func() {
		noCacheFlag = originalNoCache
		viper.Set("cache.enabled", originalEnabled)
	}()

	tests := []struct {
		name        string
		noCacheFlag bool
		configValue bool
		wantEnabled bool
	}{
		{
			name:        "cache enabled by default",
			noCacheFlag: false,
			configValue: true,
			wantEnabled: true,
		},
		{
			name:        "cache disabled by flag",
			noCacheFlag: true,
			configValue: true,
			wantEnabled: false,
		},
		{
			name:        "cache disabled by config",
			noCacheFlag: false,
			configValue: false,
			wantEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noCacheFlag = tt.noCacheFlag
			viper.Set("cache.enabled", tt.configValue)

			enabled := isCacheEnabled()
			if enabled != tt.wantEnabled {
				t.Errorf("isCacheEnabled() = %v, want %v", enabled, tt.wantEnabled)
			}
		})
	}
}

func TestGetCacheTTL(t *testing.T) {
	// Save original values
	originalTTLFlag := cacheTTLFlag
	originalTTL := viper.GetString("cache.ttl")
	defer func() {
		cacheTTLFlag = originalTTLFlag
		viper.Set("cache.ttl", originalTTL)
	}()

	tests := []struct {
		name        string
		ttlFlag     string
		configValue string
		wantTTL     time.Duration
	}{
		{
			name:        "default TTL",
			ttlFlag:     "",
			configValue: "",
			wantTTL:     15 * time.Minute,
		},
		{
			name:        "TTL from flag",
			ttlFlag:     "30m",
			configValue: "10m",
			wantTTL:     30 * time.Minute,
		},
		{
			name:        "TTL from config",
			ttlFlag:     "",
			configValue: "1h",
			wantTTL:     1 * time.Hour,
		},
		{
			name:        "invalid flag uses default",
			ttlFlag:     "invalid",
			configValue: "",
			wantTTL:     15 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheTTLFlag = tt.ttlFlag
			viper.Set("cache.ttl", tt.configValue)

			ttl := getCacheTTL()
			if ttl != tt.wantTTL {
				t.Errorf("getCacheTTL() = %v, want %v", ttl, tt.wantTTL)
			}
		})
	}
}

func TestGetCacheKey(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		domain       string
		filters      map[string]string
		wantContains []string
	}{
		{
			name:         "basic cache key",
			resourceType: "users",
			domain:       "example.com",
			filters:      nil,
			wantContains: []string{"users", "example.com", "default"},
		},
		{
			name:         "cache key with filters",
			resourceType: "groups",
			domain:       "test.com",
			filters:      map[string]string{"inactive": "true"},
			wantContains: []string{"groups", "test.com"},
		},
		{
			name:         "group members key",
			resourceType: "group-members",
			domain:       "team@example.com",
			filters:      nil,
			wantContains: []string{"group-members", "team@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := getCacheKey(tt.resourceType, tt.domain, tt.filters)

			for _, contains := range tt.wantContains {
				if !containsString(key, contains) {
					t.Errorf("getCacheKey() = %v, want to contain %v", key, contains)
				}
			}

			// Verify key ends with .json
			if !containsString(key, ".json") {
				t.Errorf("getCacheKey() = %v, should end with .json", key)
			}
		})
	}
}

func TestWriteAndReadCache(t *testing.T) {
	// Setup test cache directory
	testCacheDir := filepath.Join(os.TempDir(), "gac-test-cache")
	viper.Set("cache.directory", testCacheDir)
	viper.Set("cache.enabled", true)
	noCacheFlag = false

	defer func() {
		if err := os.RemoveAll(testCacheDir); err != nil {
			t.Logf("Failed to clean up test cache directory: %v", err)
		}
		viper.Set("cache.directory", "")
	}()

	tests := []struct {
		name    string
		key     string
		data    interface{}
		ttl     time.Duration
		wantErr bool
	}{
		{
			name: "write and read string data",
			key:  "test-string.json",
			data: "test data",
			ttl:  1 * time.Minute,
		},
		{
			name: "write and read map data",
			key:  "test-map.json",
			data: map[string]string{"key": "value"},
			ttl:  5 * time.Minute,
		},
		{
			name: "write and read slice data",
			key:  "test-slice.json",
			data: []string{"item1", "item2", "item3"},
			ttl:  10 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write to cache
			err := writeToCache(tt.key, tt.data, tt.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeToCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Read from cache
				cachedData, err := readFromCache(tt.key, tt.ttl)
				if err != nil {
					t.Errorf("readFromCache() error = %v", err)
					return
				}

				// Verify data is not nil
				if cachedData == nil {
					t.Errorf("readFromCache() returned nil data")
				}
			}
		})
	}
}

func TestCacheExpiration(t *testing.T) {
	// Setup test cache directory
	testCacheDir := filepath.Join(os.TempDir(), "gac-test-cache-expiry")
	viper.Set("cache.directory", testCacheDir)
	viper.Set("cache.enabled", true)
	noCacheFlag = false

	defer func() {
		if err := os.RemoveAll(testCacheDir); err != nil {
			t.Logf("Failed to clean up test cache directory: %v", err)
		}
		viper.Set("cache.directory", "")
	}()

	key := "test-expiry.json"
	data := "test data"
	ttl := 100 * time.Millisecond // Very short TTL for testing

	// Write to cache
	err := writeToCache(key, data, ttl)
	if err != nil {
		t.Fatalf("writeToCache() error = %v", err)
	}

	// Read immediately - should succeed
	_, err = readFromCache(key, ttl)
	if err != nil {
		t.Errorf("readFromCache() should succeed immediately after write, got error: %v", err)
	}

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Read after expiration - should fail
	_, err = readFromCache(key, ttl)
	if err == nil {
		t.Errorf("readFromCache() should fail after TTL expiration")
	}
}

func TestClearCache(t *testing.T) {
	// Setup test cache directory
	testCacheDir := filepath.Join(os.TempDir(), "gac-test-cache-clear")
	viper.Set("cache.directory", testCacheDir)
	viper.Set("cache.enabled", true)
	noCacheFlag = false

	defer func() {
		if err := os.RemoveAll(testCacheDir); err != nil {
			t.Logf("Failed to clean up test cache directory: %v", err)
		}
		viper.Set("cache.directory", "")
	}()

	// Create some test cache files
	if err := writeToCache("users-test.com-default.json", "users data", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}
	if err := writeToCache("groups-test.com-default.json", "groups data", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}
	if err := writeToCache("other-test.json", "other data", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}

	tests := []struct {
		name         string
		resourceType string
		wantErr      bool
		checkFiles   map[string]bool // filename -> should exist after clear
	}{
		{
			name:         "clear users cache",
			resourceType: "users",
			wantErr:      false,
			checkFiles: map[string]bool{
				"users-test.com-default.json":  false,
				"groups-test.com-default.json": true,
				"other-test.json":              true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate cache files before each test
			if err := writeToCache("users-test.com-default.json", "users data", 1*time.Hour); err != nil {
				t.Fatalf("Failed to write test cache: %v", err)
			}
			if err := writeToCache("groups-test.com-default.json", "groups data", 1*time.Hour); err != nil {
				t.Fatalf("Failed to write test cache: %v", err)
			}
			if err := writeToCache("other-test.json", "other data", 1*time.Hour); err != nil {
				t.Fatalf("Failed to write test cache: %v", err)
			}

			err := clearCache(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("clearCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check which files exist
			for filename, shouldExist := range tt.checkFiles {
				cachePath := filepath.Join(testCacheDir, filename)
				_, err := os.Stat(cachePath)
				exists := !os.IsNotExist(err)

				if exists != shouldExist {
					t.Errorf("File %s: exists = %v, want %v", filename, exists, shouldExist)
				}
			}
		})
	}
}

func TestGetCacheStats(t *testing.T) {
	// Setup test cache directory
	testCacheDir := filepath.Join(os.TempDir(), "gac-test-cache-stats")
	viper.Set("cache.directory", testCacheDir)
	viper.Set("cache.enabled", true)
	noCacheFlag = false

	defer func() {
		if err := os.RemoveAll(testCacheDir); err != nil {
			t.Logf("Failed to clean up test cache directory: %v", err)
		}
		viper.Set("cache.directory", "")
	}()

	// Create some test cache files
	if err := writeToCache("test1.json", "data1", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}
	if err := writeToCache("test2.json", "data2", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}
	if err := writeToCache("test3.json", "data3", 1*time.Hour); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}

	stats, err := getCacheStats()
	if err != nil {
		t.Fatalf("getCacheStats() error = %v", err)
	}

	if stats.EntryCount != 3 {
		t.Errorf("getCacheStats() EntryCount = %v, want 3", stats.EntryCount)
	}

	if stats.TotalSize == 0 {
		t.Errorf("getCacheStats() TotalSize should be > 0")
	}

	if !stats.CacheEnabled {
		t.Errorf("getCacheStats() CacheEnabled = false, want true")
	}

	if stats.CacheLocation != testCacheDir {
		t.Errorf("getCacheStats() CacheLocation = %v, want %v", stats.CacheLocation, testCacheDir)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"zero bytes", 0, "0 B"},
		{"bytes", 100, "100 B"},
		{"kilobytes", 1024, "1.0 KB"},
		{"megabytes", 1024 * 1024, "1.0 MB"},
		{"gigabytes", 1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("formatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
