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
	Version:           versionString(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{printf \"%s %s\\n\" .Name .Version}}")
	rootCmd.Flags().BoolP("version", "V", false, "Print the version number")
	rootCmd.AddCommand(configCmd)
}
