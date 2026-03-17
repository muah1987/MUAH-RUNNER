#!/usr/bin/env bash
# scripts/install-docker-mcp.sh — Install and configure the Docker MCP server
# for muah-runner's self-healing simulation engine.
# Platform-agnostic: Linux, macOS, Windows (Git Bash / WSL).
#
# Usage:
#   bash scripts/install-docker-mcp.sh
#   # or via muah-runner:
#   muah-runner mcp install docker

set -euo pipefail

MUAH_DIR="${MUAH_DIR:-.muah}"
MCP_CONFIG="${MUAH_DIR}/mcp/native/docker-mcp.json"

# ── Colour helpers ────────────────────────────────────────────────────────────
BLUE='\033[0;34m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'
info()    { echo -e "${BLUE}[docker-mcp]${NC} $*"; }
success() { echo -e "${GREEN}[docker-mcp]${NC} ✓ $*"; }
warn()    { echo -e "${YELLOW}[docker-mcp]${NC} ⚠ $*"; }
die()     { echo -e "${RED}[docker-mcp]${NC} ✗ $*" >&2; exit 1; }

# ── Check Docker ──────────────────────────────────────────────────────────────
check_docker() {
  if ! command -v docker &>/dev/null; then
    die "Docker is required. Install from https://docs.docker.com/get-docker/"
  fi
  if ! docker info &>/dev/null 2>&1; then
    die "Docker daemon is not running. Start Docker and re-run this script."
  fi
  DOCKER_VERSION="$(docker --version | awk '{print $3}' | tr -d ',')"
  success "Docker ${DOCKER_VERSION} available"
}

# ── Pull base images for simulations ─────────────────────────────────────────
pull_base_images() {
  info "Pulling muah-runner simulation base images..."
  local images=("ubuntu:22.04" "alpine:3.19" "node:20-alpine" "golang:1.24-alpine")
  for img in "${images[@]}"; do
    info "  Pulling ${img}..."
    docker pull "$img" --quiet || warn "  Failed to pull ${img} (non-fatal)"
  done
  success "Base images ready"
}

# ── Create muah-runner Docker network ────────────────────────────────────────
create_network() {
  local network="muah-runner-net"
  if docker network inspect "$network" &>/dev/null 2>&1; then
    success "Docker network '${network}' already exists"
  else
    info "Creating Docker network '${network}'..."
    docker network create "$network" --driver bridge || warn "Network creation failed (non-fatal)"
    success "Docker network '${network}' created"
  fi
}

# ── Write Docker MCP server config to .muah/ ─────────────────────────────────
write_mcp_config() {
  if [ ! -d "$MUAH_DIR" ]; then
    warn ".muah/ directory not found — run 'muah-runner init' first"
    return
  fi
  mkdir -p "$(dirname "$MCP_CONFIG")"
  cat > "$MCP_CONFIG" << MCPEOF
{
  "name": "docker-mcp",
  "type": "native",
  "status": "installed",
  "version": "$(docker --version | awk '{print $3}' | tr -d ',')",
  "capabilities": [
    "run_container",
    "simulate_failure",
    "matrix_test",
    "selfheal",
    "cleanup"
  ],
  "base_images": ["ubuntu:22.04", "alpine:3.19", "node:20-alpine", "golang:1.24-alpine"],
  "network": "muah-runner-net",
  "simulation_dir": ".muah/docker/simulations",
  "installed_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
MCPEOF
  success "${MCP_CONFIG} written"
}

# ── Write docker-compose template for simulations ─────────────────────────────
write_compose_template() {
  if [ ! -d "$MUAH_DIR" ]; then
    return
  fi
  local tmpl_dir="${MUAH_DIR}/docker/compose-templates"
  mkdir -p "$tmpl_dir"
  local tmpl="${tmpl_dir}/simulation.yml.tmpl"
  if [ -f "$tmpl" ]; then
    return
  fi
  cat > "$tmpl" << 'COMPOSEEOF'
# muah-runner simulation compose template
# Variables: {{IMAGE}}, {{COMMAND}}, {{ENV_VARS}}, {{VOLUMES}}
version: "3.9"
services:
  simulation:
    image: "{{IMAGE}}"
    command: ["sh", "-c", "{{COMMAND}}"]
    environment: {{ENV_VARS}}
    volumes: {{VOLUMES}}
    networks:
      - muah-runner-net
    labels:
      muah-runner.simulation: "true"
      muah-runner.run-id: "{{RUN_ID}}"
    mem_limit: "512m"
    cpus: "1.0"
    read_only: false

networks:
  muah-runner-net:
    external: true
COMPOSEEOF
  success "Simulation compose template written to ${tmpl}"
}

# ── Verify ────────────────────────────────────────────────────────────────────
verify() {
  info "Verifying Docker MCP setup..."
  docker run --rm alpine:3.19 echo "muah-runner docker-mcp OK" 2>/dev/null \
    && success "Docker MCP test container ran successfully" \
    || warn "Test container failed — Docker may have resource restrictions"
}

# ── Main ──────────────────────────────────────────────────────────────────────
main() {
  echo ""
  echo "  🐳 Docker MCP setup for muah-runner"
  echo "  ──────────────────────────────────────"
  echo ""
  check_docker
  pull_base_images
  create_network
  write_mcp_config
  write_compose_template
  verify
  echo ""
  success "Docker MCP setup complete!"
  echo ""
  echo "  Usage:"
  echo "    muah-runner docker status              # Check Docker health"
  echo "    muah-runner docker simulate            # Reproduce last failure"
  echo "    muah-runner docker heal                # Auto-fix last failure"
  echo "    muah-runner docker matrix              # Multi-version matrix test"
  echo ""
}

main "$@"
