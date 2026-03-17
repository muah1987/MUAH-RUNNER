package questions

import (
"testing"
)

func TestAutoSelect(t *testing.T) {
dir := t.TempDir()
asker := NewAsker(dir)
q := &Question{
ID:   "test-q",
Text: "Test question?",
Options: []Option{
{0, "Option A", "a"},
{1, "Option B", "b"},
},
Recommended: 1,
TimeoutSec:  10,
}
answer := asker.Ask(q, true)
if !answer.AutoSelected {
t.Error("expected auto-selected")
}
if answer.Value != "b" {
t.Errorf("expected 'b', got %s", answer.Value)
}
}

func TestHistorySaved(t *testing.T) {
dir := t.TempDir()
asker := NewAsker(dir)
q := &Question{
ID:          "q1",
Text:        "Question 1?",
Options:     []Option{{0, "Yes", "yes"}},
Recommended: 0,
}
asker.Ask(q, true)
history := asker.History()
if len(history) != 1 {
t.Errorf("expected 1 history entry, got %d", len(history))
}
}

func TestStandardQuestions(t *testing.T) {
qs := StandardQuestions()
if len(qs) == 0 {
t.Error("expected standard questions")
}
}
