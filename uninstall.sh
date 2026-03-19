#!/bin/sh
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

log_warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

# Detect config directory (matches Go's os.UserConfigDir)
detect_config_dir() {
    case "$(uname -s)" in
        Darwin*)    echo "$HOME/Library/Application Support";;
        Linux*)     echo "$HOME/.config";;
        CYGWIN*|MINGW*) echo "$APPDATA";;
        *)          echo "$HOME/.config";;
    esac
}

# Prompt user for confirmation
prompt() {
    if [ "$ASSUME_YES" = "1" ]; then
        return 0
    fi
    printf "%s [y/N]: " "$1"
    read -r answer
    case "$answer" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

uninstall() {
    log_info "Uninstalling menlo..."

    # Remove binary
    log_info "Removing binary from /usr/local/bin..."
    if [ -f /usr/local/bin/menlo ]; then
        if [ -w /usr/local/bin/menlo ]; then
            rm /usr/local/bin/menlo
        else
            log_warn "Need sudo to remove binary"
            sudo rm /usr/local/bin/menlo
        fi
        log_info "Binary removed"
    else
        log_info "Binary not found at /usr/local/bin/menlo"
    fi

    # Detect config directory
    CONFIG_BASE_DIR=$(detect_config_dir)
    CONFIG_DIR="$CONFIG_BASE_DIR/menlo"

    # Ask about config removal
    if [ -d "$CONFIG_DIR" ]; then
        if prompt "Remove config directory ($CONFIG_DIR)?"; then
            rm -rf "$CONFIG_DIR"
            log_info "Config directory removed"
        else
            log_info "Keeping config directory"
        fi
    else
        log_info "Config directory not found"
    fi

    log_info "Uninstall complete!"

    # Show instructions for shell completions
    case "$(uname -s)" in
        Darwin*)
            CONFIG_DIR="$HOME/Library/Application Support/menlo"
            ;;
        Linux*)
            CONFIG_DIR="$HOME/.config/menlo"
            ;;
        CYGWIN*|MINGW*)
            CONFIG_DIR="$APPDATA/menlo"
            ;;
    esac

    log_info "Shell completions are in: $CONFIG_DIR/completions/"
    log_info "Remove these lines from your shell rc file if you added them:"
    log_info "  source $CONFIG_DIR/completions/zsh   # for zsh"
    log_info "  source $CONFIG_DIR/completions/fish  # for fish"
    log_info "  source $CONFIG_DIR/completions/bash  # for bash"
}

# Check if curl is installed
if ! command -v curl >/dev/null 2>&1; then
    log_error "curl is required but not installed. Please install curl first."
    exit 1
fi

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        -h|--help)
            echo "Usage: sh uninstall.sh"
            echo ""
            echo "Options:"
            echo "  -h, --help               Show this help message"
            exit 0
            ;;
        -y|--yes)
            # Non-interactive mode - assume yes to all prompts
            export ASSUME_YES=1
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Run uninstall
uninstall