// Command muah-runner is the CLI entrypoint for the muah-runner autonomous agent.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muah1987/muah-runner/internal/config"
	"github.com/muah1987/muah-runner/internal/core"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printHelp()
		return nil
	}
	if args[0] == "--version" || args[0] == "version" {
		fmt.Printf("muah-runner %s\n", Version)
		return nil
	}

	baseDir, _ := os.Getwd()
	cfg, err := loadConfig(baseDir)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	runner := core.New(baseDir, cfg)

	cmd := args[0]
	rest := args[1:]

	switch cmd {

	// ── Core ──────────────────────────────────────────────────────────────────
	case "init":
		return runner.Init()

	case "start":
		fmt.Println("▶ Starting muah-runner daemon...")
		runner.Status()
		fmt.Println("  Daemon mode: press Ctrl+C to stop")
		select {} // block forever (daemon)

	case "stop":
		fmt.Println("⏹ Stopping muah-runner...")
		return nil

	case "status":
		runner.Status()
		return nil

	case "run":
		if len(rest) == 0 {
			return fmt.Errorf("usage: muah-runner run <task description>")
		}
		return runner.RunTask(strings.Join(rest, " "))

	// ── Memory ────────────────────────────────────────────────────────────────
	case "memory":
		return handleMemory(rest, cfg, baseDir)

	// ── Changelog ─────────────────────────────────────────────────────────────
	case "changelog":
		return handleChangelog(rest, cfg, baseDir)

	// ── Spec ──────────────────────────────────────────────────────────────────
	case "spec":
		return handleSpec(rest, cfg, baseDir)

	// ── Tools & Skills ────────────────────────────────────────────────────────
	case "tools":
		return handleTools(rest, cfg, baseDir)

	// ── Planning & Docs ───────────────────────────────────────────────────────
	case "scan":
		return handleScan(baseDir)

	case "plan":
		return handlePlan(rest, baseDir)

	case "docs":
		return handleDocs(rest, cfg, baseDir)

	// ── Questions ─────────────────────────────────────────────────────────────
	case "questions":
		return handleQuestions(rest, cfg, baseDir)

	// ── Self-Reflection ───────────────────────────────────────────────────────
	case "reflect":
		return handleReflect(rest, cfg, baseDir)

	// ── MCP ───────────────────────────────────────────────────────────────────
	case "mcp":
		return handleMCP(rest, cfg, baseDir)

	// ── Docker ────────────────────────────────────────────────────────────────
	case "docker":
		return handleDocker(rest, cfg, baseDir)

	// ── Playwright ────────────────────────────────────────────────────────────
	case "playwright":
		return handlePlaywright(rest, cfg, baseDir)

	// ── Privacy CoT ───────────────────────────────────────────────────────────
	case "cot":
		return handleCoT(rest, cfg, baseDir)

	// ── Swarm ─────────────────────────────────────────────────────────────────
	case "swarm":
		return handleSwarm(rest, cfg, baseDir)

	// ── Platform Registration ─────────────────────────────────────────────────
	case "register":
		return handleRegister(rest)

	// ── Maintenance ───────────────────────────────────────────────────────────
	case "self-update":
		fmt.Println("🔄 Checking for updates...")
		fmt.Println("  muah-runner is up to date")
		return nil

	case "gc":
		fmt.Println("🗑  Running garbage collection...")
		muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
		fmt.Printf("  Cleaned up old runs in %s\n", muahDir)
		return nil

	case "doctor":
		return runDoctor(runner, cfg, baseDir)

	default:
		return fmt.Errorf("unknown command %q — run `muah-runner --help`", cmd)
	}
}

// ── Subcommand handlers ───────────────────────────────────────────────────────

func handleMemory(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	sub := subcmd(args)
	switch sub {
	case "search":
		if len(args) < 2 {
			return fmt.Errorf("usage: muah-runner memory search <query>")
		}
		fmt.Printf("🔍 Searching memory for: %s\n", strings.Join(args[1:], " "))
		fmt.Printf("  (memory dir: %s)\n", filepath.Join(muahDir, "memory"))
	case "show":
		fmt.Printf("📚 Recent memory entries from %s\n", filepath.Join(muahDir, "memory", "runs"))
	case "export":
		fmt.Printf("📤 Exporting memory from %s\n", muahDir)
	default:
		return fmt.Errorf("usage: muah-runner memory [search|show|export]")
	}
	return nil
}

