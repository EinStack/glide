package clients

import (
	"encoding/json"
	"time"
)

type ClientConfig struct {
	Timeout *time.Duration `yaml:"timeout,omitempty" json:"timeout" swaggertype:"primitive,string"`
}

func DefaultClientConfig() *ClientConfig {
	defaultTimeout := 10 * time.Second

	return &ClientConfig{
		Timeout: &defaultTimeout,
	}
}

func (c *ClientConfig) MarshalJSON() ([]byte, error) {
	type Alias ClientConfig
	return json.Marshal(&struct {
		Timeout string `json:"timeout,omitempty"`
		*Alias
	}{
		Timeout: c.Timeout.String(),
		Alias:   (*Alias)(c),
	})
}

func (c *ClientConfig) MarshalYAML() ([]byte, error) {
	type Alias ClientConfig
	return json.Marshal(&struct {
		Timeout string `yaml:"timeout,omitempty"`
		*Alias
	}{
		Timeout: c.Timeout.String(),
		Alias:   (*Alias)(c),
	})
}
