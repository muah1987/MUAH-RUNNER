package toollearn

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
)

type Lesson struct {
Pattern   string     `json:"pattern"`
Do        string     `json:"do"`
DontDo    string     `json:"dont_do"`
ErrorClass ErrorClass `json:"error_class"`
Examples  []string   `json:"examples"`
}

type Teacher struct {
muahDir string
lessons []Lesson
}

func NewTeacher(muahDir string) *Teacher {
t := &Teacher{muahDir: muahDir}
t.load()
return t
}

func (t *Teacher) lessonsPath() string {
return filepath.Join(t.muahDir, "tool_use", "teaching", "lessons.json")
}

func (t *Teacher) load() {
data, err := os.ReadFile(t.lessonsPath())
if err != nil {
return
}
json.Unmarshal(data, &t.lessons)
}

func (t *Teacher) GenerateLesson(record FailureRecord) Lesson {
strategy := GetRecoveryStrategy(record.Class)
lesson := Lesson{
Pattern:    fmt.Sprintf("When using %s and getting %s", record.Tool, record.Class),
Do:         strategy.Strategy,
DontDo:     fmt.Sprintf("Don't retry immediately on %s", record.Class),
ErrorClass: record.Class,
Examples:   []string{fmt.Sprintf("Error: %s", record.Error)},
}
t.lessons = append(t.lessons, lesson)
t.save()
return lesson
}

func (t *Teacher) save() {
os.MkdirAll(filepath.Dir(t.lessonsPath()), 0755)
data, _ := json.MarshalIndent(t.lessons, "", "  ")
os.WriteFile(t.lessonsPath(), data, 0644)
}

func (t *Teacher) AllLessons() []Lesson {
return t.lessons
}
