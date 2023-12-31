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
	User             interface{}      `yaml:"user,omitempty"`
	Seed             interface{}      `yaml:"seed,omitempty"`
	Tools            []string         `yaml:"tools,omitempty"`
	ToolChoice       interface{}      `yaml:"tool_choice,omitempty"`
	ResponseFormat   interface{}      `yaml:"response_format,omitempty"` // TODO: should this be a part of the chat request API?
	// Stream           bool             `json:"stream,omitempty"` // TODO: we are not supporting this at the moment
}

// Defaults
// Temperature:      0.8,
// TopP:             1,
// MaxTokens:        100,
// N:                1,
// StopWords:        []string{},
// Stream:           false,
// FrequencyPenalty: 0,
// PresencePenalty:  0,
// LogitBias:        nil,
// User:             nil,
// Seed:             nil,
// Tools:            []string{},
// ToolChoice:       nil,
// ResponseFormat:   nil,

// defaultChatModel = "gpt-3.5-turbo"
// defaultEndpoint  = "/chat/completions"

type Config struct {
	BaseURL       string        `yaml:"base_url"`
	ChatEndpoint  string        `yaml:"chat_endpoint"`
	Model         string        `yaml:"model"`
	APIKey        fields.Secret `yaml:"api_key" validate:"required"`
	DefaultParams *Params       `yaml:"default_params,omitempty"`
}

// https://api.openai.com/v1
