#!/usr/bin/env bash
# scripts/setup.sh — One-command bootstrap for muah-runner on any platform.
# Works on Linux, macOS, and Windows (via Git Bash / WSL).
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/muah1987/muah-runner/main/scripts/setup.sh | bash
#   # or locally:
#   bash scripts/setup.sh

set -euo pipefail

REPO="muah1987/muah-runner"
BINARY="muah-runner"
INSTALL_DIR="${MUAH_INSTALL_DIR:-/usr/local/bin}"
VERSION="${MUAH_VERSION:-latest}"

# ── Colour helpers ────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()    { echo -e "${BLUE}[muah-runner]${NC} $*"; }
success() { echo -e "${GREEN}[muah-runner]${NC} ✓ $*"; }
warn()    { echo -e "${YELLOW}[muah-runner]${NC} ⚠ $*"; }
die()     { echo -e "${RED}[muah-runner]${NC} ✗ $*" >&2; exit 1; }

# ── Detect OS/arch ────────────────────────────────────────────────────────────
detect_os_arch() {
  OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
  ARCH="$(uname -m)"
  case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) die "Unsupported architecture: $ARCH" ;;
  esac
  case "$OS" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) die "Unsupported OS: $OS" ;;
  esac
  info "Detected platform: ${OS}/${ARCH}"
}

# ── Check dependencies ────────────────────────────────────────────────────────
check_deps() {
  local missing=()
  for cmd in curl git; do
    command -v "$cmd" &>/dev/null || missing+=("$cmd")
  done
  if [ ${#missing[@]} -gt 0 ]; then
    die "Missing required tools: ${missing[*]}. Please install them and re-run."
  fi
}

# ── Install Go (if not present) ───────────────────────────────────────────────
ensure_go() {
  if command -v go &>/dev/null; then
    GO_VERSION="$(go version 2>/dev/null | awk '{print $3}')"
    success "Go already installed: ${GO_VERSION}"
    return
  fi
  info "Go not found — installing Go 1.24..."
  local GO_VERSION_NUM="1.24.1"
  local tarball="go${GO_VERSION_NUM}.${OS}-${ARCH}.tar.gz"
  local url="https://go.dev/dl/${tarball}"
  local tmpdir
  tmpdir="$(mktemp -d)"
  info "Downloading ${url}..."
  curl -fsSL "$url" -o "${tmpdir}/${tarball}" || die "Failed to download Go"
  if [ "$OS" = "linux" ] || [ "$OS" = "darwin" ]; then
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "${tmpdir}/${tarball}"
    export PATH="/usr/local/go/bin:$PATH"
    echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc || true
    echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.zshrc 2>/dev/null || true
  fi
  rm -rf "$tmpdir"
  success "Go ${GO_VERSION_NUM} installed"
}

# ── Build or download muah-runner binary ─────────────────────────────────────
install_binary() {
  if command -v go &>/dev/null; then
    info "Building muah-runner from source..."
    local tmpdir
    tmpdir="$(mktemp -d)"
    git clone --depth=1 "https://github.com/${REPO}.git" "${tmpdir}/muah-runner" 2>/dev/null \
      || die "Failed to clone repository"
    (cd "${tmpdir}/muah-runner" && go build -o "${BINARY}" ./cmd/muah-runner/) \
      || die "Build failed"
    if [ -w "$INSTALL_DIR" ]; then
      mv "${tmpdir}/muah-runner/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
      sudo mv "${tmpdir}/muah-runner/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi
    rm -rf "$tmpdir"
    success "muah-runner installed to ${INSTALL_DIR}/${BINARY}"
  else
    die "Go is required to build muah-runner. Install Go first or set GOPATH."
  fi
}

# ── Verify installation ───────────────────────────────────────────────────────
verify() {
  if command -v muah-runner &>/dev/null; then
    success "muah-runner $(muah-runner --version 2>/dev/null || echo '(version unknown)') is ready"
  elif [ -f "${INSTALL_DIR}/${BINARY}" ]; then
    success "${INSTALL_DIR}/${BINARY} installed (add ${INSTALL_DIR} to PATH)"
  else
    warn "Binary not found in PATH — add ${INSTALL_DIR} to your PATH"
  fi
}

# ── Optional: bootstrap .muah/ in current dir ─────────────────────────────────
maybe_init() {
  if [ -f "go.mod" ] || [ -f "package.json" ] || [ -f "Makefile" ] || [ -f ".git/config" ]; then
    echo ""
    info "Looks like a project repo. Bootstrap .muah/ now? [Y/n] (auto-yes in 10s)"
    local answer
    read -t 10 -r answer || answer="y"
    answer="${answer:-y}"
    if [[ "$answer" =~ ^[Yy]$ ]]; then
      "${INSTALL_DIR}/${BINARY}" init || warn "init failed — run 'muah-runner init' manually"
    fi
  fi
}

# ── Main ──────────────────────────────────────────────────────────────────────
main() {
  echo ""
  echo "  🚀 muah-runner setup"
  echo "  ────────────────────────────────────"
  echo ""

  check_deps
  detect_os_arch
  ensure_go
  install_binary
  verify

  echo ""
  success "Setup complete!"
  echo ""
  echo "  Next steps:"
  echo "    muah-runner init              # Bootstrap .muah/ in your project"
  echo "    muah-runner run 'my task'     # Run a task"
  echo "    muah-runner status            # Check health"
  echo "    muah-runner --help            # Full CLI reference"
  echo ""

  maybe_init
}

main "$@"
