package config

import (
	"glide/pkg/api"
	"glide/pkg/telemetry"
)

// Config is a general top-level Glide configuration
type Config struct {
	Telemetry *telemetry.Config `mapstructure:"telemetry"`
	API       *api.Config       `mapstructure:"api"`
}
