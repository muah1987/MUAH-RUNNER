package sensor

import "fmt"

type HealthStatus struct {
OK      bool
Details map[string]string
}

func CheckHealth() HealthStatus {
status := HealthStatus{
OK:      true,
Details: make(map[string]string),
}
env := DetectEnvironment()
tc := DetectToolchain()
status.Details["platform"] = env.Platform
status.Details["os"] = env.OS
if tc.Go == "" {
status.OK = false
status.Details["go"] = "not found"
} else {
status.Details["go"] = "ok"
}
if tc.Docker != "" {
status.Details["docker"] = "ok"
} else {
status.Details["docker"] = "not found"
}
return status
}

func (h HealthStatus) String() string {
s := "Health: "
if h.OK {
s += "OK\n"
} else {
s += "DEGRADED\n"
}
for k, v := range h.Details {
s += fmt.Sprintf("  %s: %s\n", k, v)
}
return s
}
