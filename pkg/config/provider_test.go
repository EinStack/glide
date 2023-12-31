package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigProvider_NonExistingConfigFile(t *testing.T) {
	_, err := NewProvider().Load("./testdata/doesntexist.yaml")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestConfigProvider_NonYAMLConfigFile(t *testing.T) {
	_, err := NewProvider().Load("./testdata/provider.broken.yaml")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "unable to parse config file")
}
