package docker

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type SimResult struct {
Success   bool   `json:"success"`
Output    string `json:"output"`
FixWorked bool   `json:"fix_worked"`
}

type Simulation struct {
ID         string    `json:"id"`
Reason     string    `json:"reason"`
Image      string    `json:"image"`
Hypothesis string    `json:"hypothesis"`
Steps      []string  `json:"steps"`
Result     SimResult `json:"result"`
CreatedAt  time.Time `json:"created_at"`
}

type Simulator struct {
muahDir string
engine  *Engine
}

func NewSimulator(muahDir string) *Simulator {
return &Simulator{
muahDir: muahDir,
engine:  NewEngine(),
}
}

func (s *Simulator) simsDir() string {
return filepath.Join(s.muahDir, "docker", "simulations")
}

func (s *Simulator) Run(sim *Simulation) error {
sim.ID = fmt.Sprintf("sim-%s", time.Now().Format("20060102-150405"))
sim.CreatedAt = time.Now()

if !s.engine.IsAvailable() {
sim.Result = SimResult{
Success: false,
Output:  "Docker not available for simulation",
}
} else {
// Execute simulation steps
for _, step := range sim.Steps {
out, err := s.engine.Run(sim.Image, []string{"/bin/sh", "-c", step}, RunOptions{})
if err != nil {
sim.Result = SimResult{
Success: false,
Output:  fmt.Sprintf("Step failed: %s\nOutput: %s", step, out),
}
return s.save(sim)
}
sim.Result.Output += out + "\n"
}
sim.Result.Success = true
sim.Result.FixWorked = true
}

return s.save(sim)
}

func (s *Simulator) save(sim *Simulation) error {
dir := filepath.Join(s.simsDir(), sim.ID)
os.MkdirAll(dir, 0755)
data, err := json.MarshalIndent(sim, "", "  ")
if err != nil {
return err
}
return os.WriteFile(filepath.Join(dir, "result.json"), data, 0644)
}
