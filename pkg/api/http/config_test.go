package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPConfig_DefaultConfig(t *testing.T) {
	config := DefaultServerConfig()

	require.NotNil(t, config.Address())
	require.NotNil(t, config.ToServer())
}
