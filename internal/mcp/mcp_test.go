package mcp

import "testing"

func TestDiscovery(t *testing.T) {
dir := t.TempDir()
d := NewDiscovery(dir)
servers, err := d.Discover()
if err != nil {
t.Fatalf("discover error: %v", err)
}
if len(servers) == 0 {
t.Error("expected discovered servers")
}
}

func TestConnectionManager(t *testing.T) {
dir := t.TempDir()
cm := NewConnectionManager(dir)
server := MCPServer{Name: "test-mcp", URL: "mcp://test", Type: "native"}
conn, err := cm.Connect(server)
if err != nil {
t.Fatalf("connect error: %v", err)
}
if !conn.Active {
t.Error("expected active connection")
}
if !cm.IsConnected("test-mcp") {
t.Error("expected to be connected")
}
cm.Disconnect("test-mcp")
if cm.IsConnected("test-mcp") {
t.Error("expected disconnected after disconnect")
}
}

func TestRegistry(t *testing.T) {
r := NewMCPRegistry()
server := MCPServer{
Name: "test-mcp",
Tools: []MCPTool{{Name: "do_thing", Description: "Does a thing"}},
}
r.RegisterServer(server)
tool, err := r.GetTool("test-mcp", "do_thing")
if err != nil {
t.Fatalf("get tool error: %v", err)
}
if tool.Name != "do_thing" {
t.Errorf("wrong tool name: %s", tool.Name)
}
}
