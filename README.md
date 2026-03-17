# muah-runner

> **NOT** a traditional CI runner.
> A **persistent, memory-driven autonomous agent** you drop into any system — GitHub Actions, GitLab CI, Claude, Gemini, Jenkins, or a bare terminal — and it maintains its own memory, changelogs, specs, and tooling across every run.

---

## Philosophy

| Old runner | muah-runner |
|---|---|
| Runs a job, forgets it | Remembers every run forever |
| Configured once | Learns and improves every run |
| Lives in CI | Lives anywhere |
| You manage it | It manages itself |
| Retries blindly | Classifies errors, teaches itself |

muah-runner has a **persistent brain** (the `.muah/` directory), a mandatory **self-reflection phase** after every task, and a **Vana'diel swarm** of AI agents that collaborate under a single party leader.

---

## Quick Start (< 5 minutes)

```bash
# 1. Install
curl -fsSL https://raw.githubusercontent.com/muah1987/muah-runner/main/scripts/setup.sh | bash

# 2. Bootstrap your repo
cd my-project
muah-runner init

# 3. Run a task
muah-runner run "add unit tests to the auth package"

# 4. Check what it did
muah-runner memory show
muah-runner reflect show
muah-runner changelog show
```

### From source

```bash
git clone https://github.com/muah1987/muah-runner.git
cd muah-runner
make build
./muah-runner init
```

