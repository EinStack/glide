package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/r3labs/sse/v2"
	"glide/pkg/providers/clients"

	"go.uber.org/zap"

	"glide/pkg/api/schemas"
)

var streamDoneMarker = []byte("[DONE]")

func (c *Client) SupportChatStream() bool {
	return true
}

func (c *Client) ChatStream(ctx context.Context, req *schemas.ChatRequest) <-chan *clients.ChatStreamResult {
	streamResultC := make(chan *clients.ChatStreamResult)

	go c.streamChat(ctx, req, streamResultC)

	return streamResultC
}

func (c *Client) streamChat(ctx context.Context, request *schemas.ChatRequest, resultC chan *clients.ChatStreamResult) {
	// Create a new chat request
	resp, err := c.initChatStream(ctx, request)

	defer close(resultC)

	if err != nil {
		resultC <- clients.NewChatStreamResult(nil, err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultC <- clients.NewChatStreamResult(nil, c.handleChatReqErrs(resp))
	}

	reader := sse.NewEventStreamReader(resp.Body, 4096) // TODO: should we expose maxBufferSize?

	var completionChunk ChatCompletionChunk

	for {
		started_at := time.Now()
		rawEvent, err := reader.ReadEvent()
		chunkLatency := time.Since(started_at)

		if err != nil {
			if err == io.EOF {
				c.tel.L().Debug("Chat stream is over", zap.String("provider", c.Provider()))

				return
			}

			c.tel.L().Warn(
				"Chat stream is unexpectedly disconnected",
				zap.String("provider", c.Provider()),
			)

			resultC <- clients.NewChatStreamResult(nil, clients.ErrProviderUnavailable)

			return
		}

		c.tel.L().Debug(
			"Raw chat stream chunk",
			zap.String("provider", c.Provider()),
			zap.ByteString("rawChunk", rawEvent),
		)

		event, err := clients.ParseSSEvent(rawEvent)

		if bytes.Equal(event.Data, streamDoneMarker) {
			return
		}

		if err != nil {
			resultC <- clients.NewChatStreamResult(nil, fmt.Errorf("failed to parse chat stream message: %v", err))
			return
		}

		if !event.HasContent() {
			c.tel.L().Debug(
				"Received an empty message in chat stream, skipping it",
				zap.String("provider", c.Provider()),
				zap.Any("msg", event),
			)

			continue
		}

		err = json.Unmarshal(event.Data, &completionChunk)
		if err != nil {
			resultC <- clients.NewChatStreamResult(nil, fmt.Errorf("failed to unmarshal chat stream chunk: %v", err))
			return
		}

		c.tel.L().Debug(
			"Chat response chunk",
			zap.String("provider", c.Provider()),
			zap.Any("chunk", completionChunk),
		)

		// TODO: use objectpool here
		chatRespChunk := schemas.ChatStreamChunk{
			ID:        completionChunk.ID,
			Created:   completionChunk.Created,
			Provider:  providerName,
			Cached:    false,
			ModelName: completionChunk.ModelName,
			ModelResponse: schemas.ModelResponse{
				SystemID: map[string]string{
					"system_fingerprint": completionChunk.SystemFingerprint,
				},
				Message: schemas.ChatMessage{
					Role:    completionChunk.Choices[0].Delta.Role,
					Content: completionChunk.Choices[0].Delta.Content,
				},
			},
			Latency: &chunkLatency,
			// TODO: Pass info if this is the final message
		}

		resultC <- clients.NewChatStreamResult(
			&chatRespChunk,
			nil,
		)
	}
}

// initChatStream establishes a new chat stream
func (c *Client) initChatStream(ctx context.Context, request *schemas.ChatRequest) (*http.Response, error) {
	chatRequest := *c.createChatRequestSchema(request)
	chatRequest.Stream = true

	rawPayload, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal openAI chat stream request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create OpenAI stream chat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(c.config.APIKey)))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.tel.L().Debug(
		"Stream chat request",
		zap.String("chatURL", c.chatURL),
		zap.Any("payload", chatRequest),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send OpenAI stream chat request: %w", err)
	}

	return resp, nil
}
