package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"glide/pkg/api/schemas"
	"go.uber.org/zap"
)

// ChatRequestSchema is an OpenAI-specific request schema
type ChatRequestSchema struct {
	Model            string              `json:"model"`
	Messages         []map[string]string `json:"messages"`
	Temperature      float64             `json:"temperature,omitempty"`
	TopP             float64             `json:"top_p,omitempty"`
	MaxTokens        int                 `json:"max_tokens,omitempty"`
	N                int                 `json:"n,omitempty"`
	StopWords        []string            `json:"stop,omitempty"`
	Stream           bool                `json:"stream,omitempty"`
	FrequencyPenalty int                 `json:"frequency_penalty,omitempty"`
	PresencePenalty  int                 `json:"presence_penalty,omitempty"`
	LogitBias        *map[int]float64    `json:"logit_bias,omitempty"`
	User             interface{}         `json:"user,omitempty"`
	Seed             interface{}         `json:"seed,omitempty"`
	Tools            []string            `json:"tools,omitempty"`
	ToolChoice       interface{}         `json:"tool_choice,omitempty"`
	ResponseFormat   interface{}         `json:"response_format,omitempty"`
}

// Chat sends a chat request to the specified OpenAI model.
func (c *Client) Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)

	// TODO: this is suspicious we do zero remapping of OpenAI response and send it back as is.
	//  Does it really work well across providers?
	chatResponse, err := c.doChatRequest(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

func (c *Client) createChatRequestSchema(request *schemas.ChatRequest) *ChatRequestSchema {
	var messages []map[string]string

	// Add items from messageHistory first
	messages = append(messages, request.MessageHistory...)

	// Add msg variable last
	messages = append(messages, request.Message)

	// Iterate through unifiedData.Params and add them to the request, otherwise leave the default value
	defaultParams := u.Params

	chatRequest := &ChatRequestSchema{
		Model:            c.config.Model,
		Messages:         messages,
		Temperature:      0.8,
		TopP:             1,
		MaxTokens:        100,
		N:                1,
		StopWords:        []string{},
		Stream:           false,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		LogitBias:        nil,
		User:             nil,
		Seed:             nil,
		Tools:            []string{},
		ToolChoice:       nil,
		ResponseFormat:   nil,
	}

	// TODO: set params

	return chatRequest
}

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequestSchema) (*schemas.ChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal openai chat request payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create openai chat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.telemetry.Logger.Debug(
		"openai chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send openai chat request: %w", err)
	}

	defer resp.Body.Close() // TODO: handle this error

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.telemetry.Logger.Error("failed to read openai chat response", zap.Error(err))
		}

		// TODO: Handle failure conditions

		c.telemetry.Logger.Error(
			"openai chat request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
		)

		// TODO: return errors
	}

	// Parse response
	var response schemas.ChatResponse

	return &response, json.NewDecoder(resp.Body).Decode(&response)
}
