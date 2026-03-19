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

    # Remove completion lines from shell rc files
    log_info "Removing shell completion lines..."

    # Zsh
    if [ -f "$HOME/.zshrc" ]; then
        sed -i '' '/menlo\/completions\//d' "$HOME/.zshrc" 2>/dev/null || \
        sed -i '/menlo\/completions\//d' "$HOME/.zshrc" 2>/dev/null || true
        log_info "Removed completion from .zshrc"
    fi

    # Bash
    if [ -f "$HOME/.bashrc" ]; then
        sed -i '' '/menlo\/completions\//d' "$HOME/.bashrc" 2>/dev/null || \
        sed -i '/menlo\/completions\//d' "$HOME/.bashrc" 2>/dev/null || true
        log_info "Removed completion from .bashrc"
    fi

    # Fish
    FISH_CONFIG="$HOME/.config/fish/config.fish"
    if [ -f "$FISH_CONFIG" ]; then
        sed -i '' '/menlo\/completions\//d' "$FISH_CONFIG" 2>/dev/null || \
        sed -i '/menlo\/completions\//d' "$FISH_CONFIG" 2>/dev/null || true
        log_info "Removed completion from config.fish"
    fi
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