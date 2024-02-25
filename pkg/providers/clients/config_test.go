package clients

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientConfig_DefaultConfig(t *testing.T) {
	config := DefaultClientConfig()

	require.NotEmpty(t, config.Timeout)
}
