package toollearn

type Monitor struct {
consecutiveFailures map[string]int
profiler            *Profiler
teacher             *Teacher
muahDir             string
}

func NewMonitor(muahDir string) *Monitor {
return &Monitor{
consecutiveFailures: make(map[string]int),
profiler:            NewProfiler(muahDir),
teacher:             NewTeacher(muahDir),
muahDir:             muahDir,
}
}

func (m *Monitor) RecordUse(aiName, tool string, success bool, errMsg string) {
key := aiName + ":" + tool
if success {
m.consecutiveFailures[key] = 0
m.profiler.RecordSuccess(aiName, tool)
} else {
m.consecutiveFailures[key]++
m.profiler.RecordFailure(aiName, tool, errMsg)
errClass := ClassifyError(errMsg)
record := FailureRecord{
Tool:  tool,
AI:    aiName,
Error: errMsg,
Class: errClass,
}
LogFailure(m.muahDir, record)
if m.consecutiveFailures[key] >= ConsecutiveFailureThreshold {
m.teacher.GenerateLesson(record)
}
}
m.profiler.Save(aiName)
}

func (m *Monitor) ShouldRecover(aiName, tool string) bool {
key := aiName + ":" + tool
return m.consecutiveFailures[key] >= ConsecutiveFailureThreshold
}
