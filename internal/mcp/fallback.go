package mcp

import (
"fmt"
"os/exec"
)

// Fallback provides CLI fallback when MCP server is unavailable
type Fallback struct{}

func NewFallback() *Fallback {
return &Fallback{}
}

func (f *Fallback) Execute(serverName, toolName string, args map[string]string) (string, error) {
switch serverName {
case "filesystem-mcp":
return f.filesystemFallback(toolName, args)
case "github-mcp":
return f.githubFallback(toolName, args)
default:
return "", fmt.Errorf("no fallback available for %s", serverName)
}
}

func (f *Fallback) filesystemFallback(tool string, args map[string]string) (string, error) {
switch tool {
case "read_file":
path := args["path"]
out, err := exec.Command("cat", path).Output()
return string(out), err
case "list_directory":
path := args["path"]
out, err := exec.Command("ls", "-la", path).Output()
return string(out), err
}
return "", fmt.Errorf("fallback: unknown filesystem tool %s", tool)
}

func (f *Fallback) githubFallback(tool string, args map[string]string) (string, error) {
switch tool {
case "list_repos":
out, err := exec.Command("gh", "repo", "list").Output()
return string(out), err
}
return "", fmt.Errorf("fallback: unknown github tool %s", tool)
}
