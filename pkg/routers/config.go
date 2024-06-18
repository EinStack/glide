package routers

import (
	"fmt"
	"github.com/EinStack/glide/pkg/routers/lang"
	"github.com/EinStack/glide/pkg/telemetry"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Config struct {
	LanguageRouters  []lang.LangRouterConfig `yaml:"language" validate:"required,dive"` // the list of language routers
	EmbeddingRouters []EmbeddingRouterConfig `yaml:"embedding" validate:"required,dive"`
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
			tel.L().Info(fmt.Sprintf("Router \"%v\" is disabled, skipping", routerConfig.ID))
			continue
		}

		tel.L().Debug("Init router", zap.String("routerID", routerConfig.ID))

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