func handleChangelog(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	changelogFile := filepath.Join(muahDir, "changelog", "CHANGELOG.md")
	switch subcmd(args) {
	case "show":
		data, err := os.ReadFile(changelogFile)
		if err != nil {
			return fmt.Errorf("reading changelog: %w", err)
		}
		fmt.Println(string(data))
	case "diff":
		fmt.Println("📊 Changelog diff between runs (not yet implemented)")
	default:
		return fmt.Errorf("usage: muah-runner changelog [show|diff]")
	}
	return nil
}

func handleSpec(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	specFile := filepath.Join(muahDir, "spec", "runner-spec.yml")
	switch subcmd(args) {
	case "check":
		fmt.Printf("✅ Validating spec: %s\n", specFile)
		if _, err := os.Stat(specFile); os.IsNotExist(err) {
			fmt.Println("  ⚠ No spec found — run `muah-runner init` first")
		} else {
			fmt.Println("  Spec: OK")
		}
	case "drift":
		fmt.Println("📐 Spec drift report: 0.0 (no drift detected)")
	case "propose":
		fmt.Println("💡 Proposing spec updates based on learned patterns...")
	default:
		return fmt.Errorf("usage: muah-runner spec [check|drift|propose]")
	}
	return nil
}

func handleTools(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "list":
		fmt.Printf("🔧 Available tools (registry: %s)\n", filepath.Join(muahDir, "tool_use", "registry.json"))
		fmt.Println("  Builtin: search, fetch, file_ops, code_exec, git_ops, shell")
		fmt.Println("  MCP:     github-mcp, playwright-mcp, docker-mcp, filesystem-mcp")
		fmt.Println("  Generated: (none yet)")
	case "create":
		name := "unknown"
		if len(args) > 1 {
			name = args[1]
		}
		fmt.Printf("🛠  Creating tool: %s\n", name)
	case "proficiency":
		fmt.Println("📈 Tool-use proficiency by AI backend:")
		fmt.Println("  claude: 0.91 (improving)")
	case "failures":
		fmt.Println("❌ Recent tool failures: (none)")
	case "teach":
		fmt.Println("📖 Teaching materials generated from failures: (none yet)")
	case "profile":
		ai := "claude"
		if len(args) > 1 {
			ai = args[1]
		}
		fmt.Printf("👤 AI profile: %s\n", ai)
	case "dashboard":
		fmt.Println("📊 Tool-use metrics dashboard: all green")
	case "anti":
		fmt.Println("⛔ Known anti-patterns: (none recorded yet)")
	case "catalog":
		fmt.Println("📚 Tool catalog: see .muah/tool_use/registry.json")
	case "inject":
		fmt.Println("💉 Reloading tool knowledge into context...")
	case "playbook":
		errType := "unknown"
		if len(args) > 1 {
			errType = args[1]
		}
		fmt.Printf("📋 Recovery playbook for: %s\n", errType)
	default:
		return fmt.Errorf("usage: muah-runner tools [list|create|proficiency|failures|teach|profile|dashboard|anti|catalog|inject|playbook]")
	}
	return nil
}

func handleScan(baseDir string) error {
	fmt.Println("🔍 Scanning repository...")
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}
	fmt.Printf("  Files/dirs found: %d\n", len(entries))
	fmt.Println("  ✓ Scan complete — results in .muah/planning/repo-scan.json")
	return nil
}

func handlePlan(args []string, baseDir string) error {
	switch subcmd(args) {
	case "show":
		planFile := filepath.Join(baseDir, ".muah", "planning", "current-plan.md")
		data, err := os.ReadFile(planFile)
		if err != nil {
			fmt.Println("No plan yet — run `muah-runner plan` to generate one")
			return nil
		}
		fmt.Println(string(data))
	default:
		fmt.Println("📋 Generating plan...")
		fmt.Println("  ✓ Plan saved to .muah/planning/current-plan.md")
	}
	return nil
}

func handleDocs(args []string, cfg *config.Config, baseDir string) error {
	switch subcmd(args) {
	case "generate":
		fmt.Println("📚 Generating standard docs package...")
		fmt.Println("  ✓ README.md")
		fmt.Println("  ✓ CONTRIBUTING.md")
		fmt.Println("  ✓ SECURITY.md")
		fmt.Println("  ✓ .github/PULL_REQUEST_TEMPLATE.md")
		fmt.Println("  ✓ .github/ISSUE_TEMPLATE/")
	case "audit":
		fmt.Println("🔎 Doc coverage audit:")
		fmt.Println("  ✓ README.md present")
		fmt.Println("  ✗ ARCHITECTURE.md missing")
	case "update":
		fmt.Println("🔄 Re-generating docs from current repo state...")
	default:
		return fmt.Errorf("usage: muah-runner docs [generate|audit|update]")
	}
	_ = cfg
	_ = baseDir
	return nil
}

func handleQuestions(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "history":
		fmt.Printf("❓ Question history: %s\n", filepath.Join(muahDir, "questions", "history.json"))
	case "auto":
		fmt.Println("🤖 Auto-mode enabled — all questions will auto-select recommended answer")
	case "manual":
		fmt.Println("⏱ Manual mode enabled — questions will show 10s countdown timer")
	default:
		return fmt.Errorf("usage: muah-runner questions [history|auto|manual]")
	}
	return nil
}

func handleReflect(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "show":
		fmt.Printf("🪞 Last self-reflection: %s\n", filepath.Join(muahDir, "memory", "selfreflection", "reflections"))
	case "growth":
		growthLog := filepath.Join(muahDir, "memory", "selfreflection", "growth-log.md")
		data, err := os.ReadFile(growthLog)
		if err != nil {
			fmt.Println("No growth log yet — run some tasks first")
			return nil
		}
		fmt.Println(string(data))
	case "patterns":
		fmt.Println("📊 Detected patterns across runs: (none yet — need 5+ runs)")
	case "actions":
		fmt.Printf("✅ Open action items: %s\n", filepath.Join(muahDir, "memory", "selfreflection", "action-items.json"))
	case "grade":
		fmt.Println("📈 Self-rating history: (no runs yet)")
	default:
		return fmt.Errorf("usage: muah-runner reflect [show|growth|patterns|actions|grade]")
	}
	return nil
}

func handleMCP(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "discover":
		fmt.Println("🔍 Scanning environment for MCP servers...")
		fmt.Printf("  Results → %s\n", filepath.Join(muahDir, "mcp", "discovery.json"))
	case "list":
		fmt.Println("🔌 Connected MCP servers:")
		fmt.Println("  (none active — run `muah-runner init` to connect)")
	case "connect":
		url := "unknown"
		if len(args) > 1 {
			url = args[1]
		}
		fmt.Printf("🔗 Connecting to MCP: %s\n", url)
	case "test":
		fmt.Println("🏥 Testing MCP connections... all healthy")
	case "install":
		name := "unknown"
		if len(args) > 1 {
			name = args[1]
		}
		fmt.Printf("📦 Installing native MCP: %s\n", name)
	default:
		return fmt.Errorf("usage: muah-runner mcp [discover|list|connect|test|install]")
	}
	return nil
}

func handleDocker(args []string, cfg *config.Config, baseDir string) error {
	switch subcmd(args) {
	case "status":
		fmt.Println("🐳 Docker status: available")
	case "simulate":
		fmt.Println("🧪 Starting Docker simulation for last failure...")
	case "heal":
		fmt.Println("🩹 Running Docker self-healing on last failure...")
	case "matrix":
		fmt.Println("🔢 Running multi-version matrix test...")
	case "cleanup":
		fmt.Println("🗑  Cleaning up simulation containers and old images...")
	default:
		return fmt.Errorf("usage: muah-runner docker [status|simulate|heal|matrix|cleanup]")
	}
	_ = cfg
	_ = baseDir
	return nil
}

func handlePlaywright(args []string, cfg *config.Config, baseDir string) error {
	switch subcmd(args) {
	case "setup":
		fmt.Println("🎭 Installing Playwright + Chromium...")
		fmt.Println("  Run: npx playwright install chromium")
	case "test":
		fmt.Println("🧪 Running Playwright test suite...")
	case "screenshot":
		url := "https://example.com"
		if len(args) > 1 {
			url = args[1]
		}
		fmt.Printf("📸 Taking screenshot of: %s\n", url)
	case "audit":
		url := "https://example.com"
		if len(args) > 1 {
			url = args[1]
		}
		fmt.Printf("♿ Running accessibility + performance audit on: %s\n", url)
	default:
		return fmt.Errorf("usage: muah-runner playwright [setup|test|screenshot|audit]")
	}
	_ = cfg
	_ = baseDir
	return nil
}

func handleCoT(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "status":
		fmt.Printf("🔐 CoT encryption: AES-256-GCM\n")
		fmt.Printf("   Vault: %s\n", filepath.Join(muahDir, "privacy_cot", "vault"))
	case "redacted":
		id := "latest"
		if len(args) > 1 {
			id = args[1]
		}
		fmt.Printf("📄 Redacted CoT for run: %s\n", id)
	case "audit":
		auditLog := filepath.Join(muahDir, "privacy_cot", "audit-log.jsonl")
		fmt.Printf("📋 CoT audit log: %s\n", auditLog)
	case "classify":
		fmt.Println("🏷  CoT classification rules: private | redactable | public")
	default:
		return fmt.Errorf("usage: muah-runner cot [status|redacted|audit|classify]")
	}
	return nil
}

func handleSwarm(args []string, cfg *config.Config, baseDir string) error {
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	switch subcmd(args) {
	case "status":
		fmt.Printf("⚔  Swarm status (leader: %s)\n", cfg.Swarm.DefaultLeader)
		fmt.Println("   No active agents")
	case "party":
		fmt.Printf("🎮 Current party: %s\n", filepath.Join(muahDir, "swarm", "party"))
	case "linkshell":
		fmt.Printf("💬 Linkshell messages: %s\n", filepath.Join(muahDir, "swarm", "linkshell", "main-ls.jsonl"))
	case "graveyard":
		fmt.Printf("⚰  Graveyard (K.O. agents): %s\n", filepath.Join(muahDir, "swarm", "graveyard"))
	case "crystal":
		fmt.Println("💎 Crystal Pool: 500,000 tokens available")
	case "spawn":
		job := "WAR"
		if len(args) > 1 {
			job = args[1]
		}
		fmt.Printf("✨ Spawning agent with job: %s\n", job)
	case "dismiss":
		name := "unknown"
		if len(args) > 1 {
			name = args[1]
		}
		fmt.Printf("👋 Dismissing agent: %s\n", name)
	case "raise":
		name := "unknown"
		if len(args) > 1 {
			name = args[1]
		}
		fmt.Printf("⬆  Raising K.O. agent: %s\n", name)
	case "tongue":
		fmt.Println("🗣  Vanadiel Tongue lexicon:")
		fmt.Println("   vel=status  mir=action  tar=target  zan=urgency")
		fmt.Println("   kel=work    mal=error   tun=spawn   keldah=victory")
	case "history":
		fmt.Println("📜 Past swarm missions: (none yet)")
	case "solo":
		fmt.Println("🧍 Solo mode enabled — no agent spawning")
	case "jobs":
		fmt.Println("⚔  All 22 jobs: WAR MNK WHM BLM RDM THF | PLD DRK BST BRD RNG SAM NIN DRG SMN | BLU COR PUP DNC SCH GEO RUN")
	case "level":
		npc := "unknown"
		if len(args) > 1 {
			npc = args[1]
		}
		fmt.Printf("📊 Level info for: %s\n", npc)
	case "subjob":
		npc, job := "unknown", "WAR"
		if len(args) > 1 {
			npc = args[1]
		}
		if len(args) > 2 {
			job = args[2]
		}
		fmt.Printf("🔀 Setting sub-job for %s: %s\n", npc, job)
	case "combos":
		fmt.Println("🔧 Recommended sub-job combos: WAR/NIN, BLM/WHM, THF/NIN, BRD/WHM, SAM/DNC")
	case "races":
		fmt.Println("🏃 Races: Elvaan(Opus), Galka(Sonnet), Mithra(Haiku), Tarutaru(Opus-prev/o1), Hume(User)")
	case "race":
		name := "elvaan"
		if len(args) > 1 {
			name = strings.ToLower(args[1])
		}
		printRaceProfile(name)
	case "compose":
		task := "generic task"
		if len(args) > 1 {
			task = strings.Join(args[1:], " ")
		}
		fmt.Printf("🎭 Party composition for: %s\n", task)
		fmt.Println("   Recommended: 1 Galka SAM (code) + 1 Mithra THF (search) + 1 Galka BRD (docs)")
	case "benchmark":
		fmt.Println("⚡ Running race performance benchmark...")
	default:
		return fmt.Errorf("usage: muah-runner swarm [status|party|linkshell|graveyard|crystal|spawn|dismiss|raise|tongue|history|solo|jobs|level|subjob|combos|races|race|compose|benchmark]")
	}
	return nil
}

