package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPAddr    string `yaml:"http_addr"`
	DBPath      string `yaml:"db_path"`
	Timezone    string `yaml:"timezone"`
	PublishTime string `yaml:"publish_time"`
}

// Load applies precedence: env > file > defaults.
func Load() Config {
	cfg := fromFile()
	// defaults
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "poppo.db"
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "Local"
	}
	if cfg.PublishTime == "" {
		cfg.PublishTime = "08:00"
	}

	// env overrides
	if v := os.Getenv("PP_HTTP_ADDR"); v != "" {
		cfg.HTTPAddr = v
	}
	if v := os.Getenv("PP_DB_PATH"); v != "" {
		cfg.DBPath = v
	}
	if v := os.Getenv("PP_TZ"); v != "" {
		cfg.Timezone = v
	}
	if v := os.Getenv("PP_PUBLISH_TIME"); v != "" {
		cfg.PublishTime = v
	}
	return cfg
}

func fromFile() Config {
	var cfg Config
	path := os.Getenv("PP_CONFIG")
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".config", "poppo-press", "config.yaml")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = yaml.Unmarshal(b, &cfg)
	return cfg
}
