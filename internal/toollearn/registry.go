package toollearn

import (
"encoding/json"
"os"
"path/filepath"
)

type ToolEntry struct {
Name        string   `json:"name"`
Description string   `json:"description"`
Commands    []string `json:"commands"`
Category    string   `json:"category"`
Enabled     bool     `json:"enabled"`
}

type Registry struct {
muahDir string
tools   map[string]ToolEntry
}

func NewRegistry(muahDir string) *Registry {
r := &Registry{
muahDir: muahDir,
tools:   make(map[string]ToolEntry),
}
r.load()
return r
}

func (r *Registry) registryPath() string {
return filepath.Join(r.muahDir, "tool_use", "registry.json")
}

func (r *Registry) load() {
data, err := os.ReadFile(r.registryPath())
if err != nil {
return
}
json.Unmarshal(data, &r.tools)
}

func (r *Registry) Save() error {
os.MkdirAll(filepath.Dir(r.registryPath()), 0755)
data, err := json.MarshalIndent(r.tools, "", "  ")
if err != nil {
return err
}
return os.WriteFile(r.registryPath(), data, 0644)
}

func (r *Registry) Register(t ToolEntry) {
r.tools[t.Name] = t
}

func (r *Registry) Get(name string) (ToolEntry, bool) {
t, ok := r.tools[name]
return t, ok
}

func (r *Registry) All() []ToolEntry {
var list []ToolEntry
for _, t := range r.tools {
list = append(list, t)
}
return list
}
