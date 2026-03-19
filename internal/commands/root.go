package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "menlo",
	Short:             "Menlo CLI - A CLI tool for Menlo",
	Long:              `A CLI tool for Menlo research and development.`,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
}