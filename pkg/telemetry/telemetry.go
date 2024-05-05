package telemetry

import (
	"context"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

var Resource = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceName("glide"),
)

type Config struct {
	LogConfig *LogConfig `yaml:"logging" validate:"required"`
}

type Telemetry struct {
	Config *Config
	Logger *zap.Logger
	// TODO: add OTEL meter, tracer
}

func (t Telemetry) L() *zap.Logger {
	return t.Logger
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
	spanExporter, err := autoexport.NewSpanExporter(context.Background())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(Resource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(spanExporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	metricsReader, err := autoexport.NewMetricReader(context.Background())
	if err != nil {
		return nil, err
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			metricsReader,
		),
		sdkmetric.WithResource(Resource),
	)

	otel.SetMeterProvider(provider)

	return &Telemetry{
		Config: cfg,
		Logger: logger,
	}, nil
}

func NewLoggerMock() *zap.Logger {
	return zap.NewNop()
}

// NewTelemetryMock returns Telemetry object with NoOp loggers, meters, tracers
func NewTelemetryMock() *Telemetry {
	return &Telemetry{
		Config: DefaultConfig(),
		Logger: NewLoggerMock(),
	}
}
