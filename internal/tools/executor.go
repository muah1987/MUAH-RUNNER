package tools

import "fmt"

// Executor dispatches tool calls to the appropriate backend.
// All tools that interact with the host system require an MCP connection;
// direct execution is not supported to prevent command-injection attacks.
type Executor struct {
registry *Registry
}

func NewExecutor() *Executor {
return &Executor{registry: NewRegistry()}
}

// Execute runs the named tool with the supplied arguments.
// Shell execution and file I/O are intentionally delegated to MCP servers
// rather than performed directly; this prevents command-injection
// vulnerabilities and ensures all side-effects are audited through the MCP
// tool-use pipeline.
func (e *Executor) Execute(name string, args map[string]string) (string, error) {
_, err := e.registry.Get(name)
if err != nil {
return "", err
}
switch name {
case "shell":
// Direct shell execution is disabled to prevent command injection.
// Connect a shell-mcp server and route commands through it instead.
return "", fmt.Errorf("tool %q requires a shell-mcp connection; direct execution is disabled", name)
case "file_read":
importPath := args["path"]
_ = importPath
return "", fmt.Errorf("use filesystem-mcp for file operations")
default:
return "", fmt.Errorf("tool %s requires MCP connection", name)
}
}
