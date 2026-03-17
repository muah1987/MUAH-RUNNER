package tools

import "testing"

func TestRegistry(t *testing.T) {
r := NewRegistry()
all := r.All()
if len(all) == 0 {
t.Error("expected builtin tools")
}
_, err := r.Get("shell")
if err != nil {
t.Fatalf("expected shell tool: %v", err)
}
}