### Docker

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/muah1987/muah-runner:latest init
```

---

## How It Works

Every `muah-runner run <task>` executes **six phases** in order:

| Phase | Name | What happens |
|-------|------|-------------|
| **-1** | Tool-Use Bootstrap | Loads tool catalog, AI profiles, past lessons into context |
| **0** | Repo Scan | Detects tech stack, platform, gaps, missing docs |
| **1** | Plan | Generates roadmap, asks you to approve (10-second timer) |
| **2** | Docs | Generates standard docs package **before** any code changes |
| **3** | Execute | Runs the actual task using MCP tools + swarm agents |
| **4** | Review & Close | Doc audit, PR creation, run summary |
| **5** | Self-Reflection | **MANDATORY** 5-question interview → grade → growth log |

Phase 5 is **never skipped**. It is how muah-runner gets smarter over time.

---

## The `.muah/` Directory

Every project gets a persistent brain:

```
.muah/
├── memory/          # Every run, action, lesson — forever
├── changelog/       # Tamper-evident append-only changelog
├── spec/            # Operational specs and quality gates
├── tool_use/        # AI proficiency profiles, failure logs, teaching materials
├── privacy_cot/     # AES-256-GCM encrypted chain-of-thought vault
├── swarm/           # Agent party state, jobs, Linkshell messages
├── mcp/             # MCP server discovery and active connections
├── docker/          # Simulation records and compose templates
├── planning/        # Repo scan results, roadmap, current plan
├── docs/            # Generated docs manifest
├── questions/       # Question history and auto-select defaults
├── runs/            # Run artifacts
└── config/          # Runtime config and secrets
```

> Add `.muah/` to `.gitignore` — it's your local brain, not source code.

---

## Platform Support

muah-runner auto-detects the platform from environment variables. The same `.muah/` works everywhere:

| Platform | Detection | Token env var |
|---|---|---|
| GitHub Actions | `GITHUB_ACTIONS=true` | `GITHUB_TOKEN` |
| GitLab CI | `GITLAB_CI=true` | `GITLAB_TOKEN` |
| CircleCI | `CIRCLECI=true` | `CIRCLE_TOKEN` |
| Bitbucket | `BITBUCKET_BUILD_NUMBER` | `BITBUCKET_TOKEN` |
| Jenkins | `JENKINS_URL` | `JENKINS_TOKEN` |
| Azure DevOps | `TF_BUILD=True` | `AZURE_TOKEN` |
| Travis CI | `TRAVIS=true` | `TRAVIS_TOKEN` |
| TeamCity | `TEAMCITY_VERSION` | — |
| Drone CI | `DRONE=true` | `DRONE_TOKEN` |
| Claude | `CLAUDE_CONVERSATION_ID` | `ANTHROPIC_API_KEY` |
| Gemini | `GEMINI_API_KEY` | `GEMINI_API_KEY` |
| Standalone | (default) | — |

---

## CLI Reference

```
muah-runner <command> [subcommand] [arguments]
```

### Core

| Command | Description |
|---|---|
| `init` | Bootstrap `.muah/`, scan repo, generate plan, write docs |
| `run <task>` | Execute a task through full Phase -1→5 pipeline |
| `status` | Health check: platform, MCPs, Docker, Playwright, memory |
| `start` | Start daemon mode (connects all MCPs) |
| `stop` | Graceful shutdown |
| `doctor` | Full diagnostics: MCP, Docker, Playwright, swarm, memory |
| `gc` | Garbage collect old runs, artifacts, containers |
| `self-update` | Pull and apply updates |

### Memory

```bash
muah-runner memory show           # Show recent run memory
muah-runner memory search <query> # Full-text search across all runs
muah-runner memory export         # Export as portable JSON
```

### Changelog

```bash
muah-runner changelog show   # Show running changelog
muah-runner changelog diff   # Diff between two runs
```

### Spec

```bash
muah-runner spec check    # Validate current state against specs
muah-runner spec drift    # Show spec drift score
muah-runner spec propose  # Propose spec updates from learned patterns
```

### Tools

```bash
muah-runner tools list          # All available tools (builtin + MCP + generated)
muah-runner tools proficiency   # Per-AI tool-use proficiency scores
muah-runner tools failures      # Recent failures and recovery actions
muah-runner tools dashboard     # Full metrics dashboard
muah-runner tools teach         # Generate teaching materials from failures
muah-runner tools catalog       # Full catalog with AI-specific notes
muah-runner tools inject        # Reload tool knowledge into context
muah-runner tools playbook <e>  # Recovery playbook for error type
```

### Planning & Docs

```bash
muah-runner scan              # Scan repo: stack, gaps, platform
muah-runner plan              # Generate/update plan
muah-runner plan show         # Show current plan
muah-runner docs generate     # Generate full standard docs package
muah-runner docs audit        # Check doc coverage, flag missing/stale
muah-runner docs update       # Re-generate from current repo state
```

### Self-Reflection

```bash
muah-runner reflect show      # Last self-reflection entry
muah-runner reflect growth    # Growth log (improvement over time)
muah-runner reflect patterns  # Detected patterns across runs
muah-runner reflect grade     # Self-rating history and trend
muah-runner reflect actions   # Open action items from reflections
```

### MCP

```bash
muah-runner mcp discover        # Scan environment for MCP servers
muah-runner mcp list            # Show connected MCP servers + tools
muah-runner mcp connect <url>   # Manually connect to an MCP server
muah-runner mcp test            # Health-check all MCP connections
muah-runner mcp install <name>  # Install native MCP (playwright, docker…)
```

### Docker

```bash
muah-runner docker status     # Show Docker + running containers
muah-runner docker simulate   # Reproduce last failure in a container
muah-runner docker heal       # Auto-diagnose and fix last failure
muah-runner docker matrix     # Multi-version matrix test
muah-runner docker cleanup    # Remove simulation containers
```

### Playwright

```bash
muah-runner playwright setup              # Install Playwright + browsers
muah-runner playwright test               # Run Playwright test suite
muah-runner playwright screenshot <url>   # Quick screenshot
muah-runner playwright audit <url>        # Accessibility + performance audit
```

### Swarm (Vana'diel Protocol)

```bash
muah-runner swarm status           # Party roster with HP/MP
muah-runner swarm races            # All races + AI model mappings
muah-runner swarm race <name>      # Profile for Elvaan/Galka/Mithra/Tarutaru/Hume
muah-runner swarm jobs             # All 22 job classes
muah-runner swarm compose <task>   # Preview party composition (dry run)
muah-runner swarm spawn <job>      # Manually spawn an agent
muah-runner swarm dismiss <name>   # Dismiss an agent
muah-runner swarm linkshell        # Recent Linkshell messages (translated)
muah-runner swarm tongue           # Vanadiel Tongue lexicon
muah-runner swarm crystal          # Crystal Pool token budget
```

### Platform Registration

```bash
muah-runner register github   # Register as GitHub Actions self-hosted runner
muah-runner register claude   # Set up Claude adapter
muah-runner register gemini   # Set up Gemini adapter
```

---

## Vana'diel Swarm Protocol

muah-runner uses a multi-agent system inspired by Final Fantasy XI. Each AI model maps to a **race** with distinct stat profiles:

| Race | Model | HP | MP | Speed | Best jobs |
|------|-------|----|----|-------|-----------|
| **Elvaan** | Claude Opus (latest) | ×1.5 | ×1.5 | ×0.6 | BLM, SCH, SMN, PLD |
| **Galka** | Claude Sonnet (latest) | ×1.8 | ×1.3 | ×1.2 | WAR, MNK, SAM, DRK |
| **Mithra** | Claude Haiku (latest) | ×0.5 | ×0.5 | ×3.0 | THF, NIN, RNG, COR, DNC |
| **Tarutaru** | Claude Opus (prev) / o1 | ×0.7 | ×2.0 | ×0.5 | BLM, WHM, SMN, SCH, GEO |
| **Hume** | The User | — | — | — | Issues quests, reviews output |

> **HP** = context window tokens. **MP** = output tokens. **Death** = context exhausted.

**22 Job Classes:** WAR, MNK, WHM, BLM, RDM, THF, PLD, DRK, BST, BRD, RNG, SAM, NIN, DRG, SMN, BLU, COR, PUP, DNC, SCH, GEO, RUN

The Party Leader (**Prishe** by default) is the **only** agent that talks directly to you. All inter-agent coordination happens via **Linkshell** messages in **Vanadiel Tongue** — an internal compressed language (`vel mir tar zan kel mal tun keldah`).

```bash
muah-runner swarm tongue          # See the lexicon
muah-runner swarm compose "task"  # Preview optimal party for a task
```

---

## Security

- **CoT is always encrypted** — chain-of-thought stored at rest with AES-256-GCM
- **Secrets never in plaintext** — `.muah/config/secrets.enc` only
- **Never commit `.muah/`** — add to `.gitignore`
- **CoT classification:** `private` (never shown) | `redactable` (sanitized) | `public` (always shown)

```bash
muah-runner cot status        # Encryption status
muah-runner cot redacted <id> # View sanitized CoT for a run
```

---

## Configuration

`muah-runner.yml` in your project root (or `.muah/config/muah-runner.yml`):

```yaml
runner:
  name: "muah-runner-01"
  muah_dir: ".muah"
  concurrency: 4
  heartbeat_interval: 30

