package openai

import (
	"glide/pkg/config/fields"
)

// Params defines OpenAI-specific model params with the specific validation of values
// TODO: Add validations
type Params struct {
	Temperature      float64          `yaml:"temperature,omitempty"`
	TopP             float64          `yaml:"top_p,omitempty"`
	MaxTokens        int              `yaml:"max_tokens,omitempty"`
	N                int              `yaml:"n,omitempty"`
	StopWords        []string         `yaml:"stop,omitempty"`
	FrequencyPenalty int              `yaml:"frequency_penalty,omitempty"`
	PresencePenalty  int              `yaml:"presence_penalty,omitempty"`
	LogitBias        *map[int]float64 `yaml:"logit_bias,omitempty"`
	User             *string          `yaml:"user,omitempty"`
	Seed             *int             `yaml:"seed,omitempty"`
	Tools            []string         `yaml:"tools,omitempty"`
	ToolChoice       interface{}      `yaml:"tool_choice,omitempty"`
	ResponseFormat   interface{}      `yaml:"response_format,omitempty"` // TODO: should this be a part of the chat request API?
	// Stream           bool             `json:"stream,omitempty"` // TODO: we are not supporting this at the moment
}

func DefaultParams() Params {
	return Params{
		Temperature: 0.8,
		TopP:        1,
		MaxTokens:   100,
		N:           1,
		StopWords:   []string{},
		Tools:       []string{},
	}
}

func (p *Params) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = DefaultParams()

	type plain Params // to avoid recursion

	return unmarshal((*plain)(p))
}

type Config struct {
	BaseURL       string        `yaml:"base_url"`
	ChatEndpoint  string        `yaml:"chat_endpoint"`
	Model         string        `yaml:"model"`
	APIKey        fields.Secret `yaml:"api_key" validate:"required"`
	DefaultParams *Params       `yaml:"default_params,omitempty"`
}

// DefaultConfig for OpenAI models
func DefaultConfig() *Config {
	defaultParams := DefaultParams()

	return &Config{
		BaseURL:       "https://api.openai.com/v1",
		ChatEndpoint:  "/chat/completions",
		Model:         "gpt-3.5-turbo",
		DefaultParams: &defaultParams,
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = *DefaultConfig()

	type plain Config // to avoid recursion

	return unmarshal((*plain)(c))
}
