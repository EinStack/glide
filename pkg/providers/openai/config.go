package openai

import (
	"glide/pkg/providers"
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
	User             interface{}      `yaml:"user,omitempty"`
	Seed             interface{}      `yaml:"seed,omitempty"`
	Tools            []string         `yaml:"tools,omitempty"`
	ToolChoice       interface{}      `yaml:"tool_choice,omitempty"`
	ResponseFormat   interface{}      `yaml:"response_format,omitempty"` // TODO: should this be a part of the chat request API?
	// Stream           bool             `json:"stream,omitempty"` // TODO: we are not supporting this at the moment
}

type Config struct {
	Model         string           `yaml:"model"`
	APIKey        providers.Secret `yaml:"api_key" validate:"required"`
	DefaultParams *Params          `yaml:"default_params,omitempty"`
}
