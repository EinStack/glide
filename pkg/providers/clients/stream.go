package clients

import (
	"github.com/EinStack/glide/pkg/api/schemas"
)

type ChatStream interface {
	Open() error
	Recv() (*schemas.ChatStreamChunk, error)
	Close() error
}

type ChatStreamResult struct {
	chunk *schemas.ChatStreamChunk
	err   error
}

func (r *ChatStreamResult) Chunk() *schemas.ChatStreamChunk {
	return r.chunk
}

func (r *ChatStreamResult) Error() error {
	return r.err
}

func NewChatStreamResult(chunk *schemas.ChatStreamChunk, err error) *ChatStreamResult {
	return &ChatStreamResult{
		chunk: chunk,
		err:   err,
	}
}
