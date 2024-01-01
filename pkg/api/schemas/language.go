package schemas

// UnifiedChatRequest defines Glide's Chat Request Schema unified across all language models
type UnifiedChatRequest struct {
	Message        ChatMessage   `json:"message"`
	MessageHistory []ChatMessage `json:"messageHistory"`
}

// UnifiedChatResponse defines Glide's Chat Response Schema unified across all language models
type UnifiedChatResponse struct {
	ID               string           `json:"id,omitempty"`
	Created          float64          `json:"created,omitempty"`
	Provider         string           `json:"provider,omitempty"`
	Router           string           `json:"router,omitempty"`
	Model            string           `json:"model,omitempty"`
	Cached           bool             `json:"cached,omitempty"`
	ProviderResponse ProviderResponse `json:"provider_response,omitempty"`
}

// ProviderResponse contains data from the chosen provider
type ProviderResponse struct {
	ResponseId map[string]string `json:"response_id,omitempty"`
	Message    ChatMessage       `json:"message"`
	TokenCount TokenCount        `json:"token_count"`
}

type TokenCount struct {
	PromptTokens   float64`json:"prompt_tokens"`
	ResponseTokens float64 `json:"response_tokens"`
	TotalTokens    float64 `json:"total_tokens"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct {
	// The role of the author of this message. One of system, user, or assistant.
	Role string `json:"role"`
	// The content of the message.
	Content string `json:"content"`
	// The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,
	// with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`
}

// ChatChoice is a choice in a chat response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	CompletionTokens float64 `json:"completion_tokens,omitempty"`
	PromptTokens     float64 `json:"prompt_tokens,omitempty"`
	TotalTokens      float64 `json:"total_tokens,omitempty"`
}
