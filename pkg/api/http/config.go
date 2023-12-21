package http

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
)

type ServerConfig struct {
	// TODO: add fields
}

func (cfg *ServerConfig) ToServer() *server.Hertz {
	// TODO: do real server build based on provided config
	return server.Default(
		server.WithIdleTimeout(1*time.Second),
		server.WithHostPorts("127.0.0.1:9099"),
		server.WithMaxRequestBodySize(20<<20),
		server.WithTransport(netpoll.NewTransporter),
	)
}
