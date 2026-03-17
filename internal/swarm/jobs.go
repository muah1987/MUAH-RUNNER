package swarm

type Job struct {
Code        string
Name        string
Category    string // standard, advanced, expert
Description string
Abilities   []string
Unlocked    bool
}

var Jobs = map[string]Job{
// Standard
"WAR": {Code: "WAR", Name: "Warrior", Category: "standard", Description: "Brute force task executor", Abilities: []string{"provoke", "berserk"}, Unlocked: true},
"MNK": {Code: "MNK", Name: "Monk", Category: "standard", Description: "Physical task specialist", Abilities: []string{"focus", "boost"}, Unlocked: true},
"WHM": {Code: "WHM", Name: "White Mage", Category: "standard", Description: "Recovery and support specialist", Abilities: []string{"cure", "protect", "shell"}, Unlocked: true},
"BLM": {Code: "BLM", Name: "Black Mage", Category: "standard", Description: "Destructive code analysis", Abilities: []string{"fire", "blizzard", "thunder"}, Unlocked: true},
"RDM": {Code: "RDM", Name: "Red Mage", Category: "standard", Description: "Versatile hybrid agent", Abilities: []string{"convert", "refresh"}, Unlocked: true},
"THF": {Code: "THF", Name: "Thief", Category: "standard", Description: "Information gathering specialist", Abilities: []string{"steal", "sneak"}, Unlocked: true},
// Advanced
"PLD": {Code: "PLD", Name: "Paladin", Category: "advanced", Description: "Security and compliance guardian", Abilities: []string{"holy", "cover"}},
"DRK": {Code: "DRK", Name: "Dark Knight", Category: "advanced", Description: "Aggressive refactoring agent", Abilities: []string{"drain", "absorb"}},
"BST": {Code: "BST", Name: "Beastmaster", Category: "advanced", Description: "Tool and subprocess handler", Abilities: []string{"call_beast", "reward"}},
"BRD": {Code: "BRD", Name: "Bard", Category: "advanced", Description: "Documentation and communication", Abilities: []string{"minuet", "madrigal"}},
"RNG": {Code: "RNG", Name: "Ranger", Category: "advanced", Description: "Long-range scanning and monitoring", Abilities: []string{"eagle_eye_shot"}},
"SAM": {Code: "SAM", Name: "Samurai", Category: "advanced", Description: "Precision code execution", Abilities: []string{"meditate", "tachi"}},
"NIN": {Code: "NIN", Name: "Ninja", Category: "advanced", Description: "Stealth operations and secrets", Abilities: []string{"ninjutsu", "utsusemi"}},
"DRG": {Code: "DRG", Name: "Dragoon", Category: "advanced", Description: "High-altitude architecture analysis", Abilities: []string{"spirit_link", "super_jump"}},
"SMN": {Code: "SMN", Name: "Summoner", Category: "advanced", Description: "External service integrator", Abilities: []string{"astral_flow", "avatar"}},
// Expert
"BLU": {Code: "BLU", Name: "Blue Mage", Category: "expert", Description: "Learns and replicates patterns", Abilities: []string{"blue_magic", "chain_affinity"}},
"COR": {Code: "COR", Name: "Corsair", Category: "expert", Description: "Probabilistic planning and risk", Abilities: []string{"quick_draw", "random_deal"}},
"PUP": {Code: "PUP", Name: "Puppetmaster", Category: "expert", Description: "Automated sub-agent controller", Abilities: []string{"activate", "maneuver"}},
"DNC": {Code: "DNC", Name: "Dancer", Category: "expert", Description: "Adaptive UI/UX automation", Abilities: []string{"samba", "waltz"}},
"SCH": {Code: "SCH", Name: "Scholar", Category: "expert", Description: "Deep research and analysis", Abilities: []string{"light_arts", "dark_arts"}},
"GEO": {Code: "GEO", Name: "Geomancer", Category: "expert", Description: "Environment and context mapping", Abilities: []string{"geomancy", "indi"}},
"RUN": {Code: "RUN", Name: "Rune Fencer", Category: "expert", Description: "Adaptive defense and counter", Abilities: []string{"rune_enhancement", "lunge"}},
}

func GetJob(code string) (Job, bool) {
j, ok := Jobs[code]
return j, ok
}

func AllJobs() []Job {
var list []Job
for _, j := range Jobs {
list = append(list, j)
}
return list
}
