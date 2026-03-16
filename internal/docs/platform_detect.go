package docs

import "os"
import "path/filepath"

func DetectPlatform(rootDir string) string {
if _, err := os.Stat(filepath.Join(rootDir, ".github")); err == nil {
return "github"
}
if _, err := os.Stat(filepath.Join(rootDir, ".gitlab-ci.yml")); err == nil {
return "gitlab"
}
return "generic"
}
