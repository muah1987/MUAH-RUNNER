package swarm

import "fmt"

type ScalingPolicy struct {
MaxPartySize   int
MinPartySize   int
ScaleUpThreshold  float64
ScaleDownThreshold float64
}

func DefaultScalingPolicy() ScalingPolicy {
return ScalingPolicy{
MaxPartySize:      6,
MinPartySize:      1,
ScaleUpThreshold:  0.8,
ScaleDownThreshold: 0.2,
}
}

func ShouldScaleUp(ws *WorldState, policy ScalingPolicy) (bool, string) {
s := GetStatus(ws)
if s.AliveAgents >= policy.MaxPartySize {
return false, "at max party size"
}
alive := ws.Party.AllMembers()
engaged := 0
for _, a := range alive {
if a.Status == StatusEngaged {
engaged++
}
}
if s.AliveAgents == 0 {
return true, "no alive agents"
}
load := float64(engaged) / float64(s.AliveAgents)
if load >= policy.ScaleUpThreshold {
return true, fmt.Sprintf("load %.0f%% exceeds threshold", load*100)
}
return false, "load within bounds"
}
