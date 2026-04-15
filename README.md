# menlo

A CLI tool for accessing Menlo Robot.

## Quick Install

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/menloresearch/cli/release/install.sh)"

```

Or build from source:

```bash
git clone https://github.com/menloresearch/cli
cd cli
go build -o menlo ./cmd/menlo
```

## Quick Start

```bash
# Initialize the CLI (sets up API key and default robot)
menlo init

# Or configure manually
menlo config apikey <your-api-key>
menlo config default-robot <robot-id>
```

## Commands

### menlo init

Initialize menlo with your API key and select a default robot.

```bash
menlo init
```

### menlo robot

Manage robots.

#### List all robots

```bash
menlo robot list
```

#### Show robot status

```bash
menlo robot status                         # Use default robot
menlo robot status --robot-id <robot-id>  # Use specific robot
```

#### Join WebRTC session

```bash
menlo robot session                         # Use default robot
menlo robot session --robot-id <robot-id>  # Use specific robot
```

#### Set default robot

```bash
menlo robot connect              # Interactive selection
menlo robot connect <robot-id>   # Set directly
```

Same as `menlo config default-robot`. Sets the default robot for the CLI.

### menlo config

Manage configuration.

#### Set API key

```bash
menlo config apikey <your-api-key>
```

#### Show current API key

```bash
menlo config apikey
```

#### Set default robot

```bash
menlo config default-robot <robot-id>
```

Or interactively:

```bash
menlo config default-robot
```

## Configuration

Configuration is stored in:
- macOS: `~/Library/Application Support/menlo/config.yaml`
- Linux: `~/.config/menlo/config.yaml`
- Windows: `%APPDATA%\menlo\config.yaml`


# Uninstall
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/menloresearch/cli/release/uninstall.sh)"
```