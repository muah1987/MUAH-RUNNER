package config

import (
"encoding/json"
"os"
)

type Secrets struct {
GithubToken    string `json:"github_token,omitempty"`
AnthropicKey   string `json:"anthropic_key,omitempty"`
GeminiKey      string `json:"gemini_key,omitempty"`
}

func LoadSecrets(path string) (*Secrets, error) {
s := &Secrets{}
data, err := os.ReadFile(path)
if err != nil {
if os.IsNotExist(err) {
return s, nil
}
return nil, err
}
if err := json.Unmarshal(data, s); err != nil {
return nil, err
}
return s, nil
}
