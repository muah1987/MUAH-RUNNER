package planning

import (
"encoding/json"
"os"
"path/filepath"
"strings"
)

type DocsCoverage struct {
HasReadme      bool `json:"has_readme"`
HasLicense     bool `json:"has_license"`
HasContributing bool `json:"has_contributing"`
HasChangelog   bool `json:"has_changelog"`
HasSecurity    bool `json:"has_security"`
}

type RepoScan struct {
TechStack    []string     `json:"tech_stack"`
Languages    []string     `json:"languages"`
Features     []string     `json:"features"`
Gaps         []string     `json:"gaps"`
Platform     string       `json:"platform"`
DocsCoverage DocsCoverage `json:"docs_coverage"`
RootDir      string       `json:"root_dir"`
}

func ScanRepo(rootDir string) (*RepoScan, error) {
scan := &RepoScan{RootDir: rootDir}

scan.Languages = detectLanguages(rootDir)
scan.TechStack = detectTechStack(rootDir)
scan.Platform = detectPlatform(rootDir)
scan.DocsCoverage = checkDocs(rootDir)
scan.Gaps = detectGaps(scan)
scan.Features = detectFeatures(rootDir)

return scan, nil
}

func detectLanguages(dir string) []string {
langs := make(map[string]bool)
exts := map[string]string{
".go":   "Go",
".py":   "Python",
".js":   "JavaScript",
".ts":   "TypeScript",
".rs":   "Rust",
".java": "Java",
".rb":   "Ruby",
".sh":   "Shell",
}
filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
if err != nil || info.IsDir() {
return nil
}
ext := filepath.Ext(path)
if lang, ok := exts[ext]; ok {
langs[lang] = true
}
return nil
})
var result []string
for l := range langs {
result = append(result, l)
}
return result
}

func detectTechStack(dir string) []string {
var stack []string
checks := map[string]string{
"go.mod":           "Go Modules",
"package.json":     "Node.js",
"requirements.txt": "Python pip",
"Cargo.toml":       "Rust Cargo",
"pom.xml":          "Maven",
"Dockerfile":       "Docker",
"docker-compose.yml": "Docker Compose",
"Makefile":         "Make",
".github":          "GitHub Actions",
}
for file, tech := range checks {
if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
stack = append(stack, tech)
}
}
return stack
}

func detectPlatform(dir string) string {
if _, err := os.Stat(filepath.Join(dir, ".github")); err == nil {
return "github"
}
if _, err := os.Stat(filepath.Join(dir, ".gitlab-ci.yml")); err == nil {
return "gitlab"
}
return "standalone"
}

func checkDocs(dir string) DocsCoverage {
cov := DocsCoverage{}
files, _ := filepath.Glob(filepath.Join(dir, "README*"))
cov.HasReadme = len(files) > 0
files, _ = filepath.Glob(filepath.Join(dir, "LICENSE*"))
cov.HasLicense = len(files) > 0
files, _ = filepath.Glob(filepath.Join(dir, "CONTRIBUTING*"))
cov.HasContributing = len(files) > 0
files, _ = filepath.Glob(filepath.Join(dir, "CHANGELOG*"))
cov.HasChangelog = len(files) > 0
files, _ = filepath.Glob(filepath.Join(dir, "SECURITY*"))
cov.HasSecurity = len(files) > 0
return cov
}

func detectGaps(scan *RepoScan) []string {
var gaps []string
if !scan.DocsCoverage.HasReadme {
gaps = append(gaps, "Missing README")
}
if !scan.DocsCoverage.HasLicense {
gaps = append(gaps, "Missing LICENSE")
}
if !scan.DocsCoverage.HasChangelog {
gaps = append(gaps, "Missing CHANGELOG")
}
if !scan.DocsCoverage.HasContributing {
gaps = append(gaps, "Missing CONTRIBUTING guide")
}
if !scan.DocsCoverage.HasSecurity {
gaps = append(gaps, "Missing SECURITY policy")
}
return gaps
}

func detectFeatures(dir string) []string {
var features []string
// Simple feature detection based on directory names
dirs := []string{"api", "cmd", "internal", "pkg", "lib", "test", "docs", "scripts"}
for _, d := range dirs {
if _, err := os.Stat(filepath.Join(dir, d)); err == nil {
features = append(features, strings.ToUpper(d[:1])+d[1:]+" layer")
}
}
return features
}

func SaveScan(muahDir string, scan *RepoScan) error {
path := filepath.Join(muahDir, "planning", "repo-scan.json")
os.MkdirAll(filepath.Dir(path), 0755)
data, err := json.MarshalIndent(scan, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
