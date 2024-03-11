package azureopenai

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
)

func (c *Client) SupportChatStream() bool {
	return false
}

func (c *Client) ChatStream(_ context.Context, _ *schemas.ChatRequest) <-chan *clients.ChatStreamResult {
	streamResultC := make(chan *clients.ChatStreamResult)

	streamResultC <- clients.NewChatStreamResult(nil, clients.ErrChatStreamNotImplemented)
	close(streamResultC)

	return streamResultC
}
