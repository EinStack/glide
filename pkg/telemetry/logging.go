package telemetry

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	// Level is the minimum enabled logging level.
	Level zapcore.Level `yaml:"level"`

	// Encoding sets the logger's encoding. Valid values are "json", "console"
	Encoding string `yaml:"encoding"`

	// DisableCaller stops annotating logs with the calling function's file name and line number.
	// By default, all logs are annotated.
	DisableCaller bool `yaml:"disable_caller"`

	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `yaml:"disable_stacktrace"`

	// OutputPaths is a list of URLs or file paths to write logging output to.
	OutputPaths []string `yaml:"output_paths"`

	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `yaml:"initial_fields"`
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:             zap.InfoLevel,
		Encoding:          "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
		InitialFields:     make(map[string]interface{}),
	}
}

func (c *LogConfig) ToZapConfig() *zap.Config {
	zapConfig := zap.NewProductionConfig()

	if c.Encoding == "console" {
		zapConfig = zap.NewDevelopmentConfig()

		// Human-readable timestamps for console format of logs.
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// Colorized plain console logs
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig.Level = zap.NewAtomicLevelAt(c.Level)
	zapConfig.DisableCaller = c.DisableCaller
	zapConfig.DisableStacktrace = c.DisableStacktrace
	zapConfig.OutputPaths = c.OutputPaths
	zapConfig.InitialFields = c.InitialFields

	return &zapConfig
}

func NewLogger(cfg *LogConfig) (*zap.Logger, error) {
	zapConfig := cfg.ToZapConfig()

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
