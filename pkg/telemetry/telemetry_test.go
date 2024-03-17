package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTelemetry_Creation(t *testing.T) {
	_, err := NewTelemetry(DefaultConfig())
	require.NoError(t, err)
}
