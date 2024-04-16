package http

import (
	"context"
	"errors"
	"sync"

	"glide/pkg/telemetry"
	"go.uber.org/zap"

	"github.com/gofiber/contrib/websocket"
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
//	@Description	Talk to different LLM Chat APIs via unified endpoint
//	@tags			Language
//	@Param			router	path	string						true	"Router ID"
//	@Param			payload	body	schemas.ChatRequest	true	"Request Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.ChatResponse
//	@Failure		400	{object}	http.ErrorSchema
//	@Failure		404	{object}	http.ErrorSchema
//	@Router			/v1/language/{router}/chat [POST]
func LangChatHandler(routerManager *routers.RouterManager) Handler {
	return func(c *fiber.Ctx) error {
		if !c.Is("json") {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorSchema{
				Message: "Glide accepts only JSON payloads",
			})
		}

		// Unmarshal request body
		var req *schemas.ChatRequest

		err := c.BodyParser(&req)
		if err != nil {
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

func LangStreamRouterValidator(routerManager *routers.RouterManager) Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			routerID := c.Params("router")

			_, err := routerManager.GetLangRouter(routerID)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(ErrorSchema{
					Message: err.Error(),
				})
			}

			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	}
}

// LangStreamChatHandler
//
//	@id				glide-language-chat-stream
//	@Summary		Language Chat
//	@Description	Talk to different LLM Stream Chat APIs via a unified websocket endpoint
//	@tags			Language
//	@Param			router	    			path		string	true	"Router ID"
//	@Param			Connection				header		string	true	"Websocket Connection Type"
//	@Param			Upgrade 				header		string	true	"Upgrade header"
//	@Param			Sec-WebSocket-Key  		header		string	true	"Websocket Security Token"
//	@Param			Sec-WebSocket-Version  	header		string	true	"Websocket Security Token"
//	@Accept			json
//	@Success		101
//	@Failure		426
//	@Failure		404	{object}	http.ErrorSchema
//	@Router			/v1/language/{router}/chatStream [GET]
func LangStreamChatHandler(tel *telemetry.Telemetry, routerManager *routers.RouterManager) Handler {
	// TODO: expose websocket connection configs https://github.com/gofiber/contrib/tree/main/websocket
	return websocket.New(func(c *websocket.Conn) {
		routerID := c.Params("router")
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index

		var (
			err error
			wg  sync.WaitGroup
		)

		chatStreamC := make(chan *schemas.ChatStreamMessage)

		router, _ := routerManager.GetLangRouter(routerID)

		defer close(chatStreamC)
		defer c.Conn.Close()

		wg.Add(1)

		go func() {
			defer wg.Done()

			for chatStreamMsg := range chatStreamC {
				if err = c.WriteJSON(chatStreamMsg); err != nil {
					break
				}
			}
		}()

		for {
			var chatRequest schemas.ChatStreamRequest

			if err = c.ReadJSON(&chatRequest); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					tel.L().Warn("Streaming Chat connection is closed", zap.Error(err), zap.String("routerID", routerID))
				}

				tel.L().Debug("Streaming chat connection is closed by client", zap.Error(err), zap.String("routerID", routerID))

				break
			}

			// TODO: handle termination gracefully
			wg.Add(1)

			go func(chatRequest schemas.ChatStreamRequest) {
				defer wg.Done()

				router.ChatStream(context.Background(), &chatRequest, chatStreamC)
			}(chatRequest)
		}

		wg.Wait()
	})
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
