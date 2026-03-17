package adapters

import "os"

type GeminiAdapter struct{ key string }

func NewGeminiAdapter(key string) *GeminiAdapter {
if key == "" {
key = os.Getenv("GEMINI_API_KEY")
}
return &GeminiAdapter{key: key}
}

func (a *GeminiAdapter) Name() string      { return "gemini" }
func (a *GeminiAdapter) Platform() string  { return "gemini" }
func (a *GeminiAdapter) IsAvailable() bool { return a.key != "" }
func (a *GeminiAdapter) GetToken() string  { return a.key }
