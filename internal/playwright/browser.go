package playwright

import (
"fmt"
"os/exec"
)

type Browser struct {
browserType string // chromium, firefox, webkit
}

func NewBrowser(browserType string) *Browser {
if browserType == "" {
browserType = "chromium"
}
return &Browser{browserType: browserType}
}

func (b *Browser) Screenshot(url, outputPath string) error {
args := []string{
"playwright", "screenshot",
"--browser", b.browserType,
url, outputPath,
}
out, err := exec.Command("npx", args...).CombinedOutput()
if err != nil {
return fmt.Errorf("screenshot failed: %s", string(out))
}
return nil
}

func (b *Browser) Navigate(url string) error {
// For CLI usage, we just verify the URL is accessible
out, err := exec.Command("curl", "-sSI", url).Output()
if err != nil {
return fmt.Errorf("navigate failed: %v", err)
}
_ = out
return nil
}
