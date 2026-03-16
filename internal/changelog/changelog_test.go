package changelog

import (
"testing"
)

func TestAppendEntry(t *testing.T) {
dir := t.TempDir()
w := NewWriter(dir)
e, err := w.Append(Added, "Added new feature", "Details here", "run-001")
if err != nil {
t.Fatalf("append error: %v", err)
}
if e.Hash == "" {
t.Error("expected non-empty hash")
}
if e.Category != Added {
t.Errorf("expected Added, got %s", e.Category)
}
}

func TestHashChain(t *testing.T) {
dir := t.TempDir()
w := NewWriter(dir)
e1, _ := w.Append(Added, "First entry", "", "run-001")
e2, _ := w.Append(Fixed, "Second entry", "", "run-002")
if e2.PrevHash != e1.Hash {
t.Errorf("chain broken: e2.PrevHash=%s e1.Hash=%s", e2.PrevHash, e1.Hash)
}
}

func TestFormatEntry(t *testing.T) {
dir := t.TempDir()
w := NewWriter(dir)
e, _ := w.Append(Security, "Security fix", "CVE-2026-001", "run-sec")
out := FormatEntry(e)
if len(out) == 0 {
t.Error("expected non-empty formatted output")
}
}
