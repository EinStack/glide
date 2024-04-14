package providers

import (
	"glide/pkg/config/fields"
)

// ModelProvider exposes provider context
type ModelProvider interface {
	Provider() string
}

// Model represent a configured external modality-agnostic model with its routing properties and status
type Model interface {
	ID() string
	Healthy() bool
	LatencyUpdateInterval() *fields.Duration
	Weight() int
}
