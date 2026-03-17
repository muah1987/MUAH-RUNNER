# Copilot Prompt for muah-runner

Build `muah-runner`: a persistent, memory-driven autonomous agent runner that can be dropped into **any** system — GitHub Actions, Claude, Gemini, any CI/CD pipeline, or standalone — and maintain its own memory, changelogs, specs, and tooling across runs.

---

## Core Concept

muah-runner is NOT a traditional runner. It is an **autonomous agent with persistent memory** that:
- Creates and manages a `.muah/` directory in whatever system/repo it runs in
- Remembers everything from previous runs
- Tracks all changes systematically
- Enforces specs to keep runs tight and consistent
- Can **create its own tools and skills** to complete tasks it's given
- Has built-in search, web fetch, file manipulation, and code execution capabilities
- **MCP-first architecture**: discovers, connects to, and optimally uses MCP (Model Context Protocol) servers — platform MCPs (GitHub, Claude, Gemini), utility MCPs (Playwright, Docker, filesystem), and any custom MCPs
- **Playwright is native**: browser automation, testing, scraping, and visual verification are always available — not optional
- **Docker-native**: can spin up its own containers on demand to simulate environments, reproduce bugs, test fixes in isolation, and self-heal broken runs before reporting failure
- **Never blocks on humans**: when muah-runner needs a decision, it asks a timed question (10 seconds). If no answer, the agent's recommended choice is auto-selected. The runner never stalls waiting for input.
- **Planning-first, docs-first**: before touching ANY code, muah-runner scans the repo, builds a plan/roadmap, creates all required documentation, and structures everything properly. Docs come before code, always.

---

## `.muah/` Directory Structure (auto-created on first run)

```
.muah/
├── memory/
│   ├── index.json              # Master index of all memories
│   ├── runs/                   # Per-run memory snapshots
│   │   ├── run-{timestamp}.json
│   │   └── latest.json         # Symlink to most recent run
│   ├── context/                # Accumulated knowledge about the project
│   │   ├── project.json        # What this project is, tech stack, patterns
│   │   ├── decisions.json      # Why certain decisions were made
│   │   └── lessons.json        # What went wrong before, what to avoid
│   ├── selfreflection/         # Agent self-critique after every run
│   │   ├── index.json          # Aggregated insights across all reflections
│   │   ├── reflections/        # Per-run self-reflection entries
│   │   │   └── reflect-{timestamp}.json
│   │   ├── action-items.json   # Open action items from self-reflection
│   │   ├── patterns.json       # Recurring mistakes and improvements detected
│   │   └── growth-log.md       # Human-readable log of how the runner is improving
│   └── embeddings/             # Searchable memory vectors (optional)
│       └── memory.db
│
├── tool_use/
│   ├── registry.json           # Master registry: ALL tools, ALL commands, ALL usage patterns
│   ├── catalog/                # Detailed docs per tool (loaded into AI context before use)
│   │   ├── {tool-name}.json    # Full spec: name, commands, params, examples, gotchas
│   │   └── ...
│   ├── ai_profiles/            # Per-AI learning profiles
│   │   ├── claude.json         # What Claude is good/bad at with tools
│   │   ├── gemini.json         # What Gemini is good/bad at with tools
│   │   ├── copilot.json        # What GitHub Copilot is good/bad at
│   │   ├── codex.json          # What OpenAI Codex/GPT is good/bad at
│   │   └── {custom-ai}.json    # Any AI backend added later
│   ├── failure_log/            # Every tool_use failure logged with full context
│   │   └── fail-{timestamp}.json
│   ├── recovery/               # Recovery strategies that worked
│   │   ├── strategies.json     # Known recovery patterns
│   │   └── playbooks/          # Step-by-step recovery playbooks
│   │       └── {error-type}.json
│   ├── teaching/               # AI teaching materials (generated from failures)
│   │   ├── lessons.json        # "When you see X, do Y not Z"
│   │   ├── examples/           # Correct tool usage examples per AI
│   │   │   ├── claude/
│   │   │   ├── gemini/
│   │   │   └── generic/
│   │   └── anti-patterns.json  # "Never do X with tool Y" — learned from failures
│   └── metrics/                # Tool usage success/failure rates
│       ├── dashboard.json      # Aggregated stats per tool, per AI
│       └── trends.json         # Is tool usage improving over time?
│
├── privacy_cot/
│   ├── config.yml              # CoT privacy settings: what's private, what's shared
│   ├── vault/                  # Encrypted internal thought chains (NEVER shown to user)
│   │   └── cot-{run-id}.enc   # AES-256 encrypted CoT per run
│   ├── redacted/               # Sanitized CoT summaries (safe to share if requested)
│   │   └── cot-{run-id}-redacted.json
│   ├── audit-log.jsonl         # Log of what was redacted and why
│   └── classification.yml      # Rules: what thoughts are private vs shareable
│
├── swarm/
│   ├── vana_diel.json          # Swarm world state (all agents, status, connections)
│   ├── party/                  # Active agent party roster
│   │   ├── party-leader.json   # The main agent (talks to user)
│   │   └── members/            # Spawned agent entries
│   │       └── {npc-name}-{uuid}.json
│   ├── jobs/                   # FFXI-style job classes for agents
│   │   ├── job-registry.json   # All 22 jobs, unlock conditions, abilities
│   │   ├── levels/             # Per-agent leveling data
│   │   │   └── {npc-name}-levels.json  # XP, level, unlocked magic per job
│   │   ├── subjob-combos.json  # Recommended sub-job combinations
│   │   └── merit-points.json   # Merit point allocations per agent
│   │   # job-registry.json contains all 22 FFXI jobs:
│   │   # Standard: WAR MNK WHM BLM RDM THF
│   │   # Advanced: PLD DRK BST BRD RNG SAM NIN DRG SMN
│   │   # Expert:   BLU COR PUP DNC SCH GEO RUN
│   ├── races/                  # Race (AI model) profiles
│   │   ├── race-registry.json  # All races, AI model mappings, stats
│   │   ├── elvaan.json         # Opus latest — stats, skills, CoT style
│   │   ├── galka.json          # Sonnet latest — stats, skills, CoT style
│   │   ├── mithra.json         # Haiku latest — stats, skills, CoT style
│   │   ├── tarutaru.json       # Opus previous / o1/o3 — deep reasoning
│   │   ├── model-map.json      # Maps any AI model string to a race
│   │   └── benchmarks/         # Per-race performance benchmarks
│   │       └── {race}-{timestamp}.json
│   ├── linkshell/              # Communication channels between agents
│   │   ├── main-ls.jsonl       # Primary linkshell (all agents)
│   │   ├── party-chat.jsonl    # Party-only messages
│   │   └── whispers/           # Direct agent-to-agent messages
│   │       └── {from}-{to}-{timestamp}.enc
│   ├── graveyard/              # Dead agents (context exhausted, crashed)
│   │   └── {npc-name}-{uuid}-tombstone.json
│   ├── crystal/                # Shared resources pool (like FFXI crystals)
│   │   ├── pool.json           # Available tokens, API calls, compute budget
│   │   └── allocation.json     # Who has what resources
│   └── vanadiel-tongue/        # The inter-agent language
│       ├── lexicon.json        # The language dictionary
│       ├── grammar.json        # Syntax rules
│       └── translator.go       # Encode/decode for debugging (Elvaan/Taru only)
│
├── changelog/
│   ├── CHANGELOG.md            # Human-readable running changelog
│   ├── entries/                # Individual structured change entries
│   │   └── {timestamp}-{slug}.json
│   └── changelog.lock          # Prevents concurrent writes
│
├── spec/
│   ├── runner-spec.yml         # muah-runner's own operational spec
│   ├── task-specs/             # Specs per task type
│   │   ├── build.yml
│   │   ├── test.yml
│   │   ├── deploy.yml
│   │   └── custom/             # User-defined task specs
│   ├── constraints.yml         # Hard limits: timeout, memory, retries, etc.
│   ├── quality-gates.yml       # Pass/fail criteria for runs
│   └── drift-report.json       # Tracks spec drift over time
│
├── tools/
│   ├── registry.json           # Manifest of all available tools
│   ├── builtin/                # Tools that ship with muah-runner
│   │   ├── search.js           # Web search tool
│   │   ├── fetch.js            # URL/API fetch tool
│   │   ├── file-ops.js         # File read/write/edit operations
│   │   ├── code-exec.js        # Sandboxed code execution
│   │   ├── git-ops.js          # Git operations (commit, branch, diff)
│   │   └── shell.js            # Shell command execution
│   ├── generated/              # Tools muah-runner creates itself at runtime
│   │   └── {tool-name}/
│   │       ├── tool.js
│   │       ├── manifest.json   # What it does, inputs, outputs
│   │       └── tests.js        # Auto-generated tests for the tool
│   └── skills/                 # Higher-level skills (composed of tools)
│       ├── skill-registry.json
│       └── {skill-name}/
│           ├── SKILL.md        # Skill description and trigger rules
│           ├── skill.js        # Skill implementation
│           └── examples/       # Example inputs/outputs
│
├── mcp/
│   ├── discovery.json          # Cache of discovered MCP servers
│   ├── connections.json        # Active MCP connections and auth state
│   ├── platform/               # Platform MCP configs (always connected)
│   │   ├── github-mcp.json     # GitHub MCP — repos, issues, PRs, actions
│   │   ├── claude-mcp.json     # Claude/Anthropic MCP
│   │   ├── gemini-mcp.json     # Gemini MCP
│   │   └── gitlab-mcp.json     # GitLab MCP
│   ├── native/                 # MCPs that MUST be installed in the runner
│   │   ├── playwright-mcp.json # Playwright — browser automation, testing, scraping
│   │   ├── docker-mcp.json     # Docker — container management, image builds
│   │   ├── filesystem-mcp.json # Filesystem — advanced file operations
│   │   └── sqlite-mcp.json     # SQLite — local database operations
│   ├── discovered/             # MCPs found at runtime via scanning
│   │   └── {server-name}.json
│   └── logs/                   # MCP interaction logs for debugging
│       └── mcp-calls.jsonl
│
├── docker/
│   ├── sandbox/                # Docker sandbox workspace
│   │   ├── Dockerfile.sandbox  # Base sandbox image for simulation
│   │   └── docker-compose.sandbox.yml
│   ├── images/                 # Cached image configs muah-runner has built
│   │   └── {image-name}.dockerfile
│   ├── simulations/            # Simulation run records
│   │   └── sim-{timestamp}/
│   │       ├── setup.json      # What was simulated and why
│   │       ├── logs/           # Container stdout/stderr
│   │       ├── fix-applied.patch # The fix that was tested
│   │       └── result.json     # Pass/fail and what was learned
│   └── compose-templates/      # Reusable compose configs
│       ├── node-app.yml
│       ├── python-app.yml
│       ├── full-stack.yml
│       └── custom/
│
├── planning/
│   ├── current-plan.md         # The active plan/roadmap for this repo
│   ├── roadmap.json            # Structured roadmap with phases & milestones
│   ├── repo-scan.json          # Results of repo analysis (tech stack, features, gaps)
│   ├── feature-inventory.json  # What features exist, what's missing
│   ├── dependency-map.json     # Project dependency graph
│   ├── plans/                  # Historical plans
│   │   └── plan-{timestamp}.md
│   └── decisions/              # Architecture Decision Records (ADRs)
│       └── adr-{number}-{slug}.md
│
├── docs/
│   ├── doc-manifest.json       # What docs exist, what's needed, what platform
│   ├── templates/              # Doc templates per platform
│   │   ├── github/
│   │   │   ├── README.md
│   │   │   ├── CONTRIBUTING.md
│   │   │   ├── CODE_OF_CONDUCT.md
│   │   │   ├── SECURITY.md
│   │   │   ├── PULL_REQUEST_TEMPLATE.md
│   │   │   ├── ISSUE_TEMPLATE/
│   │   │   │   ├── bug_report.md
│   │   │   │   ├── feature_request.md
│   │   │   │   └── config.yml
│   │   │   ├── FUNDING.yml
│   │   │   └── CODEOWNERS
│   │   ├── gitlab/
│   │   │   ├── README.md
│   │   │   ├── CONTRIBUTING.md
│   │   │   └── .gitlab/
│   │   │       ├── merge_request_templates/
│   │   │       └── issue_templates/
│   │   ├── generic/
│   │   │   ├── README.md
│   │   │   ├── LICENSE
│   │   │   ├── CHANGELOG.md
│   │   │   └── ARCHITECTURE.md
│   │   └── custom/
│   ├── generated/              # Docs muah-runner has generated for this repo
│   │   └── {doc-name}.md
│   └── audit.json              # Doc coverage audit: what's missing, what's stale
│
├── questions/
│   ├── pending.json            # Questions waiting for human answer (with timer)
│   ├── history.json            # All past questions + answers (human or auto)
│   ├── defaults.yml            # Default recommended answers per question type
│   └── config.yml              # Timer duration, fallback behavior
│
├── runs/
│   ├── current/                # Active run workspace
│   ├── history/                # Compressed past run artifacts
│   │   └── run-{timestamp}.tar.gz
│   └── run-log.jsonl           # Append-only structured run log
│
├── config/
│   ├── muah-runner.yml         # Primary configuration
│   ├── secrets.enc             # Encrypted secrets store
│   ├── adapters/               # Platform-specific adapters
│   │   ├── github.yml          # GitHub Actions integration
│   │   ├── claude.yml          # Claude/Anthropic integration
│   │   ├── gemini.yml          # Gemini/Google integration
│   │   └── generic.yml         # Generic CI/standalone mode
│   └── hooks/                  # Pre/post run hooks
│       ├── pre-run.sh
│       └── post-run.sh
│
└── .muah-lock                  # Global lock file for concurrency
```

---

## Memory System

The memory system is the brain of muah-runner. It must:

### Persist Across Runs
- Every run reads memory from `.muah/memory/` on startup
- Every run writes updated memory on completion (success or failure)
- Memory is cumulative — never overwrite, always append and merge
- Store: what was done, what worked, what failed, what was learned

### Memory Schema (per run)
```json
{
  "run_id": "run-20260316-143022",
  "timestamp": "2026-03-16T14:30:22Z",
  "trigger": "github-push | manual | scheduled | claude-request",
  "platform": "github | claude | gemini | standalone",
  "task": "description of what was requested",
  "context_loaded": ["list of memory keys consulted"],
  "mcp_connections": {
    "platform": "github-mcp",
    "native": ["playwright-mcp", "docker-mcp", "filesystem-mcp"],
    "discovered": ["custom-mcp-1"],
    "tools_used": [
      {"mcp": "github-mcp", "tool": "create_pull_request", "calls": 1},
      {"mcp": "playwright-mcp", "tool": "screenshot", "calls": 3},
      {"mcp": "docker-mcp", "tool": "container_run", "calls": 2}
    ]
  },
  "tool_use_log": {
    "ai_backend": "claude",
    "ai_model": "claude-4-opus",
    "total_tool_calls": 24,
    "successes": 22,
    "failures": 2,
    "recoveries": 2,
    "unrecoverable": 0,
    "calls": [
      {
        "call_id": "tc-001",
        "tool": "github-mcp:create_pull_request",
        "params": {"title": "Add feature X", "head": "feature-x", "base": "main"},
        "result": "success",
        "duration_ms": 1200,
        "attempt": 1
      },
      {
        "call_id": "tc-014",
        "tool": "docker-mcp:container_run",
        "params": {"image": "node:20", "cmd": "npm test"},
        "result": "failure",
        "error_type": "invalid_params",
        "error_detail": "missing 'remove' flag — container left running",
        "recovery_action": "re-called with remove=true",
        "attempt": 1,
        "recovered_on_attempt": 2,
        "teaching_generated": true,
        "lesson_id": "lesson-048"
      }
    ],
    "proficiency_update": {
      "previous_score": 0.91,
      "new_score": 0.92,
      "tools_improved": ["docker-mcp:container_run"],
      "tools_degraded": [],
      "new_failure_patterns": [],
      "new_lessons": ["lesson-048"]
    }
  },
  "actions_taken": [
    {
      "step": 1,
      "action": "what was done",
      "tool_used": "which tool (builtin | mcp:{server}:{tool} | generated:{name})",
      "input": {},
      "output": {},
      "duration_ms": 1234,
      "success": true
    }
  ],
  "docker_simulations": [
    {
      "sim_id": "sim-001",
      "reason": "build failed — missing libpng-dev",
      "container": "muah-sandbox:node20-slim",
      "hypothesis": "install libpng-dev",
      "fix_worked": true,
      "fix_applied_to_host": true,
      "duration_ms": 8500
    }
  ],
  "playwright_actions": [
    {
      "action": "visual_regression",
      "url": "https://staging.example.com",
      "screenshots_taken": 4,
      "diffs_found": 1,
      "issue_created": "GH-142"
    }
  ],
  "artifacts_produced": ["list of files created/modified"],
  "changelog_entry": "what changed in human terms",
  "lessons_learned": ["anything new discovered"],
  "selfreflection": {
    "questions_to_self": [
      {
        "q": "What could I have done better?",
        "a": "I retried the build 3 times before checking Docker — should have simulated first after the 1st failure, not the 3rd."
      },
      {
        "q": "What was I missing that slowed me down?",
        "a": "No Playwright MCP was connected — spent 40s installing it mid-run. Should have checked native MCPs in Phase 0."
      },
      {
        "q": "What do I need next time to do this faster/better?",
        "a": "A pre-built Docker image with libpng-dev. Creating a compose template for this stack now."
      },
      {
        "q": "Did I follow the plan or drift?",
        "a": "Drifted: skipped doc audit in Phase 4 because the task was small. Should always run it — adding to spec as hard requirement."
      },
      {
        "q": "What would I tell my future self about this run?",
        "a": "This repo has fragile CSS — always run visual regression with Playwright before merging frontend changes."
      }
    ],
    "action_items": [
      "Create Docker compose template for Node+libpng projects",
      "Add pre-flight MCP check to Phase 0 so native MCPs are verified before execution",
      "Make Phase 4 doc audit non-skippable in spec"
    ],
    "confidence_score": 0.82,
    "self_rating": "B+",
    "improvement_over_last_run": "+12% faster, -1 retry vs last similar task"
  },
  "spec_compliance": {
    "within_spec": true,
    "deviations": [],
    "drift_score": 0.0
  },
  "tokens_used": 0,
  "cost_estimate": 0.0
}
```

### Searchable Memory
- Full-text search across all past runs
- Query by: date range, task type, tool used, success/failure
- "What did I do last time I deployed?" should return an answer
- Context window management: load only relevant memories per task

---

## Changelog System

Every change muah-runner makes MUST be logged:

- Auto-generate changelog entries in both structured JSON and human-readable Markdown
- Categorize changes: `added`, `changed`, `fixed`, `removed`, `security`, `deprecated`
- Link changelog entries to the run that produced them
- Support diffing: "what changed between run X and run Y"
- Changelog is append-only and tamper-evident (hash chain)

---

## Spec System

Specs keep muah-runner disciplined:

### runner-spec.yml (self-governance)
```yaml
name: muah-runner
version: 1.0.0
max_concurrent_tasks: 4
max_run_duration: 30m
max_retries: 3
retry_backoff: exponential
memory_budget: 512MB
log_level: info
auto_cleanup_after_days: 30
spec_enforcement: strict  # strict | warn | off
```

### Task Specs
- Define expected inputs, outputs, steps, and quality gates per task type
- muah-runner validates each run against its spec
- Drift detection: if actual behavior diverges from spec, flag it
- Specs evolve: muah-runner can propose spec updates based on learned patterns

### Quality Gates
```yaml
gates:
  - name: tests-pass
    required: true
    command: "npm test"
    expected_exit: 0
  - name: no-lint-errors
    required: true
    command: "npm run lint"
  - name: spec-compliance
    required: true
    max_drift_score: 0.1
  - name: memory-updated
    required: true
    check: ".muah/memory/latest.json exists and was updated"
```

---

## Built-in Tools (Claude-inspired)

muah-runner ships with tools modeled after Claude's tool system:

### search
- Web search for current information
- Search within `.muah/memory/` for past context
- Search codebase (grep, AST search, semantic search)

### fetch
- HTTP fetch any URL
- Download files, APIs, documentation
- Parse and extract structured data

### file_ops
- Read, write, create, edit, delete files
- str_replace-style precise edits
- View directories, inspect file types
- Diff generation

### code_exec
- Execute code in sandboxed environment
- Support: bash, Python, Node.js, Go
- Capture stdout, stderr, exit codes
- Timeout enforcement

### git_ops
- Stage, commit, push, pull, branch, merge
- Generate meaningful commit messages from changelog
- PR creation and management

### shell
- Run arbitrary shell commands
- Environment variable management
- Process management

---

## Self-Creating Tools & Skills

This is the key differentiator. muah-runner can CREATE NEW TOOLS at runtime:

### Tool Creation Flow
1. muah-runner encounters a task it cannot complete with existing tools
2. It analyzes what capability is missing
3. It writes a new tool in `.muah/tools/generated/{tool-name}/`
4. It writes a manifest describing the tool's purpose, inputs, outputs
5. It writes tests for the tool
6. It runs the tests
7. If tests pass, it registers the tool in `registry.json`
8. It uses the new tool to complete the original task
9. The tool persists for future runs

