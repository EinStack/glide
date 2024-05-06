package retry

import (
	"encoding/json"
	"time"
)

type ExpRetryConfig struct {
	MaxRetries     int            `yaml:"max_retries,omitempty" json:"max_retries"`
	BaseMultiplier int            `yaml:"base_multiplier,omitempty" json:"base_multiplier"`
	MinDelay       time.Duration  `yaml:"min_delay,omitempty" json:"min_delay" swaggertype:"primitive,string"`
	MaxDelay       *time.Duration `yaml:"max_delay,omitempty" json:"max_delay" swaggertype:"primitive,string"`
}

func DefaultExpRetryConfig() *ExpRetryConfig {
	maxDelay := 5 * time.Second

	return &ExpRetryConfig{
		MaxRetries:     3,
		BaseMultiplier: 2,
		MinDelay:       2 * time.Second,
		MaxDelay:       &maxDelay,
	}
}

func (e *ExpRetryConfig) MarshalJSON() ([]byte, error) {
	type Alias ExpRetryConfig
	return json.Marshal(&struct {
		MinDelay string  `json:"min_delay,omitempty"`
		MaxDelay *string `json:"max_delay,omitempty"`
		*Alias
	}{
		MinDelay: e.MinDelay.String(),
		MaxDelay: durationPtrToString(e.MaxDelay),
		Alias:    (*Alias)(e),
	})
}

func (e *ExpRetryConfig) MarshalYAML() (interface{}, error) {
	type Alias ExpRetryConfig
	return &struct {
		MinDelay string  `yaml:"min_delay,omitempty"`
		MaxDelay *string `yaml:"max_delay,omitempty"`
		*Alias
	}{
		MinDelay: e.MinDelay.String(),
		MaxDelay: durationPtrToString(e.MaxDelay),
		Alias:    (*Alias)(e),
	}, nil
}

func durationPtrToString(d *time.Duration) *string {
	if d == nil {
		return nil
	}
	s := d.String()
	return &s
}
