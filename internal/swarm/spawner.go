package swarm

import "fmt"

type SpawnRequest struct {
NPCName      string
Race         string
Job          string
Task         string
ContextTokens int
OutputTokens  int
}

func Spawn(req SpawnRequest) (*Agent, error) {
if req.ContextTokens == 0 {
req.ContextTokens = 100000
}
if req.OutputTokens == 0 {
req.OutputTokens = 20000
}
agent, err := NewAgent(req.NPCName, req.Race, req.Job, req.ContextTokens, req.OutputTokens)
if err != nil {
return nil, fmt.Errorf("spawn failed: %w", err)
}
agent.Engage(req.Task)
return agent, nil
}
