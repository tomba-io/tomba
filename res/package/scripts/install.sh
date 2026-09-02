#!/bin/sh
# Tomba CLI installer script
# Supports Linux, macOS (Darwin)
# Usage 1: curl -sSL https://raw.githubusercontent.com/tomba-io/tomba/master/res/package/scripts/install.sh | sh
# Usage 2: curl -sSL https://releases.tomba.io/install.sh | sh

set -e

REPO="tomba-io/tomba"
BINARY_NAME="tomba"
INSTALL_DIR="/usr/local/bin"

# Colors — use printf to produce real escape characters
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
CYAN=$(printf '\033[0;36m')
BOLD=$(printf '\033[1m')
NC=$(printf '\033[0m')

# Disable colors if not a terminal or NO_COLOR is set
if [ ! -t 1 ] || [ -n "${NO_COLOR:-}" ]; then
    RED="" GREEN="" YELLOW="" CYAN="" BOLD="" NC=""
fi

info() {
    printf "%s\n" "${CYAN}i${NC}  $1"
}

success() {
    printf "%s\n" "${GREEN}✓${NC}  $1"
}

warn() {
    printf "%s\n" "${YELLOW}!${NC}  $1"
}

error() {
    printf "%s\n" "${RED}x${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)  OS="linux" ;;
        Darwin*) OS="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
        *)       error "Unsupported operating system: $OS" ;;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)   ARCH="amd64" ;;
        i386|i686)       ARCH="386" ;;
        aarch64|arm64)   ARCH="arm64" ;;
        armv6l|armv7l)   ARCH="armv6" ;;
        ppc64le|ppc64)   ARCH="ppc64" ;;
        *)               error "Unsupported architecture: $ARCH" ;;
    esac
    echo "$ARCH"
}

# Get latest release version from GitHub
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
    else
        error "curl or wget is required to download tomba"
    fi

    if [ -z "$VERSION" ]; then
        error "Failed to determine latest version. Check https://github.com/${REPO}/releases"
    fi

    echo "$VERSION"
}

# Download file
download() {
    URL="$1"
    OUTPUT="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -sL "$URL" -o "$OUTPUT"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "$OUTPUT"
    else
        error "curl or wget is required"
    fi
}

banner() {
    printf "%s" "${GREEN}"
    cat << 'BANNER'

    ████████╗ ██████╗ ███╗   ███╗██████╗  █████╗    ██╗ ██████╗
    ╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗██╔══██╗   ██║██╔═══██╗
       ██║   ██║   ██║██╔████╔██║██████╔╝███████║   ██║██║   ██║
       ██║   ██║   ██║██║╚██╔╝██║██╔══██╗██╔══██║   ██║██║   ██║
       ██║   ╚██████╔╝██║ ╚═╝ ██║██████╔╝██║  ██║██╗██║╚██████╔╝
       ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝ ╚═════╝

BANNER
    printf "%s\n" "${NC}"
    printf "    %sCLI Installer%s — search or verify email addresses in seconds\n\n" "${BOLD}" "${NC}"
}

main() {
    banner

    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: ${BOLD}${OS}/${ARCH}${NC}"

    VERSION=$(get_latest_version)
    info "Latest version: ${BOLD}v${VERSION}${NC}"

    # Build download URL
    if [ "$OS" = "windows" ]; then
        FILENAME="${BINARY_NAME}_${OS}_${ARCH}.zip"
    else
        FILENAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"
    info "Downloading ${CYAN}${FILENAME}${NC}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download
    download "$DOWNLOAD_URL" "${TMP_DIR}/${FILENAME}"

    if [ ! -f "${TMP_DIR}/${FILENAME}" ]; then
        error "Download failed"
    fi

    success "Downloaded successfully"

    # Extract
    info "Extracting..."
    if [ "$OS" = "windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "${TMP_DIR}/${FILENAME}" -d "${TMP_DIR}"
        else
            error "unzip is required to extract the archive"
        fi
    else
        tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"
    fi

    if [ ! -f "${TMP_DIR}/${BINARY_NAME}" ]; then
        error "Binary not found in archive"
    fi

    success "Extracted successfully"

    # Install
    info "Installing to ${BOLD}${INSTALL_DIR}/${BINARY_NAME}${NC}..."

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        warn "Requires elevated permissions to install to ${INSTALL_DIR}"
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
            sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
        elif command -v doas >/dev/null 2>&1; then
            doas mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
            doas chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
        else
            error "sudo or doas is required to install to ${INSTALL_DIR}. Try: INSTALL_DIR=\$HOME/.local/bin sh install.sh"
        fi
    fi

    success "Installed ${BOLD}tomba v${VERSION}${NC} to ${INSTALL_DIR}/${BINARY_NAME}"

    # Verify
    if command -v tomba >/dev/null 2>&1; then
        INSTALLED_VERSION=$(tomba version 2>/dev/null || echo "unknown")
        success "Verification: ${INSTALLED_VERSION}"
    else
        warn "${INSTALL_DIR} may not be in your PATH. Add it with:"
        printf "    export PATH=\"%s:\$PATH\"\n" "$INSTALL_DIR"
    fi

    printf "\n%sGetting started:%s\n" "${BOLD}" "${NC}"
    printf "  tomba login          Sign in with your API key\n"
    printf "  tomba search -t      Search emails by domain\n"
    printf "  tomba chat           Start AI chat\n"
    printf "  tomba --help         Show all commands\n\n"
}

main "$@"
