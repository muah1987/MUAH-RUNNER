package selfreflection

import (
"encoding/json"
"os"
"path/filepath"
)

type Pattern struct {
Name        string `json:"name"`
Description string `json:"description"`
Occurrences int    `json:"occurrences"`
Suggestion  string `json:"suggestion"`
}

func DetectPatterns(entries []ReflectionEntry) []Pattern {
errorCounts := 0
driftCounts := 0
for _, e := range entries {
if !e.FollowedPlan {
driftCounts++
}
// Check if answers mention errors
for _, a := range e.Answers {
if containsWord(a, "error") || containsWord(a, "fail") {
errorCounts++
}
}
}

var patterns []Pattern
if errorCounts > 2 {
patterns = append(patterns, Pattern{
Name:        "recurring-errors",
Description: "Errors appearing frequently across runs",
Occurrences: errorCounts,
Suggestion:  "Add pre-flight checks before task execution",
})
}
if driftCounts > 2 {
patterns = append(patterns, Pattern{
Name:        "plan-drift",
Description: "Runner frequently deviates from plan",
Occurrences: driftCounts,
Suggestion:  "Improve plan specificity and add checkpoints",
})
}
return patterns
}

func containsWord(s, word string) bool {
for i := 0; i+len(word) <= len(s); i++ {
if s[i:i+len(word)] == word {
return true
}
}
return false
}

func SavePatterns(muahDir string, patterns []Pattern) error {
path := filepath.Join(muahDir, "memory", "selfreflection", "patterns.json")
os.MkdirAll(filepath.Dir(path), 0755)
data, err := json.MarshalIndent(patterns, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