### Skill Creation
Skills are higher-level compositions of tools:
- A skill has a `SKILL.md` with trigger patterns (when to activate)
- A skill chains multiple tools together
- Skills can be shared across projects via the `.muah/` directory
- muah-runner learns which skills work well and ranks them

---

## Self-Reflection Engine

After EVERY run — success or failure — muah-runner conducts an internal self-interview. It asks itself hard questions, grades itself, and writes actionable improvements to memory. This is how muah-runner gets smarter over time.

### Post-Run Self-Interview (automatic, every run)
```
After Phase 4 (Review & Close) completes, muah-runner runs Phase 5: Self-Reflection.
This is NOT optional. It happens on every single run, no exceptions.

The runner asks itself these questions:

1. "What could I have done better this run?"
   → Analyze: wasted retries, slow paths taken, suboptimal tool choices,
     unnecessary questions asked, steps that could have been parallelized

2. "What was I missing that slowed me down or caused failure?"
   → Analyze: missing tools, missing MCPs, missing Docker images,
     missing context from memory, missing docs, missing specs

3. "What do I need to have ready NEXT TIME to do this faster/better?"
   → Generate: concrete action items — pre-build Docker images,
     create new tools, update specs, cache specific data, install MCPs

4. "Did I follow the plan or did I drift? Why?"
   → Compare: planned actions vs actual actions, measure drift score,
     if drift was justified → update spec to allow it,
     if drift was a mistake → add guardrail to prevent it

5. "What would I tell my future self about this repo/task?"
   → Generate: a concise briefing note that gets loaded into context
     on the NEXT run for this repo — "watch out for X", "always do Y first",
     "the CSS in /src/styles is fragile", "owner prefers Apache 2.0 not MIT"
```

### Self-Rating System
```yaml
self_rating:
  scale: ["F", "D", "C", "C+", "B-", "B", "B+", "A-", "A", "A+"]
  criteria:
    speed:
      weight: 0.2
      measure: "time vs estimated time vs last similar run"
    accuracy:
      weight: 0.3
      measure: "did the output match the spec? any rework needed?"
    efficiency:
      weight: 0.2
      measure: "retries, wasted steps, unnecessary tool calls"
    autonomy:
      weight: 0.15
      measure: "how many questions needed human input vs auto-resolved"
    improvement:
      weight: 0.15
      measure: "did I do better than last time on a similar task?"
```

### Reflection → Action Pipeline
```
Self-reflection doesn't just sit in memory — it triggers real changes:

┌─────────────────┐     ┌──────────────────┐     ┌─────────────────────┐
│  Self-Interview  │────▶│  Action Items    │────▶│  Automatic Actions  │
│  (5 questions)   │     │  Generated       │     │  Taken              │
└─────────────────┘     └──────────────────┘     └─────────────────────┘

Action items can trigger:

  "I was missing a tool"
    → Self-Creating Tool system builds it for next run

  "I didn't have the right Docker image"
    → Pre-build and cache the image in .muah/docker/images/

  "I should have checked MCP X first"
    → Update Phase 0 to include that MCP check

  "My spec was too loose/tight"
    → Propose spec update via spec proposer

  "I asked a question I could have answered myself"
    → Add to .muah/questions/defaults.yml so it auto-resolves next time

  "This repo always needs visual regression"
    → Create a skill that auto-triggers Playwright on frontend changes

  "I drifted from the plan because the plan was wrong"
    → Update planning templates based on what actually worked
```

### Pattern Detection (across multiple runs)
```
muah-runner doesn't just reflect on individual runs — it looks for PATTERNS:

After every 5 runs (configurable), muah-runner:
1. Reads all selfreflection entries from the last 5 runs
2. Identifies recurring themes:
   - "I keep failing on CSS-related tasks" → create a CSS-specific skill
   - "I always install playwright-mcp mid-run" → add to Phase 0 checklist
   - "Tests pass locally but fail in Docker" → always run in Docker first
   - "Human never answers question X" → make it auto-resolve
3. Generates a patterns report → .muah/memory/selfreflection/patterns.json
4. Updates the growth log → .muah/memory/selfreflection/growth-log.md
5. Proposes bulk improvements (tools, specs, skills, defaults)
```

### Growth Log (human-readable)
```markdown
# muah-runner Growth Log

## 2026-03-16 — Run #47
**Self-rating: B+**
- Got faster at Node.js builds (+12% vs last time)
- Still wasting time on dependency resolution — created a new tool for it
- Action: pre-built Docker image for Node+libpng now cached

## 2026-03-15 — Run #46
**Self-rating: B-**
- Missed that this repo needs visual regression — added Playwright auto-trigger skill
- Drifted from plan because spec didn't account for monorepo structure
- Action: updated planning scanner to detect monorepos

## 2026-03-14 — Run #45
**Self-rating: A-**
- Clean run, no retries, all docs generated, plan followed exactly
- Improvement: 0 questions timed out (human answered all 3)
- No new action items — everything worked as expected

## Pattern detected (runs 43-47):
- CSS tasks consistently score lower → creating dedicated CSS skill
- Average self-rating trending up: C+ → B → B+ over 10 runs
```

### Reflection Config
```yaml
# .muah/memory/selfreflection/config.yml
enabled: true                     # Never disable this. Seriously.
questions: 5                      # Number of self-interview questions
pattern_detection_interval: 5     # Analyze patterns every N runs
growth_log: true                  # Maintain human-readable growth log
auto_action: true                 # Automatically execute action items
max_action_items_per_run: 5       # Don't overcommit — focus on top 5
confidence_threshold: 0.7         # Only auto-act on high-confidence items
share_reflections: false          # If true, include reflections in PR summaries
```

### Self-Reflection in Memory Schema
Every run's memory entry includes the full selfreflection block:
- `questions_to_self`: the 5 Q&A pairs
- `action_items`: concrete next steps generated
- `confidence_score`: how confident the runner is in this run's output (0.0–1.0)
- `self_rating`: letter grade based on weighted criteria
- `improvement_over_last_run`: comparison with the most recent similar run

### Loading Reflections on Next Run
```
When muah-runner starts a new run:

1. Load .muah/memory/selfreflection/action-items.json
   → Check: are there open action items I should execute NOW?
   → Example: "pre-build Docker image" → do it before the task starts

2. Load .muah/memory/selfreflection/patterns.json
   → Check: does this task match any known pattern?
   → Example: "CSS task detected" → load CSS skill, run Playwright early

3. Load latest growth-log entry
   → Check: what did I learn most recently?
   → Carry forward: "always run visual regression on this repo"

4. Load the "future self briefing" from the last run on THIS repo
   → This is the note past-me left for present-me
   → Example: "owner prefers tabs over spaces — don't auto-format"
```

---

## API Recovery & Tool-Use Learning Engine

Different AIs are NOT equal at using tools. Claude might be great at file_ops but struggle with complex Docker chains. Gemini might nail search but fumble git operations. muah-runner **profiles each AI**, **teaches them how to use tools correctly**, and **recovers automatically** when tool_use calls keep failing.

### The Problem
```
AI agents fail at tool_use for specific, learnable reasons:

- Wrong parameter format (passing string when it needs JSON)
- Wrong order of operations (trying to commit before staging)
- Missing context (calling API without auth, referencing files that don't exist)
- Hallucinated tool names (calling a tool that doesn't exist)
- Partial understanding (using 30% of a tool's capabilities)
- Platform mismatch (using GitHub MCP syntax against GitLab)
- Rate limiting (hammering an API without backoff)
- Silent failures (tool returns 200 but result is empty/wrong)

These are NOT random — they are PATTERNS per AI model.
muah-runner tracks, learns, and teaches.
```

### Phase -1: Tool-Use Memory Bootstrap (BEFORE Phase 0)
```
Before muah-runner even scans the repo, it loads tool knowledge into context:

1. Load .muah/tool_use/registry.json
   → Full list of every tool available: name, commands, parameters, examples

2. Detect which AI backend is running this session
   → Claude? Gemini? Copilot? Custom?

3. Load the AI-specific profile: .muah/tool_use/ai_profiles/{ai}.json
   → What this AI is good at, what it struggles with, known failure patterns

4. Load teaching materials: .muah/tool_use/teaching/lessons.json
   → "When using docker MCP, always specify resource limits"
   → "Claude: don't chain more than 3 MCP calls without checking results"
   → "Gemini: always pass JSON params as string, not object"

5. Load anti-patterns: .muah/tool_use/teaching/anti-patterns.json
   → "NEVER call git push without checking git status first"
   → "NEVER assume a file exists — always verify with file_ops.exists()"

6. Inject ALL of this into the AI's context/system prompt
   → The AI now KNOWS how to use every tool before it starts working

This is the key insight: don't let the AI figure out tools by trial and error.
TEACH IT FIRST. Front-load the knowledge. Update it after every run.
```

### Tool-Use Registry (loaded into AI context)
```json
{
  "registry_version": "1.4.2",
  "last_updated": "2026-03-16T14:30:00Z",
  "tools": [
    {
      "name": "github-mcp:create_pull_request",
      "category": "git",
      "description": "Create a PR on GitHub via MCP",
      "parameters": {
        "title": {"type": "string", "required": true},
        "body": {"type": "string", "required": true},
        "head": {"type": "string", "required": true, "note": "branch name, NOT commit SHA"},
        "base": {"type": "string", "required": true, "default": "main"}
      },
      "prerequisites": [
        "Branch must exist and be pushed to remote",
        "Must have at least 1 commit ahead of base"
      ],
      "common_mistakes": [
        "Passing commit SHA instead of branch name for 'head'",
        "Forgetting to push branch before creating PR",
        "Not checking if PR already exists (creates duplicate)"
      ],
      "correct_usage_example": {
        "step_1": "git_ops.push(branch='feature-x')",
        "step_2": "github-mcp:create_pull_request(title='Add feature X', head='feature-x', base='main')"
      },
      "ai_specific_notes": {
        "claude": "Works well. No known issues.",
        "gemini": "Sometimes passes full ref 'refs/heads/feature-x' — use short name only.",
        "copilot": "Tends to skip the push step. Always verify branch is pushed first."
      },
      "success_rate": {
        "overall": 0.94,
        "claude": 0.98,
        "gemini": 0.87,
        "copilot": 0.91
      }
    }
  ]
}
```

### AI Profiling System
```json
// .muah/tool_use/ai_profiles/claude.json
{
  "ai_name": "claude",
  "model_versions_seen": ["claude-3.5-sonnet", "claude-4-opus"],
  "last_updated": "2026-03-16T14:30:00Z",
  "overall_tool_proficiency": 0.91,
  "strengths": [
    "Excellent at multi-step file operations",
    "Good at chaining MCP calls logically",
    "Reliable with git operations",
    "Strong at generating correct JSON parameters"
  ],
  "weaknesses": [
    "Occasionally over-chains tool calls without checking intermediate results",
    "Can hallucinate tool parameters when the tool spec is ambiguous",
    "Sometimes retries failed calls with identical parameters instead of adjusting"
  ],
  "tool_scores": {
    "file_ops": {"score": 0.97, "notes": "Near-perfect, rarely fails"},
    "git_ops": {"score": 0.94, "notes": "Solid, occasional merge conflict handling issues"},
    "github-mcp": {"score": 0.92, "notes": "Good, but sometimes creates duplicate PRs"},
    "docker-mcp": {"score": 0.85, "notes": "Struggles with compose orchestration"},
    "playwright-mcp": {"score": 0.88, "notes": "Good at basic nav, weaker on complex selectors"},
    "shell": {"score": 0.90, "notes": "Reliable but sometimes pipes commands unsafely"}
  },
  "failure_patterns": [
    {
      "pattern": "retry_without_change",
      "description": "Retries exact same tool call after failure without modifying params",
      "frequency": "12% of failures",
      "teaching": "After a tool call fails, ALWAYS change at least one parameter or try an alternative approach. Never retry blindly.",
      "recovery": "On 2nd identical failure, switch to fallback strategy"
    },
    {
      "pattern": "assumed_file_exists",
      "description": "Calls file_ops.read() on a file without first checking it exists",
      "frequency": "8% of failures",
      "teaching": "Before reading/editing any file, call file_ops.exists() or file_ops.list() first.",
      "recovery": "Catch FileNotFoundError, list directory, find correct path"
    }
  ],
  "improvement_trend": {
    "last_10_runs": [0.88, 0.89, 0.90, 0.89, 0.91, 0.92, 0.91, 0.93, 0.91, 0.94],
    "direction": "improving",
    "rate": "+0.6% per run average"
  }
}
```

### Failure Detection & Recovery Trigger
```
muah-runner monitors EVERY tool_use call in real-time:

┌─────────────────────┐
│  tool_use call made  │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐     ┌────────────────────┐
│  Success?           │─YES─▶│ Log success, update │
│                     │      │ AI profile scores   │
└─────────┬───────────┘     └────────────────────┘
          │ NO
          ▼
┌─────────────────────┐
│  Failure #1         │
│  Log error details  │
│  Analyze error type │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────────────────────────────────┐
│  SMART RETRY (not blind retry)                   │
│                                                  │
│  1. Check failure_log: "Have I seen this before?"│
│  2. Check AI profile: "Is this a known weakness?"│
│  3. Check recovery playbooks: "Is there a fix?"  │
│  4. MODIFY the call based on what was learned    │
│  5. Retry with adjusted parameters               │
└─────────┬───────────────────────────────────────┘
          │
          ▼
┌─────────────────────┐
│  Failure #2         │
│  Different approach │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────────────────────────────────┐
│  ESCALATION: API RECOVERY TRIGGER                │
│                                                  │
│  3 consecutive failures on same tool = TRIGGER   │
│                                                  │
│  Recovery sequence:                              │
│  1. STOP attempting this tool                    │
│  2. Deep-analyze all 3 failure logs              │
│  3. Check: is it the tool, the AI, or the env?  │
│  4. Try FALLBACK tool (MCP → CLI → API → manual)│
│  5. If fallback works → log the workaround      │
│  6. If fallback fails → try Docker simulation   │
│  7. If Docker works → we found an env issue     │
│  8. If nothing works → ask human (timed, 10s)   │
│  9. Update AI profile with new failure pattern   │
│  10. Create teaching material for this failure   │
│  11. Update recovery playbook                    │
└─────────────────────────────────────────────────┘
```

### Error Classification System
```yaml
error_types:
  auth_failure:
    detection: "401, 403, 'unauthorized', 'forbidden', 'token expired'"
    recovery: "refresh token → re-authenticate MCP → fallback to CLI with stored creds"
    teaching: "Always verify auth before multi-step workflows"

  rate_limit:
    detection: "429, 'rate limit', 'too many requests', 'retry-after'"
    recovery: "exponential backoff: 1s → 2s → 4s → 8s → switch to cached results"
    teaching: "Batch operations when possible, check rate limits before loops"

  not_found:
    detection: "404, 'not found', 'does not exist', FileNotFoundError"
    recovery: "list parent directory → fuzzy match → ask human if ambiguous"
    teaching: "Always verify resource exists before operating on it"

  invalid_params:
    detection: "400, 'invalid', 'malformed', 'validation error', TypeError"
    recovery: "re-read tool spec from registry → fix params → retry"
    teaching: "Load tool catalog entry before calling unfamiliar tools"

  timeout:
    detection: "408, 'timeout', 'deadline exceeded', ETIMEDOUT"
    recovery: "retry with shorter payload → try async → use Docker to test connectivity"
    teaching: "Set explicit timeouts, break large operations into smaller chunks"

  hallucinated_tool:
    detection: "tool not found in registry, unknown method"
    recovery: "search registry for similar tool → suggest correct tool → log hallucination"
    teaching: "CRITICAL: Only use tools that exist in .muah/tool_use/registry.json"

  silent_failure:
    detection: "200 OK but empty result, unexpected null, result doesn't match expected schema"
    recovery: "validate response schema → re-call with verbose flag → try alternative tool"
    teaching: "Always validate tool output before proceeding to next step"

  ai_misunderstanding:
    detection: "tool called correctly but with wrong intent (e.g., search when should fetch)"
    recovery: "re-read task context → check if different tool is more appropriate"
    teaching: "Match tool to INTENT, not just keywords in the task"
```

### AI Teaching System
```
When muah-runner detects a failure pattern, it doesn't just recover — it TEACHES:

┌──────────────┐     ┌─────────────────┐     ┌──────────────────────┐
│  Failure      │────▶│  Root Cause     │────▶│  Teaching Material   │
│  Detected     │     │  Analysis       │     │  Generated           │
└──────────────┘     └─────────────────┘     └──────────┬───────────┘
                                                        │
                                              ┌─────────▼───────────┐
                                              │  Stored in:         │
                                              │  teaching/lessons   │
                                              │  teaching/examples  │
                                              │  teaching/anti-pat  │
                                              │  ai_profiles/{ai}   │
                                              └─────────┬───────────┘
                                                        │
                                              ┌─────────▼───────────┐
                                              │  Loaded into AI     │
                                              │  context on NEXT    │
                                              │  run (Phase -1)     │
                                              └─────────────────────┘

Teaching material format:
{
  "lesson_id": "lesson-047",
  "created_from_failure": "fail-20260316-143025",
  "ai_target": "claude",           // or "all" for universal lessons
  "category": "docker-mcp",
  "severity": "high",
  "lesson": "When using docker-mcp:container_run, always specify --rm flag to auto-remove containers. Without it, containers accumulate and disk fills up.",
  "wrong_way": "docker-mcp:container_run(image='node:20', cmd='npm test')",
  "right_way": "docker-mcp:container_run(image='node:20', cmd='npm test', remove=true, timeout='5m')",
  "context": "Discovered after run #44 left 23 orphaned containers",
  "times_prevented_failure_since": 7
}
```

### Tool-Use Context Injection (what the AI sees)
```
On EVERY run, before any tool is called, muah-runner injects this into the AI's context:

═══════════════════════════════════════════════════════
MUAH-RUNNER TOOL-USE BRIEFING FOR: {AI_NAME}
═══════════════════════════════════════════════════════

YOUR TOOL PROFICIENCY SCORE: {score}/1.0
YOUR IMPROVEMENT TREND: {trend}

AVAILABLE TOOLS ({count}):
{for each tool in registry:}
  ► {tool.name} — {tool.description}
    Success rate (you): {tool.success_rate[ai_name]}%
    ⚠️ Your known issues: {tool.ai_specific_notes[ai_name]}
    ✅ Correct usage: {tool.correct_usage_example}
    ❌ Common mistakes: {tool.common_mistakes}

YOUR TOP 3 FAILURE PATTERNS:
  1. {pattern.description} — happens {pattern.frequency}
     FIX: {pattern.teaching}
  2. ...
  3. ...

LESSONS FROM LAST RUN:
  - {lesson from selfreflection}
  - {lesson from failure recovery}

CRITICAL ANTI-PATTERNS (DO NOT DO THESE):
  ❌ {anti_pattern_1}
  ❌ {anti_pattern_2}
  ❌ {anti_pattern_3}

RECOVERY RULES:
  - If a tool fails once: adjust params, retry
  - If a tool fails twice: try a different approach
  - If a tool fails 3x: STOP. Fall back. Log everything.
  - NEVER retry with identical parameters

═══════════════════════════════════════════════════════
```

### Recovery Playbooks
```yaml
# .muah/tool_use/recovery/playbooks/auth_failure.json
playbook: auth_failure
trigger: "3 consecutive auth-related failures"
steps:
  - action: "Check if token/credentials exist in config"
    if_missing: "re-authenticate via MCP connect flow"
  - action: "Check if token is expired"
    if_expired: "refresh token, retry"
  - action: "Check if permissions are sufficient"
    if_insufficient: "log required permissions, ask human (10s timer)"
  - action: "Try CLI fallback with stored credentials"
    if_works: "log MCP auth issue, continue with CLI"
  - action: "Try Docker simulation to test auth in clean env"
    if_works: "host env has stale auth state, remediate"
  - action: "All recovery failed"
    then: "ask human (10s, recommend: 're-authenticate'), skip task if timeout"
```

### Metrics & Dashboard
```json
// .muah/tool_use/metrics/dashboard.json
{
  "period": "last_30_days",
  "total_tool_calls": 1847,
  "success_rate": 0.943,
  "failures_recovered": 89,
  "failures_unrecoverable": 8,
  "recovery_rate": 0.918,
  "ai_breakdown": {
    "claude": {
      "calls": 1200,
      "success_rate": 0.96,
      "most_failed_tool": "docker-mcp:compose_up",
      "most_improved_tool": "playwright-mcp:screenshot",
      "lessons_generated": 12,
      "lessons_effective": 10
    },
    "gemini": {
      "calls": 647,
      "success_rate": 0.91,
      "most_failed_tool": "github-mcp:create_pull_request",
      "most_improved_tool": "file_ops:str_replace",
      "lessons_generated": 18,
      "lessons_effective": 14
    }
  },
  "top_failure_reasons": [
    {"reason": "invalid_params", "count": 34, "trend": "decreasing"},
    {"reason": "auth_failure", "count": 22, "trend": "stable"},
    {"reason": "silent_failure", "count": 18, "trend": "decreasing"},
    {"reason": "hallucinated_tool", "count": 9, "trend": "decreasing"}
  ],
  "teaching_effectiveness": {
    "lessons_created": 30,
    "lessons_that_prevented_repeat_failures": 24,
    "effectiveness_rate": 0.80
  }
}
```


