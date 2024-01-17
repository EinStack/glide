package http

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
)

type ServerConfig struct {
	HostPort string
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		HostPort: "0.0.0.0:7685", // TODO: Should this be read from config?
	}
}

func (cfg *ServerConfig) ToServer() *server.Hertz {
	// TODO: do real server build based on provided config
	return server.Default(
		server.WithIdleTimeout(1*time.Second),
		server.WithHostPorts(cfg.HostPort),
		server.WithMaxRequestBodySize(20<<20),
		server.WithTransport(netpoll.NewTransporter),
	)
}
