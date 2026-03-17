package docs

import "testing"

func TestDetectPlatform(t *testing.T) {
platform := DetectPlatform("../../")
if platform == "" {
t.Error("expected non-empty platform")
}
}

func TestGenerator(t *testing.T) {
dir := t.TempDir()
g := NewGenerator(dir, "../../")
manifest, err := g.Generate()
if err != nil {
t.Fatalf("generate error: %v", err)
}
if len(manifest.Docs) == 0 {
t.Error("expected docs in manifest")
}
}

func TestAudit(t *testing.T) {
dir := t.TempDir()
result, err := Audit(dir, "../../")
if err != nil {
t.Fatalf("audit error: %v", err)
}
if result.Summary == "" {
t.Error("expected non-empty summary")
}
}
