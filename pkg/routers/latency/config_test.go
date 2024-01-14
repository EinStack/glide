package latency

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLatencyConfig_Default(t *testing.T) {
	config := DefaultConfig()

	require.NotEmpty(t, config)
}
