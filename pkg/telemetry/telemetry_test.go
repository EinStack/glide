package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
)

func TestTelemetry_Creation(t *testing.T) {
	_, err := NewTelemetry(DefaultConfig())
	require.NoError(t, err)
}

func TestDefaults_noopExporters(t *testing.T) {
	// By default all otel providers must be noop. Since we don't have otel setup
	// in test environment, this test ensures all providers are noop.
	mr, err := newMetricReader()
	if err != nil {
		t.Fatal(err)
	}
	// ensures we have a noop metric.ManualReader
	mr.(*metric.ManualReader).Shutdown(context.Background())

	se, err := newSpanExporter()
	if err != nil {
		t.Fatal(err)
	}
	// ensures we have a noopSpanExporter
	se.(noopSpanExporter).Shutdown(context.Background())
}
