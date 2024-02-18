package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"

	"github.com/gofiber/contrib/swagger"

	"github.com/gofiber/fiber/v2"
	_ "glide/docs" // importing docs package to include them into the binary

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
	srv.server.Use(swagger.New(swagger.Config{
		Title:    "Glide API Docs",
		BasePath: "/v1/",
		Path:     "swagger",
		FilePath: "./docs/swagger.json",
	}))

	srv.server.Use(fiberzap.New(fiberzap.Config{
		Logger: srv.telemetry.Logger,
	}))

	v1 := srv.server.Group("/v1")

	v1.Get("/language/", LangRoutersHandler(srv.routerManager))
	v1.Post("/language/:router/chat/", LangChatHandler(srv.routerManager))

	v1.Use("/language/:router/chatStream", LangStreamRouterValidator(srv.routerManager))
	v1.Get("/language/:router/chatStream", LangStreamChatHandler())

	v1.Get("/health/", HealthHandler)

	srv.server.Use(NotFoundHandler)

	return srv.server.Listen(srv.config.Address())
}

func (srv *Server) Shutdown(ctx context.Context) error {
	exitWaitTime := 5 * time.Second

	srv.telemetry.Logger.Info(
		fmt.Sprintf("Begin graceful shutdown, wait at most %d seconds...", exitWaitTime/time.Second),
	)

	c, cancel := context.WithTimeout(ctx, exitWaitTime)
	defer cancel()

	if err := srv.server.ShutdownWithContext(c); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			srv.telemetry.Logger.Info("Server closed forcefully due to shutdown timeout")
			return nil
		}

		return err
	}

	return nil
}
