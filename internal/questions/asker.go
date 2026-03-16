package questions

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type Option struct {
Index int    `json:"index"`
Label string `json:"label"`
Value string `json:"value"`
}

type Question struct {
ID          string        `json:"id"`
Text        string        `json:"text"`
Options     []Option      `json:"options"`
Recommended int           `json:"recommended"`
Timeout     time.Duration `json:"-"`
TimeoutSec  int           `json:"timeout_seconds"`
}

type Answer struct {
QuestionID string    `json:"question_id"`
OptionIndex int      `json:"option_index"`
Value       string   `json:"value"`
AutoSelected bool    `json:"auto_selected"`
Timestamp   time.Time `json:"timestamp"`
}

type Asker struct {
muahDir string
history []Answer
}

func NewAsker(muahDir string) *Asker {
a := &Asker{muahDir: muahDir}
a.loadHistory()
return a
}

func (a *Asker) historyPath() string {
return filepath.Join(a.muahDir, "questions", "history.json")
}

func (a *Asker) loadHistory() {
data, err := os.ReadFile(a.historyPath())
if err != nil {
return
}
json.Unmarshal(data, &a.history)
}

func (a *Asker) saveHistory() {
os.MkdirAll(filepath.Dir(a.historyPath()), 0755)
data, _ := json.MarshalIndent(a.history, "", "  ")
os.WriteFile(a.historyPath(), data, 0644)
}

// Ask presents a question and waits for input, auto-selecting recommended if timeout expires.
// When autoMode is true or writing to a non-terminal, always auto-selects.
func (a *Asker) Ask(q *Question, autoMode bool) Answer {
if q.Timeout == 0 {
q.Timeout = 10 * time.Second
}

var answer Answer
if autoMode {
answer = a.autoSelect(q)
} else {
answer = a.timedAsk(q)
}

a.history = append(a.history, answer)
a.saveHistory()
return answer
}

func (a *Asker) autoSelect(q *Question) Answer {
idx := q.Recommended
var value string
if idx >= 0 && idx < len(q.Options) {
value = q.Options[idx].Value
}
return Answer{
QuestionID:   q.ID,
OptionIndex:  idx,
Value:        value,
AutoSelected: true,
Timestamp:    time.Now(),
}
}

func (a *Asker) timedAsk(q *Question) Answer {
fmt.Printf("\n❓ %s\n", q.Text)
for i, opt := range q.Options {
marker := "  "
if i == q.Recommended {
marker = "✓ "
}
fmt.Printf("  %s[%d] %s\n", marker, i+1, opt.Label)
}
fmt.Printf("\nAuto-selecting in %d seconds... ", int(q.Timeout.Seconds()))

done := make(chan int, 1)
go func() {
var choice int
fmt.Scanf("%d", &choice)
done <- choice - 1
}()

select {
case idx := <-done:
if idx >= 0 && idx < len(q.Options) {
return Answer{
QuestionID:  q.ID,
OptionIndex: idx,
Value:       q.Options[idx].Value,
Timestamp:   time.Now(),
}
}
return a.autoSelect(q)
case <-time.After(q.Timeout):
fmt.Printf("\n⏱ Timed out, auto-selected: %d\n", q.Recommended+1)
return a.autoSelect(q)
}
}

func (a *Asker) History() []Answer {
return a.history
}
