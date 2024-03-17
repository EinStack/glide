package routers

import (
	"context"
	"errors"

	"glide/pkg/routers/retry"
	"go.uber.org/zap"

	"glide/pkg/providers"

	"glide/pkg/api/schemas"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
)

var (
	ErrNoModels         = errors.New("no models configured for router")
	ErrNoModelAvailable = errors.New("could not handle request because all providers are not available")
)

type LangRouter struct {
	routerID          string
	Config            *LangRouterConfig
	chatModels        []*providers.LanguageModel
	chatStreamModels  []*providers.LanguageModel
	chatRouting       routing.LangModelRouting
	chatStreamRouting routing.LangModelRouting
	retry             *retry.ExpRetry
	tel               *telemetry.Telemetry
}

func NewLangRouter(cfg *LangRouterConfig, tel *telemetry.Telemetry) (*LangRouter, error) {
	chatModels, chatStreamModels, err := cfg.BuildModels(tel)
	if err != nil {
		return nil, err
	}

	chatRouting, chatStreamRouting, err := cfg.BuildRouting(chatModels, chatStreamModels)
	if err != nil {
		return nil, err
	}

	router := &LangRouter{
		routerID:          cfg.ID,
		Config:            cfg,
		chatModels:        chatModels,
		chatStreamModels:  chatStreamModels,
		retry:             cfg.BuildRetry(),
		chatRouting:       chatRouting,
		chatStreamRouting: chatStreamRouting,
		tel:               tel,
	}

	return router, err
}

func (r *LangRouter) ID() string {
	return r.routerID
}

func (r *LangRouter) Chat(ctx context.Context, req *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	if len(r.chatModels) == 0 {
		return nil, ErrNoModels
	}

	retryIterator := r.retry.Iterator()

	for retryIterator.HasNext() {
		modelIterator := r.chatRouting.Iterator()

		for {
			model, err := modelIterator.Next()

			if errors.Is(err, routing.ErrNoHealthyModels) {
				// no healthy model in the pool. Let's retry after some time
				break
			}

			langModel := model.(providers.LangModel)

			// Check if there is an override in the request
			if req.Override != nil {
				// Override the message if the language model ID matches the override model ID
				if langModel.ID() == req.Override.Model {
					req.Message = req.Override.Message
				}
			}

			resp, err := langModel.Chat(ctx, req)
			if err != nil {
				r.tel.L().Warn(
					"Lang model failed processing chat request",
					zap.String("routerID", r.ID()),
					zap.String("modelID", langModel.ID()),
					zap.String("provider", langModel.Provider()),
					zap.Error(err),
				)

				continue
			}

			resp.RouterID = r.routerID

			return resp, nil
		}

		// no providers were available to handle the request,
		//  so we have to wait a bit with a hope there is some available next time
		r.tel.L().Warn("No healthy model found to serve chat request, wait and retry", zap.String("routerID", r.ID()))

		err := retryIterator.WaitNext(ctx)
		if err != nil {
			// something has cancelled the context
			return nil, err
		}
	}

	// if we reach this part, then we are in trouble
	r.tel.L().Error("No model was available to handle chat request", zap.String("routerID", r.ID()))

	return nil, ErrNoModelAvailable
}

func (r *LangRouter) ChatStream(
	ctx context.Context,
	req *schemas.ChatStreamRequest,
	respC chan<- *schemas.ChatStreamResult,
) {
	if len(r.chatStreamModels) == 0 {
		respC <- schemas.NewChatStreamErrorResult(&schemas.ChatStreamError{
			ID:       req.ID,
			ErrCode:  "noModels",
			Message:  ErrNoModels.Error(),
			Metadata: req.Metadata,
		})

		return
	}

	retryIterator := r.retry.Iterator()

	for retryIterator.HasNext() {
		modelIterator := r.chatStreamRouting.Iterator()

	NextModel:
		for {
			model, err := modelIterator.Next()

			if errors.Is(err, routing.ErrNoHealthyModels) {
				// no healthy model in the pool. Let's retry after some time
				break
			}

			langModel := model.(providers.LangModel)
			modelRespC, err := langModel.ChatStream(ctx, req)
			if err != nil {
				r.tel.L().Error(
					"Lang model failed to create streaming chat request",
					zap.String("routerID", r.ID()),
					zap.String("modelID", langModel.ID()),
					zap.String("provider", langModel.Provider()),
					zap.Error(err),
				)

				continue
			}

			for chunkResult := range modelRespC {
				err = chunkResult.Error()
				if err != nil {
					r.tel.L().Warn(
						"Lang model failed processing streaming chat request",
						zap.String("routerID", r.ID()),
						zap.String("modelID", langModel.ID()),
						zap.String("provider", langModel.Provider()),
						zap.Error(err),
					)

					// It's challenging to hide an error in case of streaming chat as consumer apps
					//  may have already used all chunks we streamed this far (e.g. showed them to their users like OpenAI UI does),
					//  so we cannot easily restart that process from scratch
					respC <- schemas.NewChatStreamErrorResult(&schemas.ChatStreamError{
						ID:       req.ID,
						ErrCode:  "modelUnavailable",
						Message:  err.Error(),
						Metadata: req.Metadata,
					})

					continue NextModel
				}

				respC <- schemas.NewChatStreamResult(chunkResult.Chunk())
			}

			return
		}

		// no providers were available to handle the request,
		//  so we have to wait a bit with a hope there is some available next time
		r.tel.L().Warn(
			"No healthy model found to serve streaming chat request, wait and retry",
			zap.String("routerID", r.ID()),
		)

		err := retryIterator.WaitNext(ctx)
		if err != nil {
			// something has cancelled the context
			respC <- schemas.NewChatStreamErrorResult(&schemas.ChatStreamError{
				ID:       req.ID,
				ErrCode:  "other",
				Message:  err.Error(),
				Metadata: req.Metadata,
			})

			return
		}
	}

	// if we reach this part, then we are in trouble
	r.tel.L().Error(
		"No model was available to handle streaming chat request. Try to configure more fallback models to avoid this",
		zap.String("routerID", r.ID()),
	)

	respC <- schemas.NewChatStreamErrorResult(&schemas.ChatStreamError{
		ID:       req.ID,
		ErrCode:  "allModelsUnavailable",
		Message:  ErrNoModelAvailable.Error(),
		Metadata: req.Metadata,
	})
}
