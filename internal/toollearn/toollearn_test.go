package toollearn

import "testing"

func TestClassifyError(t *testing.T) {
tests := []struct {
err   string
class ErrorClass
}{
{"401 unauthorized", AuthFailure},
{"rate limit exceeded 429", RateLimit},
{"404 not found", NotFound},
{"connection timeout", Timeout},
{"json unmarshal error", ParseError},
{"some random error", UnknownError},
}
for _, tt := range tests {
got := ClassifyError(tt.err)
if got != tt.class {
t.Errorf("ClassifyError(%q) = %s, want %s", tt.err, got, tt.class)
}
}
}

func TestProfilerSuccessFailure(t *testing.T) {
dir := t.TempDir()
p := NewProfiler(dir)
p.RecordSuccess("claude", "git_commit")
p.RecordSuccess("claude", "git_commit")
p.RecordFailure("claude", "git_commit", "permission denied")

profile := p.GetOrCreate("claude")
ts := profile.ToolScores["git_commit"]
if ts.Calls != 3 {
t.Errorf("expected 3 calls, got %d", ts.Calls)
}
if ts.Successes != 2 {
t.Errorf("expected 2 successes, got %d", ts.Successes)
}
}

func TestMonitorRecovery(t *testing.T) {
dir := t.TempDir()
m := NewMonitor(dir)
for i := 0; i < ConsecutiveFailureThreshold; i++ {
m.RecordUse("gemini", "fetch", false, "timeout error")
}
if !m.ShouldRecover("gemini", "fetch") {
t.Error("should trigger recovery after threshold failures")
}
}

func TestRegistry(t *testing.T) {
dir := t.TempDir()
r := NewRegistry(dir)
r.Register(ToolEntry{Name: "git", Description: "Git operations", Enabled: true})
entry, ok := r.Get("git")
if !ok {
t.Fatal("expected to find git tool")
}
if entry.Description != "Git operations" {
t.Errorf("wrong description: %s", entry.Description)
}
}
