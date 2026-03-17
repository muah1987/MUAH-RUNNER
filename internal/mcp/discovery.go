package mcp

import (
"encoding/json"
"os"
"path/filepath"
"time"
)

type MCPTool struct {
Name        string   `json:"name"`
Description string   `json:"description"`
InputSchema map[string]interface{} `json:"input_schema,omitempty"`
}

type MCPServer struct {
Name      string    `json:"name"`
URL       string    `json:"url"`
Type      string    `json:"type"` // platform | native | discovered
Status    string    `json:"status"` // connected | disconnected | unavailable
Tools     []MCPTool `json:"tools"`
LastSeen  time.Time `json:"last_seen"`
}

type Discovery struct {
muahDir string
servers []MCPServer
}

func NewDiscovery(muahDir string) *Discovery {
d := &Discovery{muahDir: muahDir}
d.load()
return d
}

func (d *Discovery) discoveryPath() string {
return filepath.Join(d.muahDir, "mcp", "discovery.json")
}

func (d *Discovery) load() {
data, err := os.ReadFile(d.discoveryPath())
if err != nil {
return
}
json.Unmarshal(data, &d.servers)
}

func (d *Discovery) save() error {
os.MkdirAll(filepath.Dir(d.discoveryPath()), 0755)
data, err := json.MarshalIndent(d.servers, "", "  ")
if err != nil {
return err
}
return os.WriteFile(d.discoveryPath(), data, 0644)
}

func (d *Discovery) Discover() ([]MCPServer, error) {
// Auto-detect known MCPs
candidates := []MCPServer{
{Name: "github-mcp", Type: "platform", URL: "mcp://github"},
{Name: "playwright-mcp", Type: "native", URL: "mcp://playwright"},
{Name: "docker-mcp", Type: "native", URL: "mcp://docker"},
{Name: "filesystem-mcp", Type: "native", URL: "mcp://filesystem"},
{Name: "sqlite-mcp", Type: "native", URL: "mcp://sqlite"},
}

for i := range candidates {
candidates[i].Status = "disconnected"
candidates[i].LastSeen = time.Now()
}
d.servers = candidates
d.save()
return d.servers, nil
}

func (d *Discovery) GetAll() []MCPServer {
return d.servers
}

func (d *Discovery) GetByType(serverType string) []MCPServer {
var result []MCPServer
for _, s := range d.servers {
if s.Type == serverType {
result = append(result, s)
}
}
return result
}
