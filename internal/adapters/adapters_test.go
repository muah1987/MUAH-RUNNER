package adapters

import "testing"

func TestAdapters(t *testing.T) {
gh := NewGitHubAdapter("fake-token")
if !gh.IsAvailable() {
t.Error("expected available with token")
}
gen := NewGenericAdapter("test")
if !gen.IsAvailable() {
t.Error("generic always available")
}
}
