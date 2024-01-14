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

	tests := map[string]Config{
		"all healthy": {
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
			},
		},
	}

	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			routers, err := cfg.BuildLangRouters(telemetry.NewTelemetryMock())

			require.NoError(t, err)
			require.Len(t, routers, 1)
		})
	}
}
