package swarm

type HPStats struct {
Max     int `json:"max"`
Current int `json:"current"`
}

type MPStats struct {
Max     int `json:"max"`
Current int `json:"current"`
}

// HPPercent returns HP as percentage
func (h HPStats) Percent() float64 {
if h.Max == 0 {
return 100
}
return float64(h.Current) / float64(h.Max) * 100
}

// MPPercent returns MP as percentage
func (m MPStats) Percent() float64 {
if m.Max == 0 {
return 100
}
return float64(m.Current) / float64(m.Max) * 100
}

func (h *HPStats) Deplete(amount int) {
h.Current -= amount
if h.Current < 0 {
h.Current = 0
}
}

func (h *HPStats) Restore(amount int) {
h.Current += amount
if h.Current > h.Max {
h.Current = h.Max
}
}

func (m *MPStats) Consume(amount int) {
m.Current -= amount
if m.Current < 0 {
m.Current = 0
}
}

func (m *MPStats) Restore(amount int) {
m.Current += amount
if m.Current > m.Max {
m.Current = m.Max
}
}

// NewHP creates HP based on context token budget
func NewHP(contextTokens int) HPStats {
return HPStats{Max: contextTokens, Current: contextTokens}
}

// NewMP creates MP based on output token budget
func NewMP(outputTokens int) MPStats {
return MPStats{Max: outputTokens, Current: outputTokens}
}
