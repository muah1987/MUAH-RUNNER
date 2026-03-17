package changelog

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"sync"
"time"
)

type Category string

const (
Added      Category = "added"
Changed    Category = "changed"
Fixed      Category = "fixed"
Removed    Category = "removed"
Security   Category = "security"
Deprecated Category = "deprecated"
)

type Entry struct {
ID        string    `json:"id"`
Timestamp time.Time `json:"timestamp"`
Category  Category  `json:"category"`
Summary   string    `json:"summary"`
Detail    string    `json:"detail,omitempty"`
RunID     string    `json:"run_id,omitempty"`
PrevHash  string    `json:"prev_hash"`
Hash      string    `json:"hash"`
}

type Writer struct {
muahDir  string
mu       sync.Mutex
lastHash string
}

func NewWriter(muahDir string) *Writer {
w := &Writer{muahDir: muahDir}
w.lastHash = w.loadLastHash()
return w
}

func (w *Writer) entriesDir() string {
return filepath.Join(w.muahDir, "changelog", "entries")
}

func (w *Writer) changelogPath() string {
return filepath.Join(w.muahDir, "changelog", "CHANGELOG.md")
}

func (w *Writer) loadLastHash() string {
dir := w.entriesDir()
entries, err := os.ReadDir(dir)
if err != nil || len(entries) == 0 {
return ""
}
// get last entry
last := entries[len(entries)-1]
data, err := os.ReadFile(filepath.Join(dir, last.Name()))
if err != nil {
return ""
}
var e Entry
if err := json.Unmarshal(data, &e); err != nil {
return ""
}
return e.Hash
}

func (w *Writer) Append(category Category, summary, detail, runID string) (*Entry, error) {
w.mu.Lock()
defer w.mu.Unlock()

if err := os.MkdirAll(w.entriesDir(), 0755); err != nil {
return nil, err
}

e := &Entry{
ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
Timestamp: time.Now(),
Category:  category,
Summary:   summary,
Detail:    detail,
RunID:     runID,
PrevHash:  w.lastHash,
}
e.Hash = ChainHash(w.lastHash, e)
w.lastHash = e.Hash

data, err := json.MarshalIndent(e, "", "  ")
if err != nil {
return nil, err
}
filename := fmt.Sprintf("%s-%s.json", e.Timestamp.Format("20060102-150405"), slugify(summary))
if err := os.WriteFile(filepath.Join(w.entriesDir(), filename), data, 0644); err != nil {
return nil, err
}
return e, w.appendMarkdown(e)
}

func slugify(s string) string {
if len(s) > 40 {
s = s[:40]
}
result := make([]byte, 0, len(s))
for _, c := range s {
if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
result = append(result, byte(c))
} else {
result = append(result, '-')
}
}
return string(result)
}

func (w *Writer) appendMarkdown(e *Entry) error {
f, err := os.OpenFile(w.changelogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
fmt.Fprintf(f, "\n## [%s] %s - %s\n%s\n", e.Category, e.Timestamp.Format("2006-01-02"), e.Summary, e.Detail)
return nil
}
