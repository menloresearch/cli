package commands

import (
	"fmt"
	"net/url"

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
	Use:               "action <action>",
	Short:             "Send an action to a robot",
	ValidArgs:         platform.ValidSemanticCommands,
	Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]

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

var robotSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Join a WebRTC session for a robot",
	Long: `Join a session to connect to a robot via WebRTC.
Returns an SFU endpoint and WebRTC token for connecting.

Examples:
  menlo robot session --robot-id <robot-id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		return createRobotSession(robotID)
	},
}

var robotConnectCmd = &cobra.Command{
	Use:   "connect [robot-id]",
	Short: "Set or show the default robot",
	Long: `Set or show the default robot. Same as 'menlo config default-robot'.

Examples:
  menlo robot connect              # Interactive selection
  menlo robot connect <robot-id>   # Set directly`,
	Args: cobra.MaximumNArgs(1),
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

	if robot.Status != nil {
		fmt.Printf("\nStatus:\n")
		fmt.Printf("  Battery: %d%%", robot.Status.Battery.Level)
		if robot.Status.Battery.Charging {
			fmt.Printf(" (charging)")
		}
		fmt.Printf("\n")
		fmt.Printf("  Last Updated: %s\n", robot.Status.LocalTimestamp().Format("2006-01-02 15:04:05"))
	}

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

func createRobotSession(robotID string) error {
	client, err := platform.NewClient()
	if err != nil {
		return err
	}

	session, err := client.CreateSession(robotID)
	if err != nil {
		return err
	}

	// Generate meet link
	meetURL := generateMeetLink(session.SFUEndpoint, session.WebRTCToken)

	fmt.Printf("Session created for robot %s\n\n", robotID)
	fmt.Printf("Connection Endpoint: %s\n", session.SFUEndpoint)
	fmt.Printf("Agent Token: %s\n\n", session.WebRTCToken)
	fmt.Printf("You can also paste this URL into your browser to debug the sim: %s\n", meetURL)

	return nil
}

func generateMeetLink(sfuEndpoint, token string) string {
	baseURL := "https://meet.livekit.io/custom"
	params := url.Values{}
	params.Set("liveKitUrl", sfuEndpoint)
	params.Set("token", token)
	return baseURL + "?" + params.Encode()
}

var robotSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Download latest snapshot from a robot",
	Long: `Download the latest snapshot image from a robot.
The image is saved to ~/.config/menlo/snapshot/{robot-id}/latest.jpeg

Examples:
  menlo robot snapshot
  menlo robot snapshot --robot-id <robot-id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		client, err := platform.NewClient()
		if err != nil {
			return err
		}

		path, err := client.GetSnapshot(robotID)
		if err != nil {
			return err
		}

		fmt.Printf("Snapshot saved to: %s\n", path)
		return nil
	},
}

func init() {
	robotStatusCmd.Flags().String("robot-id", "", "Robot ID")
	robotActionCmd.Flags().String("robot-id", "", "Robot ID")
	robotSessionCmd.Flags().String("robot-id", "", "Robot ID")
	robotSnapshotCmd.Flags().String("robot-id", "", "Robot ID")
	robotCmd.AddCommand(robotListCmd)
	robotCmd.AddCommand(robotStatusCmd)
	// robotCmd.AddCommand(robotActionCmd)   // disabled
	robotCmd.AddCommand(robotSessionCmd)
	// robotCmd.AddCommand(robotSnapshotCmd) // disabled
	robotCmd.AddCommand(robotConnectCmd)
	rootCmd.AddCommand(robotCmd)
}