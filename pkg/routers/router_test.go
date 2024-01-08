package routers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"glide/pkg/api/schemas"
	"glide/pkg/providers"
	"glide/pkg/routers/health"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

func TestLangRouter_Priority_ChatRequest(t *testing.T) {
	budget := health.NewErrorBudget(3, health.SEC)
	models := []*providers.LangModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}),
			*budget,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}}),
			*budget,
		),
	}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(3, 2, 1*time.Second, nil),
		routing:   routing.NewPriorityRouting(models),
		models:    models,
		telemetry: telemetry.NewTelemetryMock(),
	}

	ctx := context.Background()
	req := schemas.NewChatFromStr("tell me a dad joke")

	for i := 0; i < 2; i++ {
		resp, err := router.Chat(ctx, req)

		require.Equal(t, "first", resp.Model)
		require.Equal(t, "test_router", resp.Router)
		require.NoError(t, err)
	}
}
