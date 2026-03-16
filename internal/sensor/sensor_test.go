package sensor

import "testing"

func TestDetectEnvironment(t *testing.T) {
env := DetectEnvironment()
if env.OS == "" {
t.Error("expected non-empty OS")
}
if env.Platform == "" {
t.Error("expected non-empty platform")
}
}

func TestDetectToolchain(t *testing.T) {
tc := DetectToolchain()
if tc.Go == "" {
t.Error("expected Go to be detected in test environment")
}
}

func TestCheckHealth(t *testing.T) {
h := CheckHealth()
if h.Details == nil {
t.Error("expected details")
}
}

func TestScanMCPs(t *testing.T) {
r := ScanMCPs()
total := len(r.Available) + len(r.Missing)
if total == 0 {
t.Error("expected some MCP scan results")
}
}
