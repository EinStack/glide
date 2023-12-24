package telemetry

import "go.uber.org/zap"

type Config struct {
	LogConfig *LogConfig `mapstructure:"logging"`
	// TODO: add OTEL config
}

type Telemetry struct {
	Config *Config
	Logger *zap.Logger
	// TODO: add OTEL meter, tracer
}

func NewTelemetry(cfg *Config) (*Telemetry, error) {
	logger, err := NewLogger(cfg.LogConfig)
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		Config: cfg,
		Logger: logger,
	}, nil
}
