package routers

import "glide/pkg/providers/openai"

type Config struct {
	LanguageRouters *[]LanguageRouter `yaml:"language"`
}

type RoutingStrategy struct {
	// TODO: make it based on the routing strategies supported
}

type LanguageModel struct {
	ID        string `yaml:"id"`
	TimeoutMs *int   `yaml:"timeout_ms,omitempty"` // TODO: try to use Duration to bring more flexibility
	OpenAI    *openai.Config
	// Add other providers like
	// Cohere *cohere.Config
	// Anthropic *anthropic.Config
	// TODO: add validation to ensure only one provider is configured among available
}

type LanguageRouter struct {
	ID              string          `yaml:"id"`
	Enabled         bool            `yaml:"enabled"`
	RoutingStrategy RoutingStrategy `yaml:"strategy"`
	Models          []LanguageModel `yaml:"models"`
}

func DefaultConfig() *Config {
	return &Config{}
}
