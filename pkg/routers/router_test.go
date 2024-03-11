package routers

import (
	"context"
	"testing"
	"time"

	"glide/pkg/routers/latency"

	"glide/pkg/providers/clients"

	"github.com/stretchr/testify/require"
	"glide/pkg/api/schemas"
	"glide/pkg/providers"
	"glide/pkg/routers/health"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

func TestLangRouter_Chat_PickFistHealthy(t *testing.T) {
	budget := health.NewErrorBudget(3, health.SEC)
	latConfig := latency.DefaultConfig()

	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}}, false),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:         "test_router",
		Config:           &LangRouterConfig{},
		retry:            retry.NewExpRetry(3, 2, 1*time.Second, nil),
		chatRouting:      routing.NewPriority(models),
		chatModels:       langModels,
		chatStreamModels: langModels,
		tel:              telemetry.NewTelemetryMock(),
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

func TestLangRouter_Chat_PickThirdHealthy(t *testing.T) {
	budget := health.NewErrorBudget(1, health.SEC)
	latConfig := latency.DefaultConfig()
	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "3"}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "4"}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"third",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}, false),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	expectedModels := []string{"third", "third"}

	router := LangRouter{
		routerID:          "test_router",
		Config:            &LangRouterConfig{},
		retry:             retry.NewExpRetry(3, 2, 1*time.Second, nil),
		chatRouting:       routing.NewPriority(models),
		chatStreamRouting: routing.NewPriority(models),
		chatModels:        langModels,
		chatStreamModels:  langModels,
		tel:               telemetry.NewTelemetryMock(),
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

func TestLangRouter_Chat_SuccessOnRetry(t *testing.T) {
	budget := health.NewErrorBudget(1, health.MILLI)
	latConfig := latency.DefaultConfig()
	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "2"}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Msg: "1"}}, false),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:          "test_router",
		Config:            &LangRouterConfig{},
		retry:             retry.NewExpRetry(3, 2, 1*time.Millisecond, nil),
		chatRouting:       routing.NewPriority(models),
		chatStreamRouting: routing.NewPriority(models),
		chatModels:        langModels,
		chatStreamModels:  langModels,
		tel:               telemetry.NewTelemetryMock(),
	}

	resp, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

	require.NoError(t, err)
	require.Equal(t, "first", resp.ModelID)
	require.Equal(t, "test_router", resp.RouterID)
}

func TestLangRouter_Chat_UnhealthyModelInThePool(t *testing.T) {
	budget := health.NewErrorBudget(1, health.MIN)
	latConfig := latency.DefaultConfig()
	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &clients.ErrProviderUnavailable}, {Msg: "3"}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}, false),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:          "test_router",
		Config:            &LangRouterConfig{},
		retry:             retry.NewExpRetry(3, 2, 1*time.Millisecond, nil),
		chatRouting:       routing.NewPriority(models),
		chatModels:        langModels,
		chatStreamModels:  langModels,
		chatStreamRouting: routing.NewPriority(models),
		tel:               telemetry.NewTelemetryMock(),
	}

	for i := 0; i < 2; i++ {
		resp, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

		require.NoError(t, err)
		require.Equal(t, "second", resp.ModelID)
		require.Equal(t, "test_router", resp.RouterID)
	}
}

func TestLangRouter_Chat_AllModelsUnavailable(t *testing.T) {
	budget := health.NewErrorBudget(1, health.SEC)
	latConfig := latency.DefaultConfig()
	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}, false),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}, false),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:          "test_router",
		Config:            &LangRouterConfig{},
		retry:             retry.NewExpRetry(1, 2, 1*time.Millisecond, nil),
		chatRouting:       routing.NewPriority(models),
		chatModels:        langModels,
		chatStreamModels:  langModels,
		chatStreamRouting: routing.NewPriority(models),
		tel:               telemetry.NewTelemetryMock(),
	}

	_, err := router.Chat(context.Background(), schemas.NewChatFromStr("tell me a dad joke"))

	require.Error(t, err)
}

func TestLangRouter_ChatStream(t *testing.T) {
	budget := health.NewErrorBudget(3, health.SEC)
	latConfig := latency.DefaultConfig()

	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}, {Msg: "2"}}, true),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			providers.NewProviderMock([]providers.ResponseMock{{Msg: "1"}}, true),
			budget,
			*latConfig,
			1,
		),
	}

	models := make([]providers.Model, 0, len(langModels))
	for _, model := range langModels {
		models = append(models, model)
	}

	router := LangRouter{
		routerID:          "test_stream_router",
		Config:            &LangRouterConfig{},
		retry:             retry.NewExpRetry(3, 2, 1*time.Second, nil),
		chatRouting:       routing.NewPriority(models),
		chatModels:        langModels,
		chatStreamRouting: routing.NewPriority(models),
		chatStreamModels:  langModels,
		tel:               telemetry.NewTelemetryMock(),
	}

	ctx := context.Background()
	req := schemas.NewChatFromStr("tell me a dad joke")
	respC := make(chan *schemas.ChatStreamResult)

	defer close(respC)

	go router.ChatStream(ctx, req, respC)

	select {
	case chunkResult := <-respC:
		require.Nil(t, chunkResult.Error())
		require.NotNil(t, chunkResult.Chunk().ModelResponse.Message.Content)
	case <-time.Tick(5 * time.Second):
		t.Error("Timeout while waiting for stream chat chunk")
	}
}
