package spec

import "fmt"

func ProposeChanges(current *RunnerSpec, drifts []string) []string {
var proposals []string
for _, d := range drifts {
proposals = append(proposals, fmt.Sprintf("Fix: %s", d))
}
return proposals
}
