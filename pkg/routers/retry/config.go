package retry

import "time"

type ExpRetryConfig struct {
	MaxRetries int            `yaml:"max_retries,omitempty" json:"max_retries"`
	MinDelay   time.Duration  `yaml:"min_delay,omitempty" json:"min_delay" swaggertype:"primitive,integer"`
	MaxDelay   *time.Duration `yaml:"max_delay,omitempty" json:"max_delay" swaggertype:"primitive,integer"`
}

func DefaultExpRetryConfig() *ExpRetryConfig {
	maxDelay := 5 * time.Second

	return &ExpRetryConfig{
		MaxRetries: 3,
		MinDelay:   2 * time.Second,
		MaxDelay:   &maxDelay,
	}
}
