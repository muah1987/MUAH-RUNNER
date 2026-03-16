package planning

import "testing"

func TestScanRepo(t *testing.T) {
scan, err := ScanRepo("../../")
if err != nil {
t.Fatalf("scan error: %v", err)
}
if scan.Platform == "" {
t.Error("expected non-empty platform")
}
// Should detect Go in this repo
found := false
for _, lang := range scan.Languages {
if lang == "Go" {
found = true
}
}
if !found {
t.Error("expected Go to be detected")
}
}

func TestGeneratePlan(t *testing.T) {
scan := &RepoScan{
Platform:  "github",
Languages: []string{"Go"},
Gaps:      []string{"Missing README", "Missing LICENSE"},
}
plan := GeneratePlan(scan)
if len(plan.Phases) == 0 {
t.Error("expected phases")
}
// Should have doc phase
hasDocPhase := false
for _, p := range plan.Phases {
if p.Name == "Documentation" {
hasDocPhase = true
}
}
if !hasDocPhase {
t.Error("expected documentation phase for repos with gaps")
}
}

func TestSavePlanMarkdown(t *testing.T) {
dir := t.TempDir()
scan := &RepoScan{Platform: "github", Languages: []string{"Go"}}
plan := GeneratePlan(scan)
if err := SavePlanMarkdown(dir, plan); err != nil {
t.Fatalf("save error: %v", err)
}
}

func TestGenerateRoadmap(t *testing.T) {
scan := &RepoScan{Platform: "github"}
plan := GeneratePlan(scan)
roadmap := GenerateRoadmap(plan)
if len(roadmap.Milestones) == 0 {
t.Error("expected milestones")
}
}
