package tools

import "fmt"

type Tool struct {
Name        string
Description string
Category    string
}

type Registry struct {
tools map[string]Tool
}

func NewRegistry() *Registry {
r := &Registry{tools: make(map[string]Tool)}
r.registerBuiltins()
return r
}

func (r *Registry) registerBuiltins() {
builtins := []Tool{
{Name: "search", Description: "Web search", Category: "builtin"},
{Name: "fetch", Description: "HTTP fetch", Category: "builtin"},
{Name: "file_read", Description: "Read file", Category: "builtin"},
{Name: "file_write", Description: "Write file", Category: "builtin"},
{Name: "shell", Description: "Run shell command", Category: "builtin"},
{Name: "git_commit", Description: "Git commit", Category: "builtin"},
{Name: "code_exec", Description: "Execute code", Category: "builtin"},
}
for _, t := range builtins {
r.tools[t.Name] = t
}
}

func (r *Registry) Get(name string) (Tool, error) {
t, ok := r.tools[name]
if !ok {
return Tool{}, fmt.Errorf("tool not found: %s", name)
}
return t, nil
}

func (r *Registry) All() []Tool {
var list []Tool
for _, t := range r.tools {
list = append(list, t)
}
return list
}
