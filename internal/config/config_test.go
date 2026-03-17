package config

import (
"os"
"path/filepath"
"testing"
)

func TestDefaultConfig(t *testing.T) {
cfg := DefaultConfig()
if cfg.Runner.Name == "" {
t.Error("expected non-empty runner name")
}
if cfg.Runner.MuahDir != ".muah" {
t.Errorf("expected .muah dir, got %s", cfg.Runner.MuahDir)
}
if cfg.Server.Port != 8080 {
t.Errorf("expected port 8080, got %d", cfg.Server.Port)
}
}

func TestLoadNonExistent(t *testing.T) {
cfg, err := Load("/nonexistent/path.yml")
if err != nil {
t.Fatalf("expected no error for missing file, got %v", err)
}
if cfg == nil {
t.Fatal("expected non-nil config")
}
}

func TestSaveLoad(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "config.yml")
cfg := DefaultConfig()
cfg.Runner.Name = "test-runner"
if err := Save(cfg, path); err != nil {
t.Fatalf("save error: %v", err)
}
loaded, err := Load(path)
if err != nil {
t.Fatalf("load error: %v", err)
}
if loaded.Runner.Name != "test-runner" {
t.Errorf("expected test-runner, got %s", loaded.Runner.Name)
}
}

func TestLoadSecretsNonExistent(t *testing.T) {
s, err := LoadSecrets("/nonexistent/secrets.json")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if s == nil {
t.Fatal("expected non-nil secrets")
}
}

func TestLoadSecrets(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "secrets.json")
data := []byte(`{"github_token":"tok123"}`)
os.WriteFile(path, data, 0600)
s, err := LoadSecrets(path)
if err != nil {
t.Fatalf("error: %v", err)
}
if s.GithubToken != "tok123" {
t.Errorf("expected tok123, got %s", s.GithubToken)
}
}
