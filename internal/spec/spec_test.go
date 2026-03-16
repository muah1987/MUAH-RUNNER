package spec

import "testing"

func TestLoadDefault(t *testing.T) {
s, err := Load("/nonexistent/path.yml")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if s == nil {
t.Fatal("expected non-nil spec")
}
}

func TestValidate(t *testing.T) {
s := &RunnerSpec{Constraints: map[string]int{"timeout": 3600}}
issues := Validate(s)
if len(issues) > 0 {
t.Errorf("unexpected issues: %v", issues)
}
}
