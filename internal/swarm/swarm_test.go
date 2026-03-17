package swarm

import (
"testing"
)

func TestRaceStats(t *testing.T) {
race, ok := GetRace("Elvaan")
if !ok {
t.Fatal("expected Elvaan race")
}
if race.Stats.INT < 80 {
t.Error("Elvaan should have high INT")
}
}

func TestAllRaces(t *testing.T) {
expectedRaces := []string{"Elvaan", "Galka", "Mithra", "Tarutaru", "Hume"}
for _, name := range expectedRaces {
if _, ok := GetRace(name); !ok {
t.Errorf("missing race: %s", name)
}
}
}

func TestAllJobs(t *testing.T) {
jobs := AllJobs()
if len(jobs) != 22 {
t.Errorf("expected 22 jobs, got %d", len(jobs))
}
}

func TestAgentSpawn(t *testing.T) {
agent, err := NewAgent("TestNPC", "Galka", "WAR", 100000, 20000)
if err != nil {
t.Fatalf("spawn error: %v", err)
}
if agent.NPCName != "TestNPC" {
t.Errorf("expected TestNPC, got %s", agent.NPCName)
}
if !agent.IsAlive() {
t.Error("agent should be alive")
}
}

func TestHPMP(t *testing.T) {
hp := NewHP(10000)
mp := NewMP(5000)
hp.Deplete(3000)
if hp.Current != 7000 {
t.Errorf("expected 7000, got %d", hp.Current)
}
mp.Consume(2000)
if mp.Current != 3000 {
t.Errorf("expected 3000, got %d", mp.Current)
}
pct := hp.Percent()
if pct < 69 || pct > 71 {
t.Errorf("expected ~70%%, got %.1f%%", pct)
}
}

func TestWorldState(t *testing.T) {
ws, err := NewWorldState("Prishe")
if err != nil {
t.Fatalf("world state error: %v", err)
}
if ws.Party.Leader == nil {
t.Error("expected leader")
}
status := GetStatus(ws)
if status.TotalAgents < 1 {
t.Error("expected at least 1 agent")
}
}

func TestCrystalPool(t *testing.T) {
pool := NewCrystalPool(100000, 500)
ok := pool.Allocate(50000, 100)
if !ok {
t.Error("expected allocation to succeed")
}
ok = pool.Allocate(60000, 100)
if ok {
t.Error("expected allocation to fail (over budget)")
}
tokens, calls := pool.Remaining()
if tokens != 50000 {
t.Errorf("expected 50000 tokens remaining, got %d", tokens)
}
if calls != 400 {
t.Errorf("expected 400 calls remaining, got %d", calls)
}
}

func TestAgentKO(t *testing.T) {
agent, _ := NewAgent("Zeid", "Galka", "DRK", 1000, 500)
agent.KO()
if agent.IsAlive() {
t.Error("KO agent should not be alive")
}
if agent.HP.Current != 0 {
t.Error("KO agent HP should be 0")
}
}
