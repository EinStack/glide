package providers

import (
	"errors"
	"fmt"
	"time"

	"glide/pkg/providers/openai"
	"glide/pkg/telemetry"
)

var ErrProviderNotFound = errors.New("provider not found")

type LangModelConfig struct {
	ID          string         `yaml:"id" json:"id" validate:"required"`
	Enabled     bool           `yaml:"enabled" json:"enabled"`
	ErrorBudget string         `yaml:"error_budget" json:"error_budget"` // TODO: make a custom filed to accept budget in the following format "10/m"
	Timeout     *time.Duration `yaml:"timeout,omitempty" json:"timeout" swaggertype:"primitive,integer"`
	OpenAI      *openai.Config `yaml:"openai" json:"openai"`
	// Add other providers like
	// Cohere *cohere.Config
	// Anthropic *anthropic.Config
}

func DefaultLangModelConfig() *LangModelConfig {
	defaultTimeout := 10 * time.Second

	return &LangModelConfig{
		Enabled: true,
		Timeout: &defaultTimeout,
	}
}

func (c *LangModelConfig) ToModel(tel *telemetry.Telemetry) (LanguageModel, error) {
	if c.OpenAI != nil {
		client, err := openai.NewClient(c.OpenAI, tel)
		if err != nil {
			return nil, fmt.Errorf("error initing openai client: %v", err)
		}

		return client, nil
	}

	return nil, ErrProviderNotFound
}

func (c *LangModelConfig) validateOneProvider() error {
	providersConfigured := 0

	if c.OpenAI != nil {
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
