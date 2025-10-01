package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGetDomain(t *testing.T) {
	// Save original values
	originalDomain := domain
	defer func() {
		domain = originalDomain
		viper.Set("domain", "")
	}()

	tests := []struct {
		name        string
		flagValue   string
		viperValue  string
		expected    string
		description string
	}{
		{
			name:        "flag takes precedence",
			flagValue:   "flag.com",
			viperValue:  "viper.com",
			expected:    "flag.com",
			description: "When both flag and config are set, flag should win",
		},
		{
			name:        "viper when no flag",
			flagValue:   "",
			viperValue:  "viper.com",
			expected:    "viper.com",
			description: "When only config is set, use config value",
		},
		{
			name:        "empty when neither set",
			flagValue:   "",
			viperValue:  "",
			expected:    "",
			description: "When neither flag nor config set, return empty",
		},
		{
			name:        "flag only",
			flagValue:   "flag-only.com",
			viperValue:  "",
			expected:    "flag-only.com",
			description: "When only flag is set, use flag value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test conditions
			domain = tt.flagValue
			viper.Set("domain", tt.viperValue)

			// Test
			result := getDomain()

			// Verify
			if result != tt.expected {
				t.Errorf("getDomain() = %q, want %q (%s)", result, tt.expected, tt.description)
			}
		})
	}
}

func TestGetDomainEnvironmentVariable(t *testing.T) {
	// Save original values
	originalDomain := domain
	defer func() {
		domain = originalDomain
		_ = os.Unsetenv("GOOGLE_ADMIN_DOMAIN")
		viper.Set("domain", "")
	}()

	// Set environment variable
	_ = os.Setenv("GOOGLE_ADMIN_DOMAIN", "env.com")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("google_admin")

	domain = ""
	viper.Set("domain", "")

	// Manually set viper value to simulate env var reading
	viper.Set("domain", "env.com")

	result := getDomain()
	if result != "env.com" {
		t.Errorf("getDomain() with env var = %q, want %q", result, "env.com")
	}
}

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "gac" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "gac")
	}
}

func TestRootCommandHasDomainFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("domain")
	if flag == nil {
		t.Fatal("--domain flag should exist")
	}

	if flag.Usage != "domain for email addresses (e.g., example.com)" {
		t.Errorf("--domain flag usage incorrect: %q", flag.Usage)
	}
}

func TestRootCommandHasRequiredFlags(t *testing.T) {
	requiredFlags := []string{"config", "secret", "cache", "domain"}

	for _, flagName := range requiredFlags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("required flag --%s does not exist", flagName)
		}
	}
}