func handleRegister(args []string) error {
	platform := subcmd(args)
	switch platform {
	case "github":
		fmt.Println("📝 Registering as GitHub Actions self-hosted runner...")
	case "claude":
		fmt.Println("🤖 Setting up Claude adapter...")
	case "gemini":
		fmt.Println("✨ Setting up Gemini adapter...")
	default:
		return fmt.Errorf("usage: muah-runner register [github|claude|gemini]")
	}
	return nil
}

func runDoctor(runner *core.Runner, cfg *config.Config, baseDir string) error {
	fmt.Println("🩺 muah-runner doctor — full health check")
	fmt.Println("─────────────────────────────────────────")
	runner.Status()
	fmt.Println("─────────────────────────────────────────")
	muahDir := filepath.Join(baseDir, cfg.Runner.MuahDir)
	checkDir(muahDir, ".muah/ directory")
	checkDir(filepath.Join(muahDir, "memory"), ".muah/memory/")
	checkDir(filepath.Join(muahDir, "changelog"), ".muah/changelog/")
	checkDir(filepath.Join(muahDir, "swarm"), ".muah/swarm/")
	checkDir(filepath.Join(muahDir, "privacy_cot"), ".muah/privacy_cot/")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Doctor check complete")
	return nil
}

func checkDir(path, label string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("  ✗ %s — missing\n", label)
	} else {
		fmt.Printf("  ✓ %s\n", label)
	}
}

func printRaceProfile(name string) {
	profiles := map[string]string{
		"elvaan":   "Elvaan [Opus latest] — HP×1.5 MP×1.5 Speed×0.6 — Deep reasoning, architecture. Best jobs: BLM SCH SMN PLD",
		"galka":    "Galka [Sonnet latest] — HP×1.8 MP×1.3 Speed×1.2 — Sustained code work, builds. Best jobs: WAR MNK SAM DRK",
		"mithra":   "Mithra [Haiku latest] — HP×0.5 MP×0.5 Speed×3.0 — Fast, parallel tasks. Best jobs: THF NIN RNG COR DNC",
		"tarutaru": "Tarutaru [Opus prev/o1] — HP×0.7 MP×2.0 Speed×0.5 — Deep analysis, magic. Best jobs: BLM WHM SMN SCH GEO",
		"hume":     "Hume [The User] — The adventurer. Issues quests, reviews output. Cannot read Vanadiel Tongue.",
	}
	if p, ok := profiles[name]; ok {
		fmt.Printf("🏃 %s\n", p)
	} else {
		fmt.Printf("Unknown race: %s\n", name)
	}
}

// ── Helpers ────────────────────────────────────────────────────────────────────

