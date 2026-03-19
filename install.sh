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

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*)    echo "windows";;
        MINGW*)     echo "windows";;
        *)          echo "unsupported";;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64)     echo "amd64";;
        arm64|aarch64) echo "arm64";;
        *)          echo "amd64";;
    esac
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

# Get latest version from GitHub
get_latest_version() {
    curl -sL "https://api.github.com/repos/menloresearch/cli/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | cut -d'v' -f2
}

# Install menlo
install() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    log_info "Detected OS: $OS, Arch: $ARCH"

    # Get version (default to latest if not set)
    VERSION="${VERSION:-$(get_latest_version)}"

    if [ -z "$VERSION" ]; then
        log_warn "Could not fetch latest version, using v0.0.4"
        VERSION="v0.0.4"
    fi

    # Remove 'v' prefix if present
    VERSION="${VERSION#v}"

    # Binary name
    BINARY_NAME="menlo"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="menlo.exe"
    fi

    # Download URL
    DOWNLOAD_URL="https://github.com/menloresearch/cli/releases/download/v${VERSION}/${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"

    log_info "Downloading menlo v${VERSION}..."
    log_info "URL: $DOWNLOAD_URL"

    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Download and extract
    if curl -sL "$DOWNLOAD_URL" -o "menlo.tar.gz"; then
        tar -xzf "menlo.tar.gz" || log_error "Failed to extract archive"
    else
        log_error "Failed to download menlo"
        cd /
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    # Check if binary exists
    if [ ! -f "$BINARY_NAME" ]; then
        # Try alternative naming pattern
        DOWNLOAD_URL="https://github.com/menloresearch/cli/releases/download/v${VERSION}/menlo_${VERSION}_${OS}_${ARCH}.tar.gz"
        log_info "Trying alternate URL: $DOWNLOAD_URL"

        curl -sL "$DOWNLOAD_URL" -o "menlo.tar.gz" || {
            log_error "Failed to download. Please check if version $VERSION exists."
            cd /
            rm -rf "$TEMP_DIR"
            exit 1
        }
        tar -xzf "menlo.tar.gz"
    fi

    # Install to /usr/local/bin (may need sudo)
    log_info "Installing to /usr/local/bin..."
    if [ -w /usr/local/bin ]; then
        cp "$BINARY_NAME" /usr/local/bin/menlo
        chmod +x /usr/local/bin/menlo
    else
        log_warn "Need sudo to install to /usr/local/bin"
        sudo cp "$BINARY_NAME" /usr/local/bin/menlo
        sudo chmod +x /usr/local/bin/menlo
    fi

    # Write version to config (preserve existing config)
    CONFIG_BASE_DIR=$(detect_config_dir)
    CONFIG_DIR="$CONFIG_BASE_DIR/menlo"
    mkdir -p "$CONFIG_DIR"
    CONFIG_FILE="$CONFIG_DIR/config.yaml"
    if [ -f "$CONFIG_FILE" ]; then
        # Use awk to update version in YAML - replace existing or add new
        awk -v ver="$VERSION" '
        /^version:/ { print "version: \"" ver "\""; found=1; next }
        { print }
        END { if (!found) print "version: \"" ver "\"" }
        ' "$CONFIG_FILE" > "$CONFIG_FILE.tmp" && mv "$CONFIG_FILE.tmp" "$CONFIG_FILE"
    else
        printf 'version: "%s"\n' "$VERSION" > "$CONFIG_FILE"
    fi
    log_info "Version $VERSION written to config"

    # Install shell completions to config directory (not modifying user shell files)
    COMPLETION_DIR="$CONFIG_DIR/completions"
    mkdir -p "$COMPLETION_DIR"

    # Detect current shell
    SHELL_NAME="$(basename "$SHELL" 2>/dev/null || echo "bash")"

    case "$SHELL_NAME" in
        zsh)
            "$BINARY_NAME" completion zsh > "$COMPLETION_DIR/zsh"
            log_info "Zsh completion installed to $COMPLETION_DIR/zsh"
            log_info "Add to your .zshrc: source $COMPLETION_DIR/zsh"
            ;;
        fish)
            "$BINARY_NAME" completion fish > "$COMPLETION_DIR/fish"
            log_info "Fish completion installed to $COMPLETION_DIR/fish"
            log_info "Add to your config.fish: source $COMPLETION_DIR/fish"
            ;;
        bash)
            "$BINARY_NAME" completion bash > "$COMPLETION_DIR/bash"
            log_info "Bash completion installed to $COMPLETION_DIR/bash"
            log_info "Add to your .bashrc: source $COMPLETION_DIR/bash"
            ;;
        *)
            # Install all completions
            "$BINARY_NAME" completion bash > "$COMPLETION_DIR/bash" 2>/dev/null || true
            "$BINARY_NAME" completion zsh > "$COMPLETION_DIR/zsh" 2>/dev/null || true
            "$BINARY_NAME" completion fish > "$COMPLETION_DIR/fish" 2>/dev/null || true
            log_info "Completions installed to $COMPLETION_DIR"
            ;;
    esac

    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"

    log_info "menlo v${VERSION} installed successfully!"
    log_info "Run 'menlo init' to get started"
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
            echo "Usage: sh install.sh"
            echo ""
            echo "Options:"
            echo "  -h, --help               Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Run install
install