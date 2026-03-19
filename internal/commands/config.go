package commands

import (
	"fmt"
	"os"

	"github.com/menloresearch/menlo-cli/internal/clients/platform"
	"github.com/menloresearch/menlo-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var defaultRobotCmd = &cobra.Command{
	Use:   "default-robot [robot-id]",
	Short: "Set or show the default robot",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Interactive mode - show selection
			return runRobotSelector()
		}

		// Set robot ID directly
		robotID := args[0]
		return saveDefaultRobot(robotID)
	},
}

var apikeyCmd = &cobra.Command{
	Use:   "apikey [key]",
	Short: "Manage your API key",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show current or instructions
			cfg, err := config.Load()
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("Get your API key from: https://platform.menlo.ai/account/api-keys")
					fmt.Println("Then run: menlo-cli config apikey <your-api-key>")
					return nil
				}
				return err
			}
			if cfg.APIKey != "" {
				fmt.Printf("API key: %s\n", cfg.APIKey)
			} else {
				fmt.Println("No API key set")
				fmt.Println("Get your API key from: https://platform.menlo.ai/account/api-keys")
				fmt.Println("Then run: menlo-cli config apikey <your-api-key>")
			}
			return nil
		}

		// Set API key
		apiKey := args[0]
		return saveAPIKey(apiKey)
	},
}

func runRobotSelector() error {
	client, err := platform.NewClient()
	if err != nil {
		return err
	}

	resp, err := client.ListRobots(100, "")
	if err != nil {
		return err
	}

	robots, err := resp.Robots()
	if err != nil {
		return err
	}

	if len(robots) == 0 {
		fmt.Println("No robots found")
		return nil
	}

	// Run TUI
	p := NewRobotSelector(robots)
	if err := p.Run(); err != nil {
		return err
	}

	if p.Selected() != "" {
		return saveDefaultRobot(p.Selected())
	}
	return nil
}

func saveDefaultRobot(robotID string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfg.DefaultRobotID = robotID

	if err := config.EnsureConfigDir(); err != nil {
		return err
	}

	data, err := config.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}

	fmt.Printf("Default robot set to: %s\n", robotID)
	return nil
}

func saveAPIKey(apiKey string) error {
	cfg, err := config.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		cfg = config.DefaultConfig()
	}

	cfg.APIKey = apiKey

	if err := config.EnsureConfigDir(); err != nil {
		return err
	}

	data, err := config.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}

	fmt.Println("API key saved successfully!")
	return nil
}

func init() {
	configCmd.AddCommand(defaultRobotCmd)
	configCmd.AddCommand(apikeyCmd)
}