package http

import (
	"context"
	"fmt"
	"time"

	"glide/pkg/routers"

	"glide/pkg/telemetry"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Server struct {
	telemetry *telemetry.Telemetry
	router    *routers.Router
	server    *server.Hertz
}

func NewServer(config *ServerConfig, tel *telemetry.Telemetry, router *routers.Router) (*Server, error) {
	srv := config.ToServer()

	return &Server{
		telemetry: tel,
		router:    router,
		server:    srv,
	}, nil
}

func (srv *Server) Run() error {
	srv.server.POST("/v1/language/{}/chat/", func(c context.Context, ctx *app.RequestContext) {
		// TODO: call the lang router
	})

	srv.server.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"healthy": true})
	})

	return srv.server.Run()
}

func (srv *Server) Shutdown(_ context.Context) error {
	exitWaitTime := srv.server.GetOptions().ExitWaitTimeout

	srv.telemetry.Logger.Info(
		fmt.Sprintf("Begin graceful shutdown, wait at most %d seconds...", exitWaitTime/time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), exitWaitTime)
	defer cancel()

	return srv.server.Shutdown(ctx) //nolint:contextcheck
}