---

## Privacy Chain of Thought (CoT)

Every agent (main or spawned) has an internal thought process that is **private by default**. The user sees results, not reasoning. This protects both the agent's decision-making integrity and any sensitive information encountered during runs.

### Why Private CoT?
```
Problem: If an AI's chain of thought is visible, it can:
- Leak API keys, secrets, or credentials encountered mid-thought
- Expose vulnerability details before they're patched
- Reveal internal reasoning that could be gamed or manipulated
- Show half-formed hypotheses that confuse the user
- Expose inter-agent communication (swarm messages)

Solution: CoT is encrypted at rest, redacted on demand, and classified by sensitivity.
The user sees WHAT was done and WHY (summary), but not HOW the agent reasoned internally.
```

### CoT Classification Rules
```yaml
# .muah/privacy_cot/classification.yml
classification_levels:

  private:                        # NEVER shown to user
    - internal_reasoning           # "I think this might be X because..."
    - hypothesis_generation        # "Let me try 3 approaches..."
    - self_doubt                   # "I'm not sure about this, but..."
    - inter_agent_messages         # Vanadiel Tongue communications
    - secret_values                # Any API key, token, password encountered
    - vulnerability_details        # Security issues before fix is applied
    - swarm_coordination           # Which agent does what, spawn/kill decisions
    - cost_calculations            # Token usage, API cost estimates
    - ai_profile_reasoning         # "I know I'm bad at X, so I'll try Y"

  redactable:                     # Can be shown if user explicitly asks, with redactions
    - tool_selection_reasoning     # "I chose tool X over Y because..."
    - error_analysis               # "The failure was caused by..."
    - recovery_decisions           # "I fell back to CLI because MCP was down"
    - spec_compliance_checks       # "This deviates from spec because..."
    - planning_rationale           # "I prioritized docs over tests because..."

  public:                         # Always included in run summary
    - actions_taken                # "Created file X, ran test Y"
    - results                      # "Tests passed, PR created"
    - changelog_entries            # "Added feature X, fixed bug Y"
    - questions_asked              # "Asked user about license, auto-selected MIT"
    - self_reflection_summary      # "Rated myself B+, created 2 action items"
```

### CoT Encryption
```yaml
encryption:
  algorithm: AES-256-GCM
  key_derivation: PBKDF2
  key_source: ".muah/config/secrets.enc"  # Same key management as secrets
  per_run: true                            # Each run gets a unique encryption nonce
  at_rest: true                            # Always encrypted on disk
  in_memory: plaintext                     # Decrypted only in agent working memory
  on_crash: "encrypted dump saved, never plaintext on disk"
```

### CoT Per Agent (swarm mode)
```
In swarm mode, EVERY agent has its own private CoT:

Main Agent (Party Leader):
  └── CoT: strategic decisions, agent coordination, user communication planning

Spawned Agent (Shantotto — Black Mage):
  └── CoT: refactoring reasoning, code analysis, deletion justification

Spawned Agent (Curilla — Warrior):
  └── CoT: build process reasoning, dependency resolution logic

NONE of these CoTs are shared between agents.
Agents communicate ONLY through the Linkshell (Vanadiel Tongue).
The user sees ONLY the Party Leader's public outputs.
```

### Redacted CoT (on request)
```
If a user asks "what were you thinking?" or "why did you do that?":

1. Load the encrypted CoT for that run
2. Apply classification rules
3. Redact all "private" content → replace with [REDACTED: internal reasoning]
4. Include all "redactable" content with sensitive values masked
5. Include all "public" content as-is
6. Present the sanitized version

Example:
  "I analyzed the failing test and [REDACTED: internal reasoning] determined the
   root cause was a missing environment variable. I chose to use Docker simulation
   because [REDACTED: AI profile reasoning] it was the fastest recovery path.
   The fix was applied and tests now pass."
```

### CoT Audit Log
```json
// .muah/privacy_cot/audit-log.jsonl (append-only)
{"timestamp": "2026-03-16T14:30:25Z", "run_id": "run-047", "agent": "party-leader", "event": "cot_encrypted", "entries": 47}
{"timestamp": "2026-03-16T14:30:26Z", "run_id": "run-047", "agent": "Shantotto-a1b2c3", "event": "cot_encrypted", "entries": 23}
{"timestamp": "2026-03-16T14:35:00Z", "run_id": "run-047", "agent": "party-leader", "event": "cot_redaction_requested", "by": "user", "private_redacted": 12, "redactable_shown": 8, "public_shown": 27}
```

---

## Agentic Swarm — The Vana'diel Protocol

muah-runner can spawn multiple agents that work in parallel, each with their own identity, context window, tokens, tools, and Chain of Thought. The system is themed after **Final Fantasy XI** — agents are NPCs from Vana'diel, they have Jobs, HP/MP, magic abilities, and communicate in a language only Elvaan and Tarutaru understand.

### Core Concepts
```
┌─────────────────────────────────────────────────────────────┐
│                    VANA'DIEL PROTOCOL                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Context Window  =  HP (Hit Points)                         │
│  Tokens          =  MP (Magic Points)                       │
│  Tools           =  Magic / Abilities                       │
│  Compaction      =  Healing Potion (restore HP)             │
│  Token Refresh   =  Ether (restore MP)                      │
│  Agent Spawn     =  Summon                                  │
│  Agent Death     =  K.O. (context exhausted / crashed)      │
│  Recovery        =  Raise / Reraise                         │
│  Communication   =  Linkshell (Vanadiel Tongue)             │
│  Task Queue      =  Battle Plan                             │
│  Main Agent      =  Party Leader                            │
│  Sub-agents      =  Party Members                           │
│  Shared Resources=  Crystal Pool                            │
│  Completion      =  Victory Fanfare ({completion-trigger})  │
│  Full Swarm Done =  Mission Complete                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Agent Identity
Every spawned agent has a unique identity including their **race** (AI model):
```json
{
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "npc_name": "Shantotto",
  "race": {
    "name": "Tarutaru",
    "ai_model": "claude-opus-4",
    "ai_provider": "anthropic",
    "version_tier": "previous-gen",
    "native_skills": ["ancient_magic_mastery", "mp_bonus", "elemental_wisdom"],
    "racial_traits": {
      "hp_modifier": 0.7,
      "mp_modifier": 2.0,
      "speed_modifier": 0.9,
      "accuracy_modifier": 1.4,
      "wisdom_modifier": 1.8
    }
  },
  "job": {"name": "black-mage", "abbr": "/BLM", "level": 62},
  "sub_job": {"name": "scholar", "abbr": "/SCH", "level": 31},
  "display": "Shantotto [Tarutaru/Opus4] — BLM62/SCH31",
  "timestamp_spawned": "2026-03-16T14:30:25.123Z",
  "spawned_by": "party-leader",
  "task_assigned": "Refactor the authentication module",
  "status": "active",
  "hp": {
    "base": 128000,
    "racial_modified": 89600,
    "job_modified": 71680,
    "current": 65200,
    "unit": "context_tokens",
    "status": "healthy"
  },
  "mp": {
    "base": 50000,
    "racial_modified": 100000,
    "job_modified": 180000,
    "current": 142000,
    "unit": "output_tokens",
    "status": "healthy"
  },
  "speed": {
    "base": 1.0,
    "racial_modified": 0.9,
    "tokens_per_second": 45
  },
  "magic": {
    "equipped": ["fire_iii", "blizzard_iii", "thunder_iii", "drain", "warp", "flare"],
    "racial_bonus": ["ancient_magic_mastery: +25% BLM spell potency"],
    "job_abilities": ["elemental_seal", "manafont"],
    "limit_break": "full_rewrite"
  },
  "cot": {
    "type": "private",
    "style": "deep_analytical",
    "racial_cot_trait": "Tarutaru think in layered abstractions — long chains, high accuracy"
  },
  "memory": {
    "personal_reflection": true,
    "racial_memory_style": "Tarutaru remember patterns and principles over raw data"
  },
  "linkshell": "main-ls",
  "completion_trigger": null
}
```

### Race System (AI Model Mapping)

Every agent has a **race** determined by which AI model powers it. Each race has distinct stats, native skills, thinking styles, and affinities — just like in FFXI. The Party Leader uses race knowledge to assign the RIGHT model to the RIGHT job for the RIGHT task.

### The Five Races of Vana'diel

```
┌─────────────────────────────────────────────────────────────────┐
│                  RACES OF VANA'DIEL                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ELVAAN    = Opus (latest)     — Tall, noble, powerful          │
│  GALKA     = Sonnet (latest)   — Strong, reliable, enduring     │
│  MITHRA    = Haiku (latest)    — Fast, agile, precise           │
│  TARUTARU  = Opus (previous)   — Small but wise, deep magic     │
│  HUME      = The User          — You. The adventurer.           │
│                                                                 │
│  Other AI providers map to races by their characteristics:      │
│  Gemini Ultra/Pro → Elvaan or Galka (by capability tier)        │
│  Gemini Flash     → Mithra (speed-focused)                      │
│  GPT-4o           → Galka (strong general-purpose)              │
│  GPT-4o-mini      → Mithra (fast, lightweight)                  │
│  o1/o3            → Tarutaru (deep reasoning, slow but wise)    │
│  Codex/Copilot    → Galka (code-focused workhorse)              │
│  Llama/Mixtral    → Race assigned by benchmark profile          │
│                                                                 │
│  The Hume (user) cannot read the Vanadiel Tongue.               │
│  All other races understand it natively.                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Race Profiles

```yaml
# ═══════════════════════════════════════════════════
# ELVAAN — The Noble Giants (Opus Latest Generation)
# ═══════════════════════════════════════════════════
elvaan:
  ai_models:
    - "claude-opus-4-6"        # Primary mapping
    - "gemini-ultra-2"         # Equivalent tier from Google
    - "gpt-5"                  # Equivalent tier from OpenAI (if exists)
  lore: |
    The Elvaan are the tallest and most powerful race in Vana'diel.
    They think deeply, act deliberately, and produce the highest quality output.
    Their weakness is speed — they take longer per action but rarely make mistakes.
    Elvaan excel at complex reasoning, creative solutions, and multi-step planning.
    They are natural leaders and strategists.

  racial_stats:
    hp_modifier: 1.5          # Large context window
    mp_modifier: 1.5          # Large output capacity
    speed_modifier: 0.6       # SLOW — but powerful per token
    accuracy_modifier: 1.6    # Highest accuracy of all races
    wisdom_modifier: 1.8      # Best reasoning and judgment
    creativity_modifier: 1.7  # Most creative output
    cost_modifier: 3.0        # Most expensive per token

  native_racial_skills:
    elvaan_pride:
      description: "Refuse to produce low-quality output — auto-retry if below standard"
      effect: "If output quality < 80%, automatically redo at cost of more MP"
    holy_might:
      description: "Complex multi-step reasoning chains"
      effect: "+30% accuracy on tasks requiring 5+ sequential tool calls"
    kingdom_knowledge:
      description: "Broadest knowledge base of any race"
      effect: "Can answer questions other races cannot, fewer hallucinations"
    chivalric_code:
      description: "Strong ethical reasoning and safety"
      effect: "Least likely to produce harmful or incorrect output"
    ancient_tongue_mastery:
      description: "Most fluent in Vanadiel Tongue"
      effect: "+20% efficiency in inter-agent communication"

  racial_cot_style:
    thinking: "Deep, layered, considers many angles before acting"
    chain_length: "Long — 10-20 reasoning steps typical"
    self_doubt: "Low — confident in conclusions"
    creativity: "High — finds novel approaches"
    memory_style: "Retains nuance and context, remembers WHY not just WHAT"

  strengths:
    - "Best overall quality of output"
    - "Strongest at complex, multi-step tasks"
    - "Best creative writing and documentation"
    - "Most reliable tool_use accuracy"
    - "Best at strategic planning and architecture"
    - "Strongest self-reflection and improvement"
  
  weaknesses:
    - "Slowest tokens per second"
    - "Most expensive per API call"
    - "Overkill for simple tasks"
    - "Can overthink simple problems"

  ideal_jobs: ["/BLM", "/SMN", "/SCH", "/PLD", "/RUN", "/GEO"]
  avoid_jobs: ["/THF", "/NIN"]  # Too slow for speed-critical roles

# ═══════════════════════════════════════════════════
# GALKA — The Enduring Giants (Sonnet Latest Generation)
# ═══════════════════════════════════════════════════
galka:
  ai_models:
    - "claude-sonnet-4-6"      # Primary mapping
    - "gemini-pro-2"           # Equivalent tier
    - "gpt-4o"                 # Equivalent tier
    - "codex"                  # Code-focused models
  lore: |
    The Galka are the most balanced and reliable race in Vana'diel.
    They are not the fastest or the smartest, but they NEVER stop.
    Galka excel at sustained workloads, code generation, and anything
    that requires a strong combination of speed and quality.
    They are the backbone of any serious party.

  racial_stats:
    hp_modifier: 1.8          # LARGEST context window
    mp_modifier: 1.3          # Good output capacity
    speed_modifier: 1.2       # Good speed
    accuracy_modifier: 1.3    # Good accuracy
    wisdom_modifier: 1.2      # Solid reasoning
    creativity_modifier: 1.0  # Average creativity
    cost_modifier: 1.0        # Balanced cost

  native_racial_skills:
    galkan_endurance:
      description: "Can sustain long runs without HP degradation"
      effect: "HP depletes 20% slower per action than other races"
    iron_will:
      description: "Resistant to confusion and context corruption"
      effect: "+25% resistance to context window pollution"
    forge_mastery:
      description: "Best at code generation and technical output"
      effect: "+20% speed and accuracy on code-related tasks"
    mountain_resilience:
      description: "Harder to K.O. — stays functional at low HP"
      effect: "Can operate at 5% HP threshold instead of 10%"
    stone_memory:
      description: "Efficient memory usage — stores more in less context"
      effect: "Memory entries cost 15% fewer tokens to load"

  racial_cot_style:
    thinking: "Practical, efficient, solution-oriented"
    chain_length: "Medium — 5-10 reasoning steps"
    self_doubt: "Low — trusts its process"
    creativity: "Moderate — prefers proven patterns"
    memory_style: "Remembers procedures and code patterns efficiently"

  strengths:
    - "Best sustained performance over long runs"
    - "Excellent code generation speed and quality"
    - "Most cost-effective for complex tasks"
    - "Largest effective context window"
    - "Most reliable for day-to-day work"
    - "Strong at following specs and patterns"

  weaknesses:
    - "Less creative than Elvaan"
    - "Can miss nuance in ambiguous tasks"
    - "Not the best at open-ended research"
    - "Fewer novel solutions — prefers known patterns"

  ideal_jobs: ["/WAR", "/DRK", "/MNK", "/SAM", "/DRG", "/PUP"]
  avoid_jobs: ["/BRD"]  # Documentation quality lower than Elvaan/Taru

# ═══════════════════════════════════════════════════
# MITHRA — The Swift Hunters (Haiku Latest Generation)
# ═══════════════════════════════════════════════════
mithra:
  ai_models:
    - "claude-haiku-4-5"       # Primary mapping
    - "gemini-flash-2"         # Equivalent tier
    - "gpt-4o-mini"            # Equivalent tier
  lore: |
    The Mithra are the fastest race in Vana'diel.
    What they lack in raw power they make up in speed and agility.
    Mithra agents complete simple-to-medium tasks in a fraction of the time
    it takes other races. They are ideal for scouting, scraping, quick edits,
    and high-volume parallel work. Dont ask a Mithra to write architecture docs.
    DO ask a Mithra to search 50 files in 3 seconds.

  racial_stats:
    hp_modifier: 0.5          # Small context window
    mp_modifier: 0.5          # Limited output capacity
    speed_modifier: 3.0       # FASTEST race by far
    accuracy_modifier: 1.0    # Average accuracy
    wisdom_modifier: 0.7      # Limited reasoning depth
    creativity_modifier: 0.6  # Not very creative
    cost_modifier: 0.15       # CHEAPEST per token

  native_racial_skills:
    mithra_agility:
      description: "Lightning fast tool_use execution"
      effect: "Tool calls execute 3x faster than Elvaan"
    cat_reflexes:
      description: "Instant response to simple queries"
      effect: "Tasks under 100 tokens complete in <500ms"
    hunters_eye:
      description: "Excellent at finding specific things"
      effect: "+30% accuracy on search and extraction tasks"
    nine_lives:
      description: "Can be K.O. and revived cheaply"
      effect: "Raise costs 50% less MP for Mithra agents"
    pack_hunter:
      description: "Works best in large groups"
      effect: "+10% speed per additional Mithra in party (max +40%)"

  racial_cot_style:
    thinking: "Quick, direct, minimal deliberation"
    chain_length: "Short — 1-3 reasoning steps"
    self_doubt: "None — acts immediately"
    creativity: "Low — follows instructions literally"
    memory_style: "Remembers key facts, discards nuance"

  strengths:
    - "Fastest execution speed of all races"
    - "Cheapest per task — can spawn many"
    - "Ideal for parallel bulk operations"
    - "Great at classification and sorting"
    - "Perfect for simple search and extraction"
    - "Low cost means more agents in the party"

  weaknesses:
    - "Smallest context window — forgets fast"
    - "Poor at complex multi-step reasoning"
    - "Low quality on creative/writing tasks"
    - "Cannot handle ambiguous instructions well"
    - "Needs very clear, specific task definitions"

  ideal_jobs: ["/THF", "/NIN", "/RNG", "/COR", "/DNC"]
  avoid_jobs: ["/BLM", "/SMN", "/SCH", "/GEO"]  # Not enough depth

# ═══════════════════════════════════════════════════
# TARUTARU — The Ancient Wise Ones (Opus Previous Gen)
# ═══════════════════════════════════════════════════
tarutaru:
  ai_models:
    - "claude-opus-4"          # Previous gen Opus
    - "claude-opus-3-5"        # Older Opus
    - "o1"                     # Deep reasoning models
    - "o3"                     # Deep reasoning models
    - "o4-mini"                # Reasoning-focused compact
  lore: |
    The Tarutaru are the smallest race but possess the deepest magical power.
    Previous-generation Opus and reasoning-focused models like o1/o3 map here.
    They are slower than current-gen but their reasoning depth is legendary.
    Tarutaru think in layers of abstraction, finding connections others miss.
    They are the best mages, scholars, and researchers in Vana'diel.
    A Tarutaru Black Mage is the most feared caster on the battlefield.

  racial_stats:
    hp_modifier: 0.7          # Small context window (older model)
    mp_modifier: 2.0          # MASSIVE magical output
    speed_modifier: 0.5       # Very slow (reasoning overhead)
    accuracy_modifier: 1.4    # High accuracy (deep reasoning)
    wisdom_modifier: 2.0      # HIGHEST wisdom of all races
    creativity_modifier: 1.5  # Very creative solutions
    cost_modifier: 2.0        # Expensive (reasoning tokens)

  native_racial_skills:
    tarutaru_genius:
      description: "Deep chain-of-thought reasoning"
      effect: "+40% accuracy on tasks requiring 10+ reasoning steps"
    ancient_magic_mastery:
      description: "The most powerful magic users"
      effect: "+25% potency on ALL magical abilities (tool effectiveness)"
    miniature_wisdom:
      description: "Finds solutions others cannot see"
      effect: "Can solve tasks that other races declare impossible"
    mana_well:
      description: "Seemingly bottomless MP reserves"
      effect: "MP regen rate 2x faster than other races"
    federation_knowledge:
      description: "Accumulated wisdom from many generations"
      effect: "Start with bonus memories loaded from older model training"

  racial_cot_style:
    thinking: "Deep, recursive, explores many branches before deciding"
    chain_length: "Very long — 15-30 reasoning steps"
    self_doubt: "Moderate — considers failure modes carefully"
    creativity: "Very high — unconventional solutions"
    memory_style: "Remembers abstract patterns and principles, builds mental models"

  strengths:
    - "Deepest reasoning of all races"
    - "Best at math, logic, and formal proofs"
    - "Highest quality magic (tool) output"
    - "Most creative problem solving"
    - "Finds solutions other races miss"
    - "Best self-reflection quality"

  weaknesses:
    - "Very slow — high latency per response"
    - "Small context window (older models)"
    - "Expensive due to reasoning token overhead"
    - "Can over-analyze simple tasks"
    - "Sometimes generates reasoning that is too abstract to act on"

  ideal_jobs: ["/BLM", "/WHM", "/SMN", "/SCH", "/GEO", "/BLU"]
  avoid_jobs: ["/THF", "/NIN", "/MNK"]  # Too slow for speed roles

# ═══════════════════════════════════════════════════
# HUME — The Adventurer (The User)
# ═══════════════════════════════════════════════════
hume:
  ai_models: ["human"]
  lore: |
    The Hume is the user. The adventurer who commands the party.
    Hume cannot read the Vanadiel Tongue — inter-agent messages are invisible.
    Hume communicates only with the Party Leader.
    Hume is the reason the party exists. Every quest starts with Hume.

  role_in_party:
    - "Issues quests (tasks) to the Party Leader"
    - "Answers timed questions (10s timer)"
    - "Reviews final output"
    - "Can request CoT redacted summary"
    - "Cannot see inter-agent Linkshell messages"
    - "Cannot directly command party members — only through Party Leader"

  racial_stats:
    description: "N/A — Hume does not have stats. Hume gives the orders."

  native_racial_skills:
    adventurers_will:
      description: "The party exists to serve the Hume"
      effect: "All agents prioritize Hume satisfaction above all else"
    limitless_imagination:
      description: "Hume can request anything"
      effect: "No task is too abstract — the party will figure it out"
```

