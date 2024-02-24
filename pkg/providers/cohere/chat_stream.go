package cohere

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
)

func (c *Client) SupportChatStream() bool {
	return false
}

func (c *Client) ChatStream(ctx context.Context, request *schemas.ChatRequest, responseC chan<- schemas.ChatResponse) error {
	return clients.ErrChatStreamNotImplemented
}
