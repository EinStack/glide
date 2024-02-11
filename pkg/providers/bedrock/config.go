package bedrock

import (
	"glide/pkg/config/fields"
)

// Params defines OpenAI-specific model params with the specific validation of values
// TODO: Add validations
type Params struct {
	Temperature  float64  `yaml:"temperature,omitempty" json:"temperature"`
	TopP         float64  `yaml:"top_p,omitempty" json:"top_p"`
	MaxTokens    int      `yaml:"max_tokens,omitempty" json:"max_tokens"`
	StopSequence []string `yaml:"stop_sequences,omitempty" json:"stop"`
}

func DefaultParams() Params {
	return Params{
		Temperature:  0.8,
		TopP:         1,
		MaxTokens:    100,
		StopSequence: []string{},
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
	AccessKey     fields.Secret `yaml:"access_key" json:"-" validate:"required"`
	SecretKey     fields.Secret `yaml:"secret_key" json:"-" validate:"required"`
	DefaultParams *Params       `yaml:"defaultParams,omitempty" json:"defaultParams"`
}

// DefaultConfig for OpenAI models
func DefaultConfig() *Config {
	defaultParams := DefaultParams()

	return &Config{
		BaseURL:       "", // This needs to come from config. https://bedrock-runtime.{{AWS_Region}}.amazonaws.com/
		ChatEndpoint:  "/model",
		Model:         "amazon.titan-text-express-v1",
		DefaultParams: &defaultParams,
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = *DefaultConfig()

	type plain Config // to avoid recursion

	return unmarshal((*plain)(c))
}
