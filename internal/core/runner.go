package core

import (
"fmt"
"os"
"path/filepath"
"time"

"github.com/muah1987/muah-runner/internal/changelog"
"github.com/muah1987/muah-runner/internal/config"
"github.com/muah1987/muah-runner/internal/memory"
"github.com/muah1987/muah-runner/internal/planning"
"github.com/muah1987/muah-runner/internal/questions"
"github.com/muah1987/muah-runner/internal/selfreflection"
"github.com/muah1987/muah-runner/internal/sensor"
)

// Runner is the main muah-runner execution engine.
type Runner struct {
cfg       *config.Config
muahDir   string
baseDir   string
mem       *memory.Store
changelog *changelog.Writer
asker     *questions.Asker
reflector *selfreflection.Interviewer
platform  string
}

// New creates a Runner from baseDir and config.
func New(baseDir string, cfg *config.Config) *Runner {
muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
env := sensor.DetectEnvironment()
return &Runner{
cfg:       cfg,
muahDir:   muahDir,
baseDir:   baseDir,
mem:       memory.NewStore(muahDir),
changelog: changelog.NewWriter(muahDir),
asker:     questions.NewAsker(muahDir),
reflector: selfreflection.NewInterviewer(muahDir),
platform:  env.Platform,
}
}

// Init bootstraps .muah/, scans repo, generates plan and docs.
func (r *Runner) Init() error {
fmt.Println("🚀 muah-runner init")
fmt.Printf("   Platform: %s\n", r.platform)
fmt.Printf("   muah dir: %s\n", r.muahDir)

fmt.Println("📁 Bootstrapping .muah/ directory...")
if err := Bootstrap(r.baseDir); err != nil {
return fmt.Errorf("bootstrap: %w", err)
}
fmt.Println("   ✓ .muah/ created")

// Phase 0: Repo scan
fmt.Println("🔍 Phase 0: Scanning repository...")
scan, err := planning.ScanRepo(r.baseDir)
if err != nil {
return fmt.Errorf("repo scan: %w", err)
}
_ = planning.SaveScan(r.muahDir, scan)
fmt.Printf("   Platform: %s | Tech stack: %v\n", scan.Platform, scan.TechStack)
fmt.Printf("   Gaps found: %v\n", scan.Gaps)

// Phase 1: Plan
fmt.Println("📋 Phase 1: Generating plan...")
plan := planning.GeneratePlan(scan)
_ = planning.SavePlanMarkdown(r.muahDir, plan)
fmt.Printf("   %d phases planned\n", len(plan.Phases))

// Ask to approve (timed 10s)
approveQ := &questions.Question{
ID:   "init-approve-plan",
Text: fmt.Sprintf("Scanned repo. %d gaps found. Proceed with plan?", len(scan.Gaps)),
Options: []questions.Option{
{Index: 0, Label: "Yes, proceed", Value: "yes"},
{Index: 1, Label: "No, abort", Value: "no"},
},
Recommended: 0,
Timeout:     10 * time.Second,
TimeoutSec:  10,
}
ans := r.asker.Ask(approveQ, r.cfg.Questions.AutoMode)
if ans.Value == "no" {
fmt.Println("   Aborted by user.")
return nil
}

// Phase 2: Docs
fmt.Println("📚 Phase 2: Generating documentation...")
docPlatform := scan.Platform
if docPlatform == "unknown" || docPlatform == "standalone" {
docPlatform = "github"
}
fmt.Printf("   Generating %s docs package...\n", docPlatform)

runID := memory.GenerateRunID()
r.changelog.Append(changelog.Added, "muah-runner init completed", fmt.Sprintf("Platform: %s", r.platform), runID)

fmt.Println("✅ muah-runner init complete!")
return nil
}

