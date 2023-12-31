package providers

import (
	"context"

	"glide/pkg/api/schemas"
)

// ChatModel defines the interface a provider should fulfill to be able to serve language chat requests
type ChatModel interface {
	Chat(ctx *context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error)
}
