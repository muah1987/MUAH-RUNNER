package tools

import (
"fmt"
"os/exec"
)

type Executor struct {
registry *Registry
}

func NewExecutor() *Executor {
return &Executor{registry: NewRegistry()}
}

func (e *Executor) Execute(name string, args map[string]string) (string, error) {
_, err := e.registry.Get(name)
if err != nil {
return "", err
}
switch name {
case "shell":
cmd := args["command"]
out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
return string(out), err
case "file_read":
import_path := args["path"]
_ = import_path
return "", fmt.Errorf("use filesystem-mcp for file operations")
default:
return "", fmt.Errorf("tool %s requires MCP connection", name)
}
}
