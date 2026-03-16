package planning

import (
"encoding/json"
"os"
"path/filepath"
"time"
)

type Milestone struct {
Name      string    `json:"name"`
Target    time.Time `json:"target"`
Status    string    `json:"status"`
Features  []string  `json:"features"`
}

type Roadmap struct {
CreatedAt  time.Time   `json:"created_at"`
Milestones []Milestone `json:"milestones"`
}

func GenerateRoadmap(plan *Plan) *Roadmap {
roadmap := &Roadmap{
CreatedAt: time.Now(),
}
now := time.Now()
roadmap.Milestones = []Milestone{
{
Name:   "v0.1 - Bootstrap",
Target: now.Add(7 * 24 * time.Hour),
Status: "in_progress",
Features: []string{"Core .muah/ init", "Memory system", "Basic CLI"},
},
{
Name:   "v0.2 - Swarm",
Target: now.Add(14 * 24 * time.Hour),
Status: "planned",
Features: []string{"Vana'diel protocol", "Agent spawning", "MCP integration"},
},
{
Name:   "v1.0 - Production",
Target: now.Add(30 * 24 * time.Hour),
Status: "planned",
Features: []string{"Full feature set", "Docker self-heal", "Playwright integration"},
},
}
return roadmap
}

func SaveRoadmap(muahDir string, roadmap *Roadmap) error {
path := filepath.Join(muahDir, "planning", "roadmap.json")
os.MkdirAll(filepath.Dir(path), 0755)
data, err := json.MarshalIndent(roadmap, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
