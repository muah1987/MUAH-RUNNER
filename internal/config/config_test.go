package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	if cfg.Runner.Name != "muah-runner-01" {
		t.Errorf("expected default runner name, got %q", cfg.Runner.Name)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Executor.Type != "process" {
		t.Errorf("expected default executor type 'process', got %q", cfg.Executor.Type)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port, got %d", cfg.Server.Port)
	}
}

func TestLoadYAML(t *testing.T) {
	content := `
runner:
  name: "test-runner"
  concurrency: 2
server:
  port: 9999
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Runner.Name != "test-runner" {
		t.Errorf("expected 'test-runner', got %q", cfg.Runner.Name)
	}
	if cfg.Runner.Concurrency != 2 {
		t.Errorf("expected concurrency 2, got %d", cfg.Runner.Concurrency)
	}
	if cfg.Server.Port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.Server.Port)
	}
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("MUAH_RUNNER_NAME", "env-runner")
	t.Setenv("MUAH_SERVER_PORT", "7777")

	cfg, err := Load("/nonexistent/path/config.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Runner.Name != "env-runner" {
		t.Errorf("expected 'env-runner', got %q", cfg.Runner.Name)
	}
	if cfg.Server.Port != 7777 {
		t.Errorf("expected port 7777, got %d", cfg.Server.Port)
	}
}
