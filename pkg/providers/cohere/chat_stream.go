package cohere


import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/r3labs/sse/v2"
	"glide/pkg/providers/clients"
	"glide/pkg/telemetry"

	"go.uber.org/zap"

	"glide/pkg/api/schemas"
)


var (
	StopReason       = "stream-end"
)

// ChatStream represents cohere chat stream for a specific request
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

func (s *ChatStream) Open() error {
	resp, err := s.client.Do(s.req) //nolint:bodyclose
	fmt.Print(resp.StatusCode)
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
		fmt.Print(err)
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

		responseChunk := completionChunk

		var finishReason *schemas.FinishReason

		if responseChunk.IsFinished {
			finishReason = &schemas.Complete
			return &schemas.ChatStreamChunk{
				ID:        s.reqID,
				Provider:  providerName,
				Cached:    false,
				ModelName: "NA",
				Metadata:  s.reqMetadata,
				ModelResponse: schemas.ModelChunkResponse{
					Metadata: &schemas.Metadata{
						"generationId":        responseChunk.Response.GenerationID,
						"responseId": responseChunk.Response.ResponseID,
					},
					Message: schemas.ChatMessage{
						Role:    "model",
						Content: responseChunk.Text,
					},
					FinishReason: finishReason,
				},
			}, nil
		}

		// TODO: use objectpool here
		return &schemas.ChatStreamChunk{
			ID:        s.reqID,
			Provider:  providerName,
			Cached:    false,
			ModelName: "NA",
			Metadata:  s.reqMetadata,
			ModelResponse: schemas.ModelChunkResponse{
				Message: schemas.ChatMessage{
					Role:    "model",
					Content: responseChunk.Text,
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

	chatRequest.Message = request.Message.Content

	// Build the Cohere specific ChatHistory
	if len(request.MessageHistory) > 0 {
		chatRequest.ChatHistory = make([]ChatHistory, len(request.MessageHistory))
		for i, message := range request.MessageHistory {
			chatRequest.ChatHistory[i] = ChatHistory{
				// Copy the necessary fields from message to ChatHistory
				// For example, if ChatHistory has a field called "Text", you can do:
				Role:    message.Role,
				Message: message.Content,
				User:    "",
			}
		}
	}

	return &chatRequest
}

func (c *Client) makeStreamReq(ctx context.Context, req *schemas.ChatStreamRequest) (*http.Request, error) {
	chatRequest := c.createRequestFromStream(req)

	chatRequest.Stream = true

	fmt.Print(chatRequest)

	rawPayload, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal cohere chat stream request payload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create cohere stream chat request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(c.config.APIKey)))
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

