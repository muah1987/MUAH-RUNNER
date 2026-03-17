package swarm

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type Message struct {
From      string    `json:"from"`
To        string    `json:"to"`
Channel   string    `json:"channel"`
Content   string    `json:"content"`
Encoded   bool      `json:"encoded"`
Timestamp time.Time `json:"timestamp"`
}

type Linkshell struct {
muahDir string
}

func NewLinkshell(muahDir string) *Linkshell {
return &Linkshell{muahDir: muahDir}
}

func (ls *Linkshell) Send(msg *Message) error {
dir := filepath.Join(ls.muahDir, "swarm", "linkshell")
if err := os.MkdirAll(dir, 0755); err != nil {
return err
}
var filename string
switch msg.Channel {
case "party":
filename = "party-chat.jsonl"
default:
filename = "main-ls.jsonl"
}
data, err := json.Marshal(msg)
if err != nil {
return err
}
f, err := os.OpenFile(filepath.Join(dir, filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
fmt.Fprintf(f, "%s\n", data)
return nil
}
