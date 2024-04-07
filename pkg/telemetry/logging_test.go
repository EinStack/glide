package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogging_PlainOutputSetup(t *testing.T) {
	config := LogConfig{
		Encoding: "console",
	}
	zapConfig := config.ToZapConfig()

	require.Equal(t, "console", config.Encoding)
	require.NotNil(t, zapConfig)
	require.Equal(t, "console", zapConfig.Encoding)
}

func TestLogging_JSONOutputSetup(t *testing.T) {
	config := DefaultLogConfig()
	zapConfig := config.ToZapConfig()

	require.Equal(t, "json", config.Encoding)
	require.NotNil(t, zapConfig)
	require.Equal(t, "json", zapConfig.Encoding)
}
