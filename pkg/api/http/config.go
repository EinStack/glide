package http

import (
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/network/standard"

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
	maxReqBodySizeBytes := 4 * 1024 * 1024 // 4Mb
	readTimeout := 3 * time.Minute
	writeTimeout := 3 * time.Minute
	idleTimeout := 3 * time.Minute

	return &ServerConfig{
		Host:               "127.0.0.1",
		Port:               9099,
		IdleTimeout:        &idleTimeout,
		ReadTimeout:        &readTimeout,
		WriteTimeout:       &writeTimeout,
		MaxRequestBodySize: &maxReqBodySizeBytes,
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
