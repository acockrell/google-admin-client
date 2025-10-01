package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize client authentication configuration",
	Run:   initRunFunc,
	Long: `
Initialize client authentication configuration.

Usage
-----

  $ gac init [-f]

Overview
--------

This app uses a OAuth2 to authenticate with Google.  A client ID and key need
need to be created as part of a Google application, and the client_secret.json
file made available to this CLI.

Upon first execution, an OAuth2 authentication is performed, in which a URL is
displayed to stdout.  One must navigate to the URL and copy/paste the resulting
key back to the CLI.  This is a one-time process, as the credential is cached
~/.credentials/gac.json and used for subsequent executions.


Google Setup Instructions
-------------------------

TBD

`,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force reauthentication")

}

func initRunFunc(cmd *cobra.Command, args []string) {
	if force {
		fmt.Fprintf(os.Stderr, "removing cache file\n")
		// remove ~/.credentials/google-admin.json
		cacheFile, err := tokenCacheFile()
		if err != nil {
			exitWithError(err.Error())
		}

		err = os.Remove(cacheFile)
		if err != nil && !os.IsNotExist(err) {
			exitWithError(err.Error())
		}
	}

	_, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

}
