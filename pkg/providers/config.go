package providers

import (
	"errors"
	"fmt"
	"glide/pkg/providers/openai"
	"glide/pkg/telemetry"
	"time"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
)

type LangModelConfig struct {
	ID      string         `yaml:"id"`
	Enabled bool           `yaml:"enabled"`
	Timeout *time.Duration `yaml:"timeout,omitempty"`
	OpenAI  *openai.Config
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

func (m *LangModelConfig) validateOneProvider() error {
	providersConfigured := 0

	if m.OpenAI != nil {
		providersConfigured++
	}

	// check other providers here
	if providersConfigured == 0 {
		return fmt.Errorf("exactly one provider must be cofigured for model \"%v\", none is configured", m.ID)
	}

	if providersConfigured > 1 {
		return fmt.Errorf(
			"exactly one provider must be cofigured for model \"%v\", %v are configured",
			m.ID,
			providersConfigured,
		)
	}

	return nil
}

func (m *LangModelConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*m = *DefaultLangModelConfig()

	type plain LangModelConfig // to avoid recursion

	if err := unmarshal((*plain)(m)); err != nil {
		return err
	}

	return m.validateOneProvider()
}
