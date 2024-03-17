package retry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetryConfig_DefaultConfig(t *testing.T) {
	config := DefaultExpRetryConfig()

	require.NotNil(t, config)
}