server:
  host: "0.0.0.0"
  port: 8080

swarm:
  enabled: true
  max_party_size: 6
  default_leader: "Prishe"

questions:
  timer_seconds: 10   # Auto-select after N seconds
  auto_mode: false    # true = never wait for input

selfreflection:
  enabled: true
  pattern_detection_interval: 5  # Detect patterns every N runs

mcp:
  auto_discover: true
  require_playwright: true
  require_docker: true
```

---

## Contributing

1. Fork → `git checkout -b feature/my-feature`
2. `make build` → `make test`
3. `muah-runner init` to bootstrap `.muah/` in the repo
4. `muah-runner run "implement my feature"`
5. Submit PR — muah-runner will review its own work in Phase 4

### Project Structure

```
muah-runner/
├── cmd/muah-runner/      # CLI entrypoint
├── internal/
│   ├── core/             # Run loop + .muah/ bootstrap
│   ├── memory/           # Persistent run store
│   ├── changelog/        # Append-only hash-chained changelog
│   ├── spec/             # Spec loader, validator, drift detection
│   ├── mcp/              # MCP discovery + tool registry
│   ├── docker/           # Self-healing simulation engine
│   ├── planning/         # Repo scanner + phase planner
│   ├── docs/             # Smart doc generator
│   ├── questions/        # 10-second timed question system
│   ├── selfreflection/   # Post-run self-interview + grader
│   ├── toollearn/        # AI profiling + teaching engine
│   ├── privacy/          # AES-256-GCM CoT encryption
│   ├── swarm/            # Vana'diel Protocol + 22 job classes
│   ├── playwright/       # Browser automation
│   ├── tools/            # Builtin tools registry + executor
│   ├── adapters/         # Platform adapters
│   ├── sensor/           # Environment + platform detection
│   └── config/           # Config loader
├── api/                  # REST API + Prometheus metrics
├── templates/muah/       # Default .muah/ seed files
└── scripts/              # setup.sh, install-playwright.sh, install-docker-mcp.sh
```

---

## License

See [LICENSE](LICENSE). — My Own Runner
