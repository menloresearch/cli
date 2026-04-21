package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	skillTargetClaude   = "claude"
	skillTargetOpenCode = "opencode"
	skillTargetCodex    = "codex"
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage agent skills",
}

var skillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Menlo skill for your coding agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		selectedTarget := strings.TrimSpace(strings.ToLower(target))
		if selectedTarget == "" {
			selectedTarget, err = promptSkillTarget()
			if err != nil {
				return err
			}
		}

		installPath, err := resolveSkillInstallPath(selectedTarget)
		if err != nil {
			return err
		}

		if err := writeSkillFile(installPath, force); err != nil {
			return err
		}

		fmt.Printf("Menlo skill installed for %s\n", skillTargetLabel(selectedTarget))
		fmt.Printf("Path: %s\n", installPath)
		return nil
	},
}

var skillUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Menlo skill from your coding agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}

		selectedTarget := strings.TrimSpace(strings.ToLower(target))
		if selectedTarget == "" {
			selectedTarget, err = promptSkillTarget()
			if err != nil {
				return err
			}
		}

		installPath, err := resolveSkillInstallPath(selectedTarget)
		if err != nil {
			return err
		}

		if err := uninstallSkillFile(installPath); err != nil {
			return err
		}

		fmt.Printf("Menlo skill uninstalled for %s\n", skillTargetLabel(selectedTarget))
		fmt.Printf("Path: %s\n", installPath)
		return nil
	},
}

func promptSkillTarget() (string, error) {
	selector := NewSkillTargetSelector([]skillTargetItem{
		{target: skillTargetClaude, label: "Claude Code"},
		{target: skillTargetOpenCode, label: "OpenCode"},
		{target: skillTargetCodex, label: "Codex"},
	})

	if err := selector.Run(); err != nil {
		return "", err
	}

	if selector.Selected() == "" {
		return "", fmt.Errorf("selection cancelled")
	}

	return selector.Selected(), nil
}

func resolveSkillInstallPath(target string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var basePath string
	switch target {
	case skillTargetClaude:
		basePath = filepath.Join(homeDir, ".claude")
	case skillTargetOpenCode:
		basePath = filepath.Join(homeDir, ".opencode")
	case skillTargetCodex:
		basePath = filepath.Join(homeDir, ".codex")
	default:
		return "", fmt.Errorf("invalid target %q (use: claude, opencode, codex)", target)
	}

	info, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%s is not installed (missing base directory: %s)", skillTargetLabel(target), basePath)
		}
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("invalid base path: %s", basePath)
	}

	return filepath.Join(basePath, "skills", "menlo", "SKILL.md"), nil
}

func writeSkillFile(path string, force bool) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("skill already exists at %s (use --force to overwrite)", path)
	}

	return os.WriteFile(path, []byte(menloSkillContent), 0o644)
}

func uninstallSkillFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill file not found at %s", path)
		}
		return err
	}

	return os.Remove(path)
}

func skillTargetLabel(target string) string {
	switch target {
	case skillTargetClaude:
		return "Claude Code"
	case skillTargetOpenCode:
		return "OpenCode"
	case skillTargetCodex:
		return "Codex"
	default:
		return target
	}
}

const menloSkillContent = `---
name: menlo
description: Control Menlo robots safely from the CLI with timed movement patterns.
license: MIT
compatibility: opencode, claude, codex
metadata:
  version: "1.0.0"
---

# Menlo Robot Control

Use this skill to control Menlo robots safely via the menlo CLI.

## Purpose

This skill teaches the agent how to:
- Move a robot with basic actions
- Select/set a default robot when needed
- Execute timed movement sequences safely
- Handle errors with safe recovery

## Preconditions

Before sending movement commands, ensure:
1. menlo CLI is installed and available
2. API key is configured (menlo config apikey)
3. A target robot is available either:
   - as the current default robot, or
   - by setting one with:
     - menlo robot connect (interactive), or
     - menlo robot connect <robot-id>

## Command Reference

Movement command format:

menlo robot action <action> [--robot-id <robot-id>]

Supported actions:
- forward
- backward
- left
- right
- turn-left
- turn-right
- stop

Useful helpers:
- menlo robot list
- menlo robot status [--robot-id <robot-id>]
- menlo robot connect
- menlo robot connect <robot-id>

## Safety Rules

1. Prefer short, step-by-step actions over long assumptions.
2. Do not send stop immediately after motion; allow movement time first.
3. Use explicit motion windows (for example 300ms to 1500ms) then stop.
4. If any command fails, send stop (best effort) and report the failure.
5. If no robot is selected/configured, resolve robot first.

## Timed Motion Pattern

For motions that should visibly move the robot:
1. Send movement action.
2. Wait a short duration.
3. Send stop.

Example (conceptual):
1. menlo robot action forward
2. wait ~800ms
3. menlo robot action stop

Use shorter durations for small nudges and longer durations only when requested.

## Standard Execution Pattern

When asked to move:
1. Ensure target robot is resolvable.
2. Execute one action at a time.
3. For movement actions, include a delay before stop.
4. Report exactly what was run.

Example with explicit robot:
menlo robot action turn-left --robot-id rb_xxx
# wait ~500ms
menlo robot action stop --robot-id rb_xxx
menlo robot action forward --robot-id rb_xxx
# wait ~900ms
menlo robot action stop --robot-id rb_xxx

## Error Handling

If a command fails:
1. Attempt:
   menlo robot action stop [--robot-id <robot-id>]
2. Report:
   - Which command failed
   - Error text
   - Recovery attempted

If no default robot is configured:
1. menlo robot list
2. menlo robot connect (interactive) or menlo robot connect <robot-id>

## Notes

- Use only supported menlo robot action commands.
- Do not invent unsupported actions or payloads.
- Keep outputs concise and action-oriented.
`

func init() {
	skillInstallCmd.Flags().String("target", "", "Install target: claude, opencode, codex")
	skillInstallCmd.Flags().Bool("force", false, "Overwrite existing skill file")
	skillUninstallCmd.Flags().String("target", "", "Uninstall target: claude, opencode, codex")
	skillCmd.AddCommand(skillInstallCmd)
	skillCmd.AddCommand(skillUninstallCmd)
	rootCmd.AddCommand(skillCmd)
}
