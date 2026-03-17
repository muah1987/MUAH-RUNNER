package memory

import (
"testing"
"time"
)

func TestSaveLoad(t *testing.T) {
dir := t.TempDir()
store := NewStore(dir)
m := &RunMemory{
RunID:     "run-test-001",
Timestamp: time.Now(),
Task:      "test task",
Platform:  "github",
Status:    "success",
}
if err := store.Save(m); err != nil {
t.Fatalf("save error: %v", err)
}
loaded, err := store.Load("run-test-001")
if err != nil {
t.Fatalf("load error: %v", err)
}
if loaded.Task != "test task" {
t.Errorf("expected 'test task', got %s", loaded.Task)
}
}

func TestSearch(t *testing.T) {
runs := []*RunMemory{
{Task: "build and test", Platform: "github", Status: "success"},
{Task: "deploy", Platform: "claude", Status: "failure"},
{Task: "lint", Platform: "github", Status: "success"},
}
result := Search(runs, SearchFilter{Platform: "github"})
if len(result) != 2 {
t.Errorf("expected 2, got %d", len(result))
}
result = Search(runs, SearchFilter{Status: "failure"})
if len(result) != 1 {
t.Errorf("expected 1, got %d", len(result))
}
result = Search(runs, SearchFilter{Task: "build"})
if len(result) != 1 {
t.Errorf("expected 1, got %d", len(result))
}
}

func TestMerge(t *testing.T) {
a := []string{"foo", "bar"}
b := []string{"bar", "baz"}
result := MergeLessons(a, b)
if len(result) != 3 {
t.Errorf("expected 3 unique entries, got %d", len(result))
}
}

func TestGenerateRunID(t *testing.T) {
id := GenerateRunID()
if len(id) < 10 {
t.Errorf("run id too short: %s", id)
}
}
