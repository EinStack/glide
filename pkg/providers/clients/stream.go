package clients

import "glide/pkg/api/schemas"

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
