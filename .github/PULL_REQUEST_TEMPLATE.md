## Description

<!-- Clearly describe what this PR does and why. -->

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Refactor / cleanup
- [ ] Documentation update
- [ ] Tests only
- [ ] Other: ___

## Related issue

Closes #___

---

## Checklist

### Build & tests
- [ ] `make build` passes locally
- [ ] `make test` passes — no new failures, no skipped tests
- [ ] New code has at least one `*_test.go` test

### Code quality
- [ ] No hardcoded file paths (use `filepath.Join` throughout)
- [ ] Errors are wrapped with context (`fmt.Errorf("doing X: %w", err)`)
- [ ] No global state — dependencies passed via constructors
- [ ] No secrets or API keys committed to source

### Documentation
- [ ] `README.md` updated if public behaviour changed
- [ ] Package doc comments added/updated for new/changed exported symbols
- [ ] `CONTRIBUTING.md` still accurate

### Architecture
- [ ] New packages follow the `New*(muahDir string)` constructor convention
- [ ] New tools added to `internal/tools/registry.go`
- [ ] New platform adapters implement the `adapters.Adapter` interface
- [ ] `.muah/` seed files added to `internal/core/lifecycle.go` `writeSeedFiles()` if needed

### Session memory (AI agents only)
- [ ] Read `.github/memory/session-*.md` (newest first) at start of session
- [ ] Wrote `session-{YYYYMMDD-HHmmss}-{8hexchars}.md` at end of session

---

## Phase 5 — Self-Reflection

<!-- What did you learn from making this change? What would you do differently? -->
**Key lesson:**

**Grade (F–A+):**
