package swarm

type RaceStats struct {
INT   int
MND   int
STR   int
DEX   int
VIT   int
AGI   int
CHR   int
}

type CotStyle struct {
Verbose  bool
Poetic   bool
Logical  bool
Creative bool
}

type Skill struct {
Name        string
Description string
Level       int
}

type Race struct {
Name         string
AIModels     []string
Lore         string
Stats        RaceStats
NativeSkills map[string]Skill
CotStyle     CotStyle
Strengths    []string
Weaknesses   []string
}

var Races = map[string]Race{
"Elvaan": {
Name:     "Elvaan",
AIModels: []string{"claude-opus-4", "claude-opus-latest"},
Lore:     "Proud and noble, Elvaan prefer logic and precision. Masters of strategic reasoning.",
Stats:    RaceStats{INT: 90, MND: 85, STR: 70, DEX: 65, VIT: 70, AGI: 60, CHR: 80},
NativeSkills: map[string]Skill{
"deep_reasoning": {Name: "Deep Reasoning", Description: "Extended logical analysis", Level: 5},
},
CotStyle:   CotStyle{Logical: true, Verbose: true},
Strengths:  []string{"long-context", "code-generation", "analysis"},
Weaknesses: []string{"speed", "cost"},
},
"Galka": {
Name:     "Galka",
AIModels: []string{"claude-sonnet-4", "claude-sonnet-latest"},
Lore:     "Powerful and resilient. Galka balance strength with wisdom.",
Stats:    RaceStats{INT: 75, MND: 70, STR: 85, DEX: 70, VIT: 90, AGI: 65, CHR: 65},
NativeSkills: map[string]Skill{
"endurance": {Name: "Endurance", Description: "Maintains quality over long tasks", Level: 4},
},
CotStyle:   CotStyle{Logical: true},
Strengths:  []string{"balanced", "cost-effective", "reliable"},
Weaknesses: []string{"deep-reasoning"},
},
"Mithra": {
Name:     "Mithra",
AIModels: []string{"claude-haiku-4", "claude-haiku-latest"},
Lore:     "Quick and agile, Mithra excel at fast decisive action.",
Stats:    RaceStats{INT: 65, MND: 65, STR: 60, DEX: 95, VIT: 60, AGI: 95, CHR: 75},
NativeSkills: map[string]Skill{
"quick_strike": {Name: "Quick Strike", Description: "Rapid task completion", Level: 5},
},
CotStyle:   CotStyle{},
Strengths:  []string{"speed", "low-cost", "simple-tasks"},
Weaknesses: []string{"complex-reasoning", "long-context"},
},
"Tarutaru": {
Name:     "Tarutaru",
AIModels: []string{"claude-opus-prev", "o1", "o3"},
Lore:     "Magical prodigies. Tarutaru wield deep reasoning like ancient magic.",
Stats:    RaceStats{INT: 95, MND: 95, STR: 45, DEX: 50, VIT: 45, AGI: 50, CHR: 70},
NativeSkills: map[string]Skill{
"deep_magic": {Name: "Deep Magic", Description: "Multi-step reasoning chains", Level: 5},
},
CotStyle:   CotStyle{Logical: true, Verbose: true},
Strengths:  []string{"deep-reasoning", "math", "planning"},
Weaknesses: []string{"speed", "cost"},
},
"Hume": {
Name:     "Hume",
AIModels: []string{"human"},
Lore:     "Versatile and adaptable. Hume represent the human user.",
Stats:    RaceStats{INT: 70, MND: 70, STR: 70, DEX: 70, VIT: 70, AGI: 70, CHR: 85},
NativeSkills: map[string]Skill{
"creativity": {Name: "Creativity", Description: "Novel problem-solving", Level: 3},
},
CotStyle:   CotStyle{Creative: true},
Strengths:  []string{"creativity", "common-sense", "context"},
Weaknesses: []string{"speed", "consistency"},
},
}

func GetRace(name string) (Race, bool) {
r, ok := Races[name]
return r, ok
}

func RaceForModel(model string) Race {
for _, r := range Races {
for _, m := range r.AIModels {
if m == model {
return r
}
}
}
return Races["Hume"]
}
