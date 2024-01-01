package routers

import (
	"errors"

	"glide/pkg/telemetry"
)

var ErrRouterNotFound = errors.New("no router found with given ID")

type RouterManager struct {
	config    *Config
	telemetry *telemetry.Telemetry
	myRouter  *LangRouter // TODO: replace by list of routers
}

// NewManager creates a new instance of Router Manager that creates, holds and returns all routers
func NewManager(cfg *Config, tel *telemetry.Telemetry) (*RouterManager, error) {
	// TODO: init routers by config
	router, err := NewLangRouter(tel)
	if err != nil {
		return nil, err
	}

	return &RouterManager{
		config:    cfg,
		telemetry: tel,
		myRouter:  router,
	}, nil
}

// Get returns a router by type and ID
func (r *RouterManager) GetLangRouter(routerID string) (*LangRouter, error) {
	// TODO: implement actual logic
	if routerID != "myrouter" {
		return nil, ErrRouterNotFound
	}

	return r.myRouter, nil
}
