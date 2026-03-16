package sensor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Toolchain describes a discovered CLI tool.
type Toolchain struct {
	Name    string
	Path    string
	Version string
}

// Capabilities describes detected environment capabilities.
type Capabilities struct {
	OS           string
	Arch         string
	CPUCores     int
	MemoryMB     uint64
	DiskFreeGB   uint64
	GPUAvailable bool
	Toolchains   []Toolchain
	Labels       []string
}

// Detect auto-discovers the current environment capabilities.
func Detect() (*Capabilities, error) {
	caps := &Capabilities{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCores: runtime.NumCPU(),
	}

	caps.MemoryMB = detectMemoryMB()
	caps.DiskFreeGB = detectDiskFreeGB()
	caps.GPUAvailable = detectGPU()
	caps.Toolchains = detectToolchains()
	caps.Labels = computeLabels(caps)

	return caps, nil
}

func detectMemoryMB() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			var kb uint64
			fmt.Sscanf(strings.TrimPrefix(line, "MemTotal:"), "%d", &kb)
			return kb / 1024
		}
	}
	return 0
}

func detectDiskFreeGB() uint64 {
	out, err := exec.Command("df", "-k", "/").Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0
	}
	var kb uint64
	fmt.Sscanf(fields[3], "%d", &kb)
	return kb / (1024 * 1024)
}

func detectGPU() bool {
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		return true
	}
	if _, err := os.Stat("/dev/dri"); err == nil {
		return true
	}
	return false
}

type toolSpec struct {
	name string
	cmd  string
	args []string
}

func detectToolchains() []Toolchain {
	specs := []toolSpec{
		{name: "node", cmd: "node", args: []string{"--version"}},
		{name: "python3", cmd: "python3", args: []string{"--version"}},
		{name: "python", cmd: "python", args: []string{"--version"}},
		{name: "go", cmd: "go", args: []string{"version"}},
		{name: "rust", cmd: "rustc", args: []string{"--version"}},
		{name: "docker", cmd: "docker", args: []string{"--version"}},
		{name: "kubectl", cmd: "kubectl", args: []string{"version", "--client", "--short"}},
		{name: "helm", cmd: "helm", args: []string{"version", "--short"}},
		{name: "java", cmd: "java", args: []string{"-version"}},
		{name: "mvn", cmd: "mvn", args: []string{"--version"}},
		{name: "gradle", cmd: "gradle", args: []string{"--version"}},
		{name: "make", cmd: "make", args: []string{"--version"}},
	}

	var tools []Toolchain
	for _, spec := range specs {
		path, err := exec.LookPath(spec.cmd)
		if err != nil {
			continue
		}
		out, err := exec.Command(spec.cmd, spec.args...).CombinedOutput()
		version := ""
		if err == nil {
			lines := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)
			if len(lines) > 0 {
				version = strings.TrimSpace(lines[0])
			}
		}
		tools = append(tools, Toolchain{
			Name:    spec.name,
			Path:    path,
			Version: version,
		})
	}
	return tools
}

func computeLabels(caps *Capabilities) []string {
	labels := []string{
		fmt.Sprintf("os:%s", caps.OS),
		fmt.Sprintf("arch:%s", caps.Arch),
		fmt.Sprintf("cpu-cores:%d", caps.CPUCores),
	}

	if caps.GPUAvailable {
		labels = append(labels, "gpu:available")
	}

	for _, tc := range caps.Toolchains {
		ver := extractShortVersion(tc.Version)
		if ver != "" {
			labels = append(labels, fmt.Sprintf("%s:%s", tc.Name, ver))
		} else {
			labels = append(labels, fmt.Sprintf("%s:available", tc.Name))
		}
	}
	return labels
}

func extractShortVersion(version string) string {
	fields := strings.Fields(version)
	for _, f := range fields {
		// Handle "goX.Y.Z" format (e.g. "go1.21.0")
		if strings.HasPrefix(f, "go") && len(f) > 2 && f[2] >= '0' && f[2] <= '9' {
			parts := strings.Split(f[2:], ".")
			return "v" + parts[0]
		}
		// Strip leading "v" and check for digit
		bare := strings.TrimPrefix(f, "v")
		if len(bare) > 0 && bare[0] >= '0' && bare[0] <= '9' {
			// Strip trailing comma or other punctuation
			bare = strings.TrimRight(bare, ",;:")
			parts := strings.Split(bare, ".")
			if len(parts) >= 1 && parts[0] != "" {
				return "v" + parts[0]
			}
		}
	}
	return ""
}
