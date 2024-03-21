package clients

import "time"

type ClientConfig struct {
	Timeout *time.Duration `yaml:"timeout,omitempty" json:"timeout" swaggertype:"primitive,string"`
}

func DefaultClientConfig() *ClientConfig {
	defaultTimeout := time.Duration(10) * time.Second

	return &ClientConfig{
		Timeout: &defaultTimeout,
	}
}
