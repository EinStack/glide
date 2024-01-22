package config

import (
	"glide/pkg/api"
	"glide/pkg/routers"
	"glide/pkg/telemetry"
)

// Config is a general top-level Glide configuration
type Config struct {
	Telemetry *telemetry.Config `yaml:"telemetry" validate:"required"`
	API       *api.Config       `yaml:"api" validate:"required"`
	Routers   routers.Config    `yaml:"routers" validate:"required"`
}

func DefaultConfig() *Config {
	return &Config{
		Telemetry: telemetry.DefaultConfig(),
		API:       api.DefaultConfig(),
		// Routers should be defined by users
	}
}
