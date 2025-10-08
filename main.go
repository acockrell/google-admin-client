package main

import "github.com/acockrell/google-admin-client/cmd"

// Build information variables set by ldflags during build
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

func main() {
	// Pass build information to cmd package
	cmd.SetVersionInfo(Version, Commit, Date, BuiltBy)
	cmd.Execute()
}
