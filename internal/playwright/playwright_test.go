package playwright

import "testing"

func TestSetupIsInstalled(t *testing.T) {
s := NewSetup()
// Just ensure it doesn't panic
_ = s.IsInstalled()
}

func TestBrowserCreation(t *testing.T) {
b := NewBrowser("")
if b.browserType != "chromium" {
t.Errorf("expected chromium default, got %s", b.browserType)
}
}

func TestVisualTester(t *testing.T) {
dir := t.TempDir()
vt := NewVisualTester(dir)
if vt.muahDir != dir {
t.Error("wrong muah dir")
}
}
