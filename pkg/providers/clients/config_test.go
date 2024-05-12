package clients

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientConfig_DefaultConfig(t *testing.T) {
	config := DefaultClientConfig()

	require.NotEmpty(t, config.Timeout)
}

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	require.NotNil(t, config, "Config must not be nil")
	require.NotNil(t, config.Timeout, "Timeout must not be nil")
	require.NotNil(t, config.MaxIdleConns, "MaxIdleConns must not be nil")
	require.NotNil(t, config.MaxIdleConnsPerHost, "MaxIdleConnsPerHost must not be nil")

	// Check default timeout
	expectedTimeout := 10 * time.Second
	require.Equal(t, expectedTimeout, *config.Timeout)

	// Check MaxIdleConns
	expectedMaxIdleConns := 100
	require.Equal(t, expectedMaxIdleConns, *config.MaxIdleConns)

	// Check MaxIdleConnsPerHost
	expectedMaxIdleConnsPerHost := 2
	require.Equal(t, expectedMaxIdleConnsPerHost, *config.MaxIdleConnsPerHost)
}
