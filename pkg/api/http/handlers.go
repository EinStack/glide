package http

import (
	"context"
	"encoding/json"
	"errors"

	"glide/pkg/api/schemas"
	"glide/pkg/routers"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Handler = func(ctx context.Context, c *app.RequestContext)

// Swagger 101:
// - https://github.com/swaggo/swag/tree/master/example/celler

// LangChatHandler
//
//	@id				glide-language-chat
//	@Summary		Language Chat
//	@Description	Talk to different LLMs Chat API via unified endpoint
//	@tags			Language
//	@Param			router	path	string						true	"Router ID"
//	@Param			payload	body	schemas.UnifiedChatRequest	true	"Request Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.UnifiedChatResponse
//	@Failure		400	{object}	http.ErrorSchema
//	@Failure		404	{object}	http.ErrorSchema
//	@Router			/v1/language/{router}/chat [POST]
func LangChatHandler(routerManager *routers.RouterManager) Handler {
	return func(ctx context.Context, c *app.RequestContext) {
		// Unmarshal request body
		var req []schemas.UnifiedChatRequest

		err := json.Unmarshal(c.Request.Body(), &req)
		if err != nil {
			// Return bad request error
			c.JSON(consts.StatusBadRequest, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		// Bind JSON to request
		err = c.BindJSON(&req)
		if err != nil {
			// Return bad request error
			c.JSON(consts.StatusBadRequest, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		// Get router ID from path
		routerID := c.Param("router")
		router, err := routerManager.GetLangRouter(routerID)

		if errors.Is(err, routers.ErrRouterNotFound) {
			// Return not found error
			c.JSON(consts.StatusNotFound, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		// Chat with router
		resp, err := router.Chat(ctx, req)
		if err != nil {
			// Return internal server error
			c.JSON(consts.StatusInternalServerError, ErrorSchema{
				Message: err.Error(),
			})

			return
		}

		// Return chat response
		c.JSON(consts.StatusOK, resp)
	}
}

// LangRoutersHandler
//
//	@id				glide-language-routers
//	@Summary		Language Router List
//	@Description	Retrieve list of configured language routers and their configurations
//	@tags			Language
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	http.RouterListSchema
//	@Router			/v1/language/ [GET]
func LangRoutersHandler(routerManager *routers.RouterManager) Handler {
	return func(ctx context.Context, c *app.RequestContext) {
		configuredRouters := routerManager.GetLangRouters()
		cfgs := make([]*routers.LangRouterConfig, 0, len(configuredRouters))

		for _, router := range configuredRouters {
			cfgs = append(cfgs, router.Config)
		}

		c.JSON(consts.StatusOK, RouterListSchema{Routers: cfgs})
	}
}

// HealthHandler
//
//	@id			glide-health
//	@Summary	Gateway Health
//	@Description
//	@tags		Operations
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	http.HealthSchema
//	@Router		/v1/health/ [get]
func HealthHandler(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, HealthSchema{Healthy: true})
}
