# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

# Build-time version tag (set via --build-arg VERSION=v1.2.3)
ARG VERSION=dev

WORKDIR /build

# Cache module downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
      -ldflags "-s -w -X main.Version=${VERSION}" \
      -o muah-runner \
      ./cmd/muah-runner/

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
# Pin to a specific alpine digest family for reproducibility
FROM alpine:3.21

# ca-certificates: required for HTTPS calls (MCP servers, GitHub API, etc.)
# tzdata: required for correct timezone handling in logs and scheduled tasks
# git:    required for repo-scan phase (muah-runner init / scan)
RUN apk --no-cache add ca-certificates tzdata git

WORKDIR /app

COPY --from=builder /build/muah-runner .

# API port (REST + webhooks)
EXPOSE 8080
# Prometheus metrics port
EXPOSE 9090

# .muah/ is the persistent brain — mount a named volume here so data survives
# container restarts:
#   docker run -v muah-brain:/app/.muah  ghcr.io/muah1987/muah-runner:latest
VOLUME ["/app/.muah"]

ENTRYPOINT ["./muah-runner"]

# Default: start daemon mode.
# Override with 'init' for first-run bootstrap:
#   docker run ... muah-runner init
CMD ["start"]
