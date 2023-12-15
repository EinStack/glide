package cohere

// ProviderConfig represents the provider configuration
type CohereChatProviderConfig struct {
	Model            string      `json:"model" validate:"lowercase"`
	Message         string      `json:"message" validate:"required"`
	Stream           bool        `json:"stream" validate:"omitempty,boolean"`
	PreambleOverride string      `json:"preamble_override,omitempty"`
	ChatHistory      interface{} `json:"chat_history,omitempty" validate:"omitempty,dive"`
	PromptTruncation string      `json:"prompt_truncation,omitempty"`
	Connectors       interface{} `json:"connectors,omitempty" validate:"omitempty,dive"`
	SearchQueryOnly  bool        `json:"search_query_only,omitempty" validate:"omitempty,boolean"`
	Documents        interface{} `json:"documents,omitempty" validate:"omitempty,dive"`
	CitationQuality  string      `json:"citation_quality,omitempty"`
	Temperature      float32     `json:"temperature,omitempty" validate:"omitempty,number,gte=0"`
}

// CohereChatCompleteConfig represents the configuration for Cohere chat completion
func CohereChatDefaultConfig() CohereChatProviderConfig {
	return CohereChatProviderConfig{
		Model:            "command-light",
		Message:         "hello Cohere",
		Stream:           false,
		Temperature:      0.3,
		PreambleOverride: "", // TODO: determine how this affects the request
		ChatHistory:      nil,
		PromptTruncation: "",
		Connectors:       nil,
		SearchQueryOnly:  false,
		Documents:        nil,
		CitationQuality:  "accurate",
	}
}
