package memory

import "time"

type Action struct {
Tool      string    `json:"tool"`
Input     string    `json:"input"`
Output    string    `json:"output"`
Success   bool      `json:"success"`
Timestamp time.Time `json:"timestamp"`
}

type Reflection struct {
CouldDoBetter  string  `json:"could_do_better"`
WasMissing     string  `json:"was_missing"`
NeedNextTime   string  `json:"need_next_time"`
FollowedPlan   bool    `json:"followed_plan"`
FutureAdvice   string  `json:"future_advice"`
Grade          string  `json:"grade"`
Score          float64 `json:"score"`
}

type MCPToolUse struct {
MCP   string `json:"mcp"`
Tool  string `json:"tool"`
Calls int    `json:"calls"`
}

type MCPConnections struct {
Platform   string       `json:"platform"`
Native     []string     `json:"native"`
Discovered []string     `json:"discovered"`
ToolsUsed  []MCPToolUse `json:"tools_used"`
}

type RunMemory struct {
RunID          string         `json:"run_id"`
Timestamp      time.Time      `json:"timestamp"`
Trigger        string         `json:"trigger"`
Platform       string         `json:"platform"`
Task           string         `json:"task"`
Status         string         `json:"status"`
ContextLoaded  []string       `json:"context_loaded"`
MCPConnections MCPConnections `json:"mcp_connections"`
ActionsLog     []Action       `json:"actions_taken"`
Lessons        []string       `json:"lessons_learned"`
SelfReflection *Reflection    `json:"selfreflection,omitempty"`
Duration       float64        `json:"duration_seconds"`
Tags           []string       `json:"tags"`
}

type MemoryIndex struct {
TotalRuns  int       `json:"total_runs"`
LastRunID  string    `json:"last_run_id"`
LastUpdate time.Time `json:"last_update"`
RunIDs     []string  `json:"run_ids"`
}
