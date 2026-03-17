package toollearn

import "fmt"

// Injector injects tool guidance into AI context
type Injector struct {
teacher  *Teacher
profiler *Profiler
}

func NewInjector(muahDir string) *Injector {
return &Injector{
teacher:  NewTeacher(muahDir),
profiler: NewProfiler(muahDir),
}
}

// BuildContext builds a context string for the AI with tool guidance
func (i *Injector) BuildContext(aiName, tool string) string {
profile := i.profiler.GetOrCreate(aiName)
lessons := i.teacher.AllLessons()

ctx := fmt.Sprintf("AI Profile for %s: overall score %.1f%%\n", aiName, profile.OverallScore)
if ts, ok := profile.ToolScores[tool]; ok {
ctx += fmt.Sprintf("Tool %s: %d calls, %.1f%% success rate\n", tool, ts.Calls, ts.Score)
}
for _, lesson := range lessons {
ctx += fmt.Sprintf("Lesson: %s → Do: %s\n", lesson.Pattern, lesson.Do)
}
return ctx
}
