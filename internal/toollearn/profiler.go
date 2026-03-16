package toollearn

import (
"encoding/json"
"os"
"path/filepath"
)

type ToolScore struct {
Calls     int     `json:"calls"`
Successes int     `json:"successes"`
Failures  int     `json:"failures"`
Score     float64 `json:"score"`
}

type FailurePattern struct {
ErrorType   string `json:"error_type"`
Count       int    `json:"count"`
LastSeen    string `json:"last_seen"`
Recovery    string `json:"recovery"`
}

type AIProfile struct {
AIName           string               `json:"ai_name"`
OverallScore     float64              `json:"overall_score"`
Strengths        []string             `json:"strengths"`
Weaknesses       []string             `json:"weaknesses"`
ToolScores       map[string]ToolScore `json:"tool_scores"`
FailurePatterns  []FailurePattern     `json:"failure_patterns"`
ImprovementTrend []float64            `json:"improvement_trend"`
}

type Profiler struct {
muahDir  string
profiles map[string]*AIProfile
}

func NewProfiler(muahDir string) *Profiler {
p := &Profiler{
muahDir:  muahDir,
profiles: make(map[string]*AIProfile),
}
p.load()
return p
}

func (p *Profiler) profilesDir() string {
return filepath.Join(p.muahDir, "tool_use", "ai_profiles")
}

func (p *Profiler) load() {
entries, err := os.ReadDir(p.profilesDir())
if err != nil {
return
}
for _, e := range entries {
if e.IsDir() {
continue
}
data, err := os.ReadFile(filepath.Join(p.profilesDir(), e.Name()))
if err != nil {
continue
}
var profile AIProfile
if err := json.Unmarshal(data, &profile); err == nil {
p.profiles[profile.AIName] = &profile
}
}
}

func (p *Profiler) GetOrCreate(aiName string) *AIProfile {
if profile, ok := p.profiles[aiName]; ok {
return profile
}
profile := &AIProfile{
AIName:     aiName,
ToolScores: make(map[string]ToolScore),
}
p.profiles[aiName] = profile
return profile
}

func (p *Profiler) RecordSuccess(aiName, tool string) {
profile := p.GetOrCreate(aiName)
ts := profile.ToolScores[tool]
ts.Calls++
ts.Successes++
ts.Score = float64(ts.Successes) / float64(ts.Calls) * 100
profile.ToolScores[tool] = ts
p.recalcOverall(profile)
}

func (p *Profiler) RecordFailure(aiName, tool, errorType string) {
profile := p.GetOrCreate(aiName)
ts := profile.ToolScores[tool]
ts.Calls++
ts.Failures++
ts.Score = float64(ts.Successes) / float64(ts.Calls) * 100
profile.ToolScores[tool] = ts
p.recalcOverall(profile)
}

func (p *Profiler) recalcOverall(profile *AIProfile) {
if len(profile.ToolScores) == 0 {
return
}
total := 0.0
for _, ts := range profile.ToolScores {
total += ts.Score
}
profile.OverallScore = total / float64(len(profile.ToolScores))
}

func (p *Profiler) Save(aiName string) error {
profile, ok := p.profiles[aiName]
if !ok {
return nil
}
os.MkdirAll(p.profilesDir(), 0755)
data, err := json.MarshalIndent(profile, "", "  ")
if err != nil {
return err
}
return os.WriteFile(filepath.Join(p.profilesDir(), aiName+".json"), data, 0644)
}
