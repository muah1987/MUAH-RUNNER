package toollearn

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

const ConsecutiveFailureThreshold = 3

type RecoveryStrategy struct {
ErrorClass  ErrorClass `json:"error_class"`
Strategy    string     `json:"strategy"`
Steps       []string   `json:"steps"`
}

var DefaultStrategies = map[ErrorClass]RecoveryStrategy{
AuthFailure: {
ErrorClass: AuthFailure,
Strategy:   "re-authenticate",
Steps:      []string{"Check token expiry", "Refresh credentials", "Re-initialize MCP connection"},
},
RateLimit: {
ErrorClass: RateLimit,
Strategy:   "backoff-retry",
Steps:      []string{"Wait 60 seconds", "Reduce request frequency", "Retry with exponential backoff"},
},
NotFound: {
ErrorClass: NotFound,
Strategy:   "verify-and-skip",
Steps:      []string{"Verify resource exists", "Check path/URL", "Skip if not critical"},
},
Timeout: {
ErrorClass: Timeout,
Strategy:   "retry-with-timeout-increase",
Steps:      []string{"Increase timeout by 2x", "Retry operation", "Report if still failing"},
},
}

type FailureRecord struct {
Tool      string     `json:"tool"`
AI        string     `json:"ai"`
Error     string     `json:"error"`
Class     ErrorClass `json:"class"`
RunID     string     `json:"run_id"`
Timestamp time.Time  `json:"timestamp"`
}

func LogFailure(muahDir string, record FailureRecord) error {
dir := filepath.Join(muahDir, "tool_use", "failure_log")
os.MkdirAll(dir, 0755)
filename := fmt.Sprintf("fail-%d.json", time.Now().UnixNano())
data, err := json.MarshalIndent(record, "", "  ")
if err != nil {
return err
}
return os.WriteFile(filepath.Join(dir, filename), data, 0644)
}

func GetRecoveryStrategy(class ErrorClass) RecoveryStrategy {
if s, ok := DefaultStrategies[class]; ok {
return s
}
return RecoveryStrategy{
ErrorClass: UnknownError,
Strategy:   "manual-review",
Steps:      []string{"Review error logs", "Contact support if needed"},
}
}
