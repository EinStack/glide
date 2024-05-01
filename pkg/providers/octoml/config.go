package octoml

import (
	"github.com/EinStack/glide/pkg/config/fields"
)

// Params defines OctoML-specific model params with the specific validation of values
// TODO: Add validations
type Params struct {
	Temperature      float64  `yaml:"temperature,omitempty" json:"temperature"`
	TopP             float64  `yaml:"top_p,omitempty" json:"top_p"`
	MaxTokens        int      `yaml:"max_tokens,omitempty" json:"max_tokens"`
	StopWords        []string `yaml:"stop,omitempty" json:"stop"`
	FrequencyPenalty int      `yaml:"frequency_penalty,omitempty" json:"frequency_penalty"`
	PresencePenalty  int      `yaml:"presence_penalty,omitempty" json:"presence_penalty"`
	// Stream           bool             `json:"stream,omitempty"` // TODO: we are not supporting this at the moment
}

func DefaultParams() Params {
	return Params{
		Temperature: 1,
		TopP:        1,
		MaxTokens:   100,
		StopWords:   []string{},
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

// DefaultConfig for OctoML models
func DefaultConfig() *Config {
	defaultParams := DefaultParams()

	return &Config{
		BaseURL:       "https://text.octoai.run/v1",
		ChatEndpoint:  "/chat/completions",
		Model:         "mistral-7b-instruct-fp16",
		DefaultParams: &defaultParams,
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = *DefaultConfig()

	type plain Config // to avoid recursion

	return unmarshal((*plain)(c))
}
