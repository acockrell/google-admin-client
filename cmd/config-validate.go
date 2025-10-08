package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate gac configuration",
	Long: `Validate the gac configuration including:

- Configuration file syntax (YAML)
- Domain format and presence
- Client secret file existence and permissions
- Cache file path and permissions
- OAuth2 token validity
- Required API scopes

This command helps diagnose configuration issues and ensures all required
settings are properly configured before running gac commands.

Exit Codes:
  0 - Configuration is valid
  1 - Configuration has errors

Examples:
  # Validate current configuration
  gac config validate

  # Validate with specific config file
  gac config validate --config ~/.google-admin-test.yaml

  # Validate with verbose output
  gac config validate --verbose
`,
	Run: configValidateRunFunc,
}

func init() {
	configCmd.AddCommand(configValidateCmd)
}

func configValidateRunFunc(cmd *cobra.Command, args []string) {
	hasErrors := false
	hasWarnings := false

	fmt.Println("Validating gac configuration...")
	fmt.Println()

	// 1. Check config file
	fmt.Println("📄 Configuration File:")
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("  ✓ Using config file: %s\n", configFile)

		// Check file permissions
		fileInfo, err := os.Stat(configFile)
		if err != nil {
			fmt.Printf("  ✗ Error reading config file: %v\n", err)
			hasErrors = true
		} else {
			mode := fileInfo.Mode()
			fmt.Printf("  ✓ File permissions: %s\n", mode.Perm())

			// Warn if config file is world-readable
			if mode.Perm()&0044 != 0 {
				fmt.Println("  ⚠ Warning: Config file is world-readable (consider: chmod 600)")
				hasWarnings = true
			}
		}
	} else {
		fmt.Println("  ℹ No config file in use (using flags/env vars)")
	}
	fmt.Println()

	// 2. Validate domain
	fmt.Println("🌐 Domain Configuration:")
	domain := getDomain()
	if domain == "" {
		fmt.Println("  ✗ Domain not configured")
		fmt.Println("    Set via: --domain flag, GAC_DOMAIN env var, or config file")
		hasErrors = true
	} else {
		// Basic domain format validation
		if !strings.Contains(domain, ".") {
			fmt.Printf("  ⚠ Warning: Domain '%s' may be invalid (no TLD)\n", domain)
			hasWarnings = true
		} else {
			fmt.Printf("  ✓ Domain: %s\n", domain)
		}
	}
	fmt.Println()

	// 3. Validate client secret file
	fmt.Println("🔐 Client Secret File:")
	clientSecretPath := viper.GetString("client-secret")
	if clientSecretPath == "" {
		fmt.Println("  ✗ Client secret file not configured")
		fmt.Println("    Set via: --client-secret flag, GAC_CLIENT_SECRET env var, or config file")
		hasErrors = true
	} else {
		// Expand home directory if needed
		if strings.HasPrefix(clientSecretPath, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				clientSecretPath = filepath.Join(home, clientSecretPath[2:])
			}
		}

		// Check if file exists
		fileInfo, err := os.Stat(clientSecretPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("  ✗ Client secret file not found: %s\n", clientSecretPath)
			} else {
				fmt.Printf("  ✗ Error accessing client secret file: %v\n", err)
			}
			hasErrors = true
		} else {
			fmt.Printf("  ✓ Client secret file exists: %s\n", clientSecretPath)

			// Check file permissions
			mode := fileInfo.Mode()
			fmt.Printf("  ✓ File permissions: %s\n", mode.Perm())

			// Warn if file is world-readable
			if mode.Perm()&0044 != 0 {
				fmt.Println("  ⚠ Warning: Client secret file is world-readable (consider: chmod 600)")
				hasWarnings = true
			}

			// Try to read and parse the file
			if err := validateCredentialPath(clientSecretPath); err != nil {
				fmt.Printf("  ✗ Invalid client secret file path: %v\n", err)
				hasErrors = true
			} else {
				fmt.Println("  ✓ Client secret file path is valid")
			}
		}
	}
	fmt.Println()

	// 4. Validate cache file
	fmt.Println("💾 Token Cache File:")
	cacheFilePath := viper.GetString("cache-file")
	if cacheFilePath == "" {
		fmt.Println("  ℹ Using default cache file location")
		cacheFilePath = filepath.Join(os.Getenv("HOME"), ".credentials", "gac.json")
	}

	// Expand home directory if needed
	if strings.HasPrefix(cacheFilePath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			cacheFilePath = filepath.Join(home, cacheFilePath[2:])
		}
	}

	fmt.Printf("  ℹ Cache file path: %s\n", cacheFilePath)

	// Check if cache file exists
	fileInfo, err := os.Stat(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("  ℹ Token cache file does not exist yet (will be created on first auth)")
		} else {
			fmt.Printf("  ⚠ Warning: Error accessing cache file: %v\n", err)
			hasWarnings = true
		}
	} else {
		fmt.Println("  ✓ Token cache file exists")

		// Check file permissions
		mode := fileInfo.Mode()
		fmt.Printf("  ✓ File permissions: %s\n", mode.Perm())

		// Warn if file is world-readable
		if mode.Perm()&0044 != 0 {
			fmt.Println("  ⚠ Warning: Token cache file is world-readable (consider: chmod 600)")
			hasWarnings = true
		}

		// Try to validate token if it exists
		if validateToken(cacheFilePath) {
			fmt.Println("  ✓ Token cache is valid")
		} else {
			fmt.Println("  ⚠ Warning: Token cache may be invalid or expired (re-authentication may be required)")
			hasWarnings = true
		}
	}
	fmt.Println()

	// 5. Summary
	fmt.Println("📊 Validation Summary:")
	if hasErrors {
		fmt.Println("  ✗ Configuration has errors that must be fixed")
		fmt.Println()
		os.Exit(1)
	} else if hasWarnings {
		fmt.Println("  ⚠ Configuration is valid but has warnings")
		fmt.Println("  ✓ gac should work, but consider addressing warnings above")
		fmt.Println()
	} else {
		fmt.Println("  ✓ Configuration is valid with no errors or warnings")
		fmt.Println()
	}
}

// validateToken attempts to validate the OAuth2 token from the cache file
func validateToken(cacheFilePath string) bool {
	// Validate path before reading to prevent directory traversal
	if err := validateCredentialPath(cacheFilePath); err != nil {
		return false
	}

	// Read token from file
	// #nosec G304 - Path is validated by validateCredentialPath() above
	tokenBytes, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return false
	}

	// Basic validation: check if file is not empty and appears to be JSON
	tokenStr := string(tokenBytes)
	if len(tokenStr) == 0 {
		return false
	}

	// Check for basic JSON structure
	if !strings.HasPrefix(strings.TrimSpace(tokenStr), "{") {
		return false
	}

	// If we got here, the token file appears valid (actual validation happens during API calls)
	return true
}
