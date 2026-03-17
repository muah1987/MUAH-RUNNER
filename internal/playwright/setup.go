package playwright

import (
"fmt"
"os/exec"
)

type Setup struct{}

func NewSetup() *Setup {
return &Setup{}
}

func (s *Setup) IsInstalled() bool {
_, err := exec.LookPath("npx")
if err != nil {
return false
}
out, err := exec.Command("npx", "playwright", "--version").Output()
return err == nil && len(out) > 0
}

func (s *Setup) Install() error {
if s.IsInstalled() {
return nil
}
fmt.Println("Installing Playwright...")
return exec.Command("npx", "playwright", "install").Run()
}

func (s *Setup) InstallBrowsers() error {
return exec.Command("npx", "playwright", "install", "--with-deps").Run()
}
