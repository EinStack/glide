package http

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Server struct {
	server *server.Hertz
}

func NewServer(config *ServerConfig) (*Server, error) {
	return &Server{
		server: config.ToServer(),
	}, nil
}

func (srv *Server) Run() error {
	srv.server.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"healthy": true})
	})

	return srv.server.Run()
}

func (srv *Server) Shutdown(_ context.Context) error {
	exitWaitTime := srv.server.GetOptions().ExitWaitTimeout

	println(fmt.Sprintf("Begin graceful shutdown, wait at most %d seconds...", exitWaitTime/time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), exitWaitTime)
	defer cancel()

	return srv.server.Shutdown(ctx)
}
