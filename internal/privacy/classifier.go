package privacy

import "strings"

type Classification string

const (
Private    Classification = "private"
Redactable Classification = "redactable"
Public     Classification = "public"
)

var privateKeywords = []string{"password", "secret", "token", "key", "credential", "auth", "private"}
var redactableKeywords = []string{"email", "name", "address", "phone", "personal"}

func Classify(thought string) Classification {
lower := strings.ToLower(thought)
for _, kw := range privateKeywords {
if strings.Contains(lower, kw) {
return Private
}
}
for _, kw := range redactableKeywords {
if strings.Contains(lower, kw) {
return Redactable
}
}
return Public
}
