package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "glide/docs" // importing docs package to include them into the binary
	"time"

	"glide/pkg/routers"

	"glide/pkg/telemetry"
)

type Server struct {
	config        *ServerConfig
	telemetry     *telemetry.Telemetry
	routerManager *routers.RouterManager
	server        *fiber.App
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
	srv.server.Use(NotFoundHandler)

	v1 := srv.server.Group("/v1")

	v1.Get("/language/", LangRoutersHandler(srv.routerManager))
	v1.Post("/language/:router/chat/", LangChatHandler(srv.routerManager))

	v1.Get("/health/", HealthHandler)

	//schemaDocURL := swagger.URL(fmt.Sprintf("http://%v/v1/swagger/doc.json", srv.config.Address()))
	//v1.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, schemaDocURL))

	return srv.server.Listen(":9099") // TODO: take it from configs
}

func (srv *Server) Shutdown(ctx context.Context) error {
	exitWaitTime := 5 * time.Second

	srv.telemetry.Logger.Info(
		fmt.Sprintf("Begin graceful shutdown, wait at most %d seconds...", exitWaitTime/time.Second),
	)

	c, cancel := context.WithTimeout(ctx, exitWaitTime)
	defer cancel()

	return srv.server.ShutdownWithContext(c) //nolint:contextcheck
}
