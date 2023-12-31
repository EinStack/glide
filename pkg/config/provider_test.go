package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigProvider_NonExistingConfigFile(t *testing.T) {
	_, err := NewProvider().Load("./testdata/doesntexist.yaml")

	require.Error(t, err)
	require.ErrorContains(t, err, "no such file or directory")
}

func TestConfigProvider_NonYAMLConfigFile(t *testing.T) {
	_, err := NewProvider().Load("./testdata/provider.broken.yaml")

	require.Error(t, err)
	require.ErrorContains(t, err, "unable to parse config file")
}
