package docker

import (
"fmt"
"os/exec"
"strings"
)

type Engine struct{}

func NewEngine() *Engine {
return &Engine{}
}

func (e *Engine) IsAvailable() bool {
_, err := exec.LookPath("docker")
return err == nil
}

func (e *Engine) Run(image string, cmd []string, opts RunOptions) (string, error) {
args := []string{"run", "--rm"}
if opts.Name != "" {
args = append(args, "--name", opts.Name)
}
for k, v := range opts.Env {
args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
}
for _, vol := range opts.Volumes {
args = append(args, "-v", vol)
}
if opts.WorkDir != "" {
args = append(args, "-w", opts.WorkDir)
}
args = append(args, image)
args = append(args, cmd...)
out, err := exec.Command("docker", args...).CombinedOutput()
return string(out), err
}

func (e *Engine) Logs(containerID string) (string, error) {
out, err := exec.Command("docker", "logs", containerID).CombinedOutput()
return string(out), err
}

func (e *Engine) Stop(containerID string) error {
return exec.Command("docker", "stop", containerID).Run()
}

func (e *Engine) Remove(containerID string) error {
return exec.Command("docker", "rm", "-f", containerID).Run()
}

func (e *Engine) ImageExists(image string) bool {
out, err := exec.Command("docker", "images", "-q", image).Output()
return err == nil && strings.TrimSpace(string(out)) != ""
}

type RunOptions struct {
Name    string
Env     map[string]string
Volumes []string
WorkDir string
Detach  bool
}
