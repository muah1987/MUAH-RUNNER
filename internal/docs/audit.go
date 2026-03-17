package docs

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"strings"
)

type AuditResult struct {
Platform   string      `json:"platform"`
Missing    []DocInfo   `json:"missing"`
Stale      []DocInfo   `json:"stale"`
Present    []DocInfo   `json:"present"`
Score      float64     `json:"score"`
Summary    string      `json:"summary"`
}

func Audit(muahDir, rootDir string) (*AuditResult, error) {
g := NewGenerator(muahDir, rootDir)
manifest, err := g.Generate()
if err != nil {
return nil, err
}

result := &AuditResult{Platform: manifest.Platform}
for _, doc := range manifest.Docs {
switch doc.Status {
case "missing":
result.Missing = append(result.Missing, doc)
case "stale":
result.Stale = append(result.Stale, doc)
case "exists":
result.Present = append(result.Present, doc)
}
}

total := len(manifest.Docs)
if total > 0 {
result.Score = float64(len(result.Present)) / float64(total) * 100
}

var sb strings.Builder
sb.WriteString(fmt.Sprintf("Doc coverage: %.0f%% (%d/%d)", result.Score, len(result.Present), total))
if len(result.Missing) > 0 {
sb.WriteString(fmt.Sprintf(", %d missing", len(result.Missing)))
}
result.Summary = sb.String()

// Save audit
data, err := json.MarshalIndent(result, "", "  ")
if err != nil {
return nil, err
}
auditPath := filepath.Join(muahDir, "docs", "audit.json")
os.MkdirAll(filepath.Dir(auditPath), 0755)
os.WriteFile(auditPath, data, 0644)

return result, nil
}
