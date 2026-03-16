package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds the full application configuration.
type Config struct {
	Runner   RunnerConfig   `yaml:"runner"`
	Server   ServerConfig   `yaml:"server"`
	Executor ExecutorConfig `yaml:"executor"`
	Health   HealthConfig   `yaml:"health"`
	Webhooks WebhookConfig  `yaml:"webhooks"`
	Secrets  SecretsConfig  `yaml:"secrets"`
}

type RunnerConfig struct {
	Name              string   `yaml:"name"`
	Labels            []string `yaml:"labels"`
	WorkDir           string   `yaml:"work_dir"`
	Concurrency       int      `yaml:"concurrency"`
	HeartbeatInterval int      `yaml:"heartbeat_interval"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TLS  bool   `yaml:"tls"`
}

type ExecutorConfig struct {
	Type         string `yaml:"type"`
	DockerImage  string `yaml:"docker_image"`
	Timeout      int    `yaml:"timeout"`
	MaxRetries   int    `yaml:"max_retries"`
	RetryBackoff int    `yaml:"retry_backoff"`
}

type HealthConfig struct {
	MetricsPort   int `yaml:"metrics_port"`
	CheckInterval int `yaml:"check_interval"`
}

type WebhookConfig struct {
	Enabled bool     `yaml:"enabled"`
	URLs    []string `yaml:"urls"`
	Secret  string   `yaml:"secret"`
}

type SecretsConfig struct {
	EncryptionKey string `yaml:"encryption_key"`
}

// Defaults returns a Config populated with sensible defaults.
func Defaults() *Config {
	return &Config{
		Runner: RunnerConfig{
			Name:              "muah-runner-01",
			Labels:            []string{},
			WorkDir:           "/tmp/muah-runner",
			Concurrency:       4,
			HeartbeatInterval: 30,
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			TLS:  false,
		},
		Executor: ExecutorConfig{
			Type:         "process",
			DockerImage:  "ubuntu:22.04",
			Timeout:      3600,
			MaxRetries:   3,
			RetryBackoff: 5,
		},
		Health: HealthConfig{
			MetricsPort:   9090,
			CheckInterval: 15,
		},
		Webhooks: WebhookConfig{
			Enabled: false,
			URLs:    []string{},
			Secret:  "",
		},
		Secrets: SecretsConfig{
			EncryptionKey: "",
		},
	}
}

// Load reads a YAML config file and overrides values with environment variables.
// Returns defaults if the file is not found.
func Load(path string) (*Config, error) {
	cfg := Defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			applyEnv(cfg)
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	applyEnv(cfg)
	return cfg, nil
}

// applyEnv overrides config fields from environment variables.
func applyEnv(cfg *Config) {
	if v := os.Getenv("MUAH_RUNNER_NAME"); v != "" {
		cfg.Runner.Name = v
	}
	if v := os.Getenv("MUAH_RUNNER_WORK_DIR"); v != "" {
		cfg.Runner.WorkDir = v
	}
	if v := os.Getenv("MUAH_RUNNER_CONCURRENCY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Runner.Concurrency = n
		}
	}
	if v := os.Getenv("MUAH_SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("MUAH_SERVER_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = n
		}
	}
	if v := os.Getenv("MUAH_EXECUTOR_TYPE"); v != "" {
		cfg.Executor.Type = v
	}
	if v := os.Getenv("MUAH_HEALTH_METRICS_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Health.MetricsPort = n
		}
	}
	if v := os.Getenv("MUAH_WEBHOOK_SECRET"); v != "" {
		cfg.Webhooks.Secret = v
	}
	if v := os.Getenv("MUAH_SECRETS_ENCRYPTION_KEY"); v != "" {
		cfg.Secrets.EncryptionKey = v
	}
}
