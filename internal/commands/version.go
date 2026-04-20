package commands

import (
	"github.com/menloresearch/cli/internal/config"
)

func versionString() string {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	version := cfg.Version
	if version == "" {
		return "dev"
	}

	return version
}
