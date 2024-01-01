package routers

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/openai"
	"glide/pkg/telemetry"
)

type LangRouter struct {
	openAIClient *openai.Client // TODO: replace by actual model list
	telemetry    *telemetry.Telemetry
}

func NewLangRouter(tel *telemetry.Telemetry) (*LangRouter, error) {
	openAIClient, err := openai.NewClient(openai.DefaultConfig(), tel)
	if err != nil {
		return nil, err
	}

	return &LangRouter{
		openAIClient: openAIClient,
		telemetry:    tel,
	}, nil
}

func (r *LangRouter) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	// TODO: implement actual routing & fallback logic
	return r.openAIClient.Chat(ctx, request)
}
