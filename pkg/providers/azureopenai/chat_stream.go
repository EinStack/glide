package azureopenai

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

var (
	StopReason       = "stop"
	streamDoneMarker = []byte("[DONE]")
)

// ChatStream represents chat stream for a specific request
type ChatStream struct {
	tel         *telemetry.Telemetry
	client      *http.Client
	req         *http.Request
	reqID       string
	reqMetadata *schemas.Metadata
	resp        *http.Response
	reader      *sse.EventStreamReader
	errMapper   *ErrorMapper
}

func NewChatStream(
	tel *telemetry.Telemetry,
	client *http.Client,
	req *http.Request,
	reqID string,
	reqMetadata *schemas.Metadata,
	errMapper *ErrorMapper,
) *ChatStream {
	return &ChatStream{
		tel:         tel,
		client:      client,
		req:         req,
		reqID:       reqID,
		reqMetadata: reqMetadata,
		errMapper:   errMapper,
	}
}

// Open initializes and opens a ChatStream.
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

// Recv receives a chat stream chunk from the ChatStream and returns a ChatStreamChunk object.
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
			s.tel.L().Info(
				"EOF: [DONE] marker found in chat stream",
				zap.String("provider", providerName),
			)

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
			return nil, fmt.Errorf("failed to unmarshal AzureOpenAI chat stream chunk: %v", err)
		}

		responseChunk := completionChunk.Choices[0]

		var finishReason *schemas.FinishReason

		if responseChunk.FinishReason == StopReason {
			finishReason = &schemas.Complete
		}

		// TODO: use objectpool here
		return &schemas.ChatStreamChunk{
			ID:        s.reqID,
			Provider:  providerName,
			Cached:    false,
			ModelName: completionChunk.ModelName,
			Metadata:  s.reqMetadata,
			ModelResponse: schemas.ModelChunkResponse{
				Metadata: &schemas.Metadata{
					"response_id":        completionChunk.ID,
					"system_fingerprint": completionChunk.SystemFingerprint,
				},
				Message: schemas.ChatMessage{
					Role:    responseChunk.Delta.Role,
					Content: responseChunk.Delta.Content,
				},
				FinishReason: finishReason,
			},
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

func (c *Client) ChatStream(ctx context.Context, req *schemas.ChatStreamRequest) (clients.ChatStream, error) {
	// Create a new chat request
	httpRequest, err := c.makeStreamReq(ctx, req)
	if err != nil {
		return nil, err
	}

	return NewChatStream(
		c.tel,
		c.httpClient,
		httpRequest,
		req.ID,
		req.Metadata,
		c.errMapper,
	), nil
}

func (c *Client) createRequestFromStream(request *schemas.ChatStreamRequest) *ChatRequest {
	// TODO: consider using objectpool to optimize memory allocation
	chatRequest := *c.chatRequestTemplate // hoping to get a copy of the template

	chatRequest.Messages = make([]ChatMessage, 0, len(request.MessageHistory)+1)

	// Add items from messageHistory first and the new chat message last
	for _, message := range request.MessageHistory {
		chatRequest.Messages = append(chatRequest.Messages, ChatMessage{Role: message.Role, Content: message.Content})
	}

	chatRequest.Messages = append(chatRequest.Messages, ChatMessage{Role: request.Message.Role, Content: request.Message.Content})

	return &chatRequest
}

func (c *Client) makeStreamReq(ctx context.Context, req *schemas.ChatStreamRequest) (*http.Request, error) {
	chatRequest := c.createRequestFromStream(req)

	chatRequest.Stream = true

	rawPayload, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal AzureOpenAI chat stream request payload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create AzureOpenAI stream chat request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api-key", string(c.config.APIKey))
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Connection", "keep-alive")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.tel.L().Debug(
		"Stream chat request",
		zap.String("chatURL", c.chatURL),
		zap.Any("payload", chatRequest),
	)

	return request, nil
}
