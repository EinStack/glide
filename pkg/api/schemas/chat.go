package schemas

// ChatRequest defines Glide's Chat Request Schema unified across all language models
type ChatRequest struct {
	Message        ChatMessage          `json:"message" validate:"required"`
	MessageHistory []ChatMessage        `json:"message_history"`
	OverrideParams *OverrideChatRequest `json:"override_params,omitempty"`
}

type OverrideChatRequest struct {
	ModelID string      `json:"model_id" validate:"required"`
	Message ChatMessage `json:"message" validate:"required"`
}

func NewChatFromStr(message string) *ChatRequest {
	return &ChatRequest{
		Message: ChatMessage{
			"user",
			message,
			"glide",
		},
	}
}

// ChatResponse defines Glide's Chat Response Schema unified across all language models
type ChatResponse struct {
	ID            string        `json:"id,omitempty"`
	Created       int           `json:"created_at,omitempty"`
	Provider      string        `json:"provider_id,omitempty"`
	RouterID      string        `json:"router_id,omitempty"`
	ModelID       string        `json:"model_id,omitempty"`
	ModelName     string        `json:"model_name,omitempty"`
	Cached        bool          `json:"cached,omitempty"`
	ModelResponse ModelResponse `json:"model_response,omitempty"`
}

// ModelResponse is the unified response from the provider.

type ModelResponse struct {
	Metadata   map[string]string `json:"metadata,omitempty"`
	Message    ChatMessage       `json:"message"`
	TokenUsage TokenUsage        `json:"token_usage"`
}

type TokenUsage struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct {
	// The role of the author of this message. One of system, user, or assistant.
	Role string `json:"role" validate:"required"`
	// The content of the message.
	Content string `json:"content" validate:"required"`
	// The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,
	// with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`
}
