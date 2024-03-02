package clients

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientConfig_DefaultConfig(t *testing.T) {
	config := DefaultClientConfig()

	require.NotEmpty(t, config.Timeout)
}
