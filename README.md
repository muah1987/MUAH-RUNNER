# muah-runner

A lightweight, self-hosted CI/CD runner written in Go. Drop it on any Linux server, point it at your workload, and start running jobs.

---

## Features

- **Process & Docker executors** — run jobs as local OS processes or inside Docker containers
- **Priority job queue** — thread-safe, supports concurrent workers
- **Auto-capability detection** — discovers OS, CPU, memory, disk, GPU, and installed toolchains at startup
- **REST API** — submit, list, cancel, and inspect jobs over HTTP
- **Prometheus metrics** — CPU, memory, disk, and job counters
- **GitHub webhook handler** — receives `push`, `pull_request`, and other events with HMAC-SHA256 signature verification
- **YAML + env-var config** — full config file with environment variable overrides
- **Graceful shutdown** — drains in-flight jobs on SIGTERM

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        muah-runner                          │
│                                                             │
│  ┌──────────┐   ┌──────────┐   ┌──────────────────────┐   │
│  │  sensor  │   │  config  │   │     health/monitor   │   │
│  │ (detect) │   │  (yaml)  │   │  (cpu/mem/disk/jobs) │   │
│  └────┬─────┘   └────┬─────┘   └──────────┬───────────┘   │
│       │              │                     │               │
│  ┌────▼──────────────▼─────────────────────▼───────────┐  │
│  │                   API Server (:8080)                 │  │
│  │  GET  /health          GET  /status                  │  │
│  │  GET  /api/v1/info     GET  /metrics                 │  │
│  │  GET  /jobs            POST /jobs                    │  │
│  │  GET  /jobs/:id        DELETE /jobs/:id              │  │
│  │  POST /webhooks/github                               │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│                  ┌──────▼──────┐                           │
│                  │  Job Queue  │  (priority, thread-safe)  │
│                  └──────┬──────┘                           │
│                         │  Dequeue                         │
│            ┌────────────▼────────────┐                     │
│            │     Job Processor       │  (N workers)        │
│            │  ┌─────────────────┐    │                     │
│            │  │ProcessExecutor  │    │                     │
│            │  │ or              │    │                     │
│            │  │DockerExecutor   │    │                     │
│            │  └─────────────────┘    │                     │
│            └─────────────────────────┘                     │
└─────────────────────────────────────────────────────────────┘
```

---

## Quick Start

```bash
# Build
make build

# Run with defaults (no config file needed)
./muah-runner start

# Submit a job
curl -s -X POST http://localhost:8080/jobs \
  -H 'Content-Type: application/json' \
  -d '{"name":"hello","command":["echo","hello world"]}'

# Check health
curl http://localhost:8080/health

# View Prometheus metrics
curl http://localhost:8080/metrics
```

---

## Installation

### From source

```bash
git clone https://github.com/muah1987/muah-runner.git
cd muah-runner
make build
sudo mv muah-runner /usr/local/bin/
```

### Install script

```bash
curl -fsSL https://raw.githubusercontent.com/muah1987/muah-runner/main/scripts/install.sh | bash
```

### Docker

```bash
docker run -p 8080:8080 -p 9090:9090 \
  -v $(pwd)/muah-runner.yml:/app/muah-runner.yml:ro \
  ghcr.io/muah1987/muah-runner:latest start
```

### docker-compose

```bash
docker-compose up -d
```

---

## Configuration

Configuration is loaded from `muah-runner.yml` (default). Missing keys fall back to built-in defaults. All keys can be overridden by environment variables.

```yaml
runner:
  name: "muah-runner-01"       # Runner display name
  labels: []                   # Extra labels to advertise
  work_dir: "/tmp/muah-runner" # Working directory for job files
  concurrency: 4               # Maximum parallel jobs
  heartbeat_interval: 30       # Seconds between heartbeats

server:
  host: "0.0.0.0"   # Bind address
  port: 8080         # API listen port
  tls: false         # Enable TLS (not yet implemented)

executor:
  type: "process"          # "process" or "docker"
  docker_image: "ubuntu:22.04"  # Image for docker executor
  timeout: 3600            # Job timeout in seconds
  max_retries: 3           # Retry count on failure
  retry_backoff: 5         # Seconds between retries

health:
  metrics_port: 9090   # Prometheus metrics redirect port
  check_interval: 15   # Seconds between metric collections

webhooks:
  enabled: false   # Enable webhook processing
  urls: []         # URLs to forward events to
  secret: ""       # HMAC secret for GitHub signature verification

secrets:
  encryption_key: ""   # Reserved for future secret encryption
