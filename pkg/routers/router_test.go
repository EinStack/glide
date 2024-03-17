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
	ptesting "glide/pkg/providers/testing"
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
			ptesting.NewProviderMock([]ptesting.RespMock{{Msg: "1"}, {Msg: "2"}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewProviderMock([]ptesting.RespMock{{Msg: "1"}}),
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
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Msg: "3"}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Msg: "4"}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"third",
			ptesting.NewProviderMock([]ptesting.RespMock{{Msg: "1"}, {Msg: "2"}}),
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
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Msg: "2"}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Msg: "1"}}),
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
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &clients.ErrProviderUnavailable}, {Msg: "3"}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewProviderMock([]ptesting.RespMock{{Msg: "1"}, {Msg: "2"}}),
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
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewProviderMock([]ptesting.RespMock{{Err: &ErrNoModelAvailable}, {Err: &ErrNoModelAvailable}}),
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
			ptesting.NewStreamProviderMock([]ptesting.RespStreamMock{
				ptesting.NewRespStreamMock([]ptesting.RespMock{
					{Msg: "Bill"},
					{Msg: "Gates"},
					{Msg: "entered"},
					{Msg: "the"},
					{Msg: "bar"},
				}),
			}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewStreamProviderMock([]ptesting.RespStreamMock{
				ptesting.NewRespStreamMock([]ptesting.RespMock{
					{Msg: "Knock"},
					{Msg: "Knock"},
					{Msg: "joke"},
				}),
			}),
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

	chunks := make([]string, 0, 5)

	for range 5 {
		select { //nolint:gosimple
		case chunk := <-respC:
			require.Nil(t, chunk.Error())
			require.NotNil(t, chunk.Chunk().ModelResponse.Message.Content)

			chunks = append(chunks, chunk.Chunk().ModelResponse.Message.Content)
		}
	}

	require.Equal(t, []string{"Bill", "Gates", "entered", "the", "bar"}, chunks)
}

func TestLangRouter_ChatStream_AllModelsUnavailable(t *testing.T) {
	budget := health.NewErrorBudget(1, health.SEC)
	latConfig := latency.DefaultConfig()

	langModels := []*providers.LanguageModel{
		providers.NewLangModel(
			"first",
			ptesting.NewStreamProviderMock([]ptesting.RespStreamMock{
				ptesting.NewRespStreamMock([]ptesting.RespMock{
					{Err: &clients.ErrProviderUnavailable},
				}),
			}),
			budget,
			*latConfig,
			1,
		),
		providers.NewLangModel(
			"second",
			ptesting.NewStreamProviderMock([]ptesting.RespStreamMock{
				ptesting.NewRespStreamMock([]ptesting.RespMock{
					{Err: &clients.ErrProviderUnavailable},
				}),
			}),
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

	respC := make(chan *schemas.ChatStreamResult)
	defer close(respC)

	go router.ChatStream(context.Background(), schemas.NewChatFromStr("tell me a dad joke"), respC)

	errs := make([]string, 0, 3)

	for range 3 {
		result := <-respC
		require.Nil(t, result.Chunk())

		errs = append(errs, result.Error().Reason)
	}

	require.Equal(t, []string{"modelUnavailable", "modelUnavailable", "allModelsUnavailable"}, errs)
}
