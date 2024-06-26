package routing

import (
	"errors"

	"github.com/EinStack/glide/pkg/providers"
)

var ErrNoHealthyModels = errors.New("no healthy models found")

// Strategy defines supported routing strategies for language routers
type Strategy string

type LangModelRouting interface {
	Iterator() LangModelIterator
}

type LangModelIterator interface {
	Next() (providers.Model, error)
}
