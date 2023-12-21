package telemetry

import "go.uber.org/zap"

type Config struct {
	LogConfig *LogConfig `json:"logs" yaml:"logs"`
	// TODO: add OTEL config
}

type Telemetry struct {
	logger *zap.Logger
	// TODO: add OTEL meter, tracer
}

func NewTelemetry(cfg *Config) (*Telemetry, error) {
	logger, err := NewLogger(cfg.LogConfig)
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		logger: logger,
	}, nil
}

func (t *Telemetry) Logger() *zap.Logger {
	return t.logger
}
