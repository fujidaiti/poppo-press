package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Output struct {
	Pager string `yaml:"pager"`
}

type Config struct {
	Server   string `yaml:"server"`
	Token    string `yaml:"token"`
	Timezone string `yaml:"timezone"`
	Output   Output `yaml:"output"`
}

func Path() (string, error) {
	// Windows: %APPDATA%/Poppo Press/config.yaml
	if dir := os.Getenv("APPDATA"); dir != "" && isWindows() {
		return filepath.Join(dir, "Poppo Press", "config.yaml"), nil
	}
	// XDG if set
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "poppo-press", "config.yaml"), nil
	}
	// default: ~/.config/poppo-press/config.yaml
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("HOME not set")
	}
	return filepath.Join(home, ".config", "poppo-press", "config.yaml"), nil
}

func Save(c *Config) error {
	p, err := Path()
	if err != nil {
		return err
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	// write with 0600
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return err
	}
	// ensure perms exactly on existing file/dir
	_ = os.Chmod(dir, 0o700)
	_ = os.Chmod(p, 0o600)
	return nil
}

func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func isWindows() bool {
	// crude check based on environment; avoids importing runtime here
	return os.PathSeparator == '\\'
}
