package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigProvider_ValidConfigLoaded(t *testing.T) {
	configProvider := NewProvider()
	configProvider, err := configProvider.Load("./testdata/provider.fullconfig.yaml")
	require.NoError(t, err)

	cfg := configProvider.Get()

	langRouters := cfg.Routers.LanguageRouters

	require.Len(t, langRouters, 1)
	require.Equal(t, true, langRouters[0].Enabled)

	models := langRouters[0].Models
	require.Len(t, models, 1)
}

func TestConfigProvider_NoProvider(t *testing.T) {
	configProvider := NewProvider()
	configProvider, err := configProvider.Load("./testdata/provider.nomodelprovider.yaml")
	require.Error(t, err)
	require.ErrorContains(t, err, "none is configured")
}
