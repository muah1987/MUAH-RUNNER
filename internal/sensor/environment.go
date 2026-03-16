package sensor

import (
"os"
"os/exec"
"runtime"
)

type Environment struct {
OS       string
Arch     string
Hostname string
Platform string
}

func DetectEnvironment() Environment {
hostname, _ := os.Hostname()
platform := "standalone"
if os.Getenv("GITHUB_ACTIONS") == "true" {
platform = "github"
} else if os.Getenv("GITLAB_CI") == "true" {
platform = "gitlab"
}
return Environment{OS: runtime.GOOS, Arch: runtime.GOARCH, Hostname: hostname, Platform: platform}
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
