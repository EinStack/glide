package http

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"time"
)

type HTTPServerConfig struct {
	// TODO: add fields
}

func (cfg *HTTPServerConfig) ToServer() *server.Hertz {
	// TODO: do real server build based on provided config
	return server.Default(
		server.WithIdleTimeout(1*time.Second),
		server.WithHostPorts("127.0.0.1:9099"),
		server.WithMaxRequestBodySize(20<<20),
		server.WithTransport(netpoll.NewTransporter),
	)
}
