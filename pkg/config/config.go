package config

import "glide/pkg/telemetry"

// Config is a general top-level Glide configuration
type Config struct {
	Telemetry *telemetry.Config `json:"telemetry" yaml:"telemetry"`
}
