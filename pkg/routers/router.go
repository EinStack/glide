package routers

import (
	"context"
	"errors"
	"glide/pkg/providers"
	"glide/pkg/providers/factory"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"glide/pkg/api/schemas"
	"glide/pkg/telemetry"
)

var (
	ErrNoModels = errors.New("no models configured for router")
)

type LangRouter struct {
	config    *LangRouterConfig
	models    []providers.LanguageModel
	telemetry *telemetry.Telemetry
}

func NewLangRouter(cfg *LangRouterConfig, tel *telemetry.Telemetry) (*LangRouter, error) {
	router := &LangRouter{
		config:    cfg,
		telemetry: tel,
	}

	err := router.BuildModels(cfg.Models)

	return router, err
}

func (r *LangRouter) BuildModels(modelConfigs []providers.LangModelConfig) error {
	var errs error
	models := make([]providers.LanguageModel, 0, len(modelConfigs))

	for _, modelConfig := range modelConfigs {
		if !modelConfig.Enabled {
			r.telemetry.Logger.Info("model is disabled, skipping", zap.String("modelID", modelConfig.ID))
			continue
		}

		r.telemetry.Logger.Debug("init lang model", zap.String("modelID", modelConfig.ID))

		model, err := factory.NewModelFromConfig(modelConfig, r.telemetry)

		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		models = append(models, model)
	}

	if errs != nil {
		return errs
	}

	r.models = models

	return nil
}

func (r *LangRouter) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	if len(r.models) == 0 {
		return nil, ErrNoModels
	}

	// TODO: implement actual routing & fallback logic
	return r.models[0].Chat(ctx, request)
}
