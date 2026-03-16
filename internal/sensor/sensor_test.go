package sensor

import (
	"testing"
)

func TestDetect(t *testing.T) {
	caps, err := Detect()
	if err != nil {
		t.Fatalf("Detect() error: %v", err)
	}
	if caps.OS == "" {
		t.Error("OS should not be empty")
	}
	if caps.Arch == "" {
		t.Error("Arch should not be empty")
	}
	if caps.CPUCores <= 0 {
		t.Errorf("CPUCores should be positive, got %d", caps.CPUCores)
	}
	if len(caps.Labels) == 0 {
		t.Error("Labels should not be empty")
	}
}

func TestExtractShortVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v18.17.0", "v18"},
		{"go version go1.21.0 linux/amd64", "v1"},
		{"Docker version 24.0.5, build ced0996", "v24"},
		{"", ""},
		{"no version here", ""},
	}
	for _, tt := range tests {
		got := extractShortVersion(tt.input)
		if got != tt.want {
			t.Errorf("extractShortVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestComputeLabels(t *testing.T) {
	caps := &Capabilities{
		OS:           "linux",
		Arch:         "amd64",
		CPUCores:     4,
		GPUAvailable: true,
		Toolchains: []Toolchain{
			{Name: "docker", Version: "Docker version 24.0.5"},
		},
	}
	labels := computeLabels(caps)

	has := func(label string) bool {
		for _, l := range labels {
			if l == label {
				return true
			}
		}
		return false
	}

	if !has("os:linux") {
		t.Error("missing os:linux label")
	}
	if !has("arch:amd64") {
		t.Error("missing arch:amd64 label")
	}
	if !has("gpu:available") {
		t.Error("missing gpu:available label")
	}
	if !has("docker:v24") {
		t.Error("missing docker:v24 label")
	}
}
