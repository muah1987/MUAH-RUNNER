package memory

import (
"strings"
"time"
)

type SearchFilter struct {
Task     string
Platform string
Status   string
Tool     string
Since    *time.Time
Until    *time.Time
}

func Search(runs []*RunMemory, f SearchFilter) []*RunMemory {
var result []*RunMemory
for _, r := range runs {
if f.Task != "" && !strings.Contains(strings.ToLower(r.Task), strings.ToLower(f.Task)) {
continue
}
if f.Platform != "" && r.Platform != f.Platform {
continue
}
if f.Status != "" && r.Status != f.Status {
continue
}
if f.Tool != "" {
found := false
for _, a := range r.ActionsLog {
if strings.Contains(strings.ToLower(a.Tool), strings.ToLower(f.Tool)) {
found = true
break
}
}
if !found {
continue
}
}
if f.Since != nil && r.Timestamp.Before(*f.Since) {
continue
}
if f.Until != nil && r.Timestamp.After(*f.Until) {
continue
}
result = append(result, r)
}
return result
}
