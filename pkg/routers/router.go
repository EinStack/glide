package routers

import (
	"context"
	"errors"
	"glide/pkg/providers"
	"glide/pkg/routers/routing"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"glide/pkg/api/schemas"
	"glide/pkg/telemetry"
)

var ErrNoModels = errors.New("no models configured for router")

type LangRouter struct {
	Config    *LangRouterConfig
	routing   routing.LangModelRouting
	models    *[]providers.LanguageModel
	telemetry *telemetry.Telemetry
}

// buildModels creates LanguageModel slice out of the given config
// TODO: consider moving it on the config struct level
func buildModels(modelConfigs []providers.LangModelConfig, routerID string, tel *telemetry.Telemetry) (*[]providers.LanguageModel, error) {
	var errs error

	if len(modelConfigs) == 0 {
		return nil, ErrNoModels
	}

	models := make([]providers.LanguageModel, 0, len(modelConfigs))

	for _, modelConfig := range modelConfigs {
		if !modelConfig.Enabled {
			tel.Logger.Info(
				"model is disabled, skipping",
				zap.String("router", routerID),
				zap.String("model", modelConfig.ID),
			)

			continue
		}

		tel.Logger.Debug(
			"init lang model",
			zap.String("router", routerID),
			zap.String("model", modelConfig.ID),
		)

		model, err := modelConfig.ToModel(tel)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		models = append(models, model)
	}

	if errs != nil {
		return nil, errs
	}

	return &models, nil
}

func NewLangRouter(cfg *LangRouterConfig, tel *telemetry.Telemetry) (*LangRouter, error) {
	models, err := buildModels(cfg.Models, cfg.ID, tel)

	router := &LangRouter{
		Config:    cfg,
		models:    models,
		routing:   routing.NewPriorityRouting(models),
		telemetry: tel,
	}

	return router, err
}

func (r *LangRouter) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	if len(*r.models) == 0 {
		return nil, ErrNoModels
	}

	//maxRetries := 3 // TODO: move to configs
	modelIterator := r.routing.Iterator()

	for {
		model, err := modelIterator.Next()

		if errors.Is(err, routing.ErrNoHealthyModels) {
			// no healthy model in the pool. Let's retry after some time
			//r.telemetry.Logger.Warn("")
			break
		}

		resp, err := model.Chat()

		// TODO:
		if err != nil {

		}

		return resp, nil
	}
	// TODO: wait and retry define number of times
}
