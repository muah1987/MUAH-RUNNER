package playwright

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type VisualTest struct {
ID        string    `json:"id"`
URL       string    `json:"url"`
Baseline  string    `json:"baseline"`
Current   string    `json:"current"`
Diff      float64   `json:"diff_percent"`
Passed    bool      `json:"passed"`
Timestamp time.Time `json:"timestamp"`
}

type VisualTester struct {
muahDir string
browser *Browser
}

func NewVisualTester(muahDir string) *VisualTester {
return &VisualTester{
muahDir: muahDir,
browser: NewBrowser("chromium"),
}
}

func (v *VisualTester) screenshotsDir() string {
return filepath.Join(v.muahDir, "playwright", "screenshots")
}

func (v *VisualTester) CaptureScreenshot(name, url string) (string, error) {
dir := v.screenshotsDir()
os.MkdirAll(dir, 0755)
filename := fmt.Sprintf("%s-%d.png", name, time.Now().Unix())
path := filepath.Join(dir, filename)
if err := v.browser.Screenshot(url, path); err != nil {
return "", err
}
return path, nil
}

func (v *VisualTester) SaveTestResult(test *VisualTest) error {
os.MkdirAll(v.screenshotsDir(), 0755)
data, err := json.MarshalIndent(test, "", "  ")
if err != nil {
return err
}
filename := fmt.Sprintf("test-%s.json", test.ID)
return os.WriteFile(filepath.Join(v.screenshotsDir(), filename), data, 0644)
}
