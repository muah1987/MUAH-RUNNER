package sensor

import (
	"os"
	"os/exec"
	"runtime"
)

// Environment holds detected runtime information.
type Environment struct {
	OS       string
	Arch     string
	Hostname string
	// Platform is one of: github, gitlab, circleci, bitbucket, jenkins,
	// azure, travis, teamcity, drone, claude, gemini, ci, standalone
	Platform string
	// CIName is the human-readable name of the detected CI/platform.
	CIName string
}

// DetectEnvironment auto-detects the current execution platform from
// well-known environment variables. Detection goes from most-specific
// to most-general so nested environments resolve correctly.
func DetectEnvironment() Environment {
	hostname, _ := os.Hostname()
	platform, ciName := detectPlatform()
	return Environment{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
		Platform: platform,
		CIName:   ciName,
	}
}

// detectPlatform returns (platform-id, human-name).
func detectPlatform() (string, string) {
	type check struct {
		env      string
		value    string // empty = any non-empty value
		platform string
		name     string
	}
	checks := []check{
		{"GITHUB_ACTIONS", "true", "github", "GitHub Actions"},
		{"GITLAB_CI", "true", "gitlab", "GitLab CI"},
		{"CIRCLECI", "true", "circleci", "CircleCI"},
		{"BITBUCKET_BUILD_NUMBER", "", "bitbucket", "Bitbucket Pipelines"},
		{"JENKINS_URL", "", "jenkins", "Jenkins"},
		{"TF_BUILD", "True", "azure", "Azure DevOps"},
		{"TRAVIS", "true", "travis", "Travis CI"},
		{"TEAMCITY_VERSION", "", "teamcity", "TeamCity"},
		{"DRONE", "true", "drone", "Drone CI"},
		{"CLAUDE_CONVERSATION_ID", "", "claude", "Claude"},
		{"ANTHROPIC_API_KEY", "", "claude", "Claude"},
		{"GEMINI_API_KEY", "", "gemini", "Gemini"},
		{"GOOGLE_AI_API_KEY", "", "gemini", "Gemini"},
		{"CI", "true", "ci", "Generic CI"},
	}
	for _, c := range checks {
		v := os.Getenv(c.env)
		if v == "" {
			continue
		}
		if c.value == "" || v == c.value {
			return c.platform, c.name
		}
	}
	return "standalone", "Standalone"
}

type Toolchain struct {
Go         string
Git        string
Docker     string
Node       string
NPX        string
Playwright bool
}

func DetectToolchain() Toolchain {
tc := Toolchain{}
if out, err := exec.Command("go", "version").Output(); err == nil {
tc.Go = string(out)
}
if out, err := exec.Command("git", "--version").Output(); err == nil {
tc.Git = string(out)
}
if out, err := exec.Command("docker", "--version").Output(); err == nil {
tc.Docker = string(out)
}
if out, err := exec.Command("node", "--version").Output(); err == nil {
tc.Node = string(out)
}
if _, err := exec.LookPath("npx"); err == nil {
tc.NPX = "available"
if _, err2 := exec.Command("npx", "playwright", "--version").Output(); err2 == nil {
tc.Playwright = true
}
}
return tc
}
