package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"glide/pkg/api/schemas"

	"glide/pkg/routers"

	"glide/pkg/telemetry"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Server struct {
	telemetry     *telemetry.Telemetry
	routerManager *routers.RouterManager
	server        *server.Hertz
}

func NewServer(config *ServerConfig, tel *telemetry.Telemetry, routerManager *routers.RouterManager) (*Server, error) {
	srv := config.ToServer()

	return &Server{
		telemetry:     tel,
		routerManager: routerManager,
		server:        srv,
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
		router, err := srv.routerManager.GetLangRouter(routerID)

		if errors.Is(err, routers.ErrRouterNotFound) {
			c.JSON(consts.StatusNotFound, ErrorSchema{
				Message: err.Error(),
			})
			return
		}

		resp, err := router.Chat(ctx, req)
		if err != nil {
			// TODO: do a better handling, not everything is going to be an internal error
			c.JSON(consts.StatusInternalServerError, ErrorSchema{
				Message: err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, resp)
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
