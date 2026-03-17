package spec

import (
"encoding/json"
"os"
"path/filepath"
"time"
)

type DriftReport struct {
Timestamp time.Time `json:"timestamp"`
Drifts    []string  `json:"drifts"`
}

func CheckDrift(muahDir string, current *RunnerSpec) (*DriftReport, error) {
report := &DriftReport{Timestamp: time.Now()}
issues := Validate(current)
report.Drifts = issues
data, _ := json.MarshalIndent(report, "", "  ")
path := filepath.Join(muahDir, "spec", "drift-report.json")
os.MkdirAll(filepath.Dir(path), 0755)
os.WriteFile(path, data, 0644)
return report, nil
}
