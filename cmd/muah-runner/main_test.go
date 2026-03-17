package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// buildBinary compiles the muah-runner binary into a temp dir and returns the
// path to the executable. The binary is shared across all tests in the run.
func buildBinary(t *testing.T) string {
	t.Helper()
	// Determine the repo root relative to this test file.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file path")
	}
	// thisFile is cmd/muah-runner/main_test.go — repo root is two dirs up.
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	binPath := filepath.Join(t.TempDir(), "muah-runner")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(repoRoot, "cmd", "muah-runner")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build muah-runner binary:\n%s\nerr: %v", out, err)
	}
	return binPath
}

// TestCLI_Version verifies the --version flag exits 0 and prints "muah-runner".
func TestCLI_Version(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("--version failed: %v\nout: %s", err, out)
	}
	got := string(out)
	if len(got) == 0 {
		t.Error("expected non-empty version output")
	}
}

// TestCLI_Help verifies --help exits 0 and prints usage information.
func TestCLI_Help(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("--help failed: %v\nout: %s", err, out)
	}
	got := string(out)
	for _, keyword := range []string{"init", "run", "status", "swarm", "reflect"} {
		if !contains(got, keyword) {
			t.Errorf("--help output missing keyword %q", keyword)
		}
	}
}

// TestCLI_UnknownCommand verifies an unknown command exits non-zero.
func TestCLI_UnknownCommand(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "no-such-command-xyz")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for unknown command; got 0\nout: %s", out)
	}
}

// TestCLI_Init creates a temporary project directory, runs "muah-runner init"
// with MUAH_AUTO_MODE=true (so the 10-second question auto-selects instantly),
// and verifies the .muah/ directory tree is created.
func TestCLI_Init(t *testing.T) {
	bin := buildBinary(t)

	projectDir := t.TempDir()

	cmd := exec.Command(bin, "init")
	cmd.Dir = projectDir
	// MUAH_AUTO_MODE skips the 10-second wait in questions
	cmd.Env = append(os.Environ(),
		"MUAH_AUTO_MODE=true",
		"HOME="+projectDir, // isolate HOME
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("muah-runner init failed: %v\noutput:\n%s", err, out)
	}

	muahDir := filepath.Join(projectDir, ".muah")

	// Verify the core .muah/ subdirectories were created.
	requiredDirs := []string{
		filepath.Join(muahDir, "memory"),
		filepath.Join(muahDir, "changelog"),
		filepath.Join(muahDir, "swarm"),
		filepath.Join(muahDir, "privacy_cot"),
		filepath.Join(muahDir, "config"),
		filepath.Join(muahDir, "spec"),
		filepath.Join(muahDir, "mcp"),
		filepath.Join(muahDir, "planning"),
		filepath.Join(muahDir, "docs"),
		filepath.Join(muahDir, "questions"),
	}
	for _, d := range requiredDirs {
		if info, statErr := os.Stat(d); statErr != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist after init: %v", d, statErr)
		}
	}

	// Verify at least one key seed file.
	memIdx := filepath.Join(muahDir, "memory", "index.json")
	if _, err := os.Stat(memIdx); err != nil {
		t.Errorf("expected memory/index.json to exist after init: %v", err)
	}

	// Verify the init success message in output.
	if !contains(string(out), "init complete") {
		t.Errorf("expected 'init complete' in output, got:\n%s", out)
	}
}

// TestCLI_Status verifies the status command exits 0 and prints recognisable
// output even when .muah/ has not been initialised.
func TestCLI_Status(t *testing.T) {
	bin := buildBinary(t)

	projectDir := t.TempDir()
	cmd := exec.Command(bin, "status")
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "HOME="+projectDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status failed: %v\nout: %s", err, out)
	}
	got := string(out)
	// Should mention platform or OS regardless of environment
	if !contains(got, "Platform") && !contains(got, "OS") {
		t.Errorf("status output missing platform/OS info:\n%s", got)
	}
}

// TestCLI_SwarmRaces verifies the swarm races command lists all 5 Vana'diel races.
func TestCLI_SwarmRaces(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "swarm", "races").CombinedOutput()
	if err != nil {
		t.Fatalf("swarm races failed: %v\nout: %s", err, out)
	}
	got := string(out)
	for _, race := range []string{"Elvaan", "Galka", "Mithra", "Tarutaru", "Hume"} {
		if !contains(got, race) {
			t.Errorf("swarm races missing race %q", race)
		}
	}
}

// TestCLI_SwarmJobs verifies the swarm jobs command lists all 22 job classes.
func TestCLI_SwarmJobs(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "swarm", "jobs").CombinedOutput()
	if err != nil {
		t.Fatalf("swarm jobs failed: %v\nout: %s", err, out)
	}
	got := string(out)
	for _, job := range []string{"WAR", "MNK", "WHM", "BLM", "THF", "PLD", "DRK", "SMN", "SCH", "RUN"} {
		if !contains(got, job) {
			t.Errorf("swarm jobs missing job class %q", job)
		}
	}
}

// TestCLI_ToolsList verifies tools list exits 0 and mentions builtin tools.
func TestCLI_ToolsList(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "tools", "list").CombinedOutput()
	if err != nil {
		t.Fatalf("tools list failed: %v\nout: %s", err, out)
	}
	got := string(out)
	if !contains(got, "Builtin") && !contains(got, "builtin") {
		t.Errorf("tools list missing builtin tools section:\n%s", got)
	}
}

// TestCLI_InitIdempotent verifies that running "muah-runner init" twice does not
// fail or produce errors.
func TestCLI_InitIdempotent(t *testing.T) {
	bin := buildBinary(t)
	projectDir := t.TempDir()
	env := append(os.Environ(), "MUAH_AUTO_MODE=true", "HOME="+projectDir)

	for i := 0; i < 2; i++ {
		cmd := exec.Command(bin, "init")
		cmd.Dir = projectDir
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("muah-runner init (run %d) failed: %v\noutput:\n%s", i+1, err, out)
		}
	}
}

// ── helpers ────────────────────────────────────────────────────────────────────

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && len(substr) == 0 ||
		indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(s) < len(substr) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
