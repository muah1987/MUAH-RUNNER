package docker

import (
"fmt"
)

type HealAttempt struct {
Error      string
Image      string
Fix        string
Successful bool
}

type SelfHealer struct {
simulator *Simulator
sandbox   *Sandbox
}

func NewSelfHealer(muahDir string) *SelfHealer {
return &SelfHealer{
simulator: NewSimulator(muahDir),
sandbox:   NewSandbox(muahDir),
}
}

func (h *SelfHealer) Heal(errMsg, image, proposedFix string) (*HealAttempt, error) {
attempt := &HealAttempt{
Error: errMsg,
Image: image,
Fix:   proposedFix,
}

sim := &Simulation{
Reason:     fmt.Sprintf("Self-heal: %s", errMsg),
Image:      image,
Hypothesis: proposedFix,
Steps:      []string{proposedFix},
}

if err := h.simulator.Run(sim); err != nil {
return attempt, err
}

attempt.Successful = sim.Result.FixWorked
return attempt, nil
}
