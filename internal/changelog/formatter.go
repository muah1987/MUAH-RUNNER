package changelog

import (
"fmt"
"strings"
)

func FormatEntry(e *Entry) string {
var sb strings.Builder
sb.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(string(e.Category)), e.Summary))
if e.Detail != "" {
sb.WriteString(fmt.Sprintf("  %s\n", e.Detail))
}
sb.WriteString(fmt.Sprintf("  Time: %s | Hash: %s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Hash[:8]))
return sb.String()
}