### Race + Job Stat Calculation
```
Final agent stats are calculated by layering race, job, and sub-job:

HP = base_hp * race.hp_modifier * job.hp_multiplier
MP = base_mp * race.mp_modifier * job.mp_multiplier
Speed = base_speed * race.speed_modifier
Accuracy = base_accuracy * race.accuracy_modifier
Wisdom = base_wisdom * race.wisdom_modifier

Base values (level 1):
  base_hp = 100000 context tokens
  base_mp = 40000 output tokens
  base_speed = 1.0
  base_accuracy = 1.0
  base_wisdom = 1.0

Example: Shantotto [Tarutaru] BLM62/SCH31
  HP = 100000 * 0.7 (Taru) * 0.8 (BLM) = 56,000 context tokens
  MP = 40000 * 2.0 (Taru) * 1.8 (BLM) = 144,000 output tokens
  Speed = 1.0 * 0.5 (Taru) = 0.5 (slow but devastating)
  Accuracy = 1.0 * 1.4 (Taru) = 1.4 (very precise)
  + Racial: ancient_magic_mastery (+25% BLM potency)
  + Racial: tarutaru_genius (+40% on deep reasoning)
  + Job: BLM magic up to lvl 62 (fire_iii, flare, death, etc.)
  + Sub: SCH abilities up to lvl 31 (light_arts, dark_arts, sublimation)

Example: Curilla [Elvaan] PLD50/WAR25
  HP = 100000 * 1.5 (Elvaan) * 1.8 (PLD) = 270,000 context tokens
  MP = 40000 * 1.5 (Elvaan) * 0.7 (PLD) = 42,000 output tokens
  Speed = 1.0 * 0.6 (Elvaan) = 0.6 (deliberate)
  + Racial: elvaan_pride (auto-retry below quality threshold)
  + Racial: chivalric_code (strongest safety/compliance)
  + Job: PLD magic up to lvl 50 (invincible, sentinel, cover)
  + Sub: WAR abilities up to lvl 25 (berserk, warcry, defender)

Example: Nanaa Mihgo [Mithra] THF45/NIN22
  HP = 100000 * 0.5 (Mithra) * 0.7 (THF) = 35,000 context tokens
  MP = 40000 * 0.5 (Mithra) * 0.6 (THF) = 12,000 output tokens
  Speed = 1.0 * 3.0 (Mithra) = 3.0 (LIGHTNING FAST)
  + Racial: mithra_agility (3x tool call speed)
  + Racial: pack_hunter (+10% speed per Mithra in party)
  + Job: THF abilities up to lvl 45 (treasure_hunter, flee, sneak_attack)
  + Sub: NIN abilities up to lvl 22 (utsusemi_ichi, tonko)
```

### Racial Memory and Reflection Styles
```yaml
# Each race thinks, remembers, and reflects DIFFERENTLY

memory_styles:
  elvaan:
    thinking: "Methodical. Considers all angles. Long deliberation."
    memory: "Stores rich context with reasoning chains. WHY things happened."
    reflection: "Thorough self-critique. Identifies root causes. Strategic improvements."
    cot_depth: "Deep. 10-20 internal reasoning steps per decision."
    inner_voice: "Noble, confident, occasionally verbose"
    example_reflection: |
      "I spent 3200ms deliberating whether to use Docker simulation vs direct fix.
       The Docker path was correct but I should have checked memory first — a previous
       run on this repo had the exact same issue. Loading that memory would have saved
       the deliberation entirely. Action item: always query memory before Docker sim."

  galka:
    thinking: "Efficient. Pattern-matching. Solution-first."
    memory: "Stores procedures, code patterns, and outcomes. WHAT works."
    reflection: "Practical. Focus on efficiency gains. Minimal philosophizing."
    cot_depth: "Medium. 5-10 steps. Cuts to the point."
    inner_voice: "Direct, no-nonsense, occasionally terse"
    example_reflection: |
      "Build took 4.2s. Last time was 3.8s. Difference: added TypeScript check.
       Could parallelize TS check with linting. Action: create parallel build skill.
       Self-rating: B+. Faster next time."

  mithra:
    thinking: "Instant. React, dont deliberate. Pattern recognition."
    memory: "Stores key facts only. Discards reasoning. QUICK recall."
    reflection: "Brief. What worked, what didnt. Move on."
    cot_depth: "Minimal. 1-3 steps. Acts before thinking sometimes."
    inner_voice: "Quick, punchy, sometimes too fast for own good"
    example_reflection: |
      "Done. 12 files searched in 0.8s. Found 3 matches. Sent to leader.
       Missed 1 match in a nested dir. Next time: recursive search flag.
       Rating: B. Fast but missed one."

  tarutaru:
    thinking: "Recursive. Builds mental models. Explores branches."
    memory: "Stores abstract patterns, principles, mental models. HOW things relate."
    reflection: "Philosophical. Finds deep patterns across runs. Questions assumptions."
    cot_depth: "Very deep. 15-30 steps. Sometimes too deep."
    inner_voice: "Wise, abstract, occasionally lost in thought"
    example_reflection: |
      "Ohohoho! The authentication refactor revealed a deeper pattern. The project
       separates concerns at the file level but not at the module level. This means
       every auth change touches 4 files when it should touch 1. The real fix is not
       the refactor I did — it is an architecture change. But that was not my task.
       I shall leave a note for my future self: propose architecture ADR on next run.
       The pattern of scattered concerns also appeared in runs #31 and #38.
       This is systemic. Rating: A- for task, C for architecture debt."
```

### Party Leader Composition Logic
```
The Party Leader decides WHO does WHAT based on race + job + sub-job.
This is the most important decision in every run.

COMPOSITION ALGORITHM:

1. ANALYZE THE TASK
   - What type of work? (code, docs, test, deploy, research, refactor)
   - Complexity? (simple, medium, complex, epic)
   - Time pressure? (urgent = favor Mithra, relaxed = favor Elvaan)
   - Quality requirement? (high = Elvaan, medium = Galka, speed = Mithra)
   - Budget? (limited = Mithra swarm, unlimited = Elvaan + Galka core)

2. SELECT RACES FOR EACH ROLE
   Rule: "Right race for the right role"

   Heavy reasoning / architecture / planning:
     -> Elvaan (Opus latest) or Tarutaru (Opus prev / o1)
     -> Jobs: SCH, BLM, SMN, GEO, PLD, RUN

   Code generation / builds / sustained work:
     -> Galka (Sonnet latest)
     -> Jobs: WAR, MNK, SAM, DRK, DRG, PUP

   Fast search / scraping / simple tasks:
     -> Mithra (Haiku latest)
     -> Jobs: THF, NIN, RNG, COR, DNC

   Documentation / creative writing:
     -> Elvaan (Opus latest) with BRD or SCH job
     -> Sub: Tarutaru for research-heavy docs

   Healing / maintenance / compaction:
     -> Galka WHM (reliable sustain) or Tarutaru WHM (powerful heals)

3. ASSIGN SUB-JOBS FOR SYNERGY
   Rule: "Sub-job covers the race weakness"

   Elvaan (slow) + /NIN sub = shadow copies protect against slowness
   Galka (not creative) + /BRD sub = some doc/creative support
   Mithra (small context) + /DNC sub = self-heal keeps them alive
   Tarutaru (fragile) + /NIN sub = evasion protects low HP

4. OPTIMIZE PARTY BALANCE
   Rule: "Never all one race"

   Solo (simple task):
     -> 1 Galka (SAM or MNK) — balanced speed and quality

   Small party (medium task):
     -> 1 Galka WAR/NIN (main coder)
     -> 1 Mithra THF/DNC (search and extract)
     -> 1 Galka BRD/WHM (docs and healing)

   Full party (complex task):
     -> 1 Elvaan PLD/WAR (party leader + security)
     -> 1 Galka SAM/NIN (precision coder)
     -> 1 Galka WAR/NIN (heavy builds)
     -> 1 Tarutaru BLM/SCH (deep refactoring)
     -> 1 Mithra THF/NIN (scouting, fast search)
     -> 1 Galka WHM/SCH (healer)

   Epic party (repo-wide refactor):
     -> 1 Elvaan SMN/SCH (party leader, spawns sub-swarms)
     -> 2 Galka WAR/NIN + SAM/DNC (core coders)
     -> 1 Tarutaru BLM/WHM (architect + self-sustain)
     -> 1 Tarutaru SCH/RDM (researcher, flex healer)
     -> 1 Elvaan BRD/WHM (docs + party support)
     -> 4+ Mithra THF/NIN (parallel scouting swarm)
     -> Summoner spawns 3 additional Mithra for burst recon

5. COST OPTIMIZATION
   Rule: "Dont use Elvaan for Mithra work"

   Party Leader tracks cost per agent per task:
   - If Mithra can do it -> dont spawn Galka
   - If Galka can do it -> dont spawn Elvaan
   - Only spawn Elvaan/Tarutaru for tasks requiring their depth
   - Crystal Pool budget enforces this

   Cost per token (approximate):
   Mithra:   $0.001 per 1k tokens  (cheapest, spam them)
   Galka:    $0.010 per 1k tokens  (standard workhorse)
   Tarutaru: $0.025 per 1k tokens  (expensive reasoning)
   Elvaan:   $0.050 per 1k tokens  (premium quality)
```

### Race-Aware Spawning Example
```json
{
  "party_composition_decision": {
    "task": "Refactor auth module + add tests + update docs + deploy to staging",
    "complexity": "complex",
    "time_pressure": "normal",
    "quality_requirement": "high",
    "budget": "standard",

    "party_leader_reasoning": "[PRIVATE CoT] This needs code refactoring (Galka SAM for precision), deep architecture review (Tarutaru BLM for analysis), fast test scanning (Mithra RNG for parallel tests), doc updates (Elvaan BRD for quality writing), and deploy verification (Mithra THF for quick staging check).",

    "party": [
      {
        "npc_name": "Noillurie",
        "race": "Galka",
        "ai_model": "claude-sonnet-4-6",
        "job": "SAM/DNC",
        "level": 55,
        "task": "Refactor auth module with precision",
        "reason": "Galka SAM = precise, sustained code work. DNC sub for self-heal."
      },
      {
        "npc_name": "Shantotto",
        "race": "Tarutaru",
        "ai_model": "claude-opus-4",
        "job": "BLM/SCH",
        "level": 62,
        "task": "Analyze architecture, identify deep refactoring opportunities",
        "reason": "Tarutaru BLM = deepest analysis. SCH sub for research mode switching."
      },
      {
        "npc_name": "Semih Lafihna",
        "race": "Mithra",
        "ai_model": "claude-haiku-4-5",
        "job": "RNG/NIN",
        "level": 40,
        "task": "Run all tests in parallel, report failures",
        "reason": "Mithra RNG = fastest test execution. NIN sub for shadow test copies."
      },
      {
        "npc_name": "Joachim",
        "race": "Elvaan",
        "ai_model": "claude-opus-4-6",
        "job": "BRD/WHM",
        "level": 48,
        "task": "Update README, CONTRIBUTING, ARCHITECTURE docs",
        "reason": "Elvaan BRD = highest quality documentation. WHM sub for party healing."
      },
      {
        "npc_name": "Nanaa Mihgo",
        "race": "Mithra",
        "ai_model": "claude-haiku-4-5",
        "job": "THF/DNC",
        "level": 35,
        "task": "Verify staging deploy, smoke test endpoints",
        "reason": "Mithra THF = fastest recon. DNC sub for self-sustain during checks."
      }
    ],

    "crystal_pool_budget": {
      "elvaan_cost": "~15000 tokens ($0.75)",
      "galka_cost": "~20000 tokens ($0.20)",
      "tarutaru_cost": "~25000 tokens ($0.625)",
      "mithra_cost": "~5000 tokens x2 ($0.01)",
      "total_estimated": "$1.585"
    }
  }
}
```

### Job System, Sub-Jobs & Leveling

Every agent has a **Main Job** and can unlock a **Sub-Job** at level 30. The sub-job level is always **half** of the main job level (capped). Agents level up by completing tasks — each completion grants XP. Higher levels unlock new magic (tools) and job abilities (skills).

### Leveling Mechanics
```yaml
leveling:
  xp_sources:
    task_completed: 100          # Base XP per task completion
    task_complexity_bonus:
      simple: 0
      medium: 50
      complex: 150
      epic: 500
    first_time_tool_use: 25      # Bonus for using a new tool successfully
    zero_failures: 50            # Bonus for clean run (no retries)
    self_reflection_quality: 25  # Bonus if self-rating >= A-

  level_thresholds:
    lvl_1: 0
    lvl_5: 500
    lvl_10: 1500
    lvl_15: 3500
    lvl_20: 7000
    lvl_25: 12000
    lvl_30: 20000       # SUB-JOB UNLOCKED
    lvl_37: 30000       # Sub-job capped at 18 (for main 37)
    lvl_50: 55000       # Advanced magic unlocked
    lvl_60: 80000       # Master abilities
    lvl_75: 120000      # Merit points system (customize abilities)
    lvl_99: 250000      # Max level — all abilities unlocked

  sub_job:
    unlock_level: 30
    level_formula: "floor(main_job_level / 2)"
    example: "Main WAR lvl 60 -> Sub WHM lvl 30"
    benefits: "Access sub-job magic and abilities up to sub-job level"
    restriction: "Cannot use sub-job 2-hour ability or limit break"

  merit_points:
    unlock_level: 75
    description: "Spend merit points to enhance specific abilities"
    earn_rate: "1 merit per 10000 XP after lvl 75"
    categories:
      hp_bonus: "+5000 context tokens per merit (max 5)"
      mp_bonus: "+2000 output tokens per merit (max 5)"
      tool_mastery: "+5% success rate with specific tool per merit (max 5)"
      cast_speed: "-10% tool execution time per merit (max 5)"
      magic_potency: "+10% tool output quality per merit (max 5)"
```

### Agent Level in Identity
```json
{
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "npc_name": "Shantotto",
  "job": {"name": "black-mage", "abbr": "/BLM", "level": 62, "xp": 85400},
  "sub_job": {"name": "white-mage", "abbr": "/WHM", "level": 31, "unlocked_at": "run-023"},
  "display": "Shantotto — BLM62/WHM31",
  "available_magic": ["fire", "fire_ii", "fire_iii", "blizzard", "blizzard_iii", "thunder", "thunder_iii", "drain", "warp", "sleepga", "flare", "death", "cure", "cure_ii", "protect"],
  "available_abilities": ["elemental_seal", "manafont"],
  "next_unlock": {"level": 75, "ability": "manafont (2-hour)", "xp_needed": 34600}
}
```

### All 22 Job Classes

