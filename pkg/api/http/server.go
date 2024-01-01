package http

import (
	"context"
	"fmt"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"time"

	"glide/pkg/routers"

	"glide/pkg/telemetry"

	"github.com/cloudwego/hertz/pkg/app/server"
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
	defaultGroup := srv.server.Group("/v1")

	defaultGroup.POST("/language/:router/chat/", LangChatHandler(srv.routerManager))
	defaultGroup.GET("/health/", HealthHandler)
	defaultGroup.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("http://localhost:9099/v1/swagger/doc.json")))

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
