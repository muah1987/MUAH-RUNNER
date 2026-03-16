# GitHub Copilot Instructions for muah-runner

## What is muah-runner?

muah-runner is **NOT** a traditional CI runner. It is a **persistent, memory-driven autonomous agent** that can be dropped into any system — GitHub Actions, Claude, Gemini, any CI/CD pipeline, or standalone — and maintain its own memory, changelogs, specs, and tooling across runs.

The source of truth is `docs/muah-runner-copilot-prompt.md`. When in doubt, refer to it.

---

## 🧠 Copilot Session Memory Protocol

**This section is mandatory. Follow it at the start and end of every Copilot session.**

### On Session START — Read Memory

Before writing a single line of code, read all files in `.github/memory/`:

```
.github/memory/
├── session-{timestamp}-{uuid}.md    ← previous session snapshots
├── session-{timestamp}-{uuid}.md    ← ...
└── README.md                        ← memory format spec
```

**Steps:**
1. List all `session-*.md` files in `.github/memory/` sorted by timestamp (newest first)
2. Read the most recent file first — it contains the last session's state
3. Read older files if you need broader context (decisions, patterns, lessons)
4. Load the following from each file into your working context:
   - `## What Was Done` — completed work from that session
   - `## Decisions Made` — architectural and design choices (don't re-litigate)
   - `## Open Items` — unfinished work to continue
   - `## Lessons Learned` — what went wrong / what to avoid
   - `## Current State` — where the codebase is right now
5. Begin your session already knowing the project state — no re-exploration needed

> If `.github/memory/` is empty or missing, this is the first session. Bootstrap it by writing a session file at the end.

---

### On Session END — Write Memory

Before finishing your session, **always** write a new memory file:

**Filename format:** `session-{timestamp}-{uuid}.md`
- `{timestamp}` = UTC time in format `YYYYMMDD-HHmmss` (e.g. `20260316-143022`)
- `{uuid}` = a short random hex string, 8 characters (e.g. `a1b2c3d4`)
- Example: `session-20260316-143022-a1b2c3d4.md`

**File content template:**

```markdown
# Session Memory: {timestamp}-{uuid}

**Date:** {full ISO timestamp}
**Agent:** GitHub Copilot
**Task:** {what was asked / worked on this session}

## What Was Done
- {bullet list of completed actions}
- {files created, modified, deleted}
- {commands run and their outcomes}

## Decisions Made
- {architectural decisions and WHY}
- {library/pattern choices and rationale}
- {anything future sessions should not change without reason}

## Current State
- Build: {passing/failing}
- Tests: {X packages passing / failing}
- Missing: {what still needs to be done}
- Branch: {current git branch}
- Last commit: {short commit hash and message}

## Open Items
- [ ] {unfinished task 1}
- [ ] {unfinished task 2}

## Lessons Learned
- {what failed and why}
- {what to avoid next time}
- {patterns that worked well}

## Next Session Should
1. {first thing to do next session}
2. {second thing}
3. {third thing}
```

**Rules for writing memory:**
- Write BEFORE the session ends — not after, not on next start
- Be specific — vague entries are useless ("worked on code" ❌ → "added `cmd/muah-runner/main.go` with 47 CLI subcommands" ✓)
- Always include build/test status so the next session knows if the code compiles
- Record decisions with their **rationale** — not just what, but why
- Keep files — never delete old session files (they are the history)
- Max file size: 200 lines — be concise

---

---

## Architecture Overview

```
muah-runner/
├── cmd/muah-runner/main.go      # CLI entrypoint — all subcommands
├── internal/
│   ├── core/                    # Main run loop + .muah/ bootstrap
│   ├── memory/                  # Persistent memory across runs
│   ├── changelog/               # Append-only changelog with hash chain
│   ├── spec/                    # Spec loader, validator, drift detection
│   ├── mcp/                     # MCP discovery + unified tool registry
│   ├── docker/                  # Self-healing simulation engine
│   ├── planning/                # Repo scanner, planner, roadmap
│   ├── docs/                    # Smart doc generator
│   ├── questions/               # Timed 10s question system
│   ├── selfreflection/          # Post-run self-interview (mandatory)
│   ├── toollearn/               # AI profiling + teaching engine
│   ├── privacy/                 # AES-256-GCM CoT encryption
│   ├── swarm/                   # Vana'diel Protocol (FFXI-themed agents)
│   ├── playwright/              # Browser automation
│   ├── tools/                   # Builtin tools registry + executor
│   ├── adapters/                # Platform adapters (GitHub/Claude/Gemini/generic)
│   ├── sensor/                  # Environment + platform detection
│   └── config/                  # Config loader
├── api/                         # REST API + webhooks + Prometheus metrics
└── templates/muah/              # Default .muah/ directory seed files
```

---

## The `.muah/` Directory

Every project muah-runner touches gets a `.muah/` directory. This is the **persistent brain**:

```
.muah/
├── memory/          # Run history, context, self-reflection
├── changelog/       # Append-only tamper-evident changelog
├── spec/            # Operational specs and quality gates
├── tool_use/        # AI profiling, teaching materials, failure logs
├── privacy_cot/     # Encrypted chain-of-thought vault
├── swarm/           # Agent party state, jobs, linkshell
├── mcp/             # MCP server discovery and connections
├── docker/          # Simulation records and compose templates
├── planning/        # Repo scan, roadmap, current plan
├── docs/            # Generated docs manifest
├── questions/       # Question history and defaults
├── runs/            # Run artifacts
└── config/          # Runtime config and secrets
```

**Rules:**
- `.muah/` is created automatically on `muah-runner init`
- It works on **any** platform — GitHub, GitLab, Bitbucket, Claude, Gemini, standalone
- Never hardcode paths — always use `filepath.Join` for cross-platform compatibility
- Never store secrets in plaintext — use `internal/privacy/` encryption

---

## Run Phases (CRITICAL — always follow this order)

| Phase | Name | Description |
|-------|------|-------------|
| -1 | Tool-Use Bootstrap | Load tool catalog, AI profiles, lessons into context |
| 0 | Repo Scan | Detect tech stack, platform, gaps |
| 1 | Plan | Generate roadmap, ask user to approve |
| 2 | Docs | Generate all standard docs BEFORE any code |
| 3 | Execute | Run the actual task |
| 4 | Review & Close | Doc audit, PR, run summary |
| 5 | Self-Reflection | **MANDATORY** — 5-question interview, grade, growth log |

**Phase 5 is NEVER skipped.** It is how the runner learns.

---

## Vana'diel Swarm Protocol

Multi-agent system themed after Final Fantasy XI:

| Race | AI Model | Strengths |
|------|----------|-----------|
| **Elvaan** | Opus latest (claude-opus-4-6) | Deep reasoning, architecture, docs |
| **Galka** | Sonnet latest (claude-sonnet-4-6) | Sustained code work, builds |
| **Mithra** | Haiku latest (claude-haiku-4-5) | Fast search, parallel tasks |
| **Tarutaru** | Opus prev / o1 / o3 | Deep analysis, complex reasoning |
| **Hume** | The user | Issues quests, reviews output |

**HP** = context tokens. **MP** = output tokens. **Death** = context exhausted.

Jobs (22 total): WAR, MNK, WHM, BLM, RDM, THF, PLD, DRK, BST, BRD, RNG, SAM, NIN, DRG, SMN, BLU, COR, PUP, DNC, SCH, GEO, RUN.

The Party Leader (**Prishe** by default) is the ONLY agent that talks to the user.

---

## Platform Agnosticism

muah-runner runs on **any** platform. Platform is auto-detected:

```go
// Detection priority (sensor/environment.go):
// 1. GITHUB_ACTIONS=true     → github
// 2. GITLAB_CI=true          → gitlab
// 3. CIRCLECI=true           → circleci
// 4. BITBUCKET_BUILD_NUMBER  → bitbucket
// 5. JENKINS_URL             → jenkins
// 6. CLAUDE_CONVERSATION_ID  → claude
// 7. GEMINI_API_KEY          → gemini
// 8. (default)               → standalone
```

**Rules for platform-agnostic code:**
- Always use `filepath.Join()` — never string concatenation for paths
- Use `os.PathSeparator` when displaying paths to users
- Read from env vars for tokens/secrets, never hardcode
- The `.muah/` dir system is the portability layer — same structure everywhere
- MCP adapters abstract platform differences

---

## Key Coding Conventions

### Go Style
- Module: `github.com/muah1987/muah-runner`
- Go version: 1.24+
- Use `filepath.Join` for ALL file paths (cross-platform)
- Error handling: always wrap errors with context (`fmt.Errorf("doing X: %w", err)`)
- No global state — pass dependencies explicitly
- All packages have a `New*()` constructor that takes `muahDir string`

### File Operations
```go
// Always use filepath.Join — never string concat for paths
path := filepath.Join(s.muahDir, "memory", "runs", runID+".json")

// Always create parent dirs before writing
os.MkdirAll(filepath.Dir(path), 0755)
```

### Privacy / Secrets
```go
// Use internal/privacy for anything sensitive
enc, _ := privacy.NewEncryptor(key)
ciphertext, _ := enc.Encrypt([]byte(sensitiveData))
```

### Timed Questions
```go
// Never block indefinitely — always use timed questions
asker := questions.NewAsker(muahDir)
answer := asker.Ask(&questions.Question{
    Text:        "Create a PR with these changes?",
    Options:     []questions.Option{{Label: "Yes"}, {Label: "No"}},
    Recommended: 0,
    Timeout:     10 * time.Second,
}, cfg.Questions.AutoMode)
```

### Tool Failures
```go
// 3 strikes rule — never retry more than 3 times with same params
// Use internal/toollearn for tracking and recovery
monitor := toollearn.NewMonitor(muahDir)
monitor.RecordCall("github-mcp:create_pr", params, err)
if monitor.ShouldEscalate("github-mcp:create_pr") {
    // fall back: MCP → CLI → API → Docker → ask human
}
```

### Self-Reflection (every run, no exceptions)
```go
interviewer := selfreflection.NewInterviewer(muahDir)
entry, _ := interviewer.Conduct(runID, selfreflection.RunMetrics{
    ErrorCount:      errorCount,
    RetryCount:      retryCount,
    DurationSeconds: elapsed.Seconds(),
    ActionCount:     actionCount,
    KeyLesson:       keyLesson,
})
```

---

## Adding New Features

1. **New tool?** → Add to `internal/tools/registry.go` + implement in `internal/tools/builtin/`
2. **New platform adapter?** → Implement `adapters.Adapter` interface in `internal/adapters/`
3. **New MCP server?** → Add client to `internal/mcp/platform/` or `internal/mcp/native/`
4. **New job class?** → Add to `internal/swarm/jobs.go` with full ability set
5. **New CLI command?** → Add to `cmd/muah-runner/main.go` switch statement

---

## Testing

```bash
go test ./...          # Run all tests
go build ./...         # Verify everything compiles
go vet ./...           # Check for issues
```

Each package has `*_test.go`. Tests use `os.MkdirTemp` for `.muah/` directories — never use real filesystem paths in tests.

---

## CLI Reference (quick)

```bash
muah-runner init           # Bootstrap .muah/, scan repo, plan, generate docs
muah-runner run <task>     # Execute a task (full Phase -1 through 5)
muah-runner status         # Health check: MCPs, Docker, Playwright, swarm
muah-runner scan           # Repo scan only
muah-runner plan           # Generate/show plan
muah-runner docs generate  # Generate standard docs package
muah-runner swarm status   # Party roster with HP/MP
muah-runner reflect show   # Last self-reflection entry
muah-runner doctor         # Full system health check
```

Full CLI reference: see `docs/muah-runner-copilot-prompt.md` → CLI Commands section.

---

## Security Notes

- CoT (Chain of Thought) is **always encrypted** at rest using AES-256-GCM
- Secrets are stored in `.muah/config/secrets.enc` — never in plaintext
- Never commit `.muah/` to version control (it's in `.gitignore`)
- Never log API keys, tokens, or credentials
- CoT classification: `private` (never shown) | `redactable` (sanitized on request) | `public` (always shown)

---

## Key Principles (from spec)

1. **Memory is sacred** — never lose context, always persist across runs
2. **Changelog everything** — if muah-runner touched it, it's logged
3. **MCP-first** — use platform MCP before CLI before API
4. **Self-improvement is mandatory** — Phase 5 self-reflection, every run, no exceptions
5. **Platform agnostic** — same `.muah/` works on GitHub, GitLab, Claude, Gemini, terminal
6. **Never block on humans** — timed questions, 10 seconds, auto-select smart default
7. **Plan before code** — Phases 0→2 happen before Phase 3
8. **3 strikes, recover** — classify error, fall back, teach the AI, never retry blindly
9. **Docker is your lab** — simulate fixes in containers before applying to host
10. **Thoughts are private** — CoT encrypted at rest, redacted on demand
