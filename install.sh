#!/bin/sh
set -e

REPO="onurkerem/auto-message"
BINARY="auto-message"

get_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*) echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "unknown" ;;
    esac
}

OS=$(get_os)
ARCH=$(get_arch)

if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
    echo "Error: unsupported operating system or architecture."
    exit 1
fi

if [ "$OS" = "windows" ]; then
    EXT=".zip"
else
    EXT=".tar.gz"
fi

# Get latest release tag
TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
    echo "Error: could not determine latest version."
    exit 1
fi

FILENAME="${BINARY}_${TAG#v}_${OS}_${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${FILENAME}"

echo "Installing ${BINARY} ${TAG} for ${OS}/${ARCH}..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$DOWNLOAD_URL" -o "${TMPDIR}/${FILENAME}"

cd "$TMPDIR"
if [ "$EXT" = ".tar.gz" ]; then
    tar xzf "${FILENAME}"
else
    unzip -o "${FILENAME}"
fi

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Requires sudo to install to ${INSTALL_DIR}..."
    sudo mv "$BINARY" "${INSTALL_DIR}/${BINARY}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY}"
else
    mv "$BINARY" "${INSTALL_DIR}/${BINARY}"
    chmod +x "${INSTALL_DIR}/${BINARY}"
fi

echo "Successfully installed ${BINARY} ${TAG} to ${INSTALL_DIR}/${BINARY}"
"${BINARY}" --version
