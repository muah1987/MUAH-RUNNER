package questions

// Recommender suggests answers based on context
type Recommender struct{}

func (r *Recommender) Recommend(q *Question, context map[string]string) int {
return q.Recommended
}

// StandardQuestions returns common muah-runner questions
func StandardQuestions() []*Question {
return []*Question{
{
ID:   "init-platform",
Text: "What platform are you running on?",
Options: []Option{
{0, "GitHub Actions", "github"},
{1, "Claude", "claude"},
{2, "Standalone", "standalone"},
},
Recommended: 0,
TimeoutSec:  10,
},
{
ID:   "enable-swarm",
Text: "Enable Vana'diel swarm (multi-agent)?",
Options: []Option{
{0, "Yes", "true"},
{1, "No", "false"},
},
Recommended: 0,
TimeoutSec:  10,
},
{
ID:   "default-race",
Text: "Default AI race for spawned agents?",
Options: []Option{
{0, "Elvaan (Opus - best quality)", "Elvaan"},
{1, "Galka (Sonnet - balanced)", "Galka"},
{2, "Mithra (Haiku - fastest)", "Mithra"},
},
Recommended: 1,
TimeoutSec:  10,
},
}
}
