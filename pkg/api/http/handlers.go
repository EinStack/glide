package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"

	"glide/pkg/api/schemas"
	"glide/pkg/routers"
)

type Handler = func(c *fiber.Ctx) error

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
	return func(c *fiber.Ctx) error {
		// Unmarshal request body
		var req *schemas.UnifiedChatRequest

		err := c.BodyParser(&req)

		if err != nil {
			// Return bad request error

			return c.Status(fiber.StatusBadRequest).JSON(ErrorSchema{
				Message: err.Error(),
			})
		}

		// Get router ID from path
		routerID := c.Params("router")
		router, err := routerManager.GetLangRouter(routerID)

		if errors.Is(err, routers.ErrRouterNotFound) {
			// Return not found error
			return c.Status(fiber.StatusNotFound).JSON(ErrorSchema{
				Message: err.Error(),
			})
		}

		// Chat with router
		resp, err := router.Chat(c.Context(), req)

		if err != nil {
			// Return internal server error
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorSchema{
				Message: err.Error(),
			})
		}

		// Return chat response
		return c.Status(fiber.StatusOK).JSON(resp)
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
	return func(c *fiber.Ctx) error {
		configuredRouters := routerManager.GetLangRouters()
		cfgs := make([]*routers.LangRouterConfig, 0, len(configuredRouters))

		for _, router := range configuredRouters {
			cfgs = append(cfgs, router.Config)
		}

		return c.Status(fiber.StatusOK).JSON(RouterListSchema{Routers: cfgs})
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
func HealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(HealthSchema{Healthy: true})
}

func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(ErrorSchema{
		Message: "The route is not found",
	})
}
