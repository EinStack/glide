package cohere

// Cohere Chat Response
type ChatCompletion struct {
	Text          string                 `json:"text"`
	GenerationID  string                 `json:"generation_id"`
	ResponseID    string                 `json:"response_id"`
	TokenCount    TokenCount             `json:"token_count"`
	Citations     []Citation             `json:"citations"`
	Documents     []Documents            `json:"documents"`
	SearchQueries []SearchQuery          `json:"search_queries"`
	SearchResults []SearchResults        `json:"search_results"`
	Meta          Meta                   `json:"meta"`
	ToolInputs    map[string]interface{} `json:"tool_inputs"`
}

type TokenCount struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
	BilledTokens   int `json:"billed_tokens"`
}

type Meta struct {
	APIVersion struct {
		Version string `json:"version"`
	} `json:"api_version"`
	BilledUnits struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"billed_units"`
}

type Citation struct {
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Text       string   `json:"text"`
	DocumentID []string `json:"document_id"`
}

type Documents struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"` // TODO: This needs to be updated
}

type SearchQuery struct {
	Text         string `json:"text"`
	GenerationID string `json:"generation_id"`
}

type SearchResults struct {
	SearchQuery []SearchQueryObject  `json:"search_query"`
	Connectors  []ConnectorsResponse `json:"connectors"`
	DocumentID  []string             `json:"documentId"`
}

type SearchQueryObject struct {
	Text         string `json:"text"`
	GenerationID string `json:"generationId"`
}

type ConnectorsResponse struct {
	ID              string            `json:"id"`
	UserAccessToken string            `json:"user_access_token"`
	ContOnFail      string            `json:"continue_on_failure"`
	Options         map[string]string `json:"options"`
}

// ChatCompletionChunk represents SSEvent a chat response is broken down on chat streaming
// Ref: https://docs.cohere.com/reference/about
type ChatCompletionChunk struct {
	IsFinished bool          `json:"is_finished"`
	EventType  string        `json:"event_type"`
	Text       string        `json:"text"`
	Response   FinalResponse `json:"response,omitempty"`
}

type FinalResponse struct {
	ResponseID   string     `json:"response_id"`
	Text         string     `json:"text"`
	GenerationID string     `json:"generation_id"`
	TokenCount   TokenCount `json:"token_count"`
	Meta         Meta       `json:"meta"`
	FinishReason string     `json:"finish_reason"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	Role    string `json:"role"`
	Message string `json:"message"`
	User    string `json:"user,omitempty"`
}

// ChatRequest is a request to complete a chat completion..
type ChatRequest struct {
	Model             string        `json:"model"`
	Message           string        `json:"message"`
	Temperature       float64       `json:"temperature,omitempty"`
	PreambleOverride  string        `json:"preamble_override,omitempty"`
	ChatHistory       []ChatHistory `json:"chat_history,omitempty"`
	ConversationID    string        `json:"conversation_id,omitempty"`
	PromptTruncation  string        `json:"prompt_truncation,omitempty"`
	Connectors        []string      `json:"connectors,omitempty"`
	SearchQueriesOnly bool          `json:"search_queries_only,omitempty"`
	CitiationQuality  string        `json:"citiation_quality,omitempty"`
	Stream            bool          `json:"stream,omitempty"`
}

type Connectors struct {
	ID              string            `json:"id"`
	UserAccessToken string            `json:"user_access_token"`
	ContOnFail      string            `json:"continue_on_failure"`
	Options         map[string]string `json:"options"`
}
