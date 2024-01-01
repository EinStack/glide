package providers

import (
	"context"
	"errors"
	"glide/pkg/api/schemas"
)

var ErrProviderUnavailable = errors.New("provider is not available")

// ModelProvider defines an interface all model providers should support
type ModelProvider interface {
	Provider() string
}

// LanguageModel defines the interface a provider should fulfill to be able to serve language chat requests
type LanguageModel interface {
	Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error)
}
