package routers

import (
	"fmt"

	"glide/pkg/providers"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Config struct {
	LanguageRouters []LangRouterConfig `yaml:"language" validate:"required,min=1"` // the list of language routers
}

func (c *Config) BuildLangRouters(tel *telemetry.Telemetry) ([]*LangRouter, error) {
	seenIDs := make(map[string]bool, len(c.LanguageRouters))
	routers := make([]*LangRouter, 0, len(c.LanguageRouters))

	var errs error

	for idx, routerConfig := range c.LanguageRouters {
		if _, ok := seenIDs[routerConfig.ID]; ok {
			return nil, fmt.Errorf("ID \"%v\" is specified for more than one router while each ID should be unique", routerConfig.ID)
		}

		seenIDs[routerConfig.ID] = true

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

// TODO: how to specify other backoff strategies?
// TODO: Had to keep RoutingStrategy because of https://github.com/swaggo/swag/issues/1738
// LangRouterConfig
type LangRouterConfig struct {
	ID              string                      `yaml:"id" json:"routers" validate:"required"`                                       // Unique router ID
	Enabled         bool                        `yaml:"enabled" json:"enabled" validate:"required"`                                  // Is router enabled?
	Retry           *retry.ExpRetryConfig       `yaml:"retry" json:"retry" validate:"required"`                                      // retry when no healthy model is available to router
	RoutingStrategy routing.Strategy            `yaml:"strategy" json:"strategy" swaggertype:"primitive,string" validate:"required"` // strategy on picking the next model to serve the request
	Models          []providers.LangModelConfig `yaml:"models" json:"models" validate:"required,min=1"`                              // the list of models that could handle requests
}

// BuildModels creates LanguageModel slice out of the given config
func (c *LangRouterConfig) BuildModels(tel *telemetry.Telemetry) ([]providers.LanguageModel, error) {
	var errs error

	seenIDs := make(map[string]bool, len(c.Models))
	models := make([]providers.LanguageModel, 0, len(c.Models))

	for _, modelConfig := range c.Models {
		if _, ok := seenIDs[modelConfig.ID]; ok {
			return nil, fmt.Errorf(
				"ID \"%v\" is specified for more than one model in router \"%v\", while it should be unique in scope of that pool",
				modelConfig.ID,
				c.ID,
			)
		}

		seenIDs[modelConfig.ID] = true

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

	if len(models) == 0 {
		return nil, fmt.Errorf("router \"%v\" must have at least one active model, zero defined", c.ID)
	}

	if len(models) == 1 {
		tel.Logger.Warn(
			"router has only one active model defined. "+
				"This is not recommended for production setups. "+
				"Define at least a few models to leverage resiliency logic Glide provides",
			zap.String("router", c.ID),
		)
	}

	return models, nil
}

func (c *LangRouterConfig) BuildRetry() *retry.ExpRetry {
	retryConfig := c.Retry

	return retry.NewExpRetry(
		retryConfig.MaxRetries,
		retryConfig.BaseMultiplier,
		retryConfig.MinDelay,
		retryConfig.MaxDelay,
	)
}

func (c *LangRouterConfig) BuildRouting(models []providers.LanguageModel) (routing.LangModelRouting, error) {
	m := make([]providers.Model, 0, len(models))
	for _, model := range models {
		m = append(m, model)
	}

	switch c.RoutingStrategy {
	case routing.Priority:
		return routing.NewPriority(m), nil
	case routing.RoundRobin:
		return routing.NewRoundRobinRouting(m), nil
	case routing.WeightedRoundRobin:
		return routing.NewWeightedRoundRobin(m), nil
	case routing.LeastLatency:
		return routing.NewLeastLatencyRouting(m), nil
	}

	return nil, fmt.Errorf("routing strategy \"%v\" is not supported, please make sure there is no typo", c.RoutingStrategy)
}

func DefaultLangRouterConfig() LangRouterConfig {
	return LangRouterConfig{
		Enabled:         true,
		RoutingStrategy: routing.Priority,
		Retry:           retry.DefaultExpRetryConfig(),
	}
}

func (c *LangRouterConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultLangRouterConfig()

	type plain LangRouterConfig // to avoid recursion

	return unmarshal((*plain)(c))
}
