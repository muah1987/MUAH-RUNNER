package spec

import (
"os"
"gopkg.in/yaml.v3"
)

type RunnerSpec struct {
Version     string            `yaml:"version"`
Constraints map[string]int    `yaml:"constraints"`
QualityGates map[string]float64 `yaml:"quality_gates"`
}

func Load(path string) (*RunnerSpec, error) {
s := &RunnerSpec{
Constraints: map[string]int{"timeout": 3600, "max_retries": 3},
QualityGates: map[string]float64{"min_grade": 70},
}
data, err := os.ReadFile(path)
if err != nil {
if os.IsNotExist(err) {
return s, nil
}
return nil, err
}
return s, yaml.Unmarshal(data, s)
}
