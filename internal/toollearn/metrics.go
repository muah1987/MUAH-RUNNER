package toollearn

import (
"encoding/json"
"os"
"path/filepath"
)

type Dashboard struct {
ToolStats map[string]ToolScore `json:"tool_stats"`
AIStats   map[string]float64   `json:"ai_stats"`
}

func BuildDashboard(muahDir string) (*Dashboard, error) {
profiler := NewProfiler(muahDir)
dash := &Dashboard{
ToolStats: make(map[string]ToolScore),
AIStats:   make(map[string]float64),
}
for name, profile := range profiler.profiles {
dash.AIStats[name] = profile.OverallScore
for tool, ts := range profile.ToolScores {
existing := dash.ToolStats[tool]
existing.Calls += ts.Calls
existing.Successes += ts.Successes
existing.Failures += ts.Failures
if existing.Calls > 0 {
existing.Score = float64(existing.Successes) / float64(existing.Calls) * 100
}
dash.ToolStats[tool] = existing
}
}
return dash, nil
}

func SaveDashboard(muahDir string, dash *Dashboard) error {
path := filepath.Join(muahDir, "tool_use", "metrics", "dashboard.json")
os.MkdirAll(filepath.Dir(path), 0755)
data, err := json.MarshalIndent(dash, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
