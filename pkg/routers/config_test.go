package routers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"glide/pkg/providers"
	"glide/pkg/providers/clients"
	"glide/pkg/providers/openai"
	"glide/pkg/routers/health"
	"glide/pkg/routers/latency"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

func TestRouterConfig_BuildModels(t *testing.T) {
	defaultParams := openai.DefaultParams()

	cfg := Config{
		LanguageRouters: []LangRouterConfig{
			{
				ID:              "first_router",
				Enabled:         true,
				RoutingStrategy: routing.Priority,
				Retry:           retry.DefaultExpRetryConfig(),
				Models: []providers.LangModelConfig{
					{
						ID:          "first_model",
						Enabled:     true,
						Client:      clients.DefaultClientConfig(),
						ErrorBudget: health.DefaultErrorBudget(),
						Latency:     latency.DefaultConfig(),
						OpenAI: &openai.Config{
							APIKey:        "ABC",
							DefaultParams: &defaultParams,
						},
					},
				},
			},
			{
				ID:              "second_router",
				Enabled:         true,
				RoutingStrategy: routing.LeastLatency,
				Retry:           retry.DefaultExpRetryConfig(),
				Models: []providers.LangModelConfig{
					{
						ID:          "first_model",
						Enabled:     true,
						Client:      clients.DefaultClientConfig(),
						ErrorBudget: health.DefaultErrorBudget(),
						Latency:     latency.DefaultConfig(),
						OpenAI: &openai.Config{
							APIKey:        "ABC",
							DefaultParams: &defaultParams,
						},
					},
				},
			},
		},
	}

	routers, err := cfg.BuildLangRouters(telemetry.NewTelemetryMock())

	require.NoError(t, err)
	require.Len(t, routers, 2)
	require.Len(t, routers[0].models, 1)
	require.IsType(t, &routing.PriorityRouting{}, routers[0].routing)
	require.Len(t, routers[1].models, 1)
	require.IsType(t, &routing.LeastLatencyRouting{}, routers[1].routing)
}

func TestRouterConfig_InvalidSetups(t *testing.T) {
	defaultParams := openai.DefaultParams()

	tests := []struct {
		name   string
		config Config
	}{
		{
			"duplicated router IDs",
			Config{
				LanguageRouters: []LangRouterConfig{
					{
						ID:              "first_router",
						Enabled:         true,
						RoutingStrategy: routing.Priority,
						Retry:           retry.DefaultExpRetryConfig(),
						Models: []providers.LangModelConfig{
							{
								ID:          "first_model",
								Enabled:     true,
								Client:      clients.DefaultClientConfig(),
								ErrorBudget: health.DefaultErrorBudget(),
								Latency:     latency.DefaultConfig(),
								OpenAI: &openai.Config{
									APIKey:        "ABC",
									DefaultParams: &defaultParams,
								},
							},
						},
					},
					{
						ID:              "first_router",
						Enabled:         true,
						RoutingStrategy: routing.LeastLatency,
						Retry:           retry.DefaultExpRetryConfig(),
						Models: []providers.LangModelConfig{
							{
								ID:          "first_model",
								Enabled:     true,
								Client:      clients.DefaultClientConfig(),
								ErrorBudget: health.DefaultErrorBudget(),
								Latency:     latency.DefaultConfig(),
								OpenAI: &openai.Config{
									APIKey:        "ABC",
									DefaultParams: &defaultParams,
								},
							},
						},
					},
				},
			},
		},
		{
			"duplicated model IDs",
			Config{
				LanguageRouters: []LangRouterConfig{
					{
						ID:              "first_router",
						Enabled:         true,
						RoutingStrategy: routing.Priority,
						Retry:           retry.DefaultExpRetryConfig(),
						Models: []providers.LangModelConfig{
							{
								ID:          "first_model",
								Enabled:     true,
								Client:      clients.DefaultClientConfig(),
								ErrorBudget: health.DefaultErrorBudget(),
								Latency:     latency.DefaultConfig(),
								OpenAI: &openai.Config{
									APIKey:        "ABC",
									DefaultParams: &defaultParams,
								},
							},
							{
								ID:          "first_model",
								Enabled:     true,
								Client:      clients.DefaultClientConfig(),
								ErrorBudget: health.DefaultErrorBudget(),
								Latency:     latency.DefaultConfig(),
								OpenAI: &openai.Config{
									APIKey:        "ABC",
									DefaultParams: &defaultParams,
								},
							},
						},
					},
				},
			},
		},
		{
			"no models",
			Config{
				LanguageRouters: []LangRouterConfig{
					{
						ID:              "first_router",
						Enabled:         true,
						RoutingStrategy: routing.Priority,
						Retry:           retry.DefaultExpRetryConfig(),
						Models:          []providers.LangModelConfig{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		_, err := test.config.BuildLangRouters(telemetry.NewTelemetryMock())

		require.Error(t, err)
	}
}
