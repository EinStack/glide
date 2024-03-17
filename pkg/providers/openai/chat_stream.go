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
	"glide/pkg/telemetry"

	"go.uber.org/zap"

	"glide/pkg/api/schemas"
)

var streamDoneMarker = []byte("[DONE]")

// ChatStream represents OpenAI chat stream for a specific request
type ChatStream struct {
	tel       *telemetry.Telemetry
	client    *http.Client
	req       *http.Request
	resp      *http.Response
	reader    *sse.EventStreamReader
	errMapper *ErrorMapper
}

func NewChatStream(tel *telemetry.Telemetry, client *http.Client, req *http.Request, errMapper *ErrorMapper) *ChatStream {
	return &ChatStream{
		tel:       tel,
		client:    client,
		req:       req,
		errMapper: errMapper,
	}
}

func (s *ChatStream) Open() error {
	resp, err := s.client.Do(s.req) //nolint:bodyclose
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return s.errMapper.Map(resp)
	}

	s.resp = resp
	s.reader = sse.NewEventStreamReader(resp.Body, 4096) // TODO: should we expose maxBufferSize?

	return nil
}

func (s *ChatStream) Recv() (*schemas.ChatStreamChunk, error) {
	var completionChunk ChatCompletionChunk

	for {
		rawEvent, err := s.reader.ReadEvent()
		if err != nil {
			s.tel.L().Warn(
				"Chat stream is unexpectedly disconnected",
				zap.String("provider", providerName),
				zap.Error(err),
			)

			// if err is io.EOF, this still means that the stream is interrupted unexpectedly
			//  because the normal stream termination is done via finding out streamDoneMarker

			return nil, clients.ErrProviderUnavailable
		}

		s.tel.L().Debug(
			"Raw chat stream chunk",
			zap.String("provider", providerName),
			zap.ByteString("rawChunk", rawEvent),
		)

		event, err := clients.ParseSSEvent(rawEvent)

		if bytes.Equal(event.Data, streamDoneMarker) {
			return nil, io.EOF
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse chat stream message: %v", err)
		}

		if !event.HasContent() {
			s.tel.L().Debug(
				"Received an empty message in chat stream, skipping it",
				zap.String("provider", providerName),
				zap.Any("msg", event),
			)

			continue
		}

		err = json.Unmarshal(event.Data, &completionChunk)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal chat stream chunk: %v", err)
		}

		// TODO: use objectpool here
		return &schemas.ChatStreamChunk{
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
			// TODO: Pass info if this is the final message
		}, nil
	}
}

func (s *ChatStream) Close() error {
	if s.resp != nil {
		return s.resp.Body.Close()
	}

	return nil
}

func (c *Client) SupportChatStream() bool {
	return true
}

func (c *Client) ChatStream(ctx context.Context, req *schemas.ChatRequest) (clients.ChatStream, error) {
	// Create a new chat request
	request, err := c.makeStreamReq(ctx, req)
	if err != nil {
		return nil, err
	}

	return NewChatStream(
		c.tel,
		c.httpClient,
		request,
		c.errMapper,
	), nil
}

func (c *Client) makeStreamReq(ctx context.Context, request *schemas.ChatRequest) (*http.Request, error) {
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

	return req, nil
}
