# Contributing to muah-runner

Thank you for contributing! muah-runner is a persistent, memory-driven autonomous agent runner. Every contribution goes through the same Phase 0→5 pipeline the agent uses itself.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Conventions](#coding-conventions)
- [Testing](#testing)
- [Session Memory Protocol](#session-memory-protocol)
- [Architecture Quick Reference](#architecture-quick-reference)
- [Submitting a PR](#submitting-a-pr)

---

## Code of Conduct

Be respectful and constructive. Assume good intent. Focus on ideas, not people.

---

## Getting Started

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/muah-runner.git
cd muah-runner

# 2. Install Go 1.24+ (or use the setup script)
bash scripts/setup.sh

# 3. Bootstrap .muah/ (creates persistent brain for your dev session)
muah-runner init

# 4. Build
make build

# 5. Run tests
make test
```

### Optional dependencies

```bash
# Playwright (for browser automation tests)
bash scripts/install-playwright.sh

# Docker MCP (for simulation engine)
bash scripts/install-docker-mcp.sh
```

---

## Development Workflow

muah-runner follows a strict **phase-before-code** discipline:

| Phase | What you should do |
|-------|-------------------|
| 0 — Scan | Understand the existing code and tests in the area you're changing |
| 1 — Plan | Write a clear description of the change in your PR |
| 2 — Docs | Update any documentation affected by your change |
| 3 — Execute | Write the code, then write the tests |
| 4 — Review | Self-review the diff; check for missing tests or edge cases |
| 5 — Reflect | After your PR is merged, note what you learned |

---

## Coding Conventions

### Go style

- **Module:** `github.com/muah1987/muah-runner`
- **Go version:** 1.24+
- **File paths:** always use `filepath.Join` — never string concatenation
- **Errors:** always wrap with context: `fmt.Errorf("doing X: %w", err)`
- **No global state:** pass dependencies explicitly via constructors
- **Constructors:** every package has a `New*()` or `New*(muahDir string)` constructor

### File operations

```go
// Always use filepath.Join — never string concat
path := filepath.Join(s.muahDir, "memory", "runs", runID+".json")

// Always create parent dirs before writing
os.MkdirAll(filepath.Dir(path), 0755)
```

### Timed questions — never block indefinitely

```go
asker := questions.NewAsker(muahDir)
answer := asker.Ask(&questions.Question{
    Text:        "Create a PR with these changes?",
    Options:     []questions.Option{{Label: "Yes"}, {Label: "No"}},
    Recommended: 0,
    Timeout:     10 * time.Second,
}, cfg.Questions.AutoMode)
```

### Tool failures — 3 strikes rule

```go
monitor := toollearn.NewMonitor(muahDir)
monitor.RecordCall("github-mcp:create_pr", params, err)
if monitor.ShouldEscalate("github-mcp:create_pr") {
    // Fallback chain: MCP → CLI → API → Docker → ask human
}
```

### Self-reflection — every run, no exceptions

```go
interviewer := selfreflection.NewInterviewer(muahDir)
entry, _ := interviewer.Conduct(runID, selfreflection.RunMetrics{...})
```

### Privacy — CoT always encrypted

```go
enc, _ := privacy.NewEncryptor(key)
ciphertext, _ := enc.Encrypt([]byte(sensitiveData))
```

---

## Testing

```bash
# Run all tests with race detector
make test

# Run tests for a specific package
go test -v -race ./internal/core/...

# Build the binary
make build

# Lint (requires golangci-lint)
make lint
```

### Test conventions

- Use `t.TempDir()` for `.muah/` directories — never use real filesystem paths
- Table-driven tests for parameterised cases
- Test file names: `*_test.go` in the same package
- New packages **must** have at least one test file

```go
func TestMyFeature(t *testing.T) {
    dir := t.TempDir()
    // ... use dir as muahDir
}
```

---

## Session Memory Protocol

If you are a Copilot/AI agent contributing to this repo, you **must** follow the session memory protocol defined in `.github/copilot-instructions.md`:

1. **On session START** — read all files in `.github/memory/session-*.md` (newest first)
2. **On session END** — write a new `session-{YYYYMMDD-HHmmss}-{8hexchars}.md` file

This keeps every agent session informed without re-exploring the entire codebase.

---

## Architecture Quick Reference

```
muah-runner/
├── cmd/muah-runner/      # CLI entrypoint (main.go)
├── internal/
│   ├── core/             # Run loop + .muah/ Bootstrap()
│   ├── memory/           # Persistent run store
│   ├── changelog/        # Append-only hash-chained changelog
│   ├── spec/             # Spec loader, validator, drift detection
│   ├── mcp/              # MCP discovery + unified tool registry
│   ├── docker/           # Self-healing simulation engine
│   ├── planning/         # Repo scanner + roadmap planner
│   ├── docs/             # Smart doc generator
│   ├── questions/        # 10-second timed question system
│   ├── selfreflection/   # Post-run self-interview + grader
│   ├── toollearn/        # AI profiling + teaching engine
│   ├── privacy/          # AES-256-GCM CoT encryption
│   ├── swarm/            # Vana'diel Protocol (22 jobs, 5 races)
│   ├── playwright/       # Browser automation
│   ├── tools/            # Builtin tools registry
│   ├── adapters/         # Platform adapters (GitHub/Claude/Gemini/generic)
│   ├── sensor/           # Environment + platform detection (14 platforms)
│   └── config/           # Config loader
├── api/                  # REST API + Prometheus metrics
└── templates/muah/       # Default .muah/ seed files
```

**Adding a new package:** follow the `New*(muahDir string)` constructor convention, add a `*_test.go` with at least one test using `t.TempDir()`.

---

## Submitting a PR

1. Fork → `git checkout -b feature/my-feature`
2. Make your changes following the conventions above
3. `make build && make test` — ensure all pass
4. Update docs if you changed public APIs or behaviour
5. Open a PR with a clear description
6. Address review comments

### PR checklist (also in `.github/PULL_REQUEST_TEMPLATE.md`)

- [ ] `make build` passes
- [ ] `make test` passes — no new failures
- [ ] New code has tests (`*_test.go`)
- [ ] Docs updated (README, package doc comments)
- [ ] No hardcoded file paths (use `filepath.Join`)
- [ ] No secrets in source code
- [ ] Phase 5 self-reflection: what did you learn from this change?
