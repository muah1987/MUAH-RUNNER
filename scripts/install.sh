#!/usr/bin/env bash
set -euo pipefail

REPO="muah1987/muah-runner"
BINARY="muah-runner"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

LATEST=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
  echo "Could not determine latest release. Please download manually from:"
  echo "  https://github.com/${REPO}/releases"
  exit 1
fi

URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}-${OS}-${ARCH}"
echo "Downloading muah-runner ${LATEST} for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "${BINARY}"
chmod +x "${BINARY}"
sudo mv "${BINARY}" "${INSTALL_DIR}/"
echo "Installed muah-runner to ${INSTALL_DIR}/${BINARY}"
