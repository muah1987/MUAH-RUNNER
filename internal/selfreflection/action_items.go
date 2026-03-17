package selfreflection

import (
"encoding/json"
"os"
"path/filepath"
)

type ActionItem struct {
ID       string `json:"id"`
Source   string `json:"source_run"`
Priority string `json:"priority"`
Text     string `json:"text"`
Done     bool   `json:"done"`
}

func LoadActionItems(muahDir string) ([]ActionItem, error) {
path := filepath.Join(muahDir, "memory", "selfreflection", "action-items.json")
data, err := os.ReadFile(path)
if err != nil {
if os.IsNotExist(err) {
return nil, nil
}
return nil, err
}
var items []ActionItem
if err := json.Unmarshal(data, &items); err != nil {
return nil, err
}
return items, nil
}

func SaveActionItems(muahDir string, items []ActionItem) error {
path := filepath.Join(muahDir, "memory", "selfreflection", "action-items.json")
os.MkdirAll(filepath.Dir(path), 0755)
data, err := json.MarshalIndent(items, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