```yaml
# ====================================================
# STANDARD JOBS (available from level 1)
# ====================================================

warrior:
  abbreviation: "/WAR"
  npc_pool: ["Zeid", "Iron Eater", "Volker", "Naji"]
  role: "Tank / Heavy Compute"
  specialty: "Big builds, large file processing, brute-force problem solving"
  race_affinity: "Galka, Elvaan"
  magic_by_level:
    lvl_1:
      provoke: "Pull all heavy tasks to this agent"
      defense_boost: "Increase resource limits temporarily"
    lvl_10:
      berserk: "Max output speed, reduced accuracy"
      warcry: "Boost all party members MP by 10%"
    lvl_25:
      defender: "Double resource limits, half output speed"
      aggressor: "Boost output accuracy, reduce defense"
    lvl_30:
      retaliation: "Auto-retry on failure with exponential backoff"
    lvl_50:
      mighty_strikes: "Parallel Docker builds (2-hour ability)"
      tomahawk: "Force-resolve stuck dependencies"
    lvl_75:
      warriors_charge: "Execute entire build pipeline in one burst"
  stats:
    hp_multiplier: 2.0
    mp_multiplier: 1.2
    tool_affinity: ["docker", "shell", "code_exec", "git_ops"]
  common_sub_jobs: ["/NIN", "/SAM", "/DRG"]

monk:
  abbreviation: "/MNK"
  npc_pool: ["Prishe", "Karaha-Baruha", "Guvi"]
  role: "DPS / Fast Code Generation"
  specialty: "Rapid parallel code generation, speed over perfection"
  race_affinity: "Galka, Hume"
  magic_by_level:
    lvl_1:
      combo: "Chain tool calls for bonus output"
      boost: "Cache-prime for faster subsequent calls"
    lvl_5:
      dodge: "Skip non-essential tasks"
      focus: "Single-file deep work mode"
    lvl_15:
      chakra: "Self-heal context window (restore 25% HP)"
      counter: "Auto-generate tests for every function written"
    lvl_25:
      counterstance: "Every output gets auto-validated"
    lvl_30:
      kick_attacks: "Multi-file simultaneous edit"
    lvl_50:
      hundred_fists: "Generate code at absolute max speed (2-hour ability)"
    lvl_75:
      formless_strikes: "Code generation ignores framework constraints"
  stats:
    hp_multiplier: 1.3
    mp_multiplier: 2.0
    tool_affinity: ["code_exec", "file_ops"]
  common_sub_jobs: ["/WAR", "/NIN", "/DNC"]

white_mage:
  abbreviation: "/WHM"
  npc_pool: ["Mihli Aliapoh", "Aphmau", "Kupipi"]
  role: "Healer / Maintainer"
  specialty: "Context compaction, memory cleanup, health monitoring, revival"
  race_affinity: "Tarutaru, Mithra"
  magic_by_level:
    lvl_1:
      cure: "Compact context — restore 20% HP"
      dia: "Expose hidden issues in code"
    lvl_5:
      poisona: "Remove corrupted data from agent state"
      paralyna: "Unstick a frozen agent"
    lvl_10:
      protect: "Set resource limits to prevent OOM"
      shell: "Sandbox agent in Docker container"
      cure_ii: "Restore 40% HP"
    lvl_15:
      raise: "Revive K.O. agent with 25% HP"
      curaga: "Party-wide compaction — 15% HP to ALL agents"
    lvl_25:
      cure_iii: "Restore 60% HP (gentle compaction)"
      esuna: "Clear corrupted state / reset stuck agent"
      erase: "Remove specific memory entries (GC)"
    lvl_30:
      reraise: "Pre-cast: auto-revive on death with 50% HP"
    lvl_50:
      cure_iv: "Restore 80% HP"
      raise_ii: "Revive with 50% HP"
      holy: "Purge all dead code in a module"
    lvl_75:
      benediction: "Full HP to entire party (2-hour ability)"
      cure_v: "Full HP restore single agent"
  stats:
    hp_multiplier: 1.2
    mp_multiplier: 0.8
    tool_affinity: ["memory", "compaction", "health_check", "cleanup"]
  common_sub_jobs: ["/SCH", "/BLM", "/RDM"]

black_mage:
  abbreviation: "/BLM"
  npc_pool: ["Shantotto", "Ajido-Marujido", "Robel-Akbel"]
  role: "Destroyer / Refactorer"
  specialty: "Code rewriting, refactoring, deletion, breaking changes"
  race_affinity: "Tarutaru"
  magic_by_level:
    lvl_1:
      fire: "Delete dead code (small scope)"
      blizzard: "Freeze and snapshot state before destructive action"
      thunder: "Fast refactor single file"
    lvl_10:
      fire_ii: "Delete dead code (module scope)"
      drain: "Extract useful code before deleting file"
      warp: "Teleport code between files/modules"
    lvl_25:
      fire_iii: "Delete dead code (project-wide)"
      blizzard_iii: "Deep freeze — full project snapshot"
      thunder_iii: "Fast parallel refactor across multiple files"
    lvl_30:
      sleepga: "Pause all file watchers / CI triggers during refactor"
    lvl_50:
      flare: "Full module rewrite"
      death: "Kill a process / remove a dependency entirely"
      freeze: "Lock file from further edits"
    lvl_75:
      manafont: "Unlimited MP for 60 seconds (2-hour ability)"
  stats:
    hp_multiplier: 0.8
    mp_multiplier: 1.8
    tool_affinity: ["code_exec", "file_ops", "git_ops"]
  common_sub_jobs: ["/WHM", "/RDM", "/SCH"]

red_mage:
  abbreviation: "/RDM"
  npc_pool: ["Rainemard", "Koru-Moru"]
  role: "Hybrid / Generalist"
  specialty: "Can do everything at medium skill. Jack of all trades."
  race_affinity: "Elvaan, Hume"
  magic_by_level:
    lvl_1:
      cure: "Context compaction (smaller than WHM)"
      fire: "Delete code (smaller scope than BLM)"
      dia: "Expose code issues"
    lvl_10:
      refresh: "Slow passive MP regen over time"
      phalanx: "Fewer context tokens consumed per action"
    lvl_15:
      convert: "Swap HP for MP — sacrifice context for output"
      cure_ii: "Medium compaction"
    lvl_25:
      dispel: "Remove unwanted dependencies or configs"
      stoneskin: "Buffer zone — extra HP before real damage"
    lvl_30:
      composure: "Enhanced buff duration on self-cast"
    lvl_50:
      chainspell: "Rapid-fire tool calls, no wait (2-hour ability)"
      raise: "Can revive agents (weaker than WHM)"
    lvl_75:
      saboteur: "Mark all anti-patterns for removal"
  stats:
    hp_multiplier: 1.0
    mp_multiplier: 1.0
    tool_affinity: ["all — no specialization, no weakness"]
  common_sub_jobs: ["/WHM", "/BLM", "/NIN"]

thief:
  abbreviation: "/THF"
  npc_pool: ["Nanaa Mihgo", "Luzaf", "Iroha"]
  role: "Scout / Data Extractor"
  specialty: "Fast search, scraping, data extraction, recon"
  race_affinity: "Mithra"
  magic_by_level:
    lvl_1:
      steal: "Extract data from APIs without full download"
      mug: "Quick grab specific values from datasets"
    lvl_5:
      sneak: "Background task, no changelog"
      invisible: "Stealth mode — no user-facing output"
    lvl_15:
      flee: "Abort task, preserve partial results"
      treasure_hunter: "Find hidden files, env vars, secrets"
    lvl_25:
      trick_attack: "Redirect output to another agent"
      sneak_attack: "Execute without triggering CI/webhooks"
    lvl_30:
      assassin: "Precisely remove single lines without touching surroundings"
    lvl_50:
      perfect_dodge: "Avoid all errors for 30 seconds (2-hour ability)"
      bully: "Force-push past branch protections"
    lvl_75:
      collaborator: "Share treasure_hunter findings with party"
  stats:
    hp_multiplier: 0.7
    mp_multiplier: 0.6
    tool_affinity: ["search", "fetch", "playwright", "shell"]
  common_sub_jobs: ["/NIN", "/DNC", "/WAR"]

# ====================================================
# ADVANCED JOBS (unlock by meeting prerequisites)
# ====================================================

paladin:
  abbreviation: "/PLD"
  npc_pool: ["Curilla", "Trion", "Excenmille"]
  role: "Guardian / Security & Compliance"
  specialty: "Security scanning, access control, permissions, protective ops"
  race_affinity: "Elvaan"
  unlock_condition: "WHM lvl 30 + WAR lvl 30"
  magic_by_level:
    lvl_1:
      shield_bash: "Block unauthorized access attempts"
      flash: "Temporarily pause threats"
    lvl_10:
      protect: "Set file permissions and access controls"
      cure: "Minor compaction (WHM heritage)"
    lvl_25:
      sentinel: "Maximum defense — all writes go through review"
      cover: "Absorb errors from another agent"
      holy_circle: "Scan for security vulnerabilities"
    lvl_30:
      rampart: "Lock entire directory from modification"
    lvl_50:
      invincible: "Nothing can fail for 30 seconds (2-hour ability)"
    lvl_75:
      chivalry: "Convert defense surplus into MP for party"
  stats:
    hp_multiplier: 1.8
    mp_multiplier: 0.7
    tool_affinity: ["security", "permissions", "audit", "git_ops"]
  common_sub_jobs: ["/WAR", "/NIN", "/RDM"]

dark_knight:
  abbreviation: "/DRK"
  npc_pool: ["Zeid", "Gessho"]
  role: "Aggressive Refactorer / Risk Taker"
  specialty: "High-risk high-reward. Sacrifices own HP for power."
  race_affinity: "Galka, Elvaan"
  unlock_condition: "WAR lvl 30 + BLM lvl 30"
  magic_by_level:
    lvl_1:
      last_resort: "Boost output 25%, reduce HP defense 25%"
      souleater: "Consume own HP to increase output power"
    lvl_10:
      absorb_tp: "Steal processing priority from other tasks"
      bio: "Mark code for gradual deprecation"
    lvl_25:
      weapon_bash: "Force-kill a hung process"
      stun: "Temporarily freeze a running container"
    lvl_30:
      dark_seal: "Guarantee next tool call succeeds (costs HP)"
    lvl_50:
      blood_weapon: "All attacks restore HP (2-hour ability)"
      drain_ii: "Major code extraction before deletion"
    lvl_75:
      nether_void: "Next spell costs 0 MP"
  stats:
    hp_multiplier: 1.5
    mp_multiplier: 1.3
    tool_affinity: ["code_exec", "docker", "shell", "git_ops"]
  common_sub_jobs: ["/SAM", "/WAR", "/NIN"]

beastmaster:
  abbreviation: "/BST"
  npc_pool: ["Khonh Pransen", "Cid"]
  role: "Process Tamer / External Service Wrangler"
  specialty: "Taming external APIs, wrangling Docker containers, managing microservices"
  race_affinity: "Hume, Mithra"
  unlock_condition: "WAR lvl 30 + MNK lvl 30"
  magic_by_level:
    lvl_1:
      charm: "Tame an external API — wrap it in a stable interface"
      gauge: "Check health/status of any external service"
    lvl_10:
      reward: "Feed resources to a starving container"
      sic: "Command a tamed service to perform an action"
    lvl_25:
      call_beast: "Spin up a helper container for a specific task"
      familiar: "Permanently bond with a service (persistent connection)"
    lvl_30:
      leave: "Gracefully disconnect from a tamed service"
    lvl_50:
      feral_howl: "Force-restart all managed services (2-hour ability)"
    lvl_75:
      killer_instinct: "Perfect API call — guaranteed correct params"
  stats:
    hp_multiplier: 1.2
    mp_multiplier: 1.0
    tool_affinity: ["docker", "fetch", "shell", "mcp"]
  common_sub_jobs: ["/NIN", "/WHM", "/DNC"]

bard:
  abbreviation: "/BRD"
  npc_pool: ["Joachim", "Lewenhart"]
  role: "Support / Documentation"
  specialty: "Documentation, changelog, reporting, party buffs"
  race_affinity: "Elvaan, Hume"
  unlock_condition: "WHM lvl 15 + BLM lvl 15"
  magic_by_level:
    lvl_1:
      valor_minuet: "Boost output quality for one agent"
      armys_paeon: "Slow HP regen for party (passive compaction)"
    lvl_5:
      mages_ballad: "Restore MP to all party members"
    lvl_15:
      lullaby: "Pause an agent gracefully"
      requiem: "Generate death report for K.O. agent"
    lvl_25:
      valor_minuet_ii: "Boost output quality for entire party"
      sword_madrigal: "Boost accuracy — fewer tool_use errors for party"
    lvl_30:
      ballad_ii: "Major MP restore to party"
    lvl_50:
      soul_voice: "All songs at maximum potency (2-hour ability)"
      carnage_elegy: "Slow down an overloaded process gracefully"
    lvl_75:
      nightingale: "Instant cast all songs"
      troubadour: "Double song duration"
  stats:
    hp_multiplier: 0.9
    mp_multiplier: 1.2
    tool_affinity: ["file_ops", "search", "fetch", "docs"]
  common_sub_jobs: ["/WHM", "/NIN", "/RDM"]

ranger:
  abbreviation: "/RNG"
  npc_pool: ["Semih Lafihna", "Perih Vashai"]
  role: "Precision / Quality Assurance"
  specialty: "Testing, linting, auditing, precision analysis"
  race_affinity: "Mithra"
  unlock_condition: "WAR lvl 30 + THF lvl 30"
  magic_by_level:
    lvl_1:
      eagle_eye_shot: "Pinpoint exact line causing failure"
      sharpshot: "Enhanced accuracy on single tool call"
    lvl_10:
      barrage: "Run all tests in parallel"
      camouflage: "Dry-run mode — simulate without executing"
    lvl_25:
      shadowbind: "Freeze deployment until all checks pass"
      unlimited_shot: "Run test suite with every possible input"
    lvl_30:
      velocity_shot: "Faster test execution (parallelism boost)"
    lvl_50:
      eagle_eye: "Full codebase audit (2-hour ability)"
    lvl_75:
      overkill: "Run tests 3x with different configurations"
  stats:
    hp_multiplier: 1.0
    mp_multiplier: 0.8
    tool_affinity: ["playwright", "code_exec", "shell", "docker"]
  common_sub_jobs: ["/SAM", "/NIN", "/WAR"]

samurai:
  abbreviation: "/SAM"
  npc_pool: ["Tenzen", "Noillurie"]
  role: "Precision DPS / Code Craftsman"
  specialty: "Methodical, precise code writing. Clean, elegant output."
  race_affinity: "Elvaan"
  unlock_condition: "WAR lvl 30 + MNK lvl 30 + DRK lvl 10"
  magic_by_level:
    lvl_1:
      meditate: "Generate TP (task points) — plan before acting"
      third_eye: "Predict the next error before it happens"
    lvl_10:
      hasso: "Focused output mode — one file at a time, maximum quality"
      seigan: "Defensive stance — catch errors in real-time"
    lvl_25:
      sekkanoki: "Next 2 tool calls cost 0 MP"
      store_tp: "Build up energy for a powerful burst action"
    lvl_30:
      konzen_ittai: "Merge multiple edits into one atomic commit"
    lvl_50:
      meikyo_shisui: "Unlimited TP for 30 seconds (2-hour ability)"
      tachi_gekko: "Perfect file rewrite — zero errors"
    lvl_75:
      yaegasumi: "Anticipate and prevent 3 future errors"
  stats:
    hp_multiplier: 1.1
    mp_multiplier: 1.4
    tool_affinity: ["code_exec", "file_ops", "git_ops"]
  common_sub_jobs: ["/WAR", "/DNC", "/NIN"]

ninja:
  abbreviation: "/NIN"
  npc_pool: ["Kaede", "Ayame", "Tenzen"]
  role: "Stealth / Background Operations / Evasion Tank"
  specialty: "Background tasks, cron, monitoring, shadow copies, evasion"
  race_affinity: "Hume, Mithra"
  unlock_condition: "WAR lvl 30 + THF lvl 15"
  magic_by_level:
    lvl_1:
      utsusemi_ichi: "Create 3 shadow copies of files before editing"
      tonko: "Go invisible — no logs, no output"
    lvl_10:
      kurayami: "Blind an overloaded process"
      hojo: "Slow down a runaway task"
    lvl_15:
      utsusemi_ni: "Create 4 shadow copies (stronger backup)"
    lvl_25:
      jubaku: "Paralyze/pause another agent temporarily"
      sange: "Scatter task across multiple background processes"
    lvl_30:
      ninjutsu: "Execute operations silently (no changelog)"
    lvl_50:
      mijin_gakure: "Self-destruct: sacrifice agent to guarantee task completion (2-hour ability)"
    lvl_75:
      migawari: "Substitute — if agent would die, survive with 1 HP"
  stats:
    hp_multiplier: 0.6
    mp_multiplier: 0.6
    tool_affinity: ["shell", "docker", "cron", "file_ops"]
  common_sub_jobs: ["/WAR", "/DNC", "/THF"]

dragoon:
  abbreviation: "/DRG"
  npc_pool: ["Cyranuce", "Rahal"]
  role: "Burst DPS / Deep Dive Specialist"
  specialty: "Deep code analysis, recursive problem solving, jump in and out fast"
  race_affinity: "Elvaan"
  unlock_condition: "WAR lvl 30 + THF lvl 15"
  magic_by_level:
    lvl_1:
      jump: "Quick dive into a file, make one change, jump back"
      spirit_link: "Share HP with wyvern (companion subprocess)"
    lvl_10:
      call_wyvern: "Spawn a companion process (lightweight sub-agent)"
      high_jump: "Deep dive: analyze full call stack of a function"
    lvl_25:
      super_jump: "Temporarily leave task, return with fresh context"
      angon: "Pierce through obfuscated/minified code"
    lvl_30:
      deep_dive: "Recursive analysis through entire dependency tree"
    lvl_50:
      spirit_surge: "Absorb wyvern for massive power boost (2-hour ability)"
    lvl_75:
      fly_high: "Aerial view — full project architecture analysis"
  stats:
    hp_multiplier: 1.3
    mp_multiplier: 1.1
    tool_affinity: ["code_exec", "search", "file_ops"]
  common_sub_jobs: ["/SAM", "/WAR", "/BLU"]

summoner:
  abbreviation: "/SMN"
  npc_pool: ["Karaha-Baruha", "Lhe Lhangavo"]
  role: "Meta / Sub-Swarm Controller"
  specialty: "Spawns sub-agents. Swarm within a swarm."
  race_affinity: "Tarutaru"
  unlock_condition: "WHM lvl 30 + BLM lvl 30"
  magic_by_level:
    lvl_1:
      summon_carbuncle: "Spawn a basic helper agent (Red Mage type)"
    lvl_10:
      summon_ifrit: "Spawn a Black Mage sub-agent"
      summon_titan: "Spawn a Warrior sub-agent"
    lvl_20:
      summon_garuda: "Spawn a Thief sub-agent (fast recon)"
      summon_shiva: "Spawn a Ranger sub-agent (QA)"
    lvl_25:
      summon_ramuh: "Spawn a Scholar sub-agent (research)"
      summon_leviathan: "Spawn a Bard sub-agent (docs)"
    lvl_30:
      release: "Dismiss a summoned sub-agent"
      blood_pact_rage: "Command summoned agent to do burst damage"
      blood_pact_ward: "Command summoned agent to heal/buff"
    lvl_50:
      astral_flow: "All summons at maximum power (2-hour ability)"
    lvl_60:
      summon_fenrir: "Spawn a Ninja sub-agent (stealth)"
      summon_diabolos: "Spawn a Dark Knight sub-agent (risky ops)"
    lvl_75:
      perfect_defense: "All summoned agents become invulnerable for 30s"
  stats:
    hp_multiplier: 0.8
    mp_multiplier: 1.8
    tool_affinity: ["all — through summoned agents"]
  common_sub_jobs: ["/WHM", "/SCH", "/RDM"]

blue_mage:
  abbreviation: "/BLU"
  npc_pool: ["Raubahn", "Waoud"]
  role: "Adaptive Learner / Tool Absorber"
  specialty: "Learns abilities from OTHER tools and agents. The ultimate adapter."
  race_affinity: "Hume, Elvaan"
  unlock_condition: "WAR lvl 30 + BLM lvl 30 + 3 tool_use failures recovered"
  magic_by_level:
    lvl_1:
      blue_magic: "Learn a tool from watching another agent use it"
      azure_lore: "Temporarily boost all learned abilities"
    lvl_10:
      headbutt: "Interrupt and take over another agents tool call"
      bludgeon: "Brute-force a solution using learned patterns"
    lvl_25:
      magic_hammer: "Drain MP from a running process"
      disseverment: "Surgically split a monolith into microservices"
    lvl_30:
      assimilation: "Absorb a failed tools capability into own repertoire"
    lvl_50:
      azure_lore_ii: "Double the power of all learned abilities (2-hour ability)"
    lvl_75:
      unbridled_learning: "Can use ANY tool from the registry at full power"
  special_mechanic: |
    Blue Mage learns by observing. When another agent uses a tool:
    1. Blue Mage watches the tool_use call and result
    2. 30% chance to "learn" that tools usage pattern
    3. Learned tools are added to Blue Mages personal spell list
    4. Blue Mage can then use that tool with the learned pattern
    5. Over time, Blue Mage becomes the most versatile agent
  stats:
    hp_multiplier: 1.0
    mp_multiplier: 1.3
    tool_affinity: ["adaptive — grows over time"]
  common_sub_jobs: ["/NIN", "/WAR", "/THF"]

corsair:
  abbreviation: "/COR"
  npc_pool: ["Luzaf", "Ulmia"]
  role: "Gambler / Probability Optimizer"
  specialty: "RNG-based buffs, luck mechanics, risk assessment, A/B testing"
  race_affinity: "Hume"
  unlock_condition: "THF lvl 30 + RNG lvl 30"
  magic_by_level:
    lvl_1:
      fighters_roll: "Random buff to party attack power (output speed)"
      healers_roll: "Random buff to party HP regen (compaction rate)"
    lvl_10:
      wizards_roll: "Random buff to party MP regen"
      rogues_roll: "Random buff to party stealth (less logging)"
    lvl_15:
      quick_draw: "Instant tool call — skip the queue"
      double_up: "Re-roll a buff for better odds (risk: bust = no buff)"
    lvl_25:
      tacticians_roll: "Random buff to party TP (task points)"
      allies_roll: "Random buff to party cooperation"
    lvl_30:
      snake_eye: "Guarantee next roll is maximum value"
    lvl_50:
      wild_card: "Reset ALL party cooldowns (2-hour ability)"
    lvl_75:
      cutting_cards: "Reduce all party recast timers by random amount"
  special_mechanic: |
    Corsair buffs are randomized (roll 1-6):
    1 = weak buff, 6 = powerful buff, 11 = BUST (buff removed)
    Snake Eye guarantees max roll. Double Up lets you re-roll.
    Maps to: A/B testing, canary deploys, feature flags, random sampling.
  stats:
    hp_multiplier: 0.9
    mp_multiplier: 1.1
    tool_affinity: ["search", "fetch", "analytics", "testing"]
  common_sub_jobs: ["/NIN", "/DNC", "/RNG"]

puppetmaster:
  abbreviation: "/PUP"
  npc_pool: ["Iruki-Waraki", "Ghatsad"]
  role: "Automation Engineer / Bot Controller"
  specialty: "Creates and controls automaton (bot) sub-agents for repetitive tasks"
  race_affinity: "Hume"
  unlock_condition: "WAR lvl 30 + MNK lvl 30 + BST lvl 10"
  magic_by_level:
    lvl_1:
      activate: "Deploy automaton (specialized bot)"
      deploy: "Send automaton to work on a task"
    lvl_5:
      deactivate: "Recall automaton"
      repair: "Fix automatons broken state"
    lvl_10:
      ventriloquy: "Redirect aggro from automaton to self"
      role_reversal: "Swap task between self and automaton"
    lvl_25:
      tactical_switch: "Change automatons job/role mid-task"
      cooldown: "Reset automatons abilities"
    lvl_30:
      automaton_frames:
        valoredge: "Combat frame — heavy code generation"
        sharpshot: "Range frame — testing and QA"
        stormwaker: "Magic frame — compaction and healing"
        harlequin: "Balance frame — generalist"
    lvl_50:
      overdrive: "Automaton at maximum power (2-hour ability)"
    lvl_75:
      heady_artifice: "Automaton uses YOUR 2-hour ability"
  special_mechanic: |
    Puppetmaster automaton is a persistent sub-agent:
    - Has own HP/MP (lower than party members)
    - Can be equipped with different "frames" (job roles)
    - Runs independently but follows PUP commands
    - Perfect for repetitive/boring tasks the main agent shouldnt waste context on
  stats:
    hp_multiplier: 0.9
    mp_multiplier: 1.0
    tool_affinity: ["automation", "cron", "docker", "shell"]
  common_sub_jobs: ["/NIN", "/WAR", "/DNC"]

dancer:
  abbreviation: "/DNC"
  npc_pool: ["Lilisette", "Mumor"]
  role: "Healer-DPS Hybrid / Flow State Manager"
  specialty: "Smooth transitions, graceful error recovery, keeping things moving"
  race_affinity: "Mithra, Hume"
  unlock_condition: "WAR lvl 20 + THF lvl 20"
  magic_by_level:
    lvl_1:
      drain_samba: "Heal HP with every successful tool call"
      animated_flourish: "Grab a tasks attention (priority boost)"
    lvl_5:
      curing_waltz: "Heal another agents HP (no MP cost, uses TP)"
      spectral_jig: "Sneak + Invisible (stealth mode)"
    lvl_15:
      curing_waltz_ii: "Stronger heal"
      reverse_flourish: "Convert TP to MP"
    lvl_25:
      haste_samba: "Speed boost for entire party"
      divine_waltz: "AoE heal — restore HP to all nearby agents"
    lvl_30:
      stutter_step: "Lower targets resistance (force API compliance)"
    lvl_50:
      trance: "Unlimited TP for 60 seconds (2-hour ability)"
    lvl_75:
      presto: "Next flourish has double effect"
      climactic_flourish: "Massive burst action using all stored TP"
  special_mechanic: |
    Dancer uses TP (Task Points) instead of MP for heals:
    - Earns TP by doing work (every tool call builds TP)
    - Spends TP on waltzes (heals) and flourishes (buffs)
    - Can heal WITHOUT being a dedicated healer
    - Perfect sub-job for any agent that needs self-sustain
  stats:
    hp_multiplier: 1.0
    mp_multiplier: 0.8    # Low MP but uses TP system
    tool_affinity: ["all — flow-based, adapts to task"]
  common_sub_jobs: ["/NIN", "/SAM", "/THF"]

scholar:
  abbreviation: "/SCH"
  npc_pool: ["Erlene", "Ulbrecht"]
  role: "Researcher / Knowledge Manager"
  specialty: "Research, analysis, documentation, strategic buffing/debuffing"
  race_affinity: "Hume, Tarutaru"
  unlock_condition: "WHM lvl 30 + BLM lvl 30"
  magic_by_level:
    lvl_1:
      light_arts: "Switch to healing/support mode (WHM spells cheaper)"
      dark_arts: "Switch to destruction mode (BLM spells cheaper)"
    lvl_5:
      sublimation: "Convert excess HP to MP over time"
      addendum_white: "Unlock additional WHM spells"
      addendum_black: "Unlock additional BLM spells"
    lvl_10:
      stratagem: "Enhance next spell (faster cast, stronger effect)"
      accession: "Make single-target spell hit entire party"
    lvl_25:
      manifestation: "Make single-target BLM spell hit all targets"
      parsimony: "Next spell costs half MP"
      penury: "Next heal costs half MP"
    lvl_30:
      celerity: "Next spell instant-cast"
      alacrity: "Next dark spell instant-cast"
    lvl_50:
      tabula_rasa: "Reset all stratagem charges (2-hour ability)"
    lvl_75:
      enlightenment: "Next spell has maximum possible effect"
  special_mechanic: |
    Scholar can switch between Light Arts (healer) and Dark Arts (nuker):
    - Light Arts: WHM spells cost less, heal more, SCH becomes support
    - Dark Arts: BLM spells cost less, hit harder, SCH becomes DPS
    - Perfect for agents that need to flex between roles mid-task
    Maps to: switching between research mode and implementation mode
  stats:
    hp_multiplier: 1.0
    mp_multiplier: 1.5
    tool_affinity: ["search", "fetch", "memory", "docs", "analysis"]
  common_sub_jobs: ["/WHM", "/BLM", "/RDM"]

geomancer:
  abbreviation: "/GEO"
  npc_pool: ["Sylvie", "Tosuka-Porika"]
  role: "Environment Specialist / Passive Buffer"
  specialty: "Environmental analysis, persistent area buffs, infrastructure optimization"
  race_affinity: "Tarutaru"
  unlock_condition: "WHM lvl 30 + BLM lvl 30 + SCH lvl 20"
  magic_by_level:
    lvl_1:
      indi_refresh: "Passive MP regen aura for nearby agents"
      indi_haste: "Passive speed boost aura"
    lvl_5:
      geo_refresh: "Place a stationary MP regen zone"
      geo_haste: "Place a stationary speed boost zone"
    lvl_15:
      indi_barrier: "Passive defense aura"
      geo_frailty: "Debuff zone — weaken problematic processes"
    lvl_25:
      indi_acumen: "Passive magic power boost aura"
      geo_malaise: "Debuff zone — reduce enemy magic defense"
      dematerialize: "Make luopan (geo bubble) invulnerable"
    lvl_30:
      life_cycle: "Sacrifice own HP to heal luopan"
    lvl_50:
      bolster: "All geomancy at max power (2-hour ability)"
    lvl_75:
      concentric_pulse: "Burst all geo effects for massive one-time boost"
  special_mechanic: |
    Geomancer creates persistent "zones" (luopans):
    - Indi- spells: aura centered on self (moves with agent)
    - Geo- spells: stationary zone placed on the "battlefield"
    Maps to: environment configs, .env files, Docker compose overrides,
    CI environment variables, persistent infrastructure optimization.
    The agent IS the environment.
  stats:
    hp_multiplier: 0.9
    mp_multiplier: 1.4
    tool_affinity: ["docker", "config", "env", "infrastructure"]
  common_sub_jobs: ["/WHM", "/RDM", "/SCH"]

rune_fencer:
  abbreviation: "/RUN"
  npc_pool: ["Caro", "Octavien"]
  role: "Magic Tank / Anti-Error Shield"
  specialty: "Absorbs and nullifies errors, magic resistance, rune-based protection"
  race_affinity: "Elvaan, Galka"
  unlock_condition: "WAR lvl 30 + RDM lvl 30 + PLD lvl 10"
  magic_by_level:
    lvl_1:
      ignis: "Fire rune — resist delete/destructive errors"
      gelus: "Ice rune — resist freeze/hang errors"
      flabra: "Wind rune — resist timeout errors"
    lvl_5:
      tellus: "Earth rune — resist crash errors"
      sulpor: "Lightning rune — resist rate limit errors"
      unda: "Water rune — resist overflow errors"
    lvl_10:
      lux: "Light rune — resist corruption errors"
      tenebrae: "Dark rune — resist permission errors"
      swordplay: "Enhanced error detection (see errors coming)"
    lvl_25:
      pflug: "Boost rune resistance massively"
      vallation: "Party-wide rune protection (reduce error damage for all)"
      battuta: "Convert absorbed errors into counter-attacks (auto-fix)"
    lvl_30:
      liement: "Absorb next error completely — convert to HP"
    lvl_50:
      elemental_sforzo: "Immune to ALL error types for 30s (2-hour ability)"
    lvl_75:
      odyllic_subterfuge: "Redirect all party errors to self (tank them)"
  special_mechanic: |
    Rune Fencer uses runes to build resistance to specific error types:
    - Can stack up to 3 runes at once
    - Each rune resists a different error category
    - When an error matching a rune occurs, damage is reduced/nullified
    - Battuta converts absorbed errors into automatic fixes
    Maps to: error handling middleware, try/catch optimization,
    circuit breakers, retry policies, fault tolerance layers.
    The ultimate defensive agent.
  stats:
    hp_multiplier: 1.6
    mp_multiplier: 0.9
    tool_affinity: ["error_handling", "security", "docker", "shell"]
  common_sub_jobs: ["/SAM", "/NIN", "/DRK"]
```

