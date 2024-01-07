package strategy

import "glide/pkg/providers"

// RoutingStrategy defines supported routing strategies for language routers
type RoutingStrategy string

type Strategy interface {
	Next() (*providers.LanguageModel, error)
}
