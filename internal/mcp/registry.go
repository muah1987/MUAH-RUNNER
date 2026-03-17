package mcp

import "fmt"

// Registry unifies all MCP tool registries
type Registry struct {
tools   map[string]MCPTool
servers map[string]MCPServer
}

func NewMCPRegistry() *Registry {
return &Registry{
tools:   make(map[string]MCPTool),
servers: make(map[string]MCPServer),
}
}

func (r *Registry) RegisterServer(s MCPServer) {
r.servers[s.Name] = s
for _, t := range s.Tools {
r.tools[s.Name+":"+t.Name] = t
}
}

func (r *Registry) GetTool(serverName, toolName string) (MCPTool, error) {
key := serverName + ":" + toolName
t, ok := r.tools[key]
if !ok {
return MCPTool{}, fmt.Errorf("tool %s not found in %s", toolName, serverName)
}
return t, nil
}

func (r *Registry) AllTools() []MCPTool {
var tools []MCPTool
for _, t := range r.tools {
tools = append(tools, t)
}
return tools
}
