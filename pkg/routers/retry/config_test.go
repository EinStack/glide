package retry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryConfig_DefaultConfig(t *testing.T) {
	config := DefaultExpRetryConfig()

	require.NotNil(t, config)
}

func TestExpRetryConfig_MarshalJSON(t *testing.T) {
	minDelay := 2 * time.Second
	maxDelay := 5 * time.Minute

	config := &ExpRetryConfig{
		MaxRetries:     4,
		BaseMultiplier: 2,
		MinDelay:       minDelay,
		MaxDelay:       &maxDelay,
	}

	expectedJSON := `{
		"max_retries": 4,
		"base_multiplier": 2,
		"min_delay": "2s",
		"max_delay": "5m0s"
	}`

	marshaledJSON, err := json.MarshalIndent(config, "", "\t")
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(marshaledJSON))
}
