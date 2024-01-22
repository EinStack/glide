package http

import (
	"context"
	"fmt"
	"time"

	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "glide/docs" // importing docs package to include them into the binary

	"glide/pkg/routers"

	"glide/pkg/telemetry"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type Server struct {
	config        *ServerConfig
	telemetry     *telemetry.Telemetry
	routerManager *routers.RouterManager
	server        *server.Hertz
}

func NewServer(config *ServerConfig, tel *telemetry.Telemetry, routerManager *routers.RouterManager) (*Server, error) {
	srv := config.ToServer()

	return &Server{
		config:        config,
		telemetry:     tel,
		routerManager: routerManager,
		server:        srv,
	}, nil
}

func (srv *Server) Run() error {
	defaultGroup := srv.server.Group("/v1")

	defaultGroup.GET("/language/", LangRoutersHandler(srv.routerManager))
	defaultGroup.POST("/language/:router/chat/", LangChatHandler(srv.routerManager))

	defaultGroup.GET("/health/", HealthHandler)

	schemaDocURL := swagger.URL(fmt.Sprintf("http://%v/v1/swagger/doc.json", srv.config.Address()))
	defaultGroup.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, schemaDocURL))

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
