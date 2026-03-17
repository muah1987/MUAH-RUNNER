package docker

import (
"fmt"
"os/exec"
)

type Cleanup struct {
engine *Engine
}

func NewCleanup() *Cleanup {
return &Cleanup{engine: NewEngine()}
}

func (c *Cleanup) PruneContainers() (string, error) {
if !c.engine.IsAvailable() {
return "", fmt.Errorf("docker not available")
}
out, err := exec.Command("docker", "container", "prune", "-f").Output()
return string(out), err
}

func (c *Cleanup) PruneImages() (string, error) {
if !c.engine.IsAvailable() {
return "", fmt.Errorf("docker not available")
}
out, err := exec.Command("docker", "image", "prune", "-f").Output()
return string(out), err
}

func (c *Cleanup) PruneAll() (string, error) {
if !c.engine.IsAvailable() {
return "", fmt.Errorf("docker not available")
}
out, err := exec.Command("docker", "system", "prune", "-f").Output()
return string(out), err
}
