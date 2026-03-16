package adapters

import "os"

type GitHubAdapter struct {
token string
}

func NewGitHubAdapter(token string) *GitHubAdapter {
if token == "" {
token = os.Getenv("GITHUB_TOKEN")
}
return &GitHubAdapter{token: token}
}

func (a *GitHubAdapter) Name() string      { return "github" }
func (a *GitHubAdapter) Platform() string  { return "github" }
func (a *GitHubAdapter) IsAvailable() bool { return a.token != "" }
func (a *GitHubAdapter) GetToken() string  { return a.token }
