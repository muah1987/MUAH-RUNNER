// Package core implements the main run loop and .muah/ directory bootstrap.
// It is the entry point for every muah-runner execution.
package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// MuahDir is the name of the persistent state directory.
const MuahDir = ".muah"

// Bootstrap creates the full .muah/ directory structure if it doesn't exist.
// This is idempotent — safe to call on every startup.
func Bootstrap(baseDir string) error {
	muah := filepath.Join(baseDir, MuahDir)

	dirs := []string{
		filepath.Join(muah, "memory", "runs"),
		filepath.Join(muah, "memory", "context"),
		filepath.Join(muah, "memory", "selfreflection", "reflections"),
		filepath.Join(muah, "tool_use", "ai_profiles"),
		filepath.Join(muah, "tool_use", "failure_log"),
		filepath.Join(muah, "tool_use", "recovery", "playbooks"),
		filepath.Join(muah, "tool_use", "teaching", "examples", "claude"),
		filepath.Join(muah, "tool_use", "teaching", "examples", "gemini"),
		filepath.Join(muah, "tool_use", "teaching", "examples", "generic"),
		filepath.Join(muah, "tool_use", "metrics"),
		filepath.Join(muah, "privacy_cot", "vault"),
		filepath.Join(muah, "privacy_cot", "redacted"),
		filepath.Join(muah, "swarm", "party", "members"),
		filepath.Join(muah, "swarm", "jobs", "levels"),
		filepath.Join(muah, "swarm", "races", "benchmarks"),
		filepath.Join(muah, "swarm", "linkshell", "whispers"),
		filepath.Join(muah, "swarm", "graveyard"),
		filepath.Join(muah, "swarm", "crystal"),
		filepath.Join(muah, "swarm", "vanadiel-tongue"),
		filepath.Join(muah, "changelog", "entries"),
		filepath.Join(muah, "spec", "task-specs", "custom"),
		filepath.Join(muah, "tools", "builtin"),
		filepath.Join(muah, "tools", "generated"),
		filepath.Join(muah, "tools", "skills"),
		filepath.Join(muah, "mcp", "platform"),
		filepath.Join(muah, "mcp", "native"),
		filepath.Join(muah, "mcp", "discovered"),
		filepath.Join(muah, "mcp", "logs"),
		filepath.Join(muah, "docker", "sandbox"),
		filepath.Join(muah, "docker", "images"),
		filepath.Join(muah, "docker", "simulations"),
		filepath.Join(muah, "docker", "compose-templates"),
		filepath.Join(muah, "planning", "plans"),
		filepath.Join(muah, "planning", "decisions"),
		filepath.Join(muah, "docs", "templates", "github"),
		filepath.Join(muah, "docs", "templates", "gitlab"),
		filepath.Join(muah, "docs", "templates", "generic"),
		filepath.Join(muah, "docs", "generated"),
		filepath.Join(muah, "questions"),
		filepath.Join(muah, "runs", "current"),
		filepath.Join(muah, "runs", "history"),
		filepath.Join(muah, "config", "adapters"),
		filepath.Join(muah, "config", "hooks"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", d, err)
		}
	}

	if err := writeSeedFiles(muah); err != nil {
		return fmt.Errorf("writing seed files: %w", err)
	}

	return nil
}

