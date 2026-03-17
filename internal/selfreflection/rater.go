package selfreflection

// RunMetrics are inputs for self-reflection scoring
type RunMetrics struct {
DurationSeconds float64
ActionCount     int
ErrorCount      int
RetryCount      int
PlanDriftCount  int
MissingTools    []string
KeyLesson       string
// Scoring weights
SpeedScore       float64 // 0-100
AccuracyScore    float64 // 0-100
EfficiencyScore  float64 // 0-100
AutonomyScore    float64 // 0-100
ImprovementScore float64 // 0-100
}

// CalculateScore computes weighted overall score
// Weights: speed 20%, accuracy 30%, efficiency 20%, autonomy 15%, improvement 15%
func CalculateScore(m RunMetrics) float64 {
speed := m.SpeedScore
accuracy := m.AccuracyScore
efficiency := m.EfficiencyScore
autonomy := m.AutonomyScore
improvement := m.ImprovementScore

// Default heuristics if not provided.
// Only apply the error/retry penalty when the explicit score was not set
// (zero) AND there are actual errors/retries (> 0).  Using >= 0 would always
// trigger even when no errors occurred, overriding a legitimate score of 0.
if accuracy == 0 && m.ErrorCount > 0 {
accuracy = 100 - float64(m.ErrorCount)*10
if accuracy < 0 {
accuracy = 0
}
}
if efficiency == 0 && m.RetryCount > 0 {
efficiency = 100 - float64(m.RetryCount)*5
if efficiency < 0 {
efficiency = 0
}
}
if speed == 0 {
speed = 70 // neutral default
}
if autonomy == 0 {
autonomy = 80
}
if improvement == 0 {
improvement = 75
}

return speed*0.20 + accuracy*0.30 + efficiency*0.20 + autonomy*0.15 + improvement*0.15
}

// ScoreToGrade converts a 0-100 score to a letter grade
func ScoreToGrade(score float64) string {
switch {
case score >= 97:
return "A+"
case score >= 93:
return "A"
case score >= 90:
return "A-"
case score >= 87:
return "B+"
case score >= 83:
return "B"
case score >= 80:
return "B-"
case score >= 77:
return "C+"
case score >= 73:
return "C"
case score >= 70:
return "C-"
case score >= 67:
return "D+"
case score >= 60:
return "D"
default:
return "F"
}
}