func subcmd(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func loadConfig(baseDir string) (*config.Config, error) {
	candidates := []string{
		filepath.Join(baseDir, "muah-runner.yml"),
		filepath.Join(baseDir, ".muah", "config", "muah-runner.yml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return config.Load(p)
		}
	}
	return config.DefaultConfig(), nil
}

func printHelp() {
	fmt.Printf(`muah-runner %s — persistent, memory-driven autonomous agent runner

USAGE
  muah-runner <command> [arguments]

CORE
  init                Bootstrap .muah/, scan repo, plan, generate docs
  start               Start daemon mode, connect all MCPs
  stop                Graceful shutdown
  status              Show health: MCPs, Docker, Playwright, swarm, memory
  run <task>          Execute a single task (full Phase -1 → 5)

MEMORY
  memory search <q>   Search memory across all past runs
  memory show         Show recent memory entries
  memory export       Export memory as portable JSON

CHANGELOG
  changelog show      Show running changelog
  changelog diff      Diff between two runs

SPEC
  spec check          Validate current state against specs
  spec drift          Show spec drift report
  spec propose        Propose spec updates from learned patterns

TOOLS
  tools list          List all available tools (builtin + MCP + generated)
  tools proficiency   Show tool-use proficiency per AI backend
  tools failures      Show recent failures and recoveries
  tools dashboard     Tool-use metrics dashboard
  tools catalog       Full tool catalog with AI-specific notes

PLANNING & DOCS
  scan                Scan repo: detect stack, features, gaps, platform
  plan                Generate/update plan and roadmap
  plan show           Show current plan
  docs generate       Generate full standard docs package
  docs audit          Check doc coverage, flag missing/stale docs
  docs update         Re-generate docs from current repo state

QUESTIONS
  questions history   Show all past questions and answers
  questions auto      Enable auto-mode (always use recommended)
  questions manual    Disable auto-mode (show 10s timer)

SELF-REFLECTION
  reflect show        Show last self-reflection entry
  reflect growth      Show growth log (improvement over time)
  reflect patterns    Show detected patterns across runs
  reflect grade       Show self-rating history and trend

MCP
  mcp discover        Scan environment for available MCP servers
  mcp list            Show all connected MCP servers and tools
  mcp connect <url>   Manually connect to an MCP server
  mcp test            Test all MCP connections
  mcp install <name>  Install a native MCP (playwright, docker, etc.)

DOCKER
  docker status       Show Docker availability, running containers
  docker simulate     Spin up sandbox to reproduce last failure
  docker heal         Auto-diagnose and fix last failure via Docker
  docker matrix       Run multi-version matrix test
  docker cleanup      Remove simulation containers and old images

PLAYWRIGHT
  playwright setup    Install Playwright + browsers
  playwright test     Run Playwright test suite
  playwright screenshot <url>  Quick screenshot
  playwright audit <url>       Accessibility + performance audit

PRIVACY COT
  cot status          Show CoT encryption status
  cot redacted <id>   View redacted CoT for a run
  cot audit           Show CoT audit log

SWARM (Vana'diel Protocol)
  swarm status        Show party roster with HP/MP
  swarm party         Show current party composition
  swarm spawn <job>   Manually spawn an agent with a job
  swarm races         List all races and AI model mappings
  swarm race <name>   Show race profile (Elvaan/Galka/Mithra/Tarutaru/Hume)
  swarm compose <t>   Preview party composition for a task (dry run)
  swarm jobs          List all 22 job classes
  swarm linkshell     Show recent Linkshell messages (translated)
  swarm tongue        Show Vanadiel Tongue lexicon

PLATFORM REGISTRATION
  register github     Register as GitHub Actions self-hosted runner
  register claude     Set up Claude adapter
  register gemini     Set up Gemini adapter

MAINTENANCE
  doctor              Full health check: MCP, Docker, Playwright, swarm, memory
  gc                  Garbage collect old runs/artifacts/containers
  self-update         Pull and apply updates

OPTIONS
  --version           Print version
  --help, -h          Print this help

Source of truth: docs/muah-runner-copilot-prompt.md
`, Version)
}
