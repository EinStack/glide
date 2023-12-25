package telemetry

import "go.uber.org/zap"

type Config struct {
	LogConfig *LogConfig `yaml:"logging"`
	// TODO: add OTEL config
}

type Telemetry struct {
	Config *Config
	Logger *zap.Logger
	// TODO: add OTEL meter, tracer
}

func DefaultConfig() *Config {
	return &Config{
		LogConfig: DefaultLogConfig(),
	}
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
