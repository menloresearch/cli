package commands

import (
	"fmt"

	"github.com/menloresearch/cli/internal/clients/platform"
	"github.com/menloresearch/cli/internal/config"
	"github.com/spf13/cobra"
)

var robotCmd = &cobra.Command{
	Use:   "robot",
	Short: "Manage robots",
}

var robotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all robots",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		fmt.Println("Robots:")
		for _, r := range robots {
			fmt.Printf("  %s (%s) - %s\n", r.Name, r.Type, r.ID)
		}

		return nil
	},
}

var robotStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show robot status",
	RunE: func(cmd *cobra.Command, args []string) error {
		robotID, err := cmd.Flags().GetString("robot-id")
		if err != nil {
			return err
		}

		// If no robot ID provided, try default
		if robotID == "" {
			cfg, err := config.Load()
			if err != nil {
				if config.IsNotExist(err) {
					return fmt.Errorf("please specify a robot ID with --robot-id flag")
				}
				return err
			}
			robotID = cfg.DefaultRobotID
		}

		// If still no robot ID, ask user to select
		if robotID == "" {
			robotID, err = selectRobotInteractive()
			if err != nil {
				return err
			}
			if robotID == "" {
				return nil
			}
		}

		return showRobotStatus(robotID)
	},
}

var robotActionCmd = &cobra.Command{
	Use:   "action <action>",
	Short: "Send an action to a robot",
	Long: `Available actions:
  forward     Move the robot forward
  backward    Move the robot backward
  left        Move the robot left
  right       Move the robot right
  turn-left   Turn the robot left
  turn-right  Turn the robot right

Examples:
  menlo robot action forward
  menlo robot action left --robot-id <robot-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]

		// Validate action
		valid := false
		for _, v := range platform.ValidSemanticCommands {
			if action == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid action: %s. Valid actions: forward, backward, left, right, turn-left, turn-right", action)
		}

		robotID, err := cmd.Flags().GetString("robot-id")
		if err != nil {
			return err
		}

		// If no robot ID provided, try default
		if robotID == "" {
			cfg, err := config.Load()
			if err != nil {
				if !config.IsNotExist(err) {
					return err
				}
			}
			if cfg != nil {
				robotID = cfg.DefaultRobotID
			}
		}

		// If still no robot ID, ask user to select
		if robotID == "" {
			robotID, err = selectRobotInteractive()
			if err != nil {
				return err
			}
			if robotID == "" {
				return nil
			}
		}

		return sendRobotAction(robotID, action)
	},
}

// selectRobotInteractive prompts user to select a robot from list
func selectRobotInteractive() (string, error) {
	client, err := platform.NewClient()
	if err != nil {
		return "", err
	}

	resp, err := client.ListRobots(100, "")
	if err != nil {
		return "", err
	}

	robots, err := resp.Robots()
	if err != nil {
		return "", err
	}

	if len(robots) == 0 {
		return "", fmt.Errorf("no robots found")
	}

	selector := NewRobotSelector(robots)
	if err := selector.Run(); err != nil {
		return "", err
	}

	return selector.Selected(), nil
}

func showRobotStatus(robotID string) error {
	client, err := platform.NewClient()
	if err != nil {
		return err
	}

	robot, err := client.GetRobot(robotID)
	if err != nil {
		return err
	}

	fmt.Printf("Robot: %s\n", robot.Name)
	fmt.Printf("ID: %s\n", robot.ID)
	fmt.Printf("Model: %s\n", robot.Model)
	fmt.Printf("Type: %s\n", robot.Type)

	return nil
}

func sendRobotAction(robotID, action string) error {
	client, err := platform.NewClient()
	if err != nil {
		return err
	}

	if err := client.SendSemanticCommand(robotID, action); err != nil {
		return err
	}

	fmt.Printf("Action '%s' sent to robot %s\n", action, robotID)
	return nil
}

func init() {
	robotStatusCmd.Flags().String("robot-id", "", "Robot ID")
	robotActionCmd.Flags().String("robot-id", "", "Robot ID")
	robotCmd.AddCommand(robotListCmd)
	robotCmd.AddCommand(robotStatusCmd)
	robotCmd.AddCommand(robotActionCmd)
	rootCmd.AddCommand(robotCmd)
}