package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/r3labs/sse/v2"
	"glide/pkg/providers/clients"
	"go.uber.org/zap"

	"glide/pkg/api/schemas"
)

var (
	streamDoneMarker = []byte("[DONE]")
)

func (c *Client) SupportChatStream() bool {
	return true
}

func (c *Client) ChatStream(ctx context.Context, request *schemas.ChatRequest, responseC chan<- schemas.ChatResponse) error {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)
	chatRequest.Stream = true

	rawPayload, err := json.Marshal(chatRequest)
	if err != nil {
		return fmt.Errorf("unable to marshal openAI chat stream request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))

	if err != nil {
		return fmt.Errorf("unable to create OpenAI stream chat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(c.config.APIKey)))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.telemetry.Logger.Debug(
		"Stream chat request",
		zap.String("chatURL", c.chatURL),
		zap.Any("payload", chatRequest),
	)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send OpenAI stream chat request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			c.telemetry.Logger.Error(
				"Failed to read stream chat response error",
				zap.String("provider", c.Provider()),
				zap.Int("statusCode", resp.StatusCode),
				zap.Error(err),
			)
		}

		// TODO: handle all specter of errors like in a sync chat API
		c.telemetry.L().Error(
			"Failed to start chat stream",
			zap.String("provider", c.Provider()),
			zap.Int("statusCode", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.Any("headers", resp.Header),
		)

		return clients.ErrProviderUnavailable
	}

	reader := sse.NewEventStreamReader(resp.Body, 4096) // TODO: should we expose maxBufferSize?

	var completionChunk ChatCompletionChunk

	for {
		rawEvent, err := reader.ReadEvent()

		if err != nil {
			if err == io.EOF {
				c.telemetry.L().Debug("Chat stream is over", zap.String("provider", c.Provider()))

				return nil
			}

			c.telemetry.L().Warn(
				"Chat stream is unexpectedly interrupted by disconnection",
				zap.String("provider", c.Provider()),
			)

			return clients.ErrProviderUnavailable
		}

		c.telemetry.L().Debug("Raw chat stream chunk", zap.String("provider", c.Provider()), zap.ByteString("rawChunk", rawEvent))

		event, err := clients.ParseSSEvent(rawEvent)

		if bytes.Equal(event.Data, streamDoneMarker) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to parse chat stream message: %v", err)
		}

		if len(event.ID) > 0 || len(event.Data) > 0 || len(event.Event) > 0 || len(event.Retry) > 0 {
			// has some content

			err = json.Unmarshal(event.Data, &completionChunk)

			if err != nil {
				return fmt.Errorf("failed to unmarshal chat stream message: %v", err)
			}

			c.telemetry.L().Debug("Chat response chunk", zap.String("provider", c.Provider()), zap.Any("chunk", completionChunk))

			// TODO: use objectpool here
			chatResponse := schemas.ChatResponse{
				ID:        completionChunk.ID,
				Created:   completionChunk.Created,
				Provider:  providerName,
				Cached:    false,
				ModelName: completionChunk.ModelName,
				ModelResponse: schemas.ProviderResponse{
					SystemID: map[string]string{
						"system_fingerprint": completionChunk.SystemFingerprint,
					},
					Message: schemas.ChatMessage{
						Role:    completionChunk.Choices[0].Delta.Role,
						Content: completionChunk.Choices[0].Delta.Content,
						Name:    "",
					},
				},
				// TODO: Pass info if this is the final message
			}

			responseC <- chatResponse
			continue
		}

		c.telemetry.Logger.Debug(
			"Received an empty message in chat stream, skipping it",
			zap.String("provider", c.Provider()),
			zap.Any("msg", event),
		)
	}
}
