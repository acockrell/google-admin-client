package cmd

import "github.com/spf13/cobra"

// userCmd represents the user command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group operations",
}

func init() {
	rootCmd.AddCommand(groupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
