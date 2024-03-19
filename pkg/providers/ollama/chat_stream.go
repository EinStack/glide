package ollama

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
)

func (c *Client) SupportChatStream() bool {
	return false
}

func (c *Client) ChatStream(_ context.Context, _ *schemas.ChatStreamRequest) (clients.ChatStream, error) {
	return nil, clients.ErrChatStreamNotImplemented
}