```

### Environment variable overrides

| Variable | Config key |
|---|---|
| `MUAH_RUNNER_NAME` | `runner.name` |
| `MUAH_RUNNER_WORK_DIR` | `runner.work_dir` |
| `MUAH_RUNNER_CONCURRENCY` | `runner.concurrency` |
| `MUAH_SERVER_HOST` | `server.host` |
| `MUAH_SERVER_PORT` | `server.port` |
| `MUAH_EXECUTOR_TYPE` | `executor.type` |
| `MUAH_HEALTH_METRICS_PORT` | `health.metrics_port` |
| `MUAH_WEBHOOK_SECRET` | `webhooks.secret` |
| `MUAH_SECRETS_ENCRYPTION_KEY` | `secrets.encryption_key` |

---

## CLI

```
muah-runner <command> [options]

Commands:
  start       Start the runner
  stop        Stop the runner (prints instructions)
  status      Print current configuration and capabilities
  register    Instructions to register with GitHub Actions
  unregister  Instructions to remove from GitHub Actions
  self-update Check for updates
  version     Print version
```

### `start`

```
muah-runner start [-config <path>]

  -config string   Path to config file (default: muah-runner.yml)
```

### `status`

```
muah-runner status [-config <path>]
```

Prints runner name, work dir, OS/arch, CPU cores, labels, executor type, and API/metrics URLs.

---

## REST API

### `GET /health`

Returns current health stats.

**Response**
```json
{
  "status": "ok",
  "cpu": 12.5,
  "memory": 45.1,
  "disk": 30.0,
  "time": "2024-01-15T10:00:00Z"
}
```

---

### `GET /status`

Returns runner status and job counters.

**Response**
```json
{
  "name": "muah-runner-01",
  "labels": ["os:linux", "arch:amd64", "go:v1"],
  "jobs_total": 42,
  "jobs_running": 2,
  "jobs_failed": 1,
  "queue_len": 5
}
```

---

### `GET /api/v1/info`

Returns runner identity and capabilities.

**Response**
```json
{
  "name": "muah-runner-01",
  "version": "dev",
  "labels": ["os:linux", "arch:amd64"],
  "os": "linux",
  "arch": "amd64"
}
```

---

### `GET /jobs`

List all jobs (all statuses).

**Response** — array of job objects.

---

### `POST /jobs`

Submit a new job.

**Request**
```json
{
  "name": "my-build",
  "priority": 5,
  "command": ["make", "build"],
  "env": {"CI": "true"},
  "artifacts": ["dist/**"]
}
```

**Response** `201 Created`
```json
{
  "ID": "job-1705312800000000000",
  "Name": "my-build",
  "Priority": 5,
  "Status": "pending",
  "Command": ["make", "build"],
  ...
}
```

---

### `GET /jobs/:id`

Get a specific job by ID.

---

### `DELETE /jobs/:id`

Cancel a pending job.

**Response**
```json
{"status": "cancelled"}
```

---

### `GET /metrics`

Prometheus metrics endpoint (text format).

---

### `POST /webhooks/github`

Receives GitHub webhook events. Set `webhooks.secret` to enable HMAC-SHA256 signature verification via the `X-Hub-Signature-256` header.

**Response**
```json
{"status": "received", "event": "push"}
```

---

## Prometheus Metrics

| Metric | Type | Description |
|---|---|---|
| `muah_runner_cpu_usage` | Gauge | CPU usage percentage |
| `muah_runner_memory_usage` | Gauge | Memory usage percentage |
| `muah_runner_disk_usage` | Gauge | Disk usage percentage |
| `muah_runner_jobs_total` | Gauge | Total jobs processed |
| `muah_runner_jobs_running` | Gauge | Currently running jobs |
| `muah_runner_jobs_failed` | Gauge | Total failed jobs |

---

## Docker

### Build image

```bash
docker build -t muah-runner:latest .
```

### Run

```bash
docker run -d \
  --name muah-runner \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/muah-runner.yml:/app/muah-runner.yml:ro \
  -v muah-work:/tmp/muah-runner \
  muah-runner:latest
```

### docker-compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make changes and add tests
4. Run `make test` and ensure all tests pass
5. Run `make lint` (requires `golangci-lint`)
6. Submit a pull request

### Development

```bash
# Run tests with race detector
make test

# Build binary
make build

# Install locally
make install
```

### Project Structure

```
.
├── api/                  # HTTP server and handlers
├── cmd/muah-runner/      # Main entrypoint
├── internal/
│   ├── config/           # Configuration loading
│   ├── executor/         # Process and Docker executors
│   ├── health/           # Health monitoring and Prometheus metrics
│   ├── queue/            # Priority job queue
│   └── sensor/           # Capability detection
├── scripts/              # Install script
├── muah-runner.yml       # Default config
├── Dockerfile
└── docker-compose.yml
```

---

## License

See [LICENSE](LICENSE).

My Own Runner
