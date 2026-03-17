package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestBootstrap_CreatesDirectoryTree verifies that Bootstrap creates the
// full .muah/ directory structure and every required subdirectory.
func TestBootstrap_CreatesDirectoryTree(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	muah := filepath.Join(baseDir, MuahDir)

	requiredDirs := []string{
		filepath.Join(muah, "memory", "runs"),
		filepath.Join(muah, "memory", "context"),
		filepath.Join(muah, "memory", "selfreflection", "reflections"),
		filepath.Join(muah, "tool_use", "ai_profiles"),
		filepath.Join(muah, "tool_use", "failure_log"),
		filepath.Join(muah, "privacy_cot", "vault"),
		filepath.Join(muah, "swarm", "party", "members"),
		filepath.Join(muah, "swarm", "jobs", "levels"),
		filepath.Join(muah, "changelog", "entries"),
		filepath.Join(muah, "spec", "task-specs", "custom"),
		filepath.Join(muah, "mcp", "platform"),
		filepath.Join(muah, "mcp", "native"),
		filepath.Join(muah, "docker", "simulations"),
		filepath.Join(muah, "planning", "plans"),
		filepath.Join(muah, "docs", "generated"),
		filepath.Join(muah, "questions"),
		filepath.Join(muah, "runs", "current"),
		filepath.Join(muah, "config"),
	}

	for _, dir := range requiredDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("expected dir %s to exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}
}

// TestBootstrap_CreatesSeedFiles verifies that essential seed files are written.
func TestBootstrap_CreatesSeedFiles(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	muah := filepath.Join(baseDir, MuahDir)

	requiredFiles := []string{
		filepath.Join(muah, "memory", "index.json"),
		filepath.Join(muah, "tool_use", "registry.json"),
		filepath.Join(muah, "swarm", "vana_diel.json"),
		filepath.Join(muah, "swarm", "party", "party-leader.json"),
		filepath.Join(muah, "swarm", "jobs", "job-registry.json"),
		filepath.Join(muah, "swarm", "races", "race-registry.json"),
		filepath.Join(muah, "changelog", "CHANGELOG.md"),
		filepath.Join(muah, "spec", "runner-spec.yml"),
		filepath.Join(muah, "spec", "constraints.yml"),
		filepath.Join(muah, "spec", "quality-gates.yml"),
		filepath.Join(muah, "mcp", "discovery.json"),
		filepath.Join(muah, "planning", "repo-scan.json"),
		filepath.Join(muah, "docs", "doc-manifest.json"),
		filepath.Join(muah, "questions", "history.json"),
		filepath.Join(muah, "config", "muah-runner.yml"),
		filepath.Join(muah, "privacy_cot", "config.yml"),
		filepath.Join(muah, "privacy_cot", "classification.yml"),
	}

	for _, f := range requiredFiles {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected seed file %s to exist: %v", f, err)
		}
	}
}

// TestBootstrap_Idempotent verifies that calling Bootstrap twice does not
// overwrite existing files or return an error.
func TestBootstrap_Idempotent(t *testing.T) {
	baseDir := t.TempDir()

	// First call
	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("first Bootstrap() error: %v", err)
	}

	// Write a sentinel value into a seed file
	muah := filepath.Join(baseDir, MuahDir)
	indexPath := filepath.Join(muah, "memory", "index.json")
	sentinel := []byte(`{"total_runs":999,"last_run_id":"sentinel","run_ids":[],"last_update":"2026-01-01T00:00:00Z"}`)
	if err := os.WriteFile(indexPath, sentinel, 0644); err != nil {
		t.Fatalf("writing sentinel: %v", err)
	}

	// Second call — must not overwrite existing file
	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("second Bootstrap() error: %v", err)
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading index: %v", err)
	}
	if string(data) != string(sentinel) {
		t.Errorf("Bootstrap() overwrote existing file: got %s", data)
	}
}

