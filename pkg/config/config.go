package config

import (
	"glide/pkg/api"
	"glide/pkg/telemetry"
)

// Config is a general top-level Glide configuration
type Config struct {
	Telemetry *telemetry.Config `yaml:"telemetry"`
	API       *api.Config       `yaml:"api"`
}

func DefaultConfig() *Config {
	return &Config{
		Telemetry: telemetry.DefaultConfig(),
		API:       api.DefaultConfig(),
	}
}
