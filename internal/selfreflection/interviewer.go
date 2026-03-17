package selfreflection

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

const (
Q1 = "What could I have done better?"
Q2 = "What was I missing?"
Q3 = "What do I need next time?"
Q4 = "Did I follow the plan or drift?"
Q5 = "What should future-me know?"
)

var InterviewQuestions = []string{Q1, Q2, Q3, Q4, Q5}

type ReflectionEntry struct {
RunID       string    `json:"run_id"`
Timestamp   time.Time `json:"timestamp"`
Answers     [5]string `json:"answers"`
Grade       string    `json:"grade"`
Score       float64   `json:"score"`
FollowedPlan bool     `json:"followed_plan"`
}

type Interviewer struct {
muahDir string
}

func NewInterviewer(muahDir string) *Interviewer {
return &Interviewer{muahDir: muahDir}
}

func (i *Interviewer) reflectionsDir() string {
return filepath.Join(i.muahDir, "memory", "selfreflection", "reflections")
}

// Conduct performs an automated self-reflection based on run metrics
func (i *Interviewer) Conduct(runID string, metrics RunMetrics) (*ReflectionEntry, error) {
entry := &ReflectionEntry{
RunID:     runID,
Timestamp: time.Now(),
}

entry.Answers[0] = generateAnswer(Q1, metrics)
entry.Answers[1] = generateAnswer(Q2, metrics)
entry.Answers[2] = generateAnswer(Q3, metrics)
entry.Answers[3] = generateAnswer(Q4, metrics)
entry.Answers[4] = generateAnswer(Q5, metrics)
entry.FollowedPlan = metrics.PlanDriftCount == 0
entry.Score = CalculateScore(metrics)
entry.Grade = ScoreToGrade(entry.Score)

return entry, i.save(entry)
}

func generateAnswer(question string, m RunMetrics) string {
switch question {
case Q1:
if m.ErrorCount > 0 {
return fmt.Sprintf("Could have avoided %d errors with better pre-checks", m.ErrorCount)
}
return "Execution was efficient, minor timing improvements possible"
case Q2:
if len(m.MissingTools) > 0 {
return fmt.Sprintf("Missing tools: %v", m.MissingTools)
}
return "All required tools were available"
case Q3:
if m.RetryCount > 2 {
return "Better retry logic and exponential backoff"
}
return "Cache results of expensive operations"
case Q4:
if m.PlanDriftCount > 0 {
return fmt.Sprintf("Drifted %d times from plan due to unexpected errors", m.PlanDriftCount)
}
return "Followed the plan closely"
case Q5:
return fmt.Sprintf("Run completed in %.1fs with %d actions. Key lesson: %s",
m.DurationSeconds, m.ActionCount, m.KeyLesson)
}
return ""
}

func (i *Interviewer) save(e *ReflectionEntry) error {
if err := os.MkdirAll(i.reflectionsDir(), 0755); err != nil {
return err
}
filename := fmt.Sprintf("reflect-%s.json", e.Timestamp.Format("20060102-150405"))
data, err := json.MarshalIndent(e, "", "  ")
if err != nil {
return err
}
return os.WriteFile(filepath.Join(i.reflectionsDir(), filename), data, 0644)
}
