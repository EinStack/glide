package routers

import (
	"context"
	"errors"

	"glide/pkg/providers"

	"glide/pkg/api/schemas"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

var ErrNoModels = errors.New("no models configured for router")

type LangRouter struct {
	Config    *LangRouterConfig
	routing   routing.LangModelRouting
	models    []*providers.LangModel
	telemetry *telemetry.Telemetry
}

func NewLangRouter(cfg *LangRouterConfig, tel *telemetry.Telemetry) (*LangRouter, error) {
	models, err := cfg.BuildModels(tel)
	if err != nil {
		return nil, err
	}

	router := &LangRouter{
		Config:    cfg,
		models:    models,
		routing:   routing.NewPriorityRouting(models),
		telemetry: tel,
	}

	return router, err
}

func (r *LangRouter) ID() string {
	return r.Config.ID
}

func (r *LangRouter) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	if len(r.models) == 0 {
		return nil, ErrNoModels
	}

	// maxRetries := 3 // TODO: move to configs
	modelIterator := r.routing.Iterator()

	for {
		model, err := modelIterator.Next()

		if errors.Is(err, routing.ErrNoHealthyModels) {
			// no healthy model in the pool. Let's retry after some time
			// r.telemetry.Logger.Warn("")
			break
		}

		resp, err := model.Chat(ctx, request)
		// TODO:
		if err != nil {
		}

		return resp, nil
	}
	// TODO: wait and retry define number of times
}
