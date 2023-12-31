package http

import (
	"context"
	"fmt"
	"time"

	"glide/pkg/api/schemas"

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
	srv.server.POST("/v1/language/:router/chat/", func(ctx context.Context, c *app.RequestContext) {
		var req *schemas.UnifiedChatRequest

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(consts.StatusBadRequest, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		routerID := c.Param("router")

		// TODO: call the model router and return the unified response

		c.JSON(consts.StatusOK, utils.H{
			"message": fmt.Sprintf("%v was requested", routerID),
		})
	})

	srv.server.GET("/v1/health/", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, HealthSchema{Healthy: true})
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
