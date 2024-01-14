package routers

import (
	"context"
	"testing"
	"time"

	"glide/pkg/providers/clients"

	"github.com/stretchr/testify/require"
	"glide/pkg/api/schemas"
	"glide/pkg/providers"
	"glide/pkg/routers/health"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

func TestLangRouter_Priority_PickFistHealthy(t *testing.T) {
	budget := health.NewErrorBudget(3, health.SEC)
	langModels := []providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}}),
			*budget,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(3, 2, 1*time.Second, nil),
		routing:   routing.NewPriority(models),
		models:    langModels,
		telemetry: telemetry.NewTelemetryMock(),
	}

	ctx := context.Background()
	req := schemas.NewChatFromStr("tell me a dad joke")

	for i := 0; i < 2; i++ {
		resp, err := router.Chat(ctx, req)

		require.Equal(t, "first", resp.ModelID)
		require.Equal(t, "test_router", resp.RouterID)
		require.NoError(t, err)
	}
}

func TestLangRouter_Priority_PickThirdHealthy(t *testing.T) {
	budget := health.NewErrorBudget(1, health.SEC)
	langModels := []providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "3"}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "4"}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"third",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}),
			*budget,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	expectedModels := []string{"third", "third"}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(3, 2, 1*time.Second, nil),
		routing:   routing.NewPriority(models),
		models:    langModels,
		telemetry: telemetry.NewTelemetryMock(),
	}

	ctx := context.Background()
	req := schemas.NewChatFromStr("tell me a dad joke")

	for _, modelID := range expectedModels {
		resp, err := router.Chat(ctx, req)

		require.NoError(t, err)
		require.Equal(t, modelID, resp.ModelID)
		require.Equal(t, "test_router", resp.RouterID)
	}
}

func TestLangRouter_Priority_SuccessOnRetry(t *testing.T) {
	budget := health.NewErrorBudget(1, health.MILLI)
	langModels := []providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "2"}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "1"}}),
			*budget,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(3, 2, 1*time.Millisecond, nil),
		routing:   routing.NewPriority(models),
		models:    langModels,
		telemetry: telemetry.NewTelemetryMock(),
	}

	resp, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

	require.NoError(t, err)
	require.Equal(t, "first", resp.ModelID)
	require.Equal(t, "test_router", resp.RouterID)
}

func TestLangRouter_Priority_UnhealthyModelInThePool(t *testing.T) {
	budget := health.NewErrorBudget(1, health.MIN)
	langModels := []providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &clients.ErrProviderUnavailable}, {Msg: "3"}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}),
			*budget,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(3, 2, 1*time.Millisecond, nil),
		routing:   routing.NewPriority(models),
		models:    langModels,
		telemetry: telemetry.NewTelemetryMock(),
	}

	for i := 0; i < 2; i++ {
		resp, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

		require.NoError(t, err)
		require.Equal(t, "second", resp.ModelID)
		require.Equal(t, "test_router", resp.RouterID)
	}
}

func TestLangRouter_Priority_AllModelsUnavailable(t *testing.T) {
	budget := health.NewErrorBudget(1, health.SEC)
	langModels := []providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}),
			*budget,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}),
			*budget,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:  "test_router",
		Config:    &LangRouterConfig{},
		retry:     retry.NewExpRetry(1, 2, 1*time.Millisecond, nil),
		routing:   routing.NewPriority(models),
		models:    langModels,
		telemetry: telemetry.NewTelemetryMock(),
	}

	_, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

	require.Error(t, err)
}
