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
	}

	routers, err := cfg.BuildLangRouters(telemetry.NewTelemetryMock())

	require.NoError(t, err)
	require.Len(t, routers, 2)
	require.Len(t, routers[0].models, 1)
	require.IsType(t, routers[0].routing, &routing.PriorityRouting{})
	require.Len(t, routers[1].models, 1)
	require.IsType(t, routers[1].routing, &routing.LeastLatencyRouting{})
}
