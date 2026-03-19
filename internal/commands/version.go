package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtBy   = "goreleaser"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("menlo-cli %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  date: %s\n", date)
		fmt.Printf("  built by: %s\n", builtBy)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func init() {
	// Allow overriding version via ldflags
	if v := os.Getenv("MENLO_CLI_VERSION"); v != "" {
		version = v
	}
}