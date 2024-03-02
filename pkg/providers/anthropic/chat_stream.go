package anthropic

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
)

func (c *Client) SupportChatStream() bool {
	return false
}

func (c *Client) ChatStream(_ context.Context, _ *schemas.ChatRequest, _ chan<- schemas.ChatResponse) error {
	return clients.ErrChatStreamNotImplemented
}
