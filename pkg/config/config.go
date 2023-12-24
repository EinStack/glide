package config

import (
	"glide/pkg/api"
	"glide/pkg/telemetry"
)

// Config is a general top-level Glide configuration
type Config struct {
	Telemetry *telemetry.Config `json:"telemetry" yaml:"telemetry"`
	API       *api.Config       `json:"api" yaml:"api"`
}
