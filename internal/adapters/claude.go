package adapters

import "os"

type ClaudeAdapter struct{ key string }

func NewClaudeAdapter(key string) *ClaudeAdapter {
if key == "" {
key = os.Getenv("ANTHROPIC_API_KEY")
}
return &ClaudeAdapter{key: key}
}

func (a *ClaudeAdapter) Name() string      { return "claude" }
func (a *ClaudeAdapter) Platform() string  { return "claude" }
func (a *ClaudeAdapter) IsAvailable() bool { return a.key != "" }
func (a *ClaudeAdapter) GetToken() string  { return a.key }
