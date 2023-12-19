package telemetry

import "go.uber.org/zap"

type TelemetryConfig struct {
	LogConfig *LogConfig `json:"logs" yaml:"logs"`
	// TODO: add OTEL config
}

type Telemetry struct {
	logger *zap.Logger
	// TODO: add OTEL meter, tracer
}

func NewTelemetry(cfg *TelemetryConfig) (*Telemetry, error) {
	// TODO: gonna be read from a config file
	logConfig := NewLogConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConfig.Encoding = "console"

	logger, err := telemetry.NewLogger(logConfig)

	if err != nil {
		return nil, err
	}

	return &Telemetry{
		logger: logger,
	}
}

func (t *Telemetry) Logger() *zap.Logger {
	return t.logger
}
