package swarm

import (
"fmt"
"time"

"github.com/muah1987/muah-runner/internal/privacy"
)

type AgentStatus string

const (
StatusSpawning  AgentStatus = "spawning"
StatusBuffing   AgentStatus = "buffing"
StatusEngaged   AgentStatus = "engaged"
StatusCompleted AgentStatus = "completed"
StatusDismissed AgentStatus = "dismissed"
StatusKO        AgentStatus = "ko"
)

type Agent struct {
UUID      string      `json:"uuid"`
NPCName   string      `json:"npc_name"`
Race      Race        `json:"race"`
Job       Job         `json:"job"`
SubJob    *Job        `json:"sub_job,omitempty"`
HP        HPStats     `json:"hp"`
MP        MPStats     `json:"mp"`
Level     int         `json:"level"`
Status    AgentStatus `json:"status"`
Task      string      `json:"task"`
CreatedAt time.Time   `json:"created_at"`
CoT       *privacy.CoTEntry `json:"cot,omitempty"`
}

func NewAgent(npcName, raceName, jobCode string, contextTokens, outputTokens int) (*Agent, error) {
race, ok := GetRace(raceName)
if !ok {
return nil, fmt.Errorf("unknown race: %s", raceName)
}
job, ok := GetJob(jobCode)
if !ok {
return nil, fmt.Errorf("unknown job: %s", jobCode)
}
return &Agent{
UUID:      fmt.Sprintf("%s-%d", npcName, time.Now().UnixNano()),
NPCName:   npcName,
Race:      race,
Job:       job,
HP:        NewHP(contextTokens),
MP:        NewMP(outputTokens),
Level:     1,
Status:    StatusSpawning,
CreatedAt: time.Now(),
}, nil
}

func (a *Agent) IsAlive() bool {
return a.HP.Current > 0 && a.Status != StatusKO && a.Status != StatusDismissed
}

func (a *Agent) Engage(task string) {
a.Task = task
a.Status = StatusEngaged
}

func (a *Agent) Complete() {
a.Status = StatusCompleted
}

func (a *Agent) KO() {
a.Status = StatusKO
a.HP.Current = 0
}

func (a *Agent) Dismiss() {
a.Status = StatusDismissed
}
