package clients

import (
	"time"

	"github.com/EinStack/glide/pkg/config/fields"
)

type ClientConfig struct {
	Timeout *fields.Duration `yaml:"timeout,omitempty" json:"timeout" swaggertype:"primitive,string"`
}

func DefaultClientConfig() *ClientConfig {
	defaultTimeout := 10 * time.Second

	return &ClientConfig{
		Timeout: (*fields.Duration)(&defaultTimeout),
	}
}
