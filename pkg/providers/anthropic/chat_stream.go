package anthropic

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "glide/pkg/api/schemas"
    "glide/pkg/providers/clients"
    "glide/pkg/telemetry"
    "go.uber.org/zap"
)

func (c *Client) SupportChatStream() bool {
	return true
}

type AnthropicChatStream struct {
    tel         *telemetry.Telemetry
    client      *http.client
    request     *http.Request
    response    *http.Response
    errMapper   *ErrorMapper 
}

func NewAnthropicChatStream(tel *telemetry.Telemetry, *http.client, request *http.Request, errMapper *ErrorMapper) *AnthropicChatStream {
    return &AnthropicChatStream{
        tel:       tel,
        client:    client,
        request:   request,
        errMapper: errMapper,
    }
}

// Open makes the HTTP request using the provided http.Client to initiate the chat stream.
func (s *AnthropicChatStream) Open(ctx context.Context) error {
    resp, err := s.client.Do(s.request)
    if err != nil {
        s.tel.L().Error("Failed to open chat stream", zap.Error(err))
        // Map and return the error using errMapper, if errMapper is defined.
        return s.errMapper.Map(err)
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        s.tel.L().Warn("Unexpected status code", zap.Int("status", resp.StatusCode))
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    s.response = resp
    s.tel.L().Info("Chat stream opened successfully")
    return nil
}

// Recv listens for and decodes incoming messages from the chat stream into ChatStreamChunk objects.
func (s *AnthropicChatStream) Recv() (*schemas.ChatStreamChunk, error) {
    if s.response == nil {
        s.tel.L().Error("Attempted to receive from an unopened stream")
        return nil, fmt.Errorf("stream not opened")
    }

    decoder := json.NewDecoder(s.response.Body)
    var chunk schemas.ChatStreamChunk
    if err := decoder.Decode(&chunk); err != nil {
        if err == io.EOF {
            s.tel.L().Info("Chat stream ended")
            return nil, nil // Stream ended normally.
        }
        s.tel.L().Error("Error during stream processing", zap.Error(err))
        return nil, err // An error occurred during stream processing.
    }

    return &chunk, nil
}

// Close ensures the chat stream is properly terminated by closing the response body.
func (s *AnthropicChatStream) Close() error {
    if s.response != nil {
        s.tel.L().Info("Closing chat stream")
        return s.response.Body.Close()
    }
    return nil
}