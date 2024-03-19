package anthropic

import (
	"context"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
)

type Client struct {
	APIKey string
}

func (c *Client) SupportChatStream() bool {
	return true
}

func (c *Client) ChatStream(ctx context.Context, chatReq *schemas.ChatRequest) (clients.ChatStream, error) {
	apiURL := "https://api.anthropic.com/v1/complete"
	requestBody := map[string]interface{}{
		"model":                "claude-2",
		"prompt":               chatReq.Message, // Assuming chatReq.Message contains the user's message
		"max_tokens_to_sample": 256,
		"stream":               true,
	}
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	chatStream := &AnthropicChatStream{
		responseBody: resp.Body,
	}

	return chatStream, nil
}

type AnthropicChatStream struct {
	responseBody io.ReadCloser
}

func (s *AnthropicChatStream) Receive() (string, error) {
	decoder := json.NewDecoder(s.responseBody)
	for {
		var event map[string]interface{}
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				return "", nil // No more events, return nil error
			}
			return "", err
		}

		eventType, ok := event["type"].(string)
		if !ok {
			return "", fmt.Errorf("missing event type")
		}

		switch eventType {
		case "completion":
			completionData, ok := event["completion"].(string)
			if !ok {
				return "", fmt.Errorf("missing completion data")
			}
			return completionData, nil
		case "error":
			errorData, ok := event["error"].(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("missing error data")
			}
			errorType, _ := errorData["type"].(string)
			errorMessage, _ := errorData["message"].(string)
			return "", fmt.Errorf("error from Anthropic API: %s - %s", errorType, errorMessage)
		}
	}
}

func (s *AnthropicChatStream) Close() error {
	return s.responseBody.Close()
}