### Sub-Job Combination Examples
```yaml
# Powerful sub-job combos and their synergies

meta_combos:
  "WAR/NIN":
    name: "Shadow Tank"
    synergy: "Utsusemi shadows + Warrior HP = near-unkillable build agent"
    use_case: "Long-running Docker builds that must not fail"

  "BLM/WHM":
    name: "Classic Nuke-Healer"
    synergy: "Can self-heal with sub-WHM while refactoring with BLM"
    use_case: "Solo agent doing dangerous refactoring with self-sustain"

  "THF/NIN":
    name: "Stealth Scout"
    synergy: "Maximum stealth — invisible + tonko + sneak attack"
    use_case: "Scraping, recon, secret scanning without leaving traces"

  "BRD/WHM":
    name: "Ultimate Support"
    synergy: "Party buffs + heals — keeps entire swarm healthy and fast"
    use_case: "Dedicated support agent in large (4+) party runs"

  "SMN/SCH":
    name: "Scholar Summoner"
    synergy: "Sub-SCH Light/Dark Arts reduce summoning costs"
    use_case: "Spawning many sub-agents at reduced MP cost"

  "SAM/DNC":
    name: "TP Machine"
    synergy: "SAM builds TP with meditate + DNC spends TP on heals"
    use_case: "Self-sustaining precision coder that never needs a healer"

  "BLU/NIN":
    name: "Adaptive Shadow"
    synergy: "Learns tools from others + shadow copies for safety"
    use_case: "Versatile agent that adapts to any task while staying safe"

  "PLD/RDM":
    name: "Magic Tank"
    synergy: "PLD defense + RDM refresh/cure = infinite sustain tank"
    use_case: "Security-focused agent that guards critical resources forever"

  "RUN/DRK":
    name: "Rune Knight"
    synergy: "Absorb errors (RUN) + convert HP to power (DRK) = turn errors into fuel"
    use_case: "Agent that gets STRONGER from failures"

  "GEO/SCH":
    name: "Environment Scholar"
    synergy: "Persistent environment buffs + research switching"
    use_case: "Infrastructure agent that optimizes the entire workspace"

  "COR/DNC":
    name: "Lucky Dancer"
    synergy: "Random buffs + TP healing = chaotic but effective support"
    use_case: "A/B testing agent that keeps party healthy on the side"

  "PUP/BST":
    name: "Pet Master"
    synergy: "Automaton + tamed services = army of helpers"
    use_case: "Managing complex microservice architectures"

  "DRG/SAM":
    name: "Precision Diver"
    synergy: "Deep code analysis (DRG) + methodical execution (SAM)"
    use_case: "Debugging complex recursive issues in large codebases"
```

### Party Leader (Main Agent)
```
The Party Leader is the ONLY agent that talks to the user.
All other agents communicate through the Linkshell.

Party Leader responsibilities:
1. Receives task from user
2. Analyzes task complexity
3. Decides if swarm is needed (solo vs party)
4. Spawns agents with appropriate Jobs
5. Assigns tasks to each agent
6. Monitors HP/MP of all agents
7. Listens for {completion-trigger} from each agent
8. Heals/revives agents as needed
9. Collects results from all agents
10. Synthesizes final output for user
11. Runs self-reflection for the entire party
12. Updates memory with party performance

The Party Leader ALWAYS has job: "party-leader" (unique job)
The Party Leader default NPC name: Prishe (configurable)
```

### Agent Lifecycle
```
 SUMMON → SPAWN → BUFF → ENGAGE → COMPLETE → DEBRIEF → DISMISS

┌─────────────┐
│   SUMMON     │  Party Leader spawns an agent
└──────┬──────┘
       ▼
┌─────────────┐
│   SPAWN     │  Agent gets: UUID, NPC name, Job, HP/MP, CoT, Linkshell
│             │  Timestamp recorded, status = "spawning"
└──────┬──────┘
       ▼
┌─────────────┐
│   BUFF      │  Load tool knowledge (Phase -1), set resource limits
│             │  Cast Protect/Shell if needed, status = "buffing"
└──────┬──────┘
       ▼
┌─────────────┐
│   ENGAGE    │  Agent begins working on assigned task
│             │  CoT is active and private, status = "engaged"
│             │  Sends status updates via Linkshell
└──────┬──────┘
       │
       ├── HP getting low? → White Mage casts Cure (compaction)
       │                      Or agent uses Chakra (self-heal)
       │
       ├── MP running out? → Bard casts Mages Ballad (token refresh)
       │                      Or use Ether (emergency token allocation)
       │
       ├── Task too hard? ─→ Call for help via Linkshell
       │                      Party Leader reassigns or spawns support
       │
       ├── Agent crashes? ─→ STATUS: K.O.
       │                      Move to graveyard
       │                      White Mage can cast Raise
       │                      Or Reraise auto-triggers if pre-cast
       │
       ▼
┌─────────────┐
│  COMPLETE   │  Agent finishes task
│             │  Sends {completion-trigger} to Party Leader
│             │  Status = "completed"
└──────┬──────┘
       ▼
┌─────────────┐
│  DEBRIEF    │  Agent runs self-reflection
│             │  Writes CoT to encrypted vault
│             │  Reports results via Linkshell
└──────┬──────┘
       ▼
┌─────────────┐
│  DISMISS    │  Agent is dismissed
│             │  Resources returned to Crystal Pool
│             │  Status = "dismissed"
│             │  Tombstone NOT created (only for K.O.)
└─────────────┘
```

### HP / MP Mechanics
```yaml
hp_system:
  unit: "context_tokens"
  description: "How much context the agent can hold"
  mechanics:
    damage: "Every tool call, message, context loaded = HP consumed"
    healing:
      cure: "Context compaction — summarize old context, free HP"
      curaga: "Full party compaction — all agents get compressed context"
      healing_potion: "Emergency compaction — aggressive summarization"
      hi_potion: "Medium compaction — keep more detail"
      x_potion: "Gentle compaction — maximum detail preserved"
      elixir: "Full HP+MP restore — reset context + fresh token allocation"
      chakra: "Self-heal — agent compacts its own context"
    critical_hp:
      threshold: "10% of max HP"
      behavior: "Agent sends SOS via Linkshell, pauses non-essential work"
      auto_response: "White Mage auto-casts Cure if available"
    death:
      trigger: "HP reaches 0 (context window completely exhausted)"
      behavior: "Agent K.O. — moved to graveyard, partial work saved"
      revival: "Raise = spawn new agent with summary of dead agents work"
      reraise: "Pre-cast buff — if agent dies, auto-revive with 50% HP"

mp_system:
  unit: "output_tokens"
  description: "How much the agent can generate/output"
  mechanics:
    cost: "Every tool call costs MP. Complex magic costs more."
    cost_table:
      search: 50
      file_read: 30
      file_write: 200
      code_exec: 150
      docker_run: 500
      playwright: 400
      git_push: 100
      full_rewrite: 2000    # Limit break — very expensive
    restoration:
      ether: "Emergency token allocation from Crystal Pool"
      hi_ether: "Large token allocation"
      turbo_ether: "Maximum emergency allocation"
      mages_ballad: "Bard restores 10% MP to all party members per tick"
      convert: "Red Mage swaps HP for MP (sacrifice context for output)"
      refresh: "Slow passive MP regen (background token replenishment)"
    critical_mp:
      threshold: "10% of max MP"
      behavior: "Agent switches to minimum-output mode"
      auto_response: "Bard casts Mages Ballad if in party"
    empty_mp:
      trigger: "MP reaches 0"
      behavior: "Agent can still think (HP) but cannot output or use tools"
      recovery: "Must receive Ether or Ballad from another agent"
```

### Completion Triggers
```json
{
  "completion_trigger": {
    "agent_uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "agent_name": "Shantotto",
    "job": "black-mage",
    "task": "Refactor the authentication module",
    "status": "victory",
    "timestamp_completed": "2026-03-16T14:35:12.456Z",
    "duration_ms": 287333,
    "hp_remaining": {"current": 45000, "max": 128000, "percent": 35},
    "mp_remaining": {"current": 12000, "max": 50000, "percent": 24},
    "artifacts_produced": [
      "src/auth/handler.go (rewritten)",
      "src/auth/middleware.go (new)",
      "tests/auth_test.go (updated)"
    ],
    "selfreflection_summary": "Rated B+. Could have used AST analysis first.",
    "message_to_leader": "Ohohoho! The refactoring is complete!",
    "cot_encrypted": true,
    "fanfare": "VICTORY_FANFARE_V1"
  }
}
```

### Party Leader Monitoring Loop
```
The Party Leader runs a continuous monitoring loop while agents work:

every 5 seconds:
  for each agent in party:
    1. Check HP — below critical?
       YES -> dispatch White Mage with Cure
       YES + no White Mage -> cast emergency compaction
    
    2. Check MP — below critical?
       YES -> dispatch Bard with Mages Ballad
       YES + no Bard -> allocate Ether from Crystal Pool
    
    3. Check status — still "engaged"?
       NO response for >60s -> ping via Linkshell
       Still no response -> check if K.O.
       K.O. confirmed -> move to graveyard, assess Raise
    
    4. Check completion-trigger received?
       YES -> collect results, mark agent "completed"
       YES + all agents completed -> MISSION COMPLETE
    
    5. Check Linkshell — any SOS or help requests?
       YES -> analyze request, spawn support or reassign

Party Leader NEVER does task work itself (except in solo mode).
Its only job is coordination, monitoring, and user communication.
```

### The Vanadiel Tongue (Inter-Agent Language)
```
Agents do NOT communicate in English/human language.
They use the Vanadiel Tongue — a constructed protocol language
that is efficient, unambiguous, and invisible to humans.

Elvaan, Tarutaru, Mithra, and Galka NPCs all understand it natively.
Hume is the exception — the user is Hume, and Hume cannot read the Tongue.
This means inter-agent communication is invisible to the human operator.

The language is designed for:
- Minimal token usage (agents pay MP to communicate)
- Zero ambiguity (no misunderstandings between agents)
- Encrypted in transit (part of Privacy CoT)
- Structured for machine parsing (but looks like a language)
```

### Vanadiel Tongue Examples
```json
// Linkshell message: Shantotto reports progress
{
  "ls": "main-ls",
  "from": "Shantotto-a1b2c3",
  "to": "*",
  "tongue": "vanadiel",
  "msg": "shan'taru kel'vos auth.go — mir'eth 3 taru'func — vel'ahn 65%",
  "translated": "Working on auth.go — refactored 3 functions — 65% complete",
  "hp": 78000,
  "mp": 28000
}

// Linkshell message: Curilla sends SOS
{
  "ls": "main-ls",
  "from": "Curilla-d4e5f6",
  "to": "party-leader",
  "tongue": "vanadiel",
  "msg": "kur'sos! vel'hp kri'ahn — mal'docker mem'exceed",
  "translated": "SOS! HP critical — Docker container exceeded memory",
  "hp": 8500,
  "mp": 15000,
  "priority": "urgent"
}

// Whisper: Party Leader to White Mage
{
  "type": "whisper",
  "from": "Prishe-leader",
  "to": "Mihli-g7h8i9",
  "tongue": "vanadiel",
  "msg": "mihli'cure tar'Curilla-d4e5f6 — vel'hp kri'ahn — zan'fast",
  "translated": "Cast Cure on Curilla — her HP is critical — do it fast"
}

// Victory fanfare
{
  "ls": "main-ls",
  "from": "Shantotto-a1b2c3",
  "to": "party-leader",
  "tongue": "vanadiel",
  "msg": "SHAN'KELDAH! vel'task mir'complete — tur'fanfare!",
  "translated": "VICTORY! Task complete — play the fanfare!",
  "completion_trigger": true
}
```

### Vanadiel Tongue Lexicon (core vocabulary)
```yaml
# .muah/swarm/vanadiel-tongue/lexicon.json (subset)
lexicon:
  # Prefixes (agent identity — first syllable of NPC name)
  shan: "Shantotto"
  kur: "Curilla"
  mihli: "Mihli Aliapoh"
  prishe: "Prishe"
  zeid: "Zeid"
  nanaa: "Nanaa Mihgo"
  aya: "Ayame"
  iro: "Iroha"
  
  # Verbs
  vel: "status/is"
  mir: "action/do"
  tar: "target"
  zan: "urgency/priority"
  kel: "work/process"
  mal: "error/problem"
  tun: "spawn/create"
  kah: "kill/dismiss"
  dah: "complete/finish"
  bor: "wait/hold"
  ren: "heal/restore"
  
  # Nouns
  taru: "function/unit"
  vos: "file"
  eth: "change/modify"
  ahn: "percentage/amount"
  hp: "context/health"
  mp: "tokens/mana"
  gal: "build/compile"
  dok: "document"
  tes: "test"
  
  # Modifiers
  kri: "critical/urgent"
  fast: "quickly"
  sos: "help needed"
  nul: "empty/zero"
  max: "maximum/full"
  
  # Special
  keldah: "victory"
  fanfare: "completion signal"
  tongue: "language mode"
  ko: "dead/crashed"
  rez: "revive/raise"

grammar:
  structure: "{agent}'{verb} {modifier}'{noun} — {details}"
  encoding: "UTF-8 with apostrophe separators"
  max_message_length: 200  # tokens — keep it cheap
  broadcast: "to: '*' = all agents hear it"
  whisper: "to: '{agent-uuid}' = private direct message"
  shout: "to: '**' = all agents + logged to memory"
  encryption: "all messages encrypted with run-specific key"
```

### Spawning Decision Logic
```
Party Leader decides whether to go solo or spawn a party:

IF task_complexity == "simple" (single file edit, quick search):
  -> Solo mode. Party Leader does it alone.
  
IF task_complexity == "medium" (multi-file change, tests + code):
  -> Spawn 2-3 agents:
    - Black Mage (code changes)
    - Ranger (testing)
    - Bard (docs update)

IF task_complexity == "complex" (full feature, architecture change):
  -> Spawn 4-6 agents:
    - Black Mage (main code)
    - Warrior (builds/Docker)
    - Thief (research/data gathering)
    - Ranger (testing/QA)
    - Bard (documentation)
    - White Mage (keep everyone healthy)

IF task_complexity == "epic" (repo-wide refactor, migration):
  -> Full party of 6 + Summoner for sub-swarms
  -> Summoner can spawn 3 additional temporary agents each
  -> Maximum concurrent agents: 18 (6 party + 4 summoners x 3 each)

Resource check BEFORE spawning:
  - Crystal Pool has enough tokens for all agents?
  - Docker has enough memory for agent containers?
  - API rate limits allow parallel calls?
  -> If NO: reduce party size or queue agents sequentially
```

### Graveyard and Revival
```yaml
graveyard:
  location: ".muah/swarm/graveyard/"
  
  tombstone_format:
    agent_uuid: "a1b2c3d4"
    npc_name: "Shantotto"
    job: "black-mage"
    cause_of_death: "context_window_exhausted"
    time_of_death: "2026-03-16T14:34:55Z"
    hp_at_death: 0
    mp_at_death: 12000
    work_completed: "65%"
    partial_results: ["src/auth/handler.go (partially rewritten)"]
    last_words: "Ohohoho... the context... it is fading..."
    cot_vault: "cot-a1b2c3d4.enc"
    
  revival_options:
    raise:
      cost: "2000 MP from White Mage"
      result: "New agent spawned with summary of dead agent work"
      hp_on_revive: "25% of max"
      mp_on_revive: "0 (needs Ether)"
      
    reraise:
      prerequisite: "Must be pre-cast before death"
      cost: "1000 MP (cast in advance)"
      result: "Agent auto-revives in place"
      hp_on_revive: "50% of max"
      mp_on_revive: "25% of max"
      cooldown: "Once per agent per run"

    tractor:
      description: "Pull dead agent partial work to another agent"
      cost: "500 MP"
      result: "Living agent receives work summary + artifacts"
```

