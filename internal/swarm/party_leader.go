package swarm

import (
"sync"
"time"
)

type Party struct {
Leader  *Agent
Members []*Agent
mu      sync.RWMutex
}

func NewParty(leader *Agent) *Party {
return &Party{Leader: leader}
}

func (p *Party) AddMember(a *Agent) {
p.mu.Lock()
defer p.mu.Unlock()
p.Members = append(p.Members, a)
}

func (p *Party) RemoveMember(uuid string) {
p.mu.Lock()
defer p.mu.Unlock()
var kept []*Agent
for _, m := range p.Members {
if m.UUID != uuid {
kept = append(kept, m)
}
}
p.Members = kept
}

func (p *Party) AliveCount() int {
p.mu.RLock()
defer p.mu.RUnlock()
count := 0
for _, m := range p.Members {
if m.IsAlive() {
count++
}
}
if p.Leader != nil && p.Leader.IsAlive() {
count++
}
return count
}

func (p *Party) AllMembers() []*Agent {
p.mu.RLock()
defer p.mu.RUnlock()
all := make([]*Agent, 0, len(p.Members)+1)
if p.Leader != nil {
all = append(all, p.Leader)
}
all = append(all, p.Members...)
return all
}

type WorldState struct {
Party     *Party
CreatedAt time.Time
}

func NewWorldState(leaderName string) (*WorldState, error) {
leader, err := NewAgent(leaderName, "Elvaan", "RDM", 200000, 50000)
if err != nil {
return nil, err
}
leader.Status = StatusBuffing
return &WorldState{
Party:     NewParty(leader),
CreatedAt: time.Now(),
}, nil
}
