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
	rootCmd.PersistentFlags().StringVar(&clientSecret, "secret", clientSecret, "file containing client secret JSON")
	rootCmd.PersistentFlags().StringVar(&cacheFile, "cache", cacheFile, "file containing oauth2 credential cache")
	rootCmd.PersistentFlags().StringVar(&domain, "domain", "", "domain for email addresses (e.g., example.com)")

	// Bind domain flag to viper
	if err := viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain")); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to bind domain flag: %s\n", err)
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
