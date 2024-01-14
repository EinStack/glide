package routers

import (
	"context"
	"errors"

	"glide/pkg/routers/retry"
	"go.uber.org/zap"

	"glide/pkg/providers"

	"glide/pkg/api/schemas"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

var (
	ErrNoModels         = errors.New("no models configured for router")
	ErrNoModelAvailable = errors.New("could not handle request because all providers are not available")
)

type LangRouter struct {
	routerID  string
	Config    *LangRouterConfig
	routing   routing.LangModelRouting
	retry     *retry.ExpRetry
	models    []providers.LanguageModel
	telemetry *telemetry.Telemetry
}

func NewLangRouter(cfg *LangRouterConfig, tel *telemetry.Telemetry) (*LangRouter, error) {
	models, err := cfg.BuildModels(tel)
	if err != nil {
		return nil, err
	}

	routing, err := cfg.BuildRouting(models)
	if err != nil {
		return nil, err
	}

	router := &LangRouter{
		routerID:  cfg.ID,
		Config:    cfg,
		models:    models,
		retry:     cfg.BuildRetry(),
		routing:   routing,
		telemetry: tel,
	}

	return router, err
}

func (r *LangRouter) ID() string {
	return r.routerID
}

func (r *LangRouter) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	if len(r.models) == 0 {
		return nil, ErrNoModels
	}

	retryIterator := r.retry.Iterator()

	for retryIterator.HasNext() {
		modelIterator := r.routing.Iterator()

		for {
			model, err := modelIterator.Next()

			if errors.Is(err, routing.ErrNoHealthyModels) {
				// no healthy model in the pool. Let's retry after some time
				break
			}

			langModel := model.(providers.LanguageModel)

			resp, err := langModel.Chat(ctx, request)
			if err != nil {
				r.telemetry.Logger.Warn(
					"lang model failed processing chat request",
					zap.String("routerID", r.ID()),
					zap.String("modelID", langModel.ID()),
					zap.String("provider", langModel.Provider()),
					zap.Error(err),
				)

				continue
			}

			resp.RouterID = r.routerID

			return resp, nil
		}

		// no providers were available to handle the request,
		//  so we have to wait a bit with a hope there is some available next time
		r.telemetry.Logger.Warn("no healthy model found, wait and retry", zap.String("routerID", r.ID()))

		err := retryIterator.WaitNext(ctx)
		if err != nil {
			// something has cancelled the context
			return nil, err
		}
	}

	// if we reach this part, then we are in trouble
	r.telemetry.Logger.Error("no model was available to handle request", zap.String("routerID", r.ID()))

	return nil, ErrNoModelAvailable
}
