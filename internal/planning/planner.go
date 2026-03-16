package planning

import (
"fmt"
"os"
"path/filepath"
"strings"
"time"
)

type Phase struct {
Number      int      `json:"number"`
Name        string   `json:"name"`
Description string   `json:"description"`
Tasks       []string `json:"tasks"`
Status      string   `json:"status"` // pending, in_progress, complete
}

type Plan struct {
CreatedAt time.Time `json:"created_at"`
RepoScan  *RepoScan `json:"repo_scan"`
Phases    []Phase   `json:"phases"`
Priority  string    `json:"priority"`
}

func GeneratePlan(scan *RepoScan) *Plan {
plan := &Plan{
CreatedAt: time.Now(),
RepoScan:  scan,
Priority:  "normal",
}

// Phase -1: Bootstrap
plan.Phases = append(plan.Phases, Phase{
Number:      -1,
Name:        "Bootstrap",
Description: "Initialize .muah/ directory and core systems",
Tasks:       []string{"Create .muah/ structure", "Initialize memory", "Setup config"},
Status:      "complete",
})

// Phase 0: Assessment
plan.Phases = append(plan.Phases, Phase{
Number:      0,
Name:        "Assessment",
Description: "Scan and understand the repository",
Tasks:       []string{"Scan tech stack", "Identify gaps", "Review existing docs"},
Status:      "complete",
})

// Phase 1: Docs
if len(scan.Gaps) > 0 {
var docTasks []string
for _, gap := range scan.Gaps {
docTasks = append(docTasks, "Fix: "+gap)
}
plan.Phases = append(plan.Phases, Phase{
Number:      1,
Name:        "Documentation",
Description: "Create/update all required documentation",
Tasks:       docTasks,
Status:      "pending",
})
}

// Phase 2: Core Implementation
plan.Phases = append(plan.Phases, Phase{
Number:      2,
Name:        "Core Implementation",
Description: "Implement core features",
Tasks:       []string{"Implement main features", "Write tests", "Setup CI/CD"},
Status:      "pending",
})

// Phase 3: Quality
plan.Phases = append(plan.Phases, Phase{
Number:      3,
Name:        "Quality & Security",
Description: "Ensure quality gates pass",
Tasks:       []string{"Run linting", "Security scan", "Performance testing"},
Status:      "pending",
})

return plan
}

func SavePlanMarkdown(muahDir string, plan *Plan) error {
path := filepath.Join(muahDir, "planning", "current-plan.md")
os.MkdirAll(filepath.Dir(path), 0755)

var sb strings.Builder
sb.WriteString(fmt.Sprintf("# muah-runner Plan\n\nGenerated: %s\n\n", plan.CreatedAt.Format("2006-01-02 15:04:05")))
sb.WriteString(fmt.Sprintf("Platform: %s | Languages: %s\n\n", plan.RepoScan.Platform, strings.Join(plan.RepoScan.Languages, ", ")))

for _, phase := range plan.Phases {
sb.WriteString(fmt.Sprintf("## Phase %d: %s\n", phase.Number, phase.Name))
sb.WriteString(fmt.Sprintf("%s\n\n", phase.Description))
for _, task := range phase.Tasks {
sb.WriteString(fmt.Sprintf("- [ ] %s\n", task))
}
sb.WriteString("\n")
}

return os.WriteFile(path, []byte(sb.String()), 0644)
}
