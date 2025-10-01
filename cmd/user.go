package cmd

import (
	"crypto/rand"
	"math/big"

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
	const letterRunes = "abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

	b := make([]byte, length)
	maxIdx := big.NewInt(int64(len(letterRunes)))

	for i := range b {
		n, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			// If crypto/rand fails, this is a critical error
			panic("failed to generate secure random number: " + err.Error())
		}
		b[i] = letterRunes[n.Int64()]
	}
	return string(b)
}
