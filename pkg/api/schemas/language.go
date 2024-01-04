package schemas

// UnifiedChatRequest defines Glide's Chat Request Schema unified across all language models
type UnifiedChatRequest struct {
	Message        ChatMessage   `json:"message"`
	MessageHistory []ChatMessage `json:"messageHistory"`
}

// UnifiedChatResponse defines Glide's Chat Response Schema unified across all language models
type UnifiedChatResponse struct {
	ID               string           `json:"id,omitempty"`
	Created          string          `json:"created,omitempty"`
	Provider         string           `json:"provider,omitempty"`
	Router           string           `json:"router,omitempty"`
	Model            string           `json:"model,omitempty"`
	Cached           bool             `json:"cached,omitempty"`
	ModelResponse ProviderResponse `json:"modelResponse,omitempty"`
}

// ProviderResponse is the unified response from the provider.

type ProviderResponse struct {
	ResponseID map[string]string `json:"responseID,omitempty"`
	Message    ChatMessage       `json:"message"`
	TokenCount TokenCount        `json:"tokenCount"`
}

type TokenCount struct {
    PromptTokens     float64 `json:"promptTokens"`
    ResponseTokens   float64 `json:"responseTokens"`
    TotalTokens      float64 `json:"totalTokens"`
}


// OpenAI Chat Response
type OpenAIChatCompletion struct {
    ID                string             `json:"id"`
    Object            string             `json:"object"`
    Created           string              `json:"created"`
    Model             string             `json:"model"`
    SystemFingerprint string             `json:"systemFingerprint"`
    Choices           []Choice           `json:"choices"`
    Usage             Usage              `json:"usage"`
}

type Choice struct {
    Index         int         `json:"index"`
    Message       ChatMessage     `json:"message"`
    Logprobs      interface{} `json:"logprobs"`
    FinishReason  string      `json:"finishReason"`
}

type Usage struct {
    PromptTokens     float64 `json:"promptTokens"`
    CompletionTokens float64 `json:"completionTokens"`
    TotalTokens      float64 `json:"totalTokens"`
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