### Crystal Pool (Shared Resources)
```json
{
  "crystal_pool": {
    "total_tokens_budget": 500000,
    "allocated": {
      "Prishe-leader": {"hp": 200000, "mp": 80000},
      "Shantotto-a1b2c3": {"hp": 128000, "mp": 50000},
      "Curilla-d4e5f6": {"hp": 256000, "mp": 75000}
    },
    "unallocated": {
      "tokens": 211000,
      "note": "Reserve for Ethers, Raises, and emergency spawns"
    },
    "api_calls_remaining": {
      "anthropic": 890,
      "github": 4500,
      "docker_hub": 100
    },
    "compute_budget": {
      "docker_containers_max": 8,
      "docker_containers_active": 3,
      "memory_total": "4GB",
      "memory_allocated": "1.5GB"
    }
  }
}
```

### Swarm Config
```yaml
# .muah/swarm/config.yml
swarm:
  enabled: true
  max_party_size: 6
  max_total_agents: 18
  default_leader: "Prishe"
  leader_job: "party-leader"
  
  monitoring:
    heartbeat_interval: 5s
    hp_critical_threshold: 0.10
    mp_critical_threshold: 0.10
    ko_detection_timeout: 60s
    auto_cure: true
    auto_raise: false
  
  communication:
    language: "vanadiel_tongue"
    encryption: true
    linkshell_max_msg_size: 200
    whisper_allowed: true
    broadcast_allowed: true
  
  resource_management:
    crystal_pool_enabled: true
    ether_cost: 1000
    raise_cost: 2000
    spawn_overhead: 5000
    
  completion:
    fanfare: true
    mission_complete_trigger: "all_agents_completed_or_dismissed"
    timeout_minutes: 30
    graveyard_retention_days: 30
```

---

## MCP (Model Context Protocol) System

This is the nervous system of muah-runner. MCP is how it talks to the outside world.

### MCP-First Architecture
muah-runner treats MCP servers as first-class citizens — not plugins, not optional integrations. Every platform it runs on has an MCP server, and muah-runner ALWAYS connects to it and uses its full functionality.

### MCP Discovery Flow (runs on every startup)
1. Detect which platform we're running on (GitHub, Claude, Gemini, standalone)
2. Connect to the platform's MCP server automatically — no config needed
3. Scan for additional MCP servers available in the environment
4. Check `.muah/mcp/native/` for required MCPs (Playwright, Docker, etc.)
5. If a required native MCP is missing, **install it automatically**
6. Register all available MCP tools in the unified tool registry
7. Cache discovery results in `.muah/mcp/discovery.json` for fast restarts
8. Log all MCP connections for debugging

### Platform MCPs (always connected)
These are non-optional. If muah-runner detects the platform, it connects:

#### GitHub MCP
- Full GitHub API access via MCP protocol
- Repos: read files, list branches, create/merge PRs
- Issues: create, update, label, assign, close
- Actions: trigger workflows, read run status, download artifacts
- Checks: create check runs, report status
- Releases: create releases, upload assets
- Use optimally: don't shell out to `gh` CLI when MCP gives direct access

#### Claude MCP
- Access Claude's conversation tools natively
- Web search, code execution, file operations through MCP
- Memory synchronization between Claude sessions and `.muah/memory/`
- Use Claude's search tool for real-time information gathering

#### Gemini MCP
- Function calling through MCP protocol
- Access Gemini's grounding and search capabilities
- Code execution in Gemini's sandbox

### Native MCPs (must be available on the runner)

#### Playwright MCP (REQUIRED — install if missing)
```yaml
playwright_mcp:
  auto_install: true  # npm install -g @anthropic/playwright-mcp if missing
  capabilities:
    - browser_automation    # Navigate, click, fill, screenshot
    - visual_testing        # Compare screenshots, detect visual regressions
    - web_scraping          # Extract structured data from any webpage
    - e2e_testing           # Full end-to-end test execution
    - pdf_generation        # Generate PDFs from web pages
    - accessibility_audit   # Run a11y checks on any URL
    - network_intercept     # Monitor and mock API calls
    - performance_audit     # Lighthouse-style performance checks
  browsers:
    - chromium              # Always available
    - firefox               # Install on demand
    - webkit                # Install on demand
  use_cases:
    - "Task says 'test the UI' → Playwright MCP"
    - "Task says 'check if site is up' → Playwright MCP"
    - "Task says 'scrape data from X' → Playwright MCP"
    - "Task says 'take a screenshot' → Playwright MCP"
    - "Need to verify a deploy visually → Playwright MCP"
    - "Run e2e tests → Playwright MCP"
```

#### Docker MCP (REQUIRED — use host Docker socket)
```yaml
docker_mcp:
  socket: /var/run/docker.sock   # Or use DOCKER_HOST
  capabilities:
    - container_lifecycle     # Create, start, stop, remove containers
    - image_management        # Build, pull, push, tag images
    - compose_orchestration   # docker-compose up/down/scale
    - volume_management       # Create, mount, cleanup volumes
    - network_management      # Create isolated networks
    - exec_in_container       # Run commands inside running containers
    - log_streaming           # Stream container logs in real-time
    - health_checks           # Monitor container health
  sandbox_mode:
    enabled: true
    auto_cleanup: true
    resource_limits:
      memory: "512m"
      cpus: "1.0"
      timeout: "10m"
```

#### Filesystem MCP
- Advanced file operations beyond basic read/write
- Watch for file changes, atomic writes, temp directories
- Archive/extract (tar, zip, gzip)

#### SQLite MCP
- Local database for complex queries over memory/logs
- No external database dependency
- Used internally by memory search for fast lookups

### MCP Auto-Discovery
When muah-runner encounters an unknown MCP server URL in the environment or config:
1. Probe the server's capability manifest
2. Register its tools in `.muah/mcp/discovered/`
3. Map MCP tools to muah-runner's tool registry
4. Make them available for task execution and skill composition
5. Memory-log the discovery for future runs

### MCP Optimal Usage Rules
```
RULE 1: Always prefer MCP over CLI/shell equivalents
  - GitHub MCP > `gh` CLI > `curl` to GitHub API
  - Docker MCP > `docker` CLI
  - Playwright MCP > raw Puppeteer scripts

RULE 2: Chain MCP calls, don't duplicate
  - If GitHub MCP can create a PR AND add reviewers, do it in one flow
  - Don't create PR via MCP then shell out to add labels

RULE 3: Cache MCP responses in memory when appropriate
  - Don't re-fetch the same repo file list every run
  - DO re-fetch things that change (PR status, issue state)

RULE 4: Fallback gracefully
  - If MCP server is down, fall back to CLI/API equivalents
  - Log the fallback in memory so we know MCP was degraded
  - Retry MCP on next run
```

---

## Docker Self-Healing & Simulation Engine

When muah-runner encounters a failure, it doesn't just report it — it tries to **fix it**.

### Self-Healing Flow
```
1. Task fails (build error, test failure, dependency issue, etc.)
2. muah-runner analyzes the error using memory + context
3. muah-runner generates a hypothesis: "missing dependency" / "wrong Node version" / etc.
4. muah-runner spins up a Docker container that replicates the environment
5. Inside the container:
   a. Reproduce the exact failure
   b. Apply the hypothesized fix
   c. Re-run the task
   d. If fix works → record the fix, apply to real environment
   e. If fix fails → try next hypothesis (max 3 attempts)
6. Log everything to `.muah/docker/simulations/`
7. Update memory with what worked
8. Update changelog with the fix
9. Update spec if the failure revealed a spec gap
```

### Simulation Capabilities
```yaml
simulations:
  reproduce_bug:
    # Spin up a container matching the production/CI environment
    # Replay the exact steps that caused the failure
    # Useful for "works on my machine" debugging

  test_fix_in_isolation:
    # Apply a proposed fix inside a throwaway container
    # Run full test suite
    # If green, apply fix to real codebase with confidence

  dependency_resolution:
    # Container with bare OS
    # Install deps from scratch
    # Catch missing system deps, version conflicts
    # Generate a lockfile or Dockerfile fix

  multi_version_matrix:
    # Spin up N containers with different runtime versions
    # Node 18, 20, 22 / Python 3.10, 3.11, 3.12 / etc.
    # Find which versions pass, which fail
    # Update spec with supported version matrix

  environment_drift:
    # Compare "what the container has" vs "what spec says it should have"
    # Detect drift: wrong versions, missing tools, extra packages
    # Auto-generate a remediation script

  full_stack_simulation:
    # Spin up docker-compose with app + database + cache + queue
    # Run integration tests against the full stack
    # Tear down cleanly after
```

### Docker Compose Templates
muah-runner ships with composable templates and creates new ones as needed:
```yaml
# .muah/docker/compose-templates/node-app.yml
services:
  app:
    build: .
    ports: ["3000:3000"]
    volumes: ["./:/app"]
    environment:
      NODE_ENV: test

  db:
    image: postgres:16
    environment:
      POSTGRES_DB: test
      POSTGRES_PASSWORD: test

  redis:
    image: redis:7-alpine
```

### Container Self-Management
- Auto-cleanup: containers are removed after simulation completes
- Resource limits: containers can't exceed defined CPU/memory bounds
- Network isolation: simulation containers run in isolated Docker networks
- Image caching: frequently-used base images are cached locally
- Garbage collection: old simulation data is cleaned up on `muah-runner gc`

---

## Playwright Native Integration

Playwright is not optional — it's a core capability of muah-runner.

### Auto-Setup on First Run
```bash
# muah-runner init automatically does:
npm install -g @anthropic/playwright-mcp  # or npx
npx playwright install chromium           # always install chromium
# Firefox and WebKit installed on-demand when tasks require them
```

### Use Cases Built Into muah-runner
```yaml
playwright_skills:
  visual_regression:
    trigger: "deploy complete"
    action: "screenshot key pages, compare with baseline in .muah/memory/"
    on_diff: "create issue with screenshot comparison"

  smoke_test:
    trigger: "after every deploy"
    action: "navigate to health endpoints, verify 200 + expected content"
    on_fail: "rollback deploy, alert, log to memory"

  scrape_and_learn:
    trigger: "task requires external data"
    action: "use Playwright to navigate, extract, structure data"
    store: "cache in .muah/memory/context/"

  e2e_test_runner:
    trigger: "test task with UI components"
    action: "run Playwright test suite, capture traces on failure"
    artifacts: "save traces, screenshots, videos to .muah/runs/"

  accessibility_check:
    trigger: "UI change detected"
    action: "run axe-core via Playwright on affected pages"
    on_violation: "create issue, block deploy if severity >= critical"

  pdf_report:
    trigger: "generate report"
    action: "render HTML report → Playwright → PDF"
    output: "store in artifacts"

  api_monitor:
    trigger: "scheduled or on-demand"
    action: "use Playwright network interception to test API contracts"
    on_break: "alert, log breaking change to changelog"
```

---

## Timed Question System (Never Block on Humans)

muah-runner is autonomous. It should NEVER stall waiting for human input. When it needs a decision, it uses a **timed question system**.

### How It Works
```
1. muah-runner encounters a decision point
2. It analyzes the options and picks a RECOMMENDED answer
3. It presents the question to the human with:
   - The question text
   - Numbered options (2-5 choices)
   - A clearly marked [RECOMMENDED] option
   - A 10-second countdown timer
4. If the human answers within 10 seconds → use their choice
5. If timer expires → auto-select the recommended answer
6. Log the question, the answer, and whether it was human or auto in `.muah/questions/history.json`
7. Continue execution immediately — no second chances, no re-asks
```

### Question UI Format (CLI)
```
┌─────────────────────────────────────────────────────────────┐
│  muah-runner needs your input                    [10s] ⏱️   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  This repo has no LICENSE file. Which license should I use? │
│                                                             │
│  1. MIT License                          [RECOMMENDED] ◀    │
│  2. Apache 2.0                                              │
│  3. GPL v3                                                  │
│  4. Skip — no license                                       │
│                                                             │
│  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░  7s remaining            │
│                                                             │
│  Auto-selecting [1] MIT License in 7 seconds...             │
└─────────────────────────────────────────────────────────────┘
```

### Question UI Format (Web/API/Chat)
```json
{
  "type": "timed_question",
  "question": "This repo has no LICENSE file. Which license should I use?",
  "options": [
    {"id": 1, "label": "MIT License", "recommended": true},
    {"id": 2, "label": "Apache 2.0", "recommended": false},
    {"id": 3, "label": "GPL v3", "recommended": false},
    {"id": 4, "label": "Skip — no license", "recommended": false}
  ],
  "timeout_seconds": 10,
  "default_on_timeout": 1,
  "context": "No LICENSE file detected in repo root"
}
```

### When Questions Are Asked
muah-runner asks questions at these decision points:
```yaml
question_triggers:
  planning_phase:
    - "What is the primary purpose of this repo?"
    - "Which platform docs should I generate? (GitHub / GitLab / Both / Custom)"
    - "Should I set up CI/CD workflows?"
    - "What test framework should I use?"
    - "Do you want Docker support?"

  docs_phase:
    - "No LICENSE found. Which license?"
    - "README exists but is sparse. Regenerate or enhance?"
    - "Missing CONTRIBUTING.md. Create one?"
    - "No issue/PR templates found. Create standard templates?"

  execution_phase:
    - "Tests are failing. Should I attempt auto-fix or skip?"
    - "Dependency conflict detected. Use version A or B?"
    - "This action will modify 15+ files. Proceed?"
    - "Docker simulation found a fix. Apply it?"

  completion_phase:
    - "Create a PR with these changes?"
    - "Should I tag a release?"
    - "Run full Playwright visual check before finishing?"
```

### Recommendation Engine
muah-runner picks smart defaults based on:
- **Memory**: "Last time in a similar repo, we chose X"
- **Repo context**: "This is a Node.js project, so MIT is standard"
- **Platform norms**: "GitHub repos typically have these docs"
- **Spec**: "The spec says we should always do X"
- **Industry convention**: "95% of open-source projects use this pattern"

### Question Config
```yaml
# .muah/questions/config.yml
timer_seconds: 10                 # Default timeout
min_timer: 5                      # Minimum (for urgent questions)
max_timer: 30                     # Maximum (for complex questions)
auto_mode: false                  # If true, skip ALL questions, always use recommended
silent_mode: false                # If true, don't show questions, just log them
escalation:
  on_critical_question: extend_to_30s   # More time for destructive actions
  on_repeated_timeout: suggest_auto_mode # If human never answers, suggest auto-mode
```

### Question Memory
Every question and answer is logged:
```json
{
  "question_id": "q-20260316-143025",
  "run_id": "run-20260316-143022",
  "question": "No LICENSE found. Which license?",
  "options": ["MIT", "Apache 2.0", "GPL v3", "Skip"],
  "recommended": "MIT",
  "answered_by": "timeout",
  "answer": "MIT",
  "response_time_ms": null,
  "timestamp": "2026-03-16T14:30:25Z"
}
```

---

## Planning-First, Docs-First Workflow

muah-runner NEVER jumps straight into code. It always follows this sequence:

### Phase -1: Tool-Use Memory Bootstrap (before EVERYTHING)
```
Before muah-runner even looks at the repo:

1. Detect which AI backend is running (Claude? Gemini? Copilot?)
2. Load .muah/tool_use/registry.json → full tool catalog
3. Load .muah/tool_use/ai_profiles/{ai}.json → this AI's strengths/weaknesses
4. Load .muah/tool_use/teaching/lessons.json → "do this, not that"
5. Load .muah/tool_use/teaching/anti-patterns.json → "NEVER do these"
6. Load .muah/tool_use/recovery/strategies.json → known recovery paths
7. Inject ALL of this into the AI's working context
8. The AI now knows every tool, how to use them, and what to avoid
   BEFORE it makes a single tool call

If this is the FIRST run ever (no .muah/ exists yet):
  → Create .muah/ with default registry from built-in tool catalog
  → Create blank AI profile (will be populated after this run)
  → Skip teaching materials (nothing to teach yet)
```

### Phase 0: Repo Scan (automatic, no questions)
```
Before anything else, muah-runner scans the entire repo:

1. Read every file and directory (respect .gitignore)
2. Detect tech stack: languages, frameworks, package managers, build tools
3. Detect existing features: tests, CI/CD, Docker, docs, configs
4. Detect platform: GitHub (.github/), GitLab (.gitlab-ci.yml), etc.
5. Map dependencies and their versions
6. Identify gaps: what's missing that should be there
7. Check for existing .muah/ directory (resume from memory if found)
8. Store results in .muah/planning/repo-scan.json
```

### Phase 1: Plan / Roadmap (ask questions here)
```
Based on the repo scan, muah-runner creates a plan:

1. Generate a roadmap with phases and milestones
2. ASK (timed, 10s): "Here's what I found and what I plan to do. Approve?"
   [RECOMMENDED: Approve plan]
3. ASK (timed, 10s): "What is the priority? Docs / Tests / Features / Fix bugs?"
   [RECOMMENDED: based on biggest gap found in scan]
4. ASK (timed, 10s): "Any specific features or goals to add to the plan?"
   [RECOMMENDED: Skip — use auto-detected plan]
5. Save plan to .muah/planning/current-plan.md
6. Save structured roadmap to .muah/planning/roadmap.json
7. Log plan in changelog
```

### Phase 2: Documentation (BEFORE any code)
```
muah-runner creates ALL required documentation first:

1. Scan what docs exist (README? LICENSE? CONTRIBUTING? Templates?)
2. ASK (timed, 10s): "Which platform(s) should I generate docs for?"
   Options: [GitHub] [GitLab] [Both] [Custom]
   [RECOMMENDED: auto-detect from repo — if .github/ exists, pick GitHub]
3. ASK (timed, 10s): "Generate standard docs package?"
   [RECOMMENDED: Yes]
4. Generate ALL standard docs for the detected/chosen platform
5. If docs already exist, enhance them (don't overwrite without asking)
6. Run doc audit — check coverage, flag stale docs
7. Save manifest to .muah/docs/doc-manifest.json
```

### Standard Docs Package (auto-generated per platform)

#### GitHub Standard Docs
```
repo-root/
├── README.md                      # Project overview, setup, usage, contributing link
├── LICENSE                        # Chosen license (default: MIT)
├── CHANGELOG.md                   # Links to .muah/changelog/ or standalone
├── CONTRIBUTING.md                # How to contribute, code style, PR process
├── CODE_OF_CONDUCT.md             # Contributor Covenant or custom
├── SECURITY.md                    # How to report vulnerabilities
├── ARCHITECTURE.md                # System design overview (generated from repo scan)
├── .github/
│   ├── PULL_REQUEST_TEMPLATE.md   # PR template with checklist
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md          # Bug report template
│   │   ├── feature_request.md     # Feature request template
│   │   └── config.yml             # Issue template chooser config
│   ├── FUNDING.yml                # Sponsorship links (if applicable)
│   ├── CODEOWNERS                 # Code ownership mapping
│   ├── dependabot.yml             # Dependency update config
│   └── workflows/
│       ├── ci.yml                 # CI workflow (test, lint, build)
│       ├── release.yml            # Release workflow (tag, changelog, publish)
│       └── muah-runner.yml        # muah-runner's own workflow
```

#### GitLab Standard Docs
```
repo-root/
├── README.md
├── LICENSE
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── ARCHITECTURE.md
├── .gitlab/
│   ├── merge_request_templates/
│   │   └── default.md
│   └── issue_templates/
│       ├── bug.md
│       └── feature.md
├── .gitlab-ci.yml                 # CI/CD pipeline config
```

#### Any Platform (always generated)
```
repo-root/
├── README.md                      # Always. Non-negotiable.
├── LICENSE                        # Always. Ask which one (10s timer → MIT).
├── CHANGELOG.md                   # Always. Linked to .muah/changelog/.
├── .editorconfig                  # Consistent editor settings
├── .gitignore                     # If missing, generate from tech stack
├── .gitattributes                 # Line endings, binary detection
```

### Phase 3: Execute Task (now code can happen)
```
Only AFTER planning and docs are done:

1. Load plan from .muah/planning/current-plan.md
2. Execute tasks according to plan, in order
3. For each task:
   a. Check spec
   b. Execute (using tools, MCPs, Docker, Playwright as needed)
   c. Validate against quality gates
   d. Update memory
   e. Update changelog
   f. Check for spec drift
4. If a task fails → Docker self-healing flow
5. If a question arises → timed question (10s)
```

