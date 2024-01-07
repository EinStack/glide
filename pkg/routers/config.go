package routers

import (
	"glide/pkg/providers"
	"glide/pkg/routers/routing"
)

type Config struct {
	LanguageRouters []LangRouterConfig `yaml:"language"`
}

type LangRouterConfig struct {
	ID              string                      `yaml:"id" json:"routers" validate:"required"`
	Enabled         bool                        `yaml:"enabled" json:"enabled"`
	RoutingStrategy routing.Strategy            `yaml:"strategy" json:"strategy"`
	Models          []providers.LangModelConfig `yaml:"models" json:"models" validate:"required"`
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