// TestBootstrap_MemoryIndexValid verifies that the generated memory index is
// valid JSON with expected initial values.
func TestBootstrap_MemoryIndexValid(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	indexPath := filepath.Join(baseDir, MuahDir, "memory", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading index: %v", err)
	}

	var idx map[string]interface{}
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("memory/index.json is not valid JSON: %v\ncontent: %s", err, data)
	}

	if totalRuns, ok := idx["total_runs"]; !ok || totalRuns == nil {
		t.Error("memory/index.json missing 'total_runs'")
	}
	if runIDs, ok := idx["run_ids"]; !ok || runIDs == nil {
		t.Error("memory/index.json missing 'run_ids'")
	}
}

// TestBootstrap_JobRegistryHas22Jobs verifies the Vana'diel job registry is
// seeded with the expected 22 total jobs.
func TestBootstrap_JobRegistryHas22Jobs(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	jobPath := filepath.Join(baseDir, MuahDir, "swarm", "jobs", "job-registry.json")
	data, err := os.ReadFile(jobPath)
	if err != nil {
		t.Fatalf("reading job-registry.json: %v", err)
	}

	var reg map[string]interface{}
	if err := json.Unmarshal(data, &reg); err != nil {
		t.Fatalf("job-registry.json is not valid JSON: %v", err)
	}

	totalVal, ok := reg["total"]
	if !ok {
		t.Fatal("job-registry.json missing 'total'")
	}
	// JSON numbers unmarshal as float64
	if total, ok := totalVal.(float64); !ok || int(total) != 22 {
		t.Errorf("expected 22 total jobs, got %v", totalVal)
	}
}

// TestBootstrap_RaceRegistryHas5Races verifies the Vana'diel race registry
// contains exactly 5 races.
func TestBootstrap_RaceRegistryHas5Races(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	racePath := filepath.Join(baseDir, MuahDir, "swarm", "races", "race-registry.json")
	data, err := os.ReadFile(racePath)
	if err != nil {
		t.Fatalf("reading race-registry.json: %v", err)
	}

	var reg map[string]interface{}
	if err := json.Unmarshal(data, &reg); err != nil {
		t.Fatalf("race-registry.json is not valid JSON: %v", err)
	}

	races, ok := reg["races"].([]interface{})
	if !ok {
		t.Fatal("race-registry.json missing 'races' array")
	}
	if len(races) != 5 {
		t.Errorf("expected 5 races, got %d", len(races))
	}
}

// TestBootstrap_PartyLeaderIsPrishe verifies the default party leader is Prishe.
func TestBootstrap_PartyLeaderIsPrishe(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	leaderPath := filepath.Join(baseDir, MuahDir, "swarm", "party", "party-leader.json")
	data, err := os.ReadFile(leaderPath)
	if err != nil {
		t.Fatalf("reading party-leader.json: %v", err)
	}

	var leader map[string]interface{}
	if err := json.Unmarshal(data, &leader); err != nil {
		t.Fatalf("party-leader.json is not valid JSON: %v", err)
	}

	if name, ok := leader["npc_name"]; !ok || name != "Prishe" {
		t.Errorf("expected party leader 'Prishe', got %v", leader["npc_name"])
	}
}

// TestBootstrap_ChangelogMarkdownExists verifies the CHANGELOG.md seed file
// is a non-empty Markdown file.
func TestBootstrap_ChangelogMarkdownExists(t *testing.T) {
	baseDir := t.TempDir()

	if err := Bootstrap(baseDir); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}

	changelogPath := filepath.Join(baseDir, MuahDir, "changelog", "CHANGELOG.md")
	data, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Fatalf("reading CHANGELOG.md: %v", err)
	}
	if len(data) == 0 {
		t.Error("CHANGELOG.md is empty")
	}
	if string(data[:1]) != "#" {
		t.Errorf("CHANGELOG.md should start with '#', got: %q", string(data[:1]))
	}
}
