package privacy

import "time"

type CoTEntry struct {
RunID     string         `json:"run_id"`
Timestamp time.Time      `json:"timestamp"`
Thoughts  []Thought      `json:"thoughts"`
}

type Thought struct {
ID             string         `json:"id"`
Content        string         `json:"content"`
Classification Classification `json:"classification"`
}
