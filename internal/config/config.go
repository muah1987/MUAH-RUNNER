package config

import (
"os"

"gopkg.in/yaml.v3"
)

type Config struct {
Runner        RunnerConfig        `yaml:"runner"`
Server        ServerConfig        `yaml:"server"`
Executor      ExecutorConfig      `yaml:"executor"`
Health        HealthConfig        `yaml:"health"`
Swarm         SwarmConfig         `yaml:"swarm"`
Questions     QuestionsConfig     `yaml:"questions"`
SelfReflection SelfReflectionConfig `yaml:"selfreflection"`
MCP           MCPConfig           `yaml:"mcp"`
}

type RunnerConfig struct {
Name              string `yaml:"name"`
MuahDir           string `yaml:"muah_dir"`
Concurrency       int    `yaml:"concurrency"`
HeartbeatInterval int    `yaml:"heartbeat_interval"`
}

type ServerConfig struct {
Host string `yaml:"host"`
Port int    `yaml:"port"`
}

type ExecutorConfig struct {
Type       string `yaml:"type"`
Timeout    int    `yaml:"timeout"`
MaxRetries int    `yaml:"max_retries"`
}

type HealthConfig struct {
MetricsPort int `yaml:"metrics_port"`
}

type SwarmConfig struct {
Enabled       bool   `yaml:"enabled"`
MaxPartySize  int    `yaml:"max_party_size"`
DefaultLeader string `yaml:"default_leader"`
}

type QuestionsConfig struct {
TimerSeconds int  `yaml:"timer_seconds"`
AutoMode     bool `yaml:"auto_mode"`
}

type SelfReflectionConfig struct {
Enabled                  bool `yaml:"enabled"`
PatternDetectionInterval int  `yaml:"pattern_detection_interval"`
}

type MCPConfig struct {
AutoDiscover     bool `yaml:"auto_discover"`
RequirePlaywright bool `yaml:"require_playwright"`
RequireDocker    bool `yaml:"require_docker"`
}

func DefaultConfig() *Config {
return &Config{
Runner: RunnerConfig{
Name:              "muah-runner-01",
MuahDir:           ".muah",
Concurrency:       4,
HeartbeatInterval: 30,
},
Server: ServerConfig{
Host: "0.0.0.0",
Port: 8080,
},
Executor: ExecutorConfig{
Type:       "process",
Timeout:    3600,
MaxRetries: 3,
},
Health: HealthConfig{
MetricsPort: 9090,
},
Swarm: SwarmConfig{
Enabled:       true,
MaxPartySize:  6,
DefaultLeader: "Prishe",
},
Questions: QuestionsConfig{
TimerSeconds: 10,
AutoMode:     false,
},
SelfReflection: SelfReflectionConfig{
Enabled:                  true,
PatternDetectionInterval: 5,
},
MCP: MCPConfig{
AutoDiscover:     true,
RequirePlaywright: true,
RequireDocker:    true,
},
}
}

func Load(path string) (*Config, error) {
cfg := DefaultConfig()
data, err := os.ReadFile(path)
if err != nil {
if os.IsNotExist(err) {
return cfg, nil
}
return nil, err
}
if err := yaml.Unmarshal(data, cfg); err != nil {
return nil, err
}
return cfg, nil
}

func Save(cfg *Config, path string) error {
data, err := yaml.Marshal(cfg)
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}
