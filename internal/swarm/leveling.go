package swarm

const (
BaseXPPerLevel = 1000
MaxLevel       = 99
)

type LevelData struct {
NPCName string
Job     string
Level   int
XP      int
XPNext  int
}

func XPForLevel(level int) int {
return level * BaseXPPerLevel
}

func AddXP(data *LevelData, xp int) bool {
data.XP += xp
data.XPNext = XPForLevel(data.Level + 1)
if data.XP >= data.XPNext && data.Level < MaxLevel {
data.Level++
return true // leveled up
}
return false
}
