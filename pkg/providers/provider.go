package providers

import (
	"context"

	"glide/pkg/api/schemas"
)

// ModelProvider defines an interface all model providers should support
type ModelProvider interface {
	Provider() string
}

// LanguageModel defines the interface a provider should fulfill to be able to serve language chat requests
type LanguageModel interface {
	ID() string
	Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error)
}
