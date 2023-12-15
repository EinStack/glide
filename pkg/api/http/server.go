package http

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

type HTTPServer struct {
	server *server.Hertz
}

func NewHttpServer(config *HTTPServerConfig) (*HTTPServer, error) {
	return &HTTPServer{
		server: config.ToServer(),
	}, nil
}

func (srv *HTTPServer) Run() error {
    srv.server.GET("/health", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, utils.H{"healthy": true})
    })

    srv.server.POST("/chat", func(ctx context.Context, c *app.RequestContext) {

        // Pass the client request body to SendRequest
        resp, err := pkg.LlmRouter(c)
		
        if err != nil {
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
