package api

import "github.com/EinStack/glide/pkg/api/http"

// Config defines configuration for all API types we support (e.g. HTTP, gRPC)
type Config struct {
	HTTP *http.ServerConfig `yaml:"http" validate:"required"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTP: http.DefaultServerConfig(),
	}
}
