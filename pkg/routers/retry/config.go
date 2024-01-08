package retry

import "time"

type ExpRetryConfig struct {
	MaxRetries int
	MinDelay   time.Duration
	MaxDelay   *time.Duration
}

func DefaultExpRetryConfig() *ExpRetryConfig {
	maxDelay := 5 * time.Second

	return &ExpRetryConfig{
		MaxRetries: 3,
		MinDelay:   2 * time.Second,
		MaxDelay:   &maxDelay,
	}
}
