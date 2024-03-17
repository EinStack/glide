package schemas

type (
	Metadata     = map[string]any
	FinishReason = string
)

var Complete FinishReason = "complete"

// ChatStreamRequest defines a message that requests a new streaming chat
type ChatStreamRequest struct {
	ID             string               `json:"id" validate:"required"`
	Message        ChatMessage          `json:"message" validate:"required"`
	MessageHistory []ChatMessage        `json:"messageHistory" validate:"required"`
	Override       *OverrideChatRequest `json:"overrideMessage,omitempty"`
	Metadata       *Metadata            `json:"metadata,omitempty"`
}

func NewChatStreamFromStr(message string) *ChatStreamRequest {
	return &ChatStreamRequest{
		Message: ChatMessage{
			"human",
			message,
			"glide",
		},
	}
}

type ModelChunkResponse struct {
	Metadata     *Metadata     `json:"metadata,omitempty"`
	Message      ChatMessage   `json:"message"`
	FinishReason *FinishReason `json:"finishReason,omitempty"`
}

// ChatStreamChunk defines a message for a chunk of streaming chat response
type ChatStreamChunk struct {
	ID            string             `json:"id"`
	CreatedAt     int                `json:"createdAt"`
	Provider      string             `json:"providerId"`
	RouterID      string             `json:"routerId"`
	ModelID       string             `json:"modelId"`
	Cached        bool               `json:"cached"`
	ModelName     string             `json:"modelName"`
	Metadata      *Metadata          `json:"metadata,omitempty"`
	ModelResponse ModelChunkResponse `json:"modelResponse"`
}

type ChatStreamError struct {
	ID       string    `json:"id"`
	ErrCode  string    `json:"errCode"`
	Message  string    `json:"message"`
	Metadata *Metadata `json:"metadata,omitempty"`
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
