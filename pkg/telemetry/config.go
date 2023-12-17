package telemetry

import "go.uber.org/zap"

type TelemetryConfig struct {
	LogConfig *zap.Config `json:"logs" yaml:"logs"`
}
