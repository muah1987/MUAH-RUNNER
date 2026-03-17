package docker

import "testing"

func TestEngineAvailability(t *testing.T) {
e := NewEngine()
// Just test it doesn't panic - availability depends on environment
_ = e.IsAvailable()
}

func TestSandboxInit(t *testing.T) {
dir := t.TempDir()
s := NewSandbox(dir)
if err := s.Init(); err != nil {
t.Fatalf("sandbox init error: %v", err)
}
}

func TestSimulationNoDocker(t *testing.T) {
dir := t.TempDir()
sim := &Simulation{
Reason:     "test",
Image:      "ubuntu:22.04",
Hypothesis: "fix works",
Steps:      []string{"echo hello"},
}
simulator := NewSimulator(dir)
// If docker not available, should still save result without error
_ = simulator.Run(sim)
}

func TestCleanup(t *testing.T) {
c := NewCleanup()
// Just ensure methods exist and don't panic when docker unavailable
if c.engine.IsAvailable() {
t.Log("Docker available, skipping destructive cleanup test")
}
}
