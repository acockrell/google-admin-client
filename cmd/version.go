package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Build information (set from main package)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Version flags
var (
	versionShort bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for gac, including version number,
git commit hash, build date, and Go runtime information.

Examples:
  # Show full version information
  gac version

  # Show only version number
  gac version --short
`,
	Run: versionRunFunc,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionShort, "short", "s", false, "show version number only")
}

// SetVersionInfo sets the version information from main package
func SetVersionInfo(ver, cmt, dt, blt string) {
	version = ver
	commit = cmt
	date = dt
	builtBy = blt
}

func versionRunFunc(cmd *cobra.Command, args []string) {
	if versionShort {
		fmt.Println(version)
		return
	}

	fmt.Printf("gac version %s\n", version)
	fmt.Printf("  Commit:     %s\n", commit)
	fmt.Printf("  Built:      %s by %s\n", date, builtBy)
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
