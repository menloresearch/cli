package commands

import (
	"fmt"

	"github.com/menloresearch/menlo-cli/internal/clients/platform"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize menlo-cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	fmt.Println("Welcome to menlo-cli!")
	fmt.Println()

	// Step 1: Get API key from TUI
	p := NewAPIKeyInput()
	if err := p.Run(); err != nil {
		return err
	}

	apiKey := p.Value()
	if apiKey == "" {
		fmt.Println("Initialization cancelled")
		return nil
	}

	// Save API key temporarily
	if err := saveAPIKey(apiKey); err != nil {
		return err
	}
	fmt.Println("API key saved!")

	// Step 2: Verify by listing robots
	fmt.Println("Verifying API key...")
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
		fmt.Println("No robots found. Initialization complete!")
		return nil
	}

	fmt.Printf("Found %d robot(s)\n", len(robots))

	// Step 3: Select default robot
	fmt.Println()
	fmt.Println("Select a default robot:")

	selector := NewRobotSelector(robots)
	if err := selector.Run(); err != nil {
		return err
	}

	selected := selector.Selected()
	if selected != "" {
		if err := saveDefaultRobot(selected); err != nil {
			return err
		}
	}

	fmt.Println()
	fmt.Println("Initialization complete!")
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}