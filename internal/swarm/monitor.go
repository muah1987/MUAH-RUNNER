package swarm

import "fmt"

type SwarmStatus struct {
TotalAgents int
AliveAgents int
KOAgents    int
Leaders     int
}

func GetStatus(ws *WorldState) SwarmStatus {
all := ws.Party.AllMembers()
status := SwarmStatus{TotalAgents: len(all)}
for _, a := range all {
if a.IsAlive() {
status.AliveAgents++
} else if a.Status == StatusKO {
status.KOAgents++
}
}
if ws.Party.Leader != nil {
status.Leaders = 1
}
return status
}

func (s SwarmStatus) String() string {
return fmt.Sprintf("Agents: %d total, %d alive, %d KO", s.TotalAgents, s.AliveAgents, s.KOAgents)
}