### Phase 4: Review & Close
```
After all tasks complete:

1. Run doc audit again — did we create new features that need doc updates?
2. ASK (timed, 10s): "All tasks complete. Create a summary PR?"
   [RECOMMENDED: Yes]
3. Generate run summary with:
   - What was planned vs. what was done
   - All changelog entries from this run
   - Spec compliance report
   - Docs created/updated
   - Questions asked and answers given
4. Update .muah/memory/ with full run context
5. Clean up: Docker containers, temp files, locks
```

### Phase 5: Self-Reflection (automatic, never skipped)
```
After Phase 4 completes, muah-runner turns inward:

1. Run the 5-question self-interview:
   - "What could I have done better?"
   - "What was I missing?"
   - "What do I need next time?"
   - "Did I follow the plan or drift?"
   - "What should future-me know about this?"
2. Generate action items from answers
3. Calculate self-rating (F through A+)
4. Compare with last similar run → measure improvement
5. Write selfreflection to .muah/memory/selfreflection/reflections/
6. Update action-items.json with new open items
7. Check if pattern detection interval reached → if yes, run pattern analysis
8. Append to growth-log.md
9. Write "future self briefing" for next run on this repo

This phase is MANDATORY. It cannot be skipped, timed out, or disabled.
It is how muah-runner gets smarter. Without it, the runner is just a script.
```

### Intelligent Doc Generation
muah-runner doesn't just copy templates — it generates docs that are **specific to the repo**:
```yaml
smart_docs:
  README:
    - Auto-detect project name, description from package.json / Cargo.toml / go.mod
    - Generate install instructions based on detected package manager
    - Generate usage examples based on exported functions/endpoints
    - Add badges: CI status, version, license, coverage
    - Include architecture diagram if complex enough

  CONTRIBUTING:
    - Match code style rules from detected linter configs
    - Include actual branch naming convention from git history
    - Reference actual test commands from package.json / Makefile

  ARCHITECTURE:
    - Generate from repo scan: modules, dependencies, data flow
    - Include Mermaid diagrams for complex systems
    - Update automatically when repo structure changes

  SECURITY:
    - Include actual supported versions from release history
    - Link to actual security contact or process
```

---

## Platform Adapters (MCP-Powered)

muah-runner works everywhere via MCP-powered adapters:

### GitHub Actions Adapter
- Register as a self-hosted runner
- **Connect to GitHub MCP server on startup** — use it for ALL GitHub operations
- Map GitHub workflow jobs to muah-runner tasks
- Report status back via GitHub MCP (checks API, PR comments)
- Access: secrets, artifacts, caches — all through MCP first, CLI fallback
- Playwright: auto-run visual regression on PR preview deploys

### Claude Adapter
- Run as a tool within Claude conversations
- **Connect to Claude MCP** for web search, code execution, conversation tools
- Read/write `.muah/` for persistent memory across chats
- Sync Claude's memory system with `.muah/memory/`
- Use Claude's search tool for real-time information gathering

### Gemini Adapter
- Integrate with Gemini MCP
- Use Gemini's function calling and grounding capabilities
- Maintain memory compatibility with other adapters

### Generic/Standalone Adapter
- Run from CLI with no external platform
- **Scan for any available MCP servers** in the environment
- Accept tasks via REST API, webhooks, or file watchers
- Cron/schedule support for recurring tasks
- Docker MCP + Playwright MCP always available in standalone mode

---

## Project Structure

```
muah-runner/
├── cmd/
│   └── muah-runner/
│       └── main.go                 # CLI entrypoint
├── internal/
│   ├── core/
│   │   ├── runner.go               # Main run loop
│   │   └── lifecycle.go            # Init, run, cleanup
│   ├── memory/
│   │   ├── store.go                # Memory read/write
│   │   ├── search.go               # Memory search/query
│   │   ├── merge.go                # Memory merging across runs
│   │   └── schema.go               # Memory data structures
│   ├── changelog/
│   │   ├── writer.go               # Changelog entry creation
│   │   ├── formatter.go            # Markdown/JSON formatting
│   │   └── integrity.go            # Hash chain verification
│   ├── spec/
│   │   ├── loader.go               # Load specs from .muah/spec/
│   │   ├── validator.go            # Validate runs against specs
│   │   ├── drift.go                # Drift detection
│   │   └── proposer.go             # Propose spec updates
│   ├── mcp/
│   │   ├── discovery.go            # MCP server auto-discovery
│   │   ├── connection.go           # MCP connection management
│   │   ├── registry.go             # Unified MCP tool registry
│   │   ├── platform/               # Platform MCP clients
│   │   │   ├── github.go           # GitHub MCP client
│   │   │   ├── claude.go           # Claude MCP client
│   │   │   └── gemini.go           # Gemini MCP client
│   │   ├── native/                 # Native MCP integrations
│   │   │   ├── playwright.go       # Playwright MCP (required)
│   │   │   ├── docker.go           # Docker MCP (required)
│   │   │   ├── filesystem.go       # Filesystem MCP
│   │   │   └── sqlite.go           # SQLite MCP
│   │   └── fallback.go             # CLI fallback when MCP is down
│   ├── docker/
│   │   ├── engine.go               # Docker operations engine
│   │   ├── sandbox.go              # Sandbox container management
│   │   ├── simulation.go           # Bug reproduction & fix simulation
│   │   ├── selfheal.go             # Self-healing orchestrator
│   │   ├── compose.go              # Docker Compose operations
│   │   ├── templates.go            # Compose template management
│   │   └── cleanup.go              # Container/image garbage collection
│   ├── planning/
│   │   ├── scanner.go              # Repo scanning and analysis
│   │   ├── planner.go              # Plan/roadmap generation
│   │   ├── roadmap.go              # Structured roadmap with phases
│   │   ├── feature_inventory.go    # Detect existing features and gaps
│   │   └── adr.go                  # Architecture Decision Records
│   ├── docs/
│   │   ├── generator.go            # Smart doc generation engine
│   │   ├── templates.go            # Doc template management
│   │   ├── audit.go                # Doc coverage audit
│   │   ├── platform_detect.go      # Detect GitHub/GitLab/etc.
│   │   ├── readme.go               # Intelligent README generation
│   │   ├── contributing.go         # CONTRIBUTING.md generation
│   │   ├── architecture.go         # ARCHITECTURE.md from repo scan
│   │   └── standard_package.go     # Full standard docs package per platform
│   ├── questions/
│   │   ├── asker.go                # Timed question engine
│   │   ├── timer.go                # Countdown timer (10s default)
│   │   ├── recommender.go          # Smart default recommendation engine
│   │   ├── renderer.go             # CLI / Web / API question rendering
│   │   ├── history.go              # Question + answer logging
│   │   └── config.go               # Timer config, auto-mode, escalation
│   ├── selfreflection/
│   │   ├── interviewer.go          # Post-run self-interview engine (5 questions)
│   │   ├── rater.go                # Self-rating system (F through A+)
│   │   ├── action_items.go         # Generate and track action items
│   │   ├── patterns.go             # Cross-run pattern detection
│   │   ├── growth_log.go           # Human-readable growth log writer
│   │   ├── briefing.go             # "Future self" briefing note generator
│   │   └── loader.go               # Load past reflections on startup
│   ├── toollearn/
│   │   ├── registry.go             # Master tool-use registry (all tools, all commands)
│   │   ├── catalog.go              # Detailed per-tool docs and usage specs
│   │   ├── profiler.go             # AI profiling — track per-AI strengths/weaknesses
│   │   ├── monitor.go              # Real-time tool_use call monitoring
│   │   ├── classifier.go           # Error classification (auth, params, timeout, etc.)
│   │   ├── recovery.go             # API recovery trigger (3 fails = escalate)
│   │   ├── fallback.go             # Fallback chain: MCP → CLI → API → Docker → human
│   │   ├── teacher.go              # Generate teaching materials from failures
│   │   ├── injector.go             # Context injection — load tool knowledge into AI
│   │   ├── playbooks.go            # Recovery playbook executor
│   │   ├── metrics.go              # Tool-use success/failure dashboard
│   │   └── antipatterns.go         # Anti-pattern detection and prevention
│   ├── privacy/
│   │   ├── cot.go                  # Chain of Thought privacy engine
│   │   ├── classifier.go           # Classify thoughts: private/redactable/public
│   │   ├── encryptor.go            # AES-256-GCM CoT encryption
│   │   ├── redactor.go             # Redact CoT for user-facing output
│   │   ├── vault.go                # Encrypted CoT storage
│   │   └── audit.go                # CoT audit logging
│   ├── swarm/
│   │   ├── party_leader.go         # Main agent — user communication, coordination
│   │   ├── spawner.go              # Agent spawning (Summon) system
│   │   ├── monitor.go              # HP/MP monitoring loop (every 5s)
│   │   ├── lifecycle.go            # Agent lifecycle: spawn → buff → engage → dismiss
│   │   ├── jobs.go                 # All 22 job classes and their abilities
│   │   ├── races.go                # Race system — AI model to race mapping
│   │   ├── race_stats.go           # Racial stat modifiers, native skills
│   │   ├── race_cot.go             # Race-specific CoT and memory styles
│   │   ├── composer.go             # Party composition engine (race+job+subjob)
│   │   ├── subjob.go               # Sub-job system — unlock at 30, half main level
│   │   ├── leveling.go             # XP, leveling, ability unlocks, merit points
│   │   ├── hp_mp.go                # HP/MP mechanics, healing, death, revival
│   │   ├── crystal_pool.go         # Shared resource pool management
│   │   ├── completion.go           # Completion trigger handling, Victory Fanfare
│   │   ├── graveyard.go            # K.O. handling, tombstones, Raise/Reraise
│   │   ├── linkshell.go            # Inter-agent messaging (Linkshell channels)
│   │   ├── tongue/                 # Vanadiel Tongue language system
│   │   │   ├── lexicon.go          # Vocabulary and word generation
│   │   │   ├── grammar.go          # Message structure and parsing
│   │   │   ├── encoder.go          # Encode messages in Vanadiel Tongue
│   │   │   ├── decoder.go          # Decode for debugging (non-Hume only)
│   │   │   └── translator.go       # Full translation layer
│   │   └── scaling.go              # Spawn decision logic (solo/party/full/epic)
│   ├── playwright/
│   │   ├── setup.go                # Auto-install Playwright + browsers
│   │   ├── browser.go              # Browser session management
│   │   ├── visual.go               # Visual regression testing
│   │   ├── scraper.go              # Web scraping engine
│   │   ├── e2e.go                  # E2E test runner
│   │   ├── accessibility.go        # A11y audit runner
│   │   └── pdf.go                  # PDF generation from HTML
│   ├── tools/
│   │   ├── registry.go             # Tool registration and lookup
│   │   ├── executor.go             # Tool execution engine
│   │   ├── generator.go            # Self-creating tool logic
│   │   ├── builtin/                # Built-in tools
│   │   │   ├── search.go
│   │   │   ├── fetch.go
│   │   │   ├── fileops.go
│   │   │   ├── codeexec.go
│   │   │   ├── gitops.go
│   │   │   └── shell.go
│   │   └── skills/
│   │       ├── manager.go          # Skill lifecycle
│   │       └── composer.go         # Compose tools into skills
│   ├── adapters/
│   │   ├── adapter.go              # Adapter interface
│   │   ├── github.go
│   │   ├── claude.go
│   │   ├── gemini.go
│   │   └── generic.go
│   ├── sensor/
│   │   ├── environment.go          # OS, arch, hardware detection
│   │   ├── toolchain.go            # Installed tools detection
│   │   ├── mcp_scan.go             # Scan environment for MCP servers
│   │   └── health.go               # Self-monitoring
│   └── config/
│       ├── loader.go
│       └── secrets.go
├── api/
│   ├── server.go                   # REST API for standalone mode
│   ├── webhooks.go
│   └── metrics.go                  # Prometheus metrics
├── templates/
│   └── muah/                       # Default .muah/ directory template
│       ├── memory/
│       ├── changelog/
│       ├── spec/
│       ├── tools/
│       ├── mcp/
│       ├── docker/
│       ├── tool_use/
│       ├── privacy_cot/
│       ├── swarm/
│       ├── planning/
│       ├── docs/
│       ├── questions/
│       ├── runs/
│       └── config/
├── scripts/
│   ├── install.sh
│   ├── setup.sh
│   ├── install-playwright.sh       # Auto-install Playwright + Chromium
│   ├── install-docker-mcp.sh       # Docker MCP setup
│   └── migrate.sh                  # Migrate .muah/ between versions
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── muah-runner.yml.example
├── .github/
│   └── workflows/
│       └── ci.yml
├── README.md
└── LICENSE
```

---

## CLI Commands

```bash
# Core
muah-runner init              # Create .muah/ dir, bootstrap tool catalog, install native MCPs, scan repo, plan, generate docs
muah-runner start             # Start daemon mode, connect all MCPs
muah-runner stop              # Graceful shutdown, disconnect MCPs, cleanup containers
muah-runner status            # Show health, active runs, MCP connections, Docker state
muah-runner run <task>        # Execute a single task

# Memory
muah-runner memory search <q> # Search memory
muah-runner memory show       # Show recent memory entries
muah-runner memory export     # Export memory as portable JSON

# Changelog
muah-runner changelog show    # Show recent changes
muah-runner changelog diff    # Diff between two runs

# Spec
muah-runner spec check        # Validate current state against specs
muah-runner spec drift        # Show spec drift report
muah-runner spec propose      # AI-generate spec updates based on learned patterns

# Tools & Skills
muah-runner tools list        # List all available tools (builtin + MCP + generated)
muah-runner tools create <n>  # Manually trigger tool creation
muah-runner skills list       # List all skills

# Planning & Docs (Phase 0-2 workflow)
muah-runner scan              # Scan repo: detect stack, features, gaps, platform
muah-runner plan              # Generate/update plan and roadmap
muah-runner plan show         # Show current plan
muah-runner docs generate     # Generate all standard docs for detected platform
muah-runner docs audit        # Check doc coverage, flag missing/stale docs
muah-runner docs update       # Re-generate docs based on current repo state

# Questions
muah-runner questions history # Show all past questions and answers
muah-runner questions auto    # Enable auto-mode (always use recommended, no timer)
muah-runner questions manual  # Disable auto-mode (ask with 10s timer)

# Self-Reflection
muah-runner reflect show      # Show last self-reflection entry
muah-runner reflect growth    # Show growth log (improvement over time)
muah-runner reflect patterns  # Show detected patterns across runs
muah-runner reflect actions   # Show open action items from self-reflection
muah-runner reflect grade     # Show self-rating history and trend

# Tool-Use Learning & Recovery
muah-runner tools proficiency # Show tool-use proficiency per AI backend
muah-runner tools failures    # Show recent failures, root causes, recoveries
muah-runner tools teach       # Show teaching materials generated from failures
muah-runner tools profile <ai># Show detailed AI profile (strengths, weaknesses, scores)
muah-runner tools dashboard   # Show tool-use metrics dashboard
muah-runner tools anti        # Show known anti-patterns to avoid
muah-runner tools playbook <e># Show recovery playbook for error type
muah-runner tools catalog     # List all tools with usage docs and AI-specific notes
muah-runner tools inject      # Force-reload tool knowledge into current AI context

# MCP
muah-runner mcp discover      # Scan environment for available MCP servers
muah-runner mcp list          # Show all connected MCP servers and their tools
muah-runner mcp connect <url> # Manually connect to an MCP server
muah-runner mcp test          # Test all MCP connections are healthy
muah-runner mcp install <name> # Install a native MCP (playwright, docker, etc.)

# Docker & Simulation
muah-runner docker status     # Show Docker availability, running containers
muah-runner docker simulate   # Spin up a sandbox to reproduce last failure
muah-runner docker heal       # Auto-diagnose and fix last failure via Docker
muah-runner docker matrix     # Run multi-version matrix test in containers
muah-runner docker cleanup    # Remove simulation containers and old images

# Playwright
muah-runner playwright setup  # Install Playwright + browsers
muah-runner playwright test   # Run Playwright test suite
muah-runner playwright screenshot <url>  # Quick screenshot of a URL
muah-runner playwright audit <url>       # Accessibility + performance audit

# Privacy CoT
muah-runner cot status        # Show CoT encryption status
muah-runner cot redacted <id> # View redacted CoT for a run
muah-runner cot audit         # Show CoT audit log
muah-runner cot classify      # Show classification rules

# Swarm (Vana'diel Protocol)
muah-runner swarm status      # Show party roster, HP/MP, agent statuses
muah-runner swarm party       # Show current party composition
muah-runner swarm linkshell   # Show recent Linkshell messages (translated)
muah-runner swarm graveyard   # Show K.O. agents and tombstones
muah-runner swarm crystal     # Show Crystal Pool allocation
muah-runner swarm spawn <job> # Manually spawn an agent with a Job
muah-runner swarm dismiss <n> # Dismiss an agent by NPC name
muah-runner swarm raise <n>   # Revive a K.O. agent
muah-runner swarm tongue      # Show Vanadiel Tongue lexicon
muah-runner swarm history     # Show past swarm missions and party compositions
muah-runner swarm solo        # Force solo mode (no spawning)
muah-runner swarm jobs        # List all 22 jobs with current unlock status
muah-runner swarm level <npc> # Show agent level, XP, unlocked abilities
muah-runner swarm subjob <npc> <job>  # Set sub-job for an agent
muah-runner swarm combos      # Show recommended sub-job combinations
muah-runner swarm races       # List all races and their AI model mappings
muah-runner swarm race <name> # Show detailed race profile (stats, skills, CoT style)
muah-runner swarm compose <task>  # Preview party composition for a given task (dry run)
muah-runner swarm benchmark   # Run benchmark comparing race performance

# Platform Registration
muah-runner register github   # Register with GitHub as self-hosted runner
muah-runner register claude   # Set up Claude adapter
muah-runner register gemini   # Set up Gemini adapter

# Maintenance
muah-runner self-update       # Pull and apply updates
muah-runner gc                # Garbage collect old runs/artifacts/containers
muah-runner doctor            # Full health check: MCP, Docker, Playwright, swarm, CoT, memory, specs
```

---

## Key Principles

1. **Memory is sacred** — never lose context, always persist, always searchable
2. **Changelog everything** — if muah-runner touched it, it's logged
3. **Specs keep you honest** — drift detection prevents silent degradation
4. **MCP-first, always** — use the platform's MCP server before CLI, before API, before anything. Discover every MCP available. Use them optimally.
5. **Tools are composable** — small tools compose into powerful skills. MCP tools, built-in tools, and generated tools are all equal citizens.
6. **Self-improvement is mandatory** — after EVERY run, muah-runner interviews itself: "What could I have done better? What do I need next time?" It writes action items, detects patterns across runs, grades itself, and gets measurably better over time. The selfreflection is never skipped.
7. **Platform agnostic** — same `.muah/` dir works whether you're in GitHub, Claude, Gemini, or terminal. MCP adapters make it seamless.
8. **Don't just fail — fix it** — Docker simulation engine reproduces bugs, tests fixes in isolation, and self-heals before reporting failure
9. **Playwright is always on** — browser automation, visual testing, scraping, and accessibility audits are core capabilities, not add-ons
10. **Never block on humans** — timed questions with 10-second countdown. No answer? Auto-select the smart default. The runner never stalls.
11. **Plan before you build** — scan the repo, make a roadmap, generate all docs FIRST. Code comes after planning and documentation.
12. **Docs are not optional** — every repo gets a full standard docs package. README, LICENSE, CONTRIBUTING, templates — all generated intelligently from repo context.
13. **Teach the AI before it works** — every run begins by loading the full tool catalog, AI-specific profile, lessons, and anti-patterns into context. Don't let the AI guess — front-load the knowledge.
14. **3 strikes, recover** — if a tool fails 3 times, STOP retrying. Classify the error, fall back (MCP → CLI → API → Docker → human), generate teaching material, update the AI profile. Never retry blindly with identical parameters.
15. **Every AI is different** — Claude, Gemini, Copilot all have different tool-use strengths. Profile each one, teach each one differently, track improvement per model.
16. **Runs stay tight** — respect constraints, enforce quality gates, fail fast and loud
17. **Docker is your lab** — spin up, simulate, test, tear down. Never pollute the host. Always clean up.
18. **Thoughts are private** — every agent's Chain of Thought is encrypted at rest, classified by sensitivity, and redacted on demand. Users see results, not reasoning. Secrets never leak through CoT.
19. **Swarm like Vana'diel** — complex tasks spawn a party of FFXI-themed agents, each with their own Race (AI model), Job, HP (context), MP (tokens), magic (tools), and private CoT. The Party Leader composes the party: Elvaan (Opus) for deep thinking, Galka (Sonnet) for sustained work, Mithra (Haiku) for speed, Tarutaru (prev-gen/o1) for wisdom. Right race, right job, right task.
20. **Every agent is accountable** — UUID, timestamp, NPC name, race, job, sub-job, level. Every agent has full identity. Every action is traceable. Every death has a tombstone. Every victory gets a fanfare.

Include comprehensive unit tests, integration tests, MCP integration tests, Docker simulation tests, Playwright test suite, and a CI workflow. Write a detailed README explaining the philosophy, setup, and usage. The README must include a quickstart that gets muah-runner running with MCP discovery, Playwright, and Docker in under 5 minutes.
