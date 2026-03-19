package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrConfigNotFound = errors.New("config file not found")

func IsNotExist(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, os.ErrNotExist) || errors.Is(err, ErrConfigNotFound)
}

type Config struct {
	APIKey         string `yaml:"api_key"`
	PlatformURL    string `yaml:"platform_url"`
	DefaultRobotID string `yaml:"default_robot_id"`
}

func (c *Config) SetDefaults() {
	if c.PlatformURL == "" {
		c.PlatformURL = "https://api.menlo.ai/"
	}
}

func ConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "menlo-cli"), nil
}

func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.SetDefaults()
	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		PlatformURL: "https://api.menlo.ai/",
	}
}

func EnsureConfigDir() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func Marshal(c *Config) ([]byte, error) {
	return yaml.Marshal(c)
}