// writeSeedFiles creates default JSON/YAML seed files only if they don't exist.
func writeSeedFiles(muah string) error {
	seeds := map[string]interface{}{
		filepath.Join(muah, "memory", "index.json"): map[string]interface{}{
			"total_runs":  0,
			"last_run_id": "",
			"run_ids":     []string{},
			"last_update": time.Now().UTC(),
		},
		filepath.Join(muah, "tool_use", "registry.json"): map[string]interface{}{
			"registry_version": "1.0.0",
			"last_updated":     time.Now().UTC(),
			"tools":            []interface{}{},
		},
		filepath.Join(muah, "tool_use", "teaching", "lessons.json"): map[string]interface{}{
			"lessons": []interface{}{},
		},
		filepath.Join(muah, "tool_use", "teaching", "anti-patterns.json"): map[string]interface{}{
			"anti_patterns": []interface{}{},
		},
		filepath.Join(muah, "tool_use", "recovery", "strategies.json"): map[string]interface{}{
			"strategies": defaultRecoveryStrategies(),
		},
		filepath.Join(muah, "tool_use", "metrics", "dashboard.json"): map[string]interface{}{
			"total_tool_calls": 0,
			"success_rate":     0.0,
			"ai_breakdown":     map[string]interface{}{},
		},
		filepath.Join(muah, "swarm", "vana_diel.json"): map[string]interface{}{
			"world":         "vana_diel",
			"active_agents": []interface{}{},
			"missions":      []interface{}{},
		},
		filepath.Join(muah, "swarm", "party", "party-leader.json"): map[string]interface{}{
			"npc_name": "Prishe",
			"job":      "party-leader",
			"status":   "standby",
		},
		filepath.Join(muah, "swarm", "jobs", "job-registry.json"): defaultJobRegistry(),
		filepath.Join(muah, "swarm", "races", "race-registry.json"): defaultRaceRegistry(),
		filepath.Join(muah, "swarm", "crystal", "pool.json"): map[string]interface{}{
			"total_tokens_budget": 500000,
			"allocated":           map[string]interface{}{},
			"unallocated":         map[string]interface{}{"tokens": 500000},
		},
		filepath.Join(muah, "changelog", "CHANGELOG.md"): "# Changelog\n\nAll changes made by muah-runner are logged here.\n",
		filepath.Join(muah, "spec", "runner-spec.yml"): defaultRunnerSpec(),
		filepath.Join(muah, "spec", "constraints.yml"): defaultConstraints(),
		filepath.Join(muah, "spec", "quality-gates.yml"): defaultQualityGates(),
		filepath.Join(muah, "mcp", "discovery.json"): map[string]interface{}{
			"last_scan":      nil,
			"servers":        []interface{}{},
			"platform":       "",
			"native_checked": false,
		},
		filepath.Join(muah, "mcp", "connections.json"): map[string]interface{}{
			"active":      []interface{}{},
			"last_update": time.Now().UTC(),
		},
		filepath.Join(muah, "planning", "repo-scan.json"): map[string]interface{}{
			"scanned_at":  nil,
			"tech_stack":  []string{},
			"languages":   []string{},
			"platform":    "unknown",
			"gaps":        []string{},
		},
		filepath.Join(muah, "docs", "doc-manifest.json"): map[string]interface{}{
			"platform":    "unknown",
			"docs_exist":  []string{},
			"docs_needed": []string{},
			"last_audit":  nil,
		},
		filepath.Join(muah, "questions", "pending.json"): map[string]interface{}{
			"questions": []interface{}{},
		},
		filepath.Join(muah, "questions", "history.json"): []interface{}{},
		filepath.Join(muah, "questions", "defaults.yml"):  "# Default answers for recurring questions\n# question: recommended_answer_index\n",
		filepath.Join(muah, "questions", "config.yml"):    defaultQuestionsConfig(),
		filepath.Join(muah, "runs", "run-log.jsonl"):      "",
		filepath.Join(muah, "config", "muah-runner.yml"):  defaultRunnerConfig(),
		filepath.Join(muah, "privacy_cot", "config.yml"):  defaultCotConfig(),
		filepath.Join(muah, "privacy_cot", "classification.yml"): defaultCotClassification(),
		filepath.Join(muah, "privacy_cot", "audit-log.jsonl"):    "",
		filepath.Join(muah, ".muah-lock"): map[string]interface{}{
			"locked": false,
			"pid":    0,
		},
	}

	for path, content := range seeds {
		// Don't overwrite existing files
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func writeFile(path string, content interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var data []byte
	var err error
	switch v := content.(type) {
	case string:
		data = []byte(v)
	default:
		ext := filepath.Ext(path)
		if ext == ".yml" || ext == ".yaml" {
			data, err = yaml.Marshal(content)
		} else {
			data, err = json.MarshalIndent(content, "", "  ")
		}
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

func defaultRecoveryStrategies() []map[string]string {
	return []map[string]string{
		{"error_type": "auth_failure", "strategy": "refresh_token_then_fallback_cli"},
		{"error_type": "rate_limit", "strategy": "exponential_backoff_then_cache"},
		{"error_type": "not_found", "strategy": "list_directory_fuzzy_match"},
		{"error_type": "invalid_params", "strategy": "reload_spec_fix_params_retry"},
		{"error_type": "timeout", "strategy": "retry_shorter_then_async"},
		{"error_type": "hallucinated_tool", "strategy": "search_registry_suggest_correct"},
	}
}

func defaultJobRegistry() map[string]interface{} {
	return map[string]interface{}{
		"standard_jobs": []string{"WAR", "MNK", "WHM", "BLM", "RDM", "THF"},
		"advanced_jobs": []string{"PLD", "DRK", "BST", "BRD", "RNG", "SAM", "NIN", "DRG", "SMN"},
		"expert_jobs":   []string{"BLU", "COR", "PUP", "DNC", "SCH", "GEO", "RUN"},
		"total":         22,
	}
}

func defaultRaceRegistry() map[string]interface{} {
	return map[string]interface{}{
		"races": []map[string]interface{}{
			{"name": "Elvaan", "ai_tier": "opus_latest", "primary_model": "claude-opus-4-6"},
			{"name": "Galka", "ai_tier": "sonnet_latest", "primary_model": "claude-sonnet-4-6"},
			{"name": "Mithra", "ai_tier": "haiku_latest", "primary_model": "claude-haiku-4-5"},
			{"name": "Tarutaru", "ai_tier": "opus_previous", "primary_model": "claude-opus-4"},
			{"name": "Hume", "ai_tier": "human", "primary_model": "human"},
		},
	}
}

func defaultRunnerSpec() string {
	return `name: muah-runner
version: 1.0.0
max_concurrent_tasks: 4
max_run_duration: 30m
max_retries: 3
retry_backoff: exponential
memory_budget: 512MB
log_level: info
auto_cleanup_after_days: 30
spec_enforcement: strict
`
}

func defaultConstraints() string {
	return `# Hard limits for muah-runner
timeout_seconds: 1800
max_memory_mb: 512
max_retries: 3
max_parallel_agents: 18
max_party_size: 6
question_timeout_seconds: 10
docker_memory_limit: 512m
docker_cpu_limit: "1.0"
`
}

func defaultQualityGates() string {
	return `gates:
  - name: tests-pass
    required: true
    command: "go test ./..."
    expected_exit: 0
  - name: build-succeeds
    required: true
    command: "go build ./..."
    expected_exit: 0
  - name: spec-compliance
    required: true
    max_drift_score: 0.1
  - name: memory-updated
    required: true
    check: ".muah/memory/latest.json exists and was updated"
`
}

func defaultQuestionsConfig() string {
	return `timer_seconds: 10
min_timer: 5
max_timer: 30
auto_mode: false
silent_mode: false
escalation:
  on_critical_question: extend_to_30s
  on_repeated_timeout: suggest_auto_mode
`
}

func defaultRunnerConfig() string {
	return `runner:
  name: "muah-runner-01"
  muah_dir: ".muah"
  concurrency: 4
  heartbeat_interval: 30

server:
  host: "0.0.0.0"
  port: 8080

swarm:
  enabled: true
  max_party_size: 6
  default_leader: "Prishe"

questions:
  timer_seconds: 10
  auto_mode: false

selfreflection:
  enabled: true
  pattern_detection_interval: 5

mcp:
  auto_discover: true
  require_playwright: true
  require_docker: true
`
}

func defaultCotConfig() string {
	return `enabled: true
encryption: AES-256-GCM
per_run: true
at_rest: true
share_reflections: false
`
}

func defaultCotClassification() string {
	return `classification_levels:
  private:
    - internal_reasoning
    - hypothesis_generation
    - self_doubt
    - inter_agent_messages
    - secret_values
    - vulnerability_details
    - swarm_coordination
    - cost_calculations
    - ai_profile_reasoning

  redactable:
    - tool_selection_reasoning
    - error_analysis
    - recovery_decisions
    - spec_compliance_checks
    - planning_rationale

  public:
    - actions_taken
    - results
    - changelog_entries
    - questions_asked
    - self_reflection_summary
`
}
