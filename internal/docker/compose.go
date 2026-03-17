package docker

import (
"fmt"
"os/exec"
)

type Compose struct {
composeFile string
}

func NewCompose(composeFile string) *Compose {
return &Compose{composeFile: composeFile}
}

func (c *Compose) Up() error {
return exec.Command("docker", "compose", "-f", c.composeFile, "up", "-d").Run()
}

func (c *Compose) Down() error {
return exec.Command("docker", "compose", "-f", c.composeFile, "down").Run()
}

func (c *Compose) Logs(service string) (string, error) {
args := []string{"compose", "-f", c.composeFile, "logs"}
if service != "" {
args = append(args, service)
}
out, err := exec.Command("docker", args...).Output()
return string(out), err
}

func (c *Compose) IsAvailable() bool {
_, err := exec.LookPath("docker")
if err != nil {
return false
}
out, _ := exec.Command("docker", "compose", "version").Output()
return len(out) > 0
}

func (c *Compose) String() string {
return fmt.Sprintf("Compose[%s]", c.composeFile)
}
