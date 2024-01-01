package http

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"glide/pkg/api/schemas"
	"glide/pkg/routers"
)

type Handler = func(ctx context.Context, c *app.RequestContext)

// LangChatHandler
// @id glide-language-chat
// @Summary Language Chat
// @Description Talk to different LLMs Chat API via unified endpoint
// @tags lang
// @Accept application/json
// @Produce application/json
// @Router /v1/language/:router/chat [post]
func LangChatHandler(routerManager *routers.RouterManager) Handler {
	return func(ctx context.Context, c *app.RequestContext) {
		var req *schemas.UnifiedChatRequest

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(consts.StatusBadRequest, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		routerID := c.Param("router")
		router, err := routerManager.GetLangRouter(routerID)

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
	}
}

// HealthHandler
// @id glide-health
// @Summary Gateway Health
// @Description
// @tags operations
// @Accept application/json
// @Produce application/json
// @Router /v1/health/ [get]
func HealthHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, HealthSchema{Healthy: true})
}
