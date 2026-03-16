package docker

import (
"fmt"
"os"
"path/filepath"
)

type Sandbox struct {
muahDir string
engine  *Engine
}

func NewSandbox(muahDir string) *Sandbox {
return &Sandbox{
muahDir: muahDir,
engine:  NewEngine(),
}
}

func (s *Sandbox) sandboxDir() string {
return filepath.Join(s.muahDir, "docker", "sandbox")
}

func (s *Sandbox) Init() error {
if err := os.MkdirAll(s.sandboxDir(), 0755); err != nil {
return err
}
// Write default sandbox Dockerfile
dockerfile := `FROM ubuntu:22.04
RUN apt-get update && apt-get install -y curl git build-essential
WORKDIR /workspace
CMD ["/bin/bash"]
`
return os.WriteFile(filepath.Join(s.sandboxDir(), "Dockerfile.sandbox"), []byte(dockerfile), 0644)
}

func (s *Sandbox) Run(image, command string) (string, error) {
if !s.engine.IsAvailable() {
return "", fmt.Errorf("docker not available")
}
return s.engine.Run(image, []string{"/bin/sh", "-c", command}, RunOptions{
WorkDir: "/workspace",
})
}
