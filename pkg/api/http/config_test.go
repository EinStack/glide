package http

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPConfig_DefaultConfig(t *testing.T) {
	config := DefaultServerConfig()

	require.NotNil(t, config.Address())
	require.NotNil(t, config.ToServer())
}
