package commands

import (
	"fmt"

	"github.com/menloresearch/cli/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			// If no config, use default
			cfg = config.DefaultConfig()
		}

		version := cfg.Version
		if version == "" {
			version = "dev"
		}

		fmt.Printf("menlo %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}