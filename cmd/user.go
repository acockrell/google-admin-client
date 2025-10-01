package cmd

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "List & modify users",
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func randomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
