package providers

import (
	"errors"
	"fmt"

	"glide/pkg/providers/clients"

	"glide/pkg/routers/health"

	"glide/pkg/providers/azureopenai"
	"glide/pkg/providers/cohere"
	"glide/pkg/providers/octoml"
	"glide/pkg/providers/openai"
	"glide/pkg/telemetry"
)

var ErrProviderNotFound = errors.New("provider not found")

type LangModelConfig struct {
	ID          string                `yaml:"id" json:"id" validate:"required"` // Model instance ID (unique in scope of the router)
	Enabled     bool                  `yaml:"enabled" json:"enabled"`           // Is the model enabled?
	ErrorBudget health.ErrorBudget    `yaml:"error_budget" json:"error_budget" swaggertype:"primitive,string"`
	Client      *clients.ClientConfig `yaml:"client" json:"client"`
	OpenAI      *openai.Config        `yaml:"openai" json:"openai"`
	AzureOpenAI *azureopenai.Config   `yaml:"azureopenai" json:"azureopenai"`
	Cohere      *cohere.Config        `yaml:"cohere" json:"cohere"`
	OctoML      *octoml.Config        `yaml:"octoml" json:"octoml"`
	// Add other providers like
	// Cohere *cohere.Config
	// Anthropic *anthropic.Config
}

func DefaultLangModelConfig() *LangModelConfig {
	return &LangModelConfig{
		Enabled:     true,
		Client:      clients.DefaultClientConfig(),
		ErrorBudget: health.DefaultErrorBudget(),
	}
}

func (c *LangModelConfig) ToModel(tel *telemetry.Telemetry) (*LangModel, error) {
	if c.OpenAI != nil {
		client, err := openai.NewClient(c.OpenAI, c.Client, tel)
		if err != nil {
			return nil, fmt.Errorf("error initing openai client: %v", err)
		}

		return NewLangModel(c.ID, client, c.ErrorBudget), nil
	}

	if c.AzureOpenAI != nil {
		client, err := azureopenai.NewClient(c.AzureOpenAI, c.Client, tel)
		if err != nil {
			return nil, fmt.Errorf("error initing azureopenai client: %v", err)
		}

		return NewLangModel(c.ID, client, c.ErrorBudget), nil
	}

	if c.Cohere != nil {
		client, err := cohere.NewClient(c.Cohere, c.Client, tel)
		if err != nil {
			return nil, fmt.Errorf("error initing cohere client: %v", err)
		}

		return NewLangModel(c.ID, client, c.ErrorBudget), nil
	}

	if c.OctoML != nil {
		client, err := octoml.NewClient(c.OctoML, c.Client, tel)
		if err != nil {
			return nil, fmt.Errorf("error initing cohere client: %v", err)
		}

		return NewLangModel(c.ID, client, c.ErrorBudget), nil
	}

	return nil, ErrProviderNotFound
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
