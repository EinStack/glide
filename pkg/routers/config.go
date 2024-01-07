package routers

import (
	"glide/pkg/providers"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Config struct {
	LanguageRouters []LangRouterConfig `yaml:"language"`
}

func (c *Config) BuildLangRouters(tel *telemetry.Telemetry) ([]*LangRouter, error) {
	routers := make([]*LangRouter, 0, len(c.LanguageRouters))

	var errs error

	for idx, routerConfig := range c.LanguageRouters {
		if !routerConfig.Enabled {
			tel.Logger.Info("router is disabled, skipping", zap.String("routerID", routerConfig.ID))
			continue
		}

		tel.Logger.Debug("init router", zap.String("routerID", routerConfig.ID))

		router, err := NewLangRouter(&c.LanguageRouters[idx], tel)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		routers = append(routers, router)
	}

	if errs != nil {
		return nil, errs
	}

	return routers, nil
}

type LangRouterConfig struct {
	ID              string                      `yaml:"id" json:"routers" validate:"required"`
	Enabled         bool                        `yaml:"enabled" json:"enabled"`
	RoutingStrategy routing.Strategy            `yaml:"strategy" json:"strategy"`
	Models          []providers.LangModelConfig `yaml:"models" json:"models" validate:"required"`
}

// BuildModels creates LanguageModel slice out of the given config
func (c *LangRouterConfig) BuildModels(tel *telemetry.Telemetry) ([]providers.LanguageModel, error) {
	var errs error

	models := make([]providers.LanguageModel, 0, len(c.Models))

	for _, modelConfig := range c.Models {
		if !modelConfig.Enabled {
			tel.Logger.Info(
				"model is disabled, skipping",
				zap.String("router", c.ID),
				zap.String("model", modelConfig.ID),
			)

			continue
		}

		tel.Logger.Debug(
			"init lang model",
			zap.String("router", c.ID),
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

	return models, nil
}

func DefaultLangRouterConfig() LangRouterConfig {
	return LangRouterConfig{
		Enabled:         true,
		RoutingStrategy: routing.Priority,
	}
}

func (p *LangRouterConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = DefaultLangRouterConfig()

	type plain LangRouterConfig // to avoid recursion

	return unmarshal((*plain)(p))
}
