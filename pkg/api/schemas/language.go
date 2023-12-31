package schemas

// ChatRequest defines Glide's Chat Request Schema unified across all language models
type ChatRequest struct {
	Message []struct { // TODO: could we reuse ChatMessage?
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	MessageHistory []string `json:"messageHistory"`
}

// ChatResponse defines Glide's Chat Response Schema unified across all language models
type ChatResponse struct {
	ID      string        `json:"id,omitempty"`
	Created float64       `json:"created,omitempty"`
	Choices []*ChatChoice `json:"choices,omitempty"`
	Model   string        `json:"model,omitempty"`
	Object  string        `json:"object,omitempty"` // TODO: what does this mean "Object"?
	Usage   Usage         `json:"usage,omitempty"`
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
