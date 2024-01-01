package cohere

import (
	"glide/pkg/config/fields"
)

// Params defines OpenAI-specific model params with the specific validation of values
// TODO: Add validations
type Params struct {
	Temperature       float64       `json:"temperature,omitempty"`
	Stream            bool          `json:"stream,omitempty"` // unsupported right now
	PreambleOverride  string        `json:"preamble_override,omitempty"`
	ChatHistory       []ChatHistory `json:"chat_history,omitempty"`
	ConversationID    string        `json:"conversation_id,omitempty"`
	PromptTruncation  string        `json:"prompt_truncation,omitempty"`
	Connectors        []string      `json:"connectors,omitempty"`
	SearchQueriesOnly bool          `json:"search_queries_only,omitempty"`
	CitiationQuality  string        `json:"citiation_quality,omitempty"`
}

func DefaultParams() Params {
	return Params{
		Temperature:       0.3,
		Stream:            false,
		PreambleOverride:  "",
		ChatHistory:       nil,
		ConversationID:    "",
		PromptTruncation:  "",
		Connectors:        []string{},
		SearchQueriesOnly: false,
		CitiationQuality:  "",
	}
}

func (p *Params) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = DefaultParams()

	type plain Params // to avoid recursion

	return unmarshal((*plain)(p))
}

type Config struct {
	BaseURL       string        `yaml:"base_url" json:"baseUrl" validate:"required"`
	ChatEndpoint  string        `yaml:"chat_endpoint" json:"chatEndpoint" validate:"required"`
	Model         string        `yaml:"model" json:"model" validate:"required"`
	APIKey        fields.Secret `yaml:"api_key" json:"-" validate:"required"`
	DefaultParams *Params       `yaml:"default_params,omitempty" json:"defaultParams"`
}

// DefaultConfig for OpenAI models
func DefaultConfig() *Config {
	defaultParams := DefaultParams()

	return &Config{
		BaseURL:       "https://api.cohere.ai/v1",
		ChatEndpoint:  "/chat",
		Model:         "command-light",
		DefaultParams: &defaultParams,
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = *DefaultConfig()

	type plain Config // to avoid recursion

	return unmarshal((*plain)(c))
}
