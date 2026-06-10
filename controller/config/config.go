package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen           string `yaml:"listen"`
	Database         string `yaml:"database"`
	AdminToken       string `yaml:"admin_token"`
	AgentTokenSecret string `yaml:"agent_token_secret"`
	TLSCert          string `yaml:"tls_cert"`
	TLSKey           string `yaml:"tls_key"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Allow env var overrides
	if v := os.Getenv("ADMIN_TOKEN"); v != "" {
		cfg.AdminToken = v
	}
	if v := os.Getenv("AGENT_TOKEN_SECRET"); v != "" {
		cfg.AgentTokenSecret = v
	}
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		Listen:           "0.0.0.0:8080",
		Database:         "/var/lib/ddos/controller.db",
		AdminToken:       os.Getenv("ADMIN_TOKEN"),
		AgentTokenSecret: os.Getenv("AGENT_TOKEN_SECRET"),
	}
}
