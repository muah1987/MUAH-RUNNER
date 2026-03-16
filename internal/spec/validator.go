package spec

import "fmt"

func Validate(s *RunnerSpec) []string {
var issues []string
if s.Constraints["timeout"] <= 0 {
issues = append(issues, "timeout must be positive")
}
return issues
}

func ValidateOrFail(s *RunnerSpec) error {
issues := Validate(s)
if len(issues) > 0 {
return fmt.Errorf("spec validation failed: %v", issues)
}
return nil
}
