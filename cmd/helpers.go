package cmd

import (
	"fmt"
	"os"
)

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
