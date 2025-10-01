package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags and parameters
var (
	cfgFile string
	domain  string
)

var rootCmd = &cobra.Command{
	Use:   "gac",
	Short: "CLI for administering Google Apps users",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.google-admin.yaml)")
	rootCmd.PersistentFlags().StringVar(&clientSecret, "client-secret", clientSecret, "file containing client secret JSON")
	rootCmd.PersistentFlags().StringVar(&cacheFile, "cache-file", cacheFile, "file containing oauth2 credential cache")
	rootCmd.PersistentFlags().StringVar(&domain, "domain", "", "domain for email addresses (e.g., example.com)")

	// Maintain backward compatibility with old flag names
	rootCmd.PersistentFlags().StringVar(&clientSecret, "secret", clientSecret, "deprecated: use --client-secret instead")
	rootCmd.PersistentFlags().StringVar(&cacheFile, "cache", cacheFile, "deprecated: use --cache-file instead")
	if err := rootCmd.PersistentFlags().MarkDeprecated("secret", "use --client-secret instead"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to mark secret flag as deprecated: %s\n", err)
	}
	if err := rootCmd.PersistentFlags().MarkDeprecated("cache", "use --cache-file instead"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to mark cache flag as deprecated: %s\n", err)
	}

	// Bind flags to viper
	if err := viper.BindPFlag("client-secret", rootCmd.PersistentFlags().Lookup("client-secret")); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind client-secret flag: %s\n", err)
	}
	if err := viper.BindPFlag("cache-file", rootCmd.PersistentFlags().Lookup("cache-file")); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind cache-file flag: %s\n", err)
	}
	if err := viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain")); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind domain flag: %s\n", err)
	}

	// Bind environment variables
	// Supports both GOOGLE_ADMIN_CLIENT_SECRET and GAC_CLIENT_SECRET
	if err := viper.BindEnv("client-secret", "GAC_CLIENT_SECRET", "GOOGLE_ADMIN_CLIENT_SECRET"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind client-secret env vars: %s\n", err)
	}
	if err := viper.BindEnv("cache-file", "GAC_CACHE_FILE", "GOOGLE_ADMIN_CACHE_FILE"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind cache-file env vars: %s\n", err)
	}
	if err := viper.BindEnv("domain", "GAC_DOMAIN", "GOOGLE_ADMIN_DOMAIN"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind domain env vars: %s\n", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".google-admin") // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // adding home directory as first search path
	viper.SetEnvPrefix("google_admin")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

// getDomain returns the configured domain, with fallback to default
func getDomain() string {
	if domain != "" {
		return domain
	}
	return viper.GetString("domain")
}
