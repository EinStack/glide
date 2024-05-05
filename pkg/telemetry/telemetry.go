package telemetry

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

type Config struct {
	LogConfig *LogConfig        `yaml:"logging" validate:"required"`
	Resource  map[string]string `yaml:"resource"`
}

type Telemetry struct {
	Config *Config
	Logger *zap.Logger
}

func (t Telemetry) L() *zap.Logger {
	return t.Logger
}

func DefaultConfig() *Config {
	instance := os.Getenv("POD_NAME")
	if instance == "" {
		instance = uuid.New().String()
	}
	return &Config{
		LogConfig: DefaultLogConfig(),
		Resource: map[string]string{
			string(semconv.ServiceNameKey):       "glide",
			string(semconv.ServiceInstanceIDKey): instance,
		},
	}
}

func NewTelemetry(cfg *Config) (*Telemetry, error) {
	logger, err := NewLogger(cfg.LogConfig)
	if err != nil {
		return nil, err
	}

	resourceAttr := make([]attribute.KeyValue, 0, len(cfg.Resource))

	for k, v := range cfg.Resource {
		resourceAttr = append(resourceAttr, attribute.String(k, v))
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		resourceAttr...,
	)

	spanExporter, err := autoexport.NewSpanExporter(context.Background())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
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
		sdkmetric.WithResource(resource),
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
