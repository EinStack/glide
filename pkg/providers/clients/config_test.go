package clients

import (
    "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

func TestClientConfig_DefaultConfig(t *testing.T) {
	config := DefaultClientConfig()

	require.NotEmpty(t, config.Timeout)
}

func TestClientConfig_JSONMarshal(t *testing.T) {
    defaultConfig := DefaultClientConfig()

    expectedJSON := `{
        "timeout": "10s"
    }`

	marshaledJSON, err := json.MarshalIndent(defaultConfig, "", "\t")
    
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(marshaledJSON))
}

