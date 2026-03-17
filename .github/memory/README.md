# .github/memory — Copilot Session Memory

This directory stores persistent memory snapshots for GitHub Copilot sessions.

## Purpose

Every Copilot session **reads** from this directory at start and **writes** to it at end.
This gives Copilot continuity across sessions — it always knows what was done before,
what decisions were made, and what still needs to be done.

## File Naming Convention

```
session-{YYYYMMDD-HHmmss}-{8-char-uuid}.md
```

Examples:
- `session-20260316-143022-a1b2c3d4.md`
- `session-20260317-091500-f7e2b109.md`

## File Format

Each session file follows this template:

```markdown
# Session Memory: {timestamp}-{uuid}

**Date:** {ISO 8601 timestamp}
**Agent:** GitHub Copilot
**Task:** {what was worked on}

## What Was Done
## Decisions Made
## Current State
## Open Items
## Lessons Learned
## Next Session Should
```

## Rules

1. **Always read** the newest `session-*.md` before starting work
2. **Always write** a new `session-*.md` before ending the session
3. **Never delete** old session files — they are history
4. **Be specific** — vague entries are useless
5. **Include build/test status** in every `## Current State` section
6. Files sorted by timestamp = chronological project history

## Reading Order

Load newest → oldest until you have enough context:

```bash
ls -t .github/memory/session-*.md | head -5
```

## Writing a New Session File

Use timestamp `YYYYMMDD-HHmmss` (UTC) and a random 8-char hex UUID:

```bash
# Example filename generation (bash)
TS=$(date -u +%Y%m%d-%H%M%S)
UUID=$(cat /dev/urandom | tr -dc 'a-f0-9' | head -c 8)
echo "session-${TS}-${UUID}.md"
```
