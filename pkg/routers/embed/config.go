package embed

import (
	"github.com/EinStack/glide/pkg/providers"
	"github.com/EinStack/glide/pkg/resiliency/retry"
	"github.com/EinStack/glide/pkg/routers/routing"
)

type EmbeddingRouterConfig struct {
	ID              string                      `yaml:"id" json:"routers" validate:"required"`                                       // Unique router ID
	Enabled         bool                        `yaml:"enabled" json:"enabled" validate:"required"`                                  // Is router enabled?
	Retry           *retry.ExpRetryConfig       `yaml:"retry" json:"retry" validate:"required"`                                      // retry when no healthy model is available to router
	RoutingStrategy routing.Strategy            `yaml:"strategy" json:"strategy" swaggertype:"primitive,string" validate:"required"` // strategy on picking the next model to serve the request
	Models          []providers.LangModelConfig `yaml:"models" json:"models" validate:"required,min=1,dive"`                         // the list of models that could handle requests
}
