package schemas

import "time"

type (
	Metadata     = map[string]any
	EventType    = string
	FinishReason = string
	ErrorCode    = string
)

var (
	Complete        FinishReason = "complete"
	MaxTokens       FinishReason = "max_tokens"
	ContentFiltered FinishReason = "content_filtered"
	ErrorReason     FinishReason = "error"
	OtherReason     FinishReason = "other"
)

var (
	NoModelConfigured    ErrorCode = "no_model_configured"
	ModelUnavailable     ErrorCode = "model_unavailable"
	AllModelsUnavailable ErrorCode = "all_models_unavailable"
	UnknownError         ErrorCode = "unknown_error"
)

type StreamRequestID = string

// ChatStreamRequest defines a message that requests a new streaming chat
type ChatStreamRequest struct {
	ID             StreamRequestID      `json:"id" validate:"required"`
	Message        ChatMessage          `json:"message" validate:"required"`
	MessageHistory []ChatMessage        `json:"messageHistory" validate:"required"`
	Override       *OverrideChatRequest `json:"overrideMessage,omitempty"`
	Metadata       *Metadata            `json:"metadata,omitempty"`
}

func NewChatStreamFromStr(message string) *ChatStreamRequest {
	return &ChatStreamRequest{
		Message: ChatMessage{
			"user",
			message,
			"glide",
		},
	}
}

type ModelChunkResponse struct {
	Metadata *Metadata   `json:"metadata,omitempty"`
	Message  ChatMessage `json:"message"`
}

type ChatStreamMessage struct {
	ID        StreamRequestID  `json:"id"`
	CreatedAt int              `json:"createdAt"`
	RouterID  string           `json:"routerId"`
	Metadata  *Metadata        `json:"metadata,omitempty"`
	Chunk     *ChatStreamChunk `json:"chunk,omitempty"`
	Error     *ChatStreamError `json:"error,omitempty"`
}

// ChatStreamChunk defines a message for a chunk of streaming chat response
type ChatStreamChunk struct {
	ModelID       string             `json:"modelId"`
	Provider      string             `json:"providerName"`
	ModelName     string             `json:"modelName"`
	Cached        bool               `json:"cached"`
	ModelResponse ModelChunkResponse `json:"modelResponse"`
	FinishReason  *FinishReason      `json:"finishReason,omitempty"`
}

type ChatStreamError struct {
	ErrCode ErrorCode `json:"errCode"`
	Message string    `json:"message"`
}

func NewChatStreamChunk(
	reqID StreamRequestID,
	routerID string,
	reqMetadata *Metadata,
	chunk *ChatStreamChunk,
) *ChatStreamMessage {
	return &ChatStreamMessage{
		ID:        reqID,
		RouterID:  routerID,
		CreatedAt: int(time.Now().UTC().Unix()),
		Metadata:  reqMetadata,
		Chunk:     chunk,
	}
}

func NewChatStreamError(
	reqID StreamRequestID,
	routerID string,
	errCode ErrorCode,
	errMsg string,
	reqMetadata *Metadata,
) *ChatStreamMessage {
	return &ChatStreamMessage{
		ID:        reqID,
		RouterID:  routerID,
		CreatedAt: int(time.Now().UTC().Unix()),
		Metadata:  reqMetadata,
		Error: &ChatStreamError{
			ErrCode: errCode,
			Message: errMsg,
		},
	}
}
