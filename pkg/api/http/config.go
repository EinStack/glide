package http

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"time"

	"github.com/cloudwego/hertz/pkg/common/config"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type ServerConfig struct {
	Host               string         `yaml:"host"`
	Port               int            `yaml:"port"`
	ReadTimeout        *time.Duration `yaml:"read_timeout"`
	WriteTimeout       *time.Duration `yaml:"write_timeout"`
	IdleTimeout        *time.Duration `yaml:"idle_timeout"`
	MaxRequestBodySize *int           `yaml:"max_request_body_size"`
}

func DefaultServerConfig() *ServerConfig {
	maxReqBodySize := 4 * 1024 * 1024
	readTimeout := 3 * time.Second
	writeTimeout := 3 * time.Second
	idleTimeout := 1 * time.Second

	return &ServerConfig{
		Host:               "127.0.0.1",
		Port:               9099,
		IdleTimeout:        &idleTimeout,
		ReadTimeout:        &readTimeout,
		WriteTimeout:       &writeTimeout,
		MaxRequestBodySize: &maxReqBodySize,
	}
}

func (cfg *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

func (cfg *ServerConfig) ToServer() *server.Hertz {
	// More configs are listed on https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/engine/
	serverOptions := []config.Option{
		server.WithHostPorts(cfg.Address()),
		// https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/network-lib/#choosing-appropriate-network-library
		server.WithTransport(standard.NewTransporter),
		server.WithStreamBody(true),
	}

	if cfg.IdleTimeout != nil {
		serverOptions = append(serverOptions, server.WithIdleTimeout(*cfg.IdleTimeout))
	}

	if cfg.ReadTimeout != nil {
		serverOptions = append(serverOptions, server.WithReadTimeout(*cfg.ReadTimeout))
	}

	if cfg.WriteTimeout != nil {
		serverOptions = append(serverOptions, server.WithWriteTimeout(*cfg.WriteTimeout))
	}

	if cfg.MaxRequestBodySize != nil {
		serverOptions = append(serverOptions, server.WithMaxRequestBodySize(*cfg.MaxRequestBodySize))
	}

	return server.Default(serverOptions...)
}
