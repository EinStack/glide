package providers

import (
	"errors"
	"fmt"

	"glide/pkg/routers/latency"

	"glide/pkg/providers/clients"

	"glide/pkg/routers/health"

	"glide/pkg/providers/anthropic"
	"glide/pkg/providers/azureopenai"
	"glide/pkg/providers/cohere"
	"glide/pkg/providers/octoml"
	"glide/pkg/providers/openai"
	"glide/pkg/telemetry"
)

var ErrProviderNotFound = errors.New("provider not found")

type LangModelConfig struct {
	ID          string                `yaml:"id" json:"id" validate:"required"`           // Model instance ID (unique in scope of the router)
	Enabled     bool                  `yaml:"enabled" json:"enabled" validate:"required"` // Is the model enabled?
	ErrorBudget *health.ErrorBudget   `yaml:"error_budget" json:"error_budget" swaggertype:"primitive,string"`
	Latency     *latency.Config       `yaml:"latency" json:"latency"`
	Weight      int                   `yaml:"weight" json:"weight"`
	Client      *clients.ClientConfig `yaml:"client" json:"client"`
	// Add other providers like:
	OpenAI      *openai.Config      `yaml:"openai,omitempty" json:"openai,omitempty"`
	AzureOpenAI *azureopenai.Config `yaml:"azureopenai,omitempty" json:"azureopenai,omitempty"`
	Cohere      *cohere.Config      `yaml:"cohere,omitempty" json:"cohere,omitempty"`
	OctoML      *octoml.Config      `yaml:"octoml,omitempty" json:"octoml,omitempty"`
	Anthropic   *anthropic.Config   `yaml:"anthropic,omitempty" json:"anthropic,omitempty"`
}

func DefaultLangModelConfig() *LangModelConfig {
	return &LangModelConfig{
		Enabled:     true,
		Client:      clients.DefaultClientConfig(),
		ErrorBudget: health.DefaultErrorBudget(),
		Latency:     latency.DefaultConfig(),
		Weight:      1,
	}
}

func (c *LangModelConfig) ToModel(tel *telemetry.Telemetry) (*LangModel, error) {
	client, err := c.initClient(tel)
	if err != nil {
		return nil, fmt.Errorf("error initializing client: %v", err)
	}

	return NewLangModel(c.ID, client, *c.ErrorBudget, *c.Latency, c.Weight), nil
}

// initClient initializes the language model client based on the provided configuration.
// It takes a telemetry object as input and returns a LangModelProvider and an error.
func (c *LangModelConfig) initClient(tel *telemetry.Telemetry) (LangModelProvider, error) {
	switch {
	case c.OpenAI != nil:
		return openai.NewClient(c.OpenAI, c.Client, tel)
	// case c.AzureOpenAI != nil:
	// return azureopenai.NewClient(c.AzureOpenAI, c.Client, tel)
	// case c.Cohere != nil:
	// return cohere.NewClient(c.Cohere, c.Client, tel)
	// case c.OctoML != nil:
	//	return octoml.NewClient(c.OctoML, c.Client, tel)
	// case c.Anthropic != nil:
	//	return anthropic.NewClient(c.Anthropic, c.Client, tel)
	default:
		return nil, ErrProviderNotFound
	}
}

func (c *LangModelConfig) validateOneProvider() error {
	providersConfigured := 0

	if c.OpenAI != nil {
		providersConfigured++
	}

	if c.AzureOpenAI != nil {
		providersConfigured++
	}

	if c.Cohere != nil {
		providersConfigured++
	}

	if c.OctoML != nil {
		providersConfigured++
	}

	if c.Anthropic != nil {
		providersConfigured++
	}

	// check other providers here
	if providersConfigured == 0 {
		return fmt.Errorf("exactly one provider must be cofigured for model \"%v\", none is configured", c.ID)
	}

	if providersConfigured > 1 {
		return fmt.Errorf(
			"exactly one provider must be cofigured for model \"%v\", %v are configured",
			c.ID,
			providersConfigured,
		)
	}

	return nil
}

func (c *LangModelConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = *DefaultLangModelConfig()

	type plain LangModelConfig // to avoid recursion

	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return c.validateOneProvider()
}