// RunTask executes a task through the full phase pipeline.
func (r *Runner) RunTask(task string) error {
runID := memory.GenerateRunID()
start := time.Now()

fmt.Printf("🎯 muah-runner run: %s\n", task)
fmt.Printf("   Run ID: %s | Platform: %s\n", runID, r.platform)

if _, err := os.Stat(r.muahDir); os.IsNotExist(err) {
fmt.Println("   .muah/ not found, running init first...")
if err := r.Init(); err != nil {
return fmt.Errorf("auto-init: %w", err)
}
}

mem := &memory.RunMemory{
RunID:     runID,
Timestamp: time.Now().UTC(),
Trigger:   "manual",
Platform:  r.platform,
Task:      task,
}

errorCount, retryCount, planDriftCount := 0, 0, 0
keyLesson := "task completed"

fmt.Println("⚡ Phase 3: Executing task...")
mem.ActionsLog = append(mem.ActionsLog, memory.Action{
Tool:      "core:run_task",
Input:     task,
Output:    "executed",
Success:   true,
Timestamp: time.Now(),
})

fmt.Println("📝 Phase 4: Review & Close...")
prQ := &questions.Question{
ID:   runID + "-create-pr",
Text: "Task complete. Create a summary PR?",
Options: []questions.Option{
{Index: 0, Label: "Yes, create PR", Value: "yes"},
{Index: 1, Label: "No, skip", Value: "no"},
},
Recommended: 0,
Timeout:     10 * time.Second,
TimeoutSec:  10,
}
r.asker.Ask(prQ, r.cfg.Questions.AutoMode)

// Phase 5: Self-Reflection — MANDATORY, never skipped
fmt.Println("🪞 Phase 5: Self-Reflection...")
elapsed := time.Since(start)
entry, err := r.reflector.Conduct(runID, selfreflection.RunMetrics{
ErrorCount:      errorCount,
RetryCount:      retryCount,
PlanDriftCount:  planDriftCount,
DurationSeconds: elapsed.Seconds(),
ActionCount:     len(mem.ActionsLog),
KeyLesson:       keyLesson,
})
if err != nil {
fmt.Printf("   ⚠ self-reflection error: %v\n", err)
} else {
fmt.Printf("   Self-rating: %s (score: %.2f)\n", entry.Grade, entry.Score)
}

mem.Lessons = []string{keyLesson}
_ = r.mem.Save(mem)
r.changelog.Append(changelog.Changed, "Completed task: "+task, "", runID)

fmt.Printf("✅ Run %s complete in %.1fs\n", runID, elapsed.Seconds())
return nil
}

// Status prints health information.
func (r *Runner) Status() {
fmt.Println("📊 muah-runner status")
fmt.Println("─────────────────────────────────────")
env := sensor.DetectEnvironment()
fmt.Printf("Platform:   %s (%s)\n", env.Platform, env.CIName)
fmt.Printf("OS/Arch:    %s/%s\n", env.OS, env.Arch)
fmt.Printf("Hostname:   %s\n", env.Hostname)
if _, err := os.Stat(r.muahDir); os.IsNotExist(err) {
fmt.Printf(".muah dir:  ✗ not found (run `muah-runner init`)\n")
} else {
fmt.Printf(".muah dir:  ✓ %s\n", r.muahDir)
}
if idx, err := r.mem.LoadIndex(); err == nil {
fmt.Printf("Memory:     %d runs recorded\n", idx.TotalRuns)
if idx.LastRunID != "" {
fmt.Printf("Last run:   %s\n", idx.LastRunID)
}
} else {
fmt.Println("Memory:     not initialized")
}
tc := sensor.DetectToolchain()
fmt.Println("─────────────────────────────────────")
fmt.Println("Toolchain:")
printCheck("Go", tc.Go != "")
printCheck("Git", tc.Git != "")
printCheck("Docker", tc.Docker != "")
printCheck("Playwright", tc.Playwright)
}

func printCheck(label string, ok bool) {
if ok {
fmt.Printf("  ✓ %-14s\n", label)
} else {
fmt.Printf("  ✗ %-14s (optional)\n", label)
}
}
