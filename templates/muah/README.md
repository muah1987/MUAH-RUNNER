# templates/muah/ — Default .muah/ seed files

This directory contains the default seed files that `muah-runner init` uses to
bootstrap the `.muah/` directory in any project.

## How it works

When you run `muah-runner init`, the `internal/core.Bootstrap()` function:

1. Creates the full `.muah/` directory tree (40+ directories)
2. Writes these seed files into the new `.muah/` — **only if they don't already exist**
3. Never overwrites existing files (idempotent)

## Structure mirrors `.muah/`

```
templates/muah/
├── memory/
│   └── index.json              ← Empty memory index
├── tool_use/
│   ├── registry.json           ← Empty tool registry
│   ├── recovery/strategies.json ← Default error recovery strategies
│   └── teaching/lessons.json   ← Empty lessons list
├── swarm/
│   ├── vana_diel.json          ← World state
│   ├── party/party-leader.json ← Default leader: Prishe
│   ├── jobs/job-registry.json  ← All 22 FFXI job classes
│   ├── races/race-registry.json ← 5 races mapped to AI tiers
│   └── crystal/pool.json       ← Token budget pool
├── changelog/
│   └── CHANGELOG.md            ← Empty changelog header
├── spec/
│   ├── runner-spec.yml         ← Default runner spec
│   ├── constraints.yml         ← Hard limits (timeout, memory, retries)
│   └── quality-gates.yml       ← Pass criteria (tests, build, spec)
├── mcp/
│   ├── discovery.json          ← Empty discovery state
│   └── connections.json        ← Empty connections
├── planning/
│   └── repo-scan.json          ← Empty scan results
├── docs/
│   └── doc-manifest.json       ← Empty doc manifest
├── questions/
│   ├── history.json            ← Empty question history
│   ├── defaults.yml            ← Auto-answer defaults
│   └── config.yml              ← Timer/auto-mode config
├── config/
│   └── muah-runner.yml         ← Default runtime config
└── privacy_cot/
    ├── config.yml              ← Encryption settings
    ├── classification.yml      ← private/redactable/public rules
    └── audit-log.jsonl         ← Empty audit log
```

## Adding new seed files

1. Add the file here under `templates/muah/`
2. Add the corresponding entry to `writeSeedFiles()` in `internal/core/lifecycle.go`
3. Add a test in `internal/core/core_test.go` to verify it's created by `Bootstrap()`
