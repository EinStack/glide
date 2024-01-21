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

func TestConfigProvider_ValidConfigLoaded(t *testing.T) {
	configProvider := NewProvider()
	configProvider, err := configProvider.Load("./testdata/provider.fullconfig.yaml")
	require.NoError(t, err)

	cfg := configProvider.Get()

	langRouters := cfg.Routers.LanguageRouters

	require.Len(t, langRouters, 1)
	require.True(t, langRouters[0].Enabled)

	models := langRouters[0].Models
	require.Len(t, models, 1)
}

func TestConfigProvider_InvalidConfigLoaded(t *testing.T) {
	var tests = []struct {
		name       string
		configFile string
	}{
		{"empty telemetry", "./testdata/provider.telnil.yaml"},
		{"empty logging", "./testdata/provider.loggingnil.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configProvider := NewProvider()
			configProvider, err := configProvider.Load(tt.configFile)

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to validate config file")
		})
	}
}

func TestConfigProvider_NoProvider(t *testing.T) {
	configProvider := NewProvider()
	_, err := configProvider.Load("./testdata/provider.nomodelprovider.yaml")

	require.Error(t, err)
	require.ErrorContains(t, err, "none is configured")
}
