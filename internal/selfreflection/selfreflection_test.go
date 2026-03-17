package selfreflection

import (
"testing"
)

func TestScoreToGrade(t *testing.T) {
tests := []struct {
score float64
grade string
}{
{100, "A+"},
{95, "A"},
{85, "B"},
{75, "C"},
{65, "D"},
{50, "F"},
}
for _, tt := range tests {
got := ScoreToGrade(tt.score)
if got != tt.grade {
t.Errorf("ScoreToGrade(%.0f) = %s, want %s", tt.score, got, tt.grade)
}
}
}

func TestCalculateScore(t *testing.T) {
m := RunMetrics{
SpeedScore:       80,
AccuracyScore:    90,
EfficiencyScore:  85,
AutonomyScore:    75,
ImprovementScore: 70,
}
score := CalculateScore(m)
if score < 70 || score > 100 {
t.Errorf("expected score in [70,100], got %.2f", score)
}
}

func TestCalculateScoreHeuristics(t *testing.T) {
m := RunMetrics{ErrorCount: 2, RetryCount: 1}
score := CalculateScore(m)
if score <= 0 || score > 100 {
t.Errorf("invalid score: %.2f", score)
}
}

func TestConductInterview(t *testing.T) {
dir := t.TempDir()
interviewer := NewInterviewer(dir)
m := RunMetrics{
DurationSeconds: 45.2,
ActionCount:     10,
ErrorCount:      1,
KeyLesson:       "pre-check dependencies",
}
entry, err := interviewer.Conduct("run-001", m)
if err != nil {
t.Fatalf("conduct error: %v", err)
}
if entry.Grade == "" {
t.Error("expected non-empty grade")
}
for i, a := range entry.Answers {
if a == "" {
t.Errorf("answer %d should not be empty", i)
}
}
}

func TestDetectPatterns(t *testing.T) {
entries := make([]ReflectionEntry, 5)
for i := range entries {
entries[i].FollowedPlan = false
entries[i].Answers[0] = "had errors"
}
patterns := DetectPatterns(entries)
if len(patterns) == 0 {
t.Error("expected patterns detected")
}
}
