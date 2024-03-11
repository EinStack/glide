package schemas

// ChatStreamRequest defines a message that requests a new streaming chat
type ChatStreamRequest struct {
	// TODO: implement
}

// ChatStreamChunk defines a message for a chunk of streaming chat response
type ChatStreamChunk struct {
	// TODO: modify according to the streaming chat needs
	ID            string        `json:"id,omitempty"`
	Created       int           `json:"created,omitempty"`
	Provider      string        `json:"provider,omitempty"`
	RouterID      string        `json:"router,omitempty"`
	ModelID       string        `json:"model_id,omitempty"`
	ModelName     string        `json:"model,omitempty"`
	Cached        bool          `json:"cached,omitempty"`
	ModelResponse ModelResponse `json:"modelResponse,omitempty"`
	// TODO: add chat request-specific context
}

type ChatStreamError struct {
	// TODO: add chat request-specific context
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ChatStreamResult struct {
	chunk *ChatStreamChunk
	err   *ChatStreamError
}

func (r *ChatStreamResult) Chunk() *ChatStreamChunk {
	return r.chunk
}

func (r *ChatStreamResult) Error() *ChatStreamError {
	return r.err
}

func NewChatStreamResult(chunk *ChatStreamChunk) *ChatStreamResult {
	return &ChatStreamResult{
		chunk: chunk,
		err:   nil,
	}
}

func NewChatStreamErrorResult(err *ChatStreamError) *ChatStreamResult {
	return &ChatStreamResult{
		chunk: nil,
		err:   err,
	}
}
