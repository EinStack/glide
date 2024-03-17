package telemetry

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogging_PlainOutputSetup(t *testing.T) {
	config := LogConfig{
		Encoding: "console",
	}
	zapConfig := config.ToZapConfig()

	require.Equal(t, config.Encoding, "console")
	require.NotNil(t, zapConfig)
	require.Equal(t, zapConfig.Encoding, "console")
}

func TestLogging_JSONOutputSetup(t *testing.T) {
	config := DefaultLogConfig()
	zapConfig := config.ToZapConfig()

	require.Equal(t, config.Encoding, "json")
	require.NotNil(t, zapConfig)
	require.Equal(t, zapConfig.Encoding, "json")
}
