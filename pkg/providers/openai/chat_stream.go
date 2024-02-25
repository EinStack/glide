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

func (c *Client) SupportChatStream() bool {
	return true
}

func (c *Client) ChatStream(ctx context.Context, request *schemas.ChatRequest, responseC chan<- schemas.ChatResponse) error {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)
	chatRequest.Stream = true

	rawPayload, err := json.Marshal(chatRequest)
	if err != nil {
		return fmt.Errorf("unable to marshal openai chat stream request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return fmt.Errorf("unable to create openai stream chat request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(c.config.APIKey)))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.telemetry.Logger.Debug(
		"openai stream chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", chatRequest),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send openai stream chat request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO: handle all specter of errors like in a sync chat API
		return clients.ErrProviderUnavailable
	}

	reader := sse.NewEventStreamReader(resp.Body, 4096) // TODO: should we expose maxBufferSize?

	var completionChunk ChatCompletionChunk

	for {
		rawEvent, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				// TODO: the stream is over
				// erChan <- nil
				// return
			}

			// TODO: we are disconnected

			// erChan <- err
			// return
		}

		event, err := clients.ParseSSEvent(rawEvent)
		if err != nil {
			// TODO: handle
		}

		if len(event.ID) > 0 || len(event.Data) > 0 || len(event.Event) > 0 || len(event.Retry) > 0 {
			// has some content

			err = json.Unmarshal(event.Data, &completionChunk)
			if err != nil {
				// TODO: handle
			}

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
			}

			responseC <- chatResponse
			continue
		}

		c.telemetry.Logger.Debug(
			"Received an empty event on OpenAI chat stream",
			zap.Any("event", event),
		)
	}
}
