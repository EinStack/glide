package http

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
	"glide/pkg/api"
	"sync"
	"log/slog"
	
)

type HTTPServer struct {
	server *server.Hertz
}

type ServerManager struct {
	httpServer *HTTPServer
	shutdownWG *sync.WaitGroup
}

func NewHttpServer(config *HTTPServerConfig) (*HTTPServer, error) {
	return &HTTPServer{
		server: config.ToServer(),
	}, nil
}

// This func has the routes for the server
func (srv *HTTPServer) Run() error {
    srv.server.GET("/health", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, utils.H{"healthy": true})
    })

    srv.server.POST("/chat", func(ctx context.Context, c *app.RequestContext) {

		slog.Info("POST request at /chat received")

        // Pass the client request body to SendRequest
        resp, err := api.Router(c)
		

		slog.Info("Provider response received")
		
        if err != nil {
			slog.Error("Error in Router Response: %v", err)
            c.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
            return
        }
        c.JSON(consts.StatusOK, utils.H{"response": resp})
    })

    return srv.server.Run()
}

func (srv *HTTPServer) Shutdown(_ context.Context) error {
	exitWaitTime := srv.server.GetOptions().ExitWaitTimeout

	println(fmt.Sprintf("Begin graceful shutdown, wait at most %d seconds...", exitWaitTime/time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), exitWaitTime)
	defer cancel()

	return srv.server.Shutdown(ctx)
}

func NewServerManager(httpConfig *HTTPServerConfig) (*ServerManager, error) {
	httpServer, err := NewHttpServer(httpConfig)
	// TODO: init other servers like gRPC in future

	if err != nil {
		return nil, err
	}

	return &ServerManager{
		httpServer: httpServer,
		shutdownWG: &sync.WaitGroup{},
	}, nil
}

func (mgr *ServerManager) Start() {
	if mgr.httpServer != nil {
		mgr.shutdownWG.Add(1)

		go func() {
			defer mgr.shutdownWG.Done()

			// TODO: log the error
			err := mgr.httpServer.Run()

			println(err)
		}()
	}
}

func (mgr *ServerManager) Shutdown(ctx context.Context) error {
	var err error

	if mgr.httpServer != nil {
		err = mgr.httpServer.Shutdown(ctx)
	}

	mgr.shutdownWG.Wait()

	return err
}

