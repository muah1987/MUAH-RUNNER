package memory

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type Store struct {
muahDir string
}

func NewStore(muahDir string) *Store {
return &Store{muahDir: muahDir}
}

func (s *Store) runsDir() string {
return filepath.Join(s.muahDir, "memory", "runs")
}

func (s *Store) indexPath() string {
return filepath.Join(s.muahDir, "memory", "index.json")
}

func (s *Store) Save(m *RunMemory) error {
if err := os.MkdirAll(s.runsDir(), 0755); err != nil {
return err
}
data, err := json.MarshalIndent(m, "", "  ")
if err != nil {
return err
}
path := filepath.Join(s.runsDir(), m.RunID+".json")
if err := os.WriteFile(path, data, 0644); err != nil {
return err
}
// update latest symlink
latest := filepath.Join(s.runsDir(), "latest.json")
os.Remove(latest)
os.WriteFile(latest, data, 0644)
return s.updateIndex(m)
}

func (s *Store) updateIndex(m *RunMemory) error {
idx, _ := s.LoadIndex()
if idx == nil {
idx = &MemoryIndex{}
}
idx.TotalRuns++
idx.LastRunID = m.RunID
idx.LastUpdate = time.Now()
idx.RunIDs = append(idx.RunIDs, m.RunID)
data, err := json.MarshalIndent(idx, "", "  ")
if err != nil {
return err
}
return os.WriteFile(s.indexPath(), data, 0644)
}

func (s *Store) LoadIndex() (*MemoryIndex, error) {
data, err := os.ReadFile(s.indexPath())
if err != nil {
return nil, err
}
var idx MemoryIndex
if err := json.Unmarshal(data, &idx); err != nil {
return nil, err
}
return &idx, nil
}

func (s *Store) Load(runID string) (*RunMemory, error) {
path := filepath.Join(s.runsDir(), runID+".json")
data, err := os.ReadFile(path)
if err != nil {
return nil, err
}
var m RunMemory
if err := json.Unmarshal(data, &m); err != nil {
return nil, err
}
return &m, nil
}

func (s *Store) LoadAll() ([]*RunMemory, error) {
idx, err := s.LoadIndex()
if err != nil {
return nil, err
}
var runs []*RunMemory
for _, id := range idx.RunIDs {
m, err := s.Load(id)
if err != nil {
continue
}
runs = append(runs, m)
}
return runs, nil
}

func GenerateRunID() string {
return fmt.Sprintf("run-%s", time.Now().Format("20060102-150405"))
}
