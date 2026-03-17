package mcp

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type Connection struct {
Server    MCPServer `json:"server"`
ConnectedAt time.Time `json:"connected_at"`
Active    bool      `json:"active"`
}

type ConnectionManager struct {
muahDir     string
connections map[string]*Connection
}

func NewConnectionManager(muahDir string) *ConnectionManager {
cm := &ConnectionManager{
muahDir:     muahDir,
connections: make(map[string]*Connection),
}
cm.load()
return cm
}

func (cm *ConnectionManager) connectionsPath() string {
return filepath.Join(cm.muahDir, "mcp", "connections.json")
}

func (cm *ConnectionManager) load() {
data, err := os.ReadFile(cm.connectionsPath())
if err != nil {
return
}
json.Unmarshal(data, &cm.connections)
}

func (cm *ConnectionManager) save() {
os.MkdirAll(filepath.Dir(cm.connectionsPath()), 0755)
data, _ := json.MarshalIndent(cm.connections, "", "  ")
os.WriteFile(cm.connectionsPath(), data, 0644)
}

func (cm *ConnectionManager) Connect(server MCPServer) (*Connection, error) {
conn := &Connection{
Server:      server,
ConnectedAt: time.Now(),
Active:      true,
}
cm.connections[server.Name] = conn
cm.save()
return conn, nil
}

func (cm *ConnectionManager) Disconnect(serverName string) {
if conn, ok := cm.connections[serverName]; ok {
conn.Active = false
cm.save()
}
}

func (cm *ConnectionManager) IsConnected(serverName string) bool {
conn, ok := cm.connections[serverName]
return ok && conn.Active
}

func (cm *ConnectionManager) ActiveConnections() []*Connection {
var active []*Connection
for _, c := range cm.connections {
if c.Active {
active = append(active, c)
}
}
return active
}

func (cm *ConnectionManager) Test(serverName string) error {
if !cm.IsConnected(serverName) {
return fmt.Errorf("server %s is not connected", serverName)
}
return nil
}
