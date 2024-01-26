package openai

import (
	"glide/pkg/config/fields"
)

// Params defines OpenAI-specific model params with the specific validation of values
// TODO: Add validations
type Params struct {
	Temperature      float64          `yaml:"temperature,omitempty" json:"temperature"`
	TopP             float64          `yaml:"top_p,omitempty" json:"top_p"`
	MaxTokens        int              `yaml:"max_tokens,omitempty" json:"max_tokens"`
	N                int              `yaml:"n,omitempty" json:"n"`
	StopWords        []string         `yaml:"stop,omitempty" json:"stop"`
	FrequencyPenalty int              `yaml:"frequency_penalty,omitempty" json:"frequency_penalty"`
	PresencePenalty  int              `yaml:"presence_penalty,omitempty" json:"presence_penalty"`
	LogitBias        *map[int]float64 `yaml:"logit_bias,omitempty" json:"logit_bias"`
	User             *string          `yaml:"user,omitempty" json:"user"`
	Seed             *int             `yaml:"seed,omitempty" json:"seed"`
	Tools            []string         `yaml:"tools,omitempty" json:"tools"`
	ToolChoice       interface{}      `yaml:"tool_choice,omitempty" json:"tool_choice"`
	ResponseFormat   interface{}      `yaml:"response_format,omitempty" json:"response_format"` // TODO: should this be a part of the chat request API?
	Stream           bool             `json:"stream,omitempty"`                                 // TODO: we are not supporting this at the moment
}

func DefaultParams() Params {
	return Params{
		Temperature: 0.8,
		TopP:        1,
		MaxTokens:   100,
		N:           1,
		StopWords:   []string{},
		Tools:       []string{},
		Stream:      true,
	}
}

func (p *Params) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = DefaultParams()

	type plain Params // to avoid recursion

	return unmarshal((*plain)(p))
}

type Config struct {
	BaseURL       string        `yaml:"baseUrl" json:"baseUrl" validate:"required"`
	ChatEndpoint  string        `yaml:"chatEndpoint" json:"chatEndpoint" validate:"required"`
	Model         string        `yaml:"model" json:"model" validate:"required"`
	APIKey        fields.Secret `yaml:"api_key" json:"-" validate:"required"`
	DefaultParams *Params       `yaml:"defaultParams,omitempty" json:"defaultParams"`
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
