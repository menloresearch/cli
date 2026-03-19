# menlo-cli

A CLI tool for Menlo research and development.

## Installation

```bash
go install github.com/menloresearch/menlo-cli@latest
```

Or build from source:

```bash
git clone https://github.com/menloresearch/menlo-cli
cd menlo-cli
go build -o menlo-cli ./cmd/menlo-cli
```

## Quick Start

```bash
# Initialize the CLI (sets up API key and default robot)
menlo-cli init

# Or configure manually
menlo-cli config apikey <your-api-key>
menlo-cli config default-robot <robot-id>
```

## Commands

### menlo-cli init

Initialize menlo-cli with your API key and select a default robot.

```bash
menlo-cli init
```

### menlo-cli robot

Manage robots.

#### List all robots

```bash
menlo-cli robot list
```

#### Show robot status

```bash
menlo-cli robot status                         # Use default robot
menlo-cli robot status --robot-id <robot-id>  # Use specific robot
```

#### Send action to robot

```bash
menlo-cli robot action forward                 # Use default robot
menlo-cli robot action left --robot-id <id>   # Use specific robot
```

Available actions:
- `forward` - Move the robot forward
- `backward` - Move the robot backward
- `left` - Move the robot left
- `right` - Move the robot right
- `turn-left` - Turn the robot left
- `turn-right` - Turn the robot right

### menlo-cli config

Manage configuration.

#### Set API key

```bash
menlo-cli config apikey <your-api-key>
```

#### Show current API key

```bash
menlo-cli config apikey
```

#### Set default robot

```bash
menlo-cli config default-robot <robot-id>
```

Or interactively:

```bash
menlo-cli config default-robot
```

## Configuration

Configuration is stored in:
- macOS: `~/Library/Application Support/menlo-cli/config.yaml`
- Linux: `~/.config/menlo-cli/config.yaml`
- Windows: `%APPDATA%\menlo-cli\config.yaml`

## Shell Completion

```bash
# Bash
menlo-cli completion bash > /etc/bash_completion.d/menlo-cli

# Zsh
menlo-cli completion zsh > "${fpath[1]}/_menlo-cli"

# Fish
menlo-cli completion fish > ~/.config/fish/completions/menlo-cli.fish
```