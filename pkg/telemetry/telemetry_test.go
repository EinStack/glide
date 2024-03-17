package telemetry

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTelemetry_Creation(t *testing.T) {
	_, err := NewTelemetry(DefaultConfig())
	require.NoError(t, err)
}
