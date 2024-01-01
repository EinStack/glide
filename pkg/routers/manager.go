package routers

import (
	"errors"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"glide/pkg/telemetry"
)

var ErrRouterNotFound = errors.New("no router found with given ID")

type RouterManager struct {
	config      *Config
	telemetry   *telemetry.Telemetry
	langRouters *map[string]*LangRouter
}

// NewManager creates a new instance of Router Manager that creates, holds and returns all routers
func NewManager(cfg *Config, tel *telemetry.Telemetry) (*RouterManager, error) {
	manager := RouterManager{
		config:    cfg,
		telemetry: tel,
	}

	err := manager.BuildRouters(cfg.LanguageRouters)

	return &manager, err
}

func (r *RouterManager) BuildRouters(routerConfigs []LangRouterConfig) error {
	routers := make(map[string]*LangRouter, len(routerConfigs))

	var errs error

	for idx, routerConfig := range routerConfigs {
		if !routerConfig.Enabled {
			r.telemetry.Logger.Info("router is disabled, skipping", zap.String("routerID", routerConfig.ID))
			continue
		}

		r.telemetry.Logger.Debug("init router", zap.String("routerID", routerConfig.ID))

		router, err := NewLangRouter(&routerConfigs[idx], r.telemetry)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		routers[routerConfig.ID] = router
	}

	if errs != nil {
		return errs
	}

	r.langRouters = &routers

	return nil
}

// GetLangRouter returns a router by type and ID
func (r *RouterManager) GetLangRouter(routerID string) (*LangRouter, error) {
	if router, found := (*r.langRouters)[routerID]; found {
		return router, nil
	}

	return nil, ErrRouterNotFound
}
