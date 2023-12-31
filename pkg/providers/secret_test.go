package providers

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecret_OpaqueOnMarshaling(t *testing.T) {
	name := "OpenAI"
	secretValue := "ABCDE123"

	config := struct {
		APIKey Secret `json:"api_key"`
		Name   string `json:"name"`
	}{
		APIKey: Secret(secretValue),
		Name:   name,
	}

	rawConfig, err := yaml.Marshal(config)
	require.NoError(t, err)

	rawConfigStr := string(rawConfig)

	assert.NotContains(t, rawConfigStr, secretValue)
	assert.Contains(t, rawConfigStr, maskedSecret)
	assert.Contains(t, rawConfigStr, name)
}
