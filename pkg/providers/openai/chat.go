package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/EinStack/glide/pkg/api/schemas"
	"go.uber.org/zap"
)

// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		Model:            cfg.Model,
		Temperature:      cfg.DefaultParams.Temperature,
		TopP:             cfg.DefaultParams.TopP,
		MaxTokens:        cfg.DefaultParams.MaxTokens,
		N:                cfg.DefaultParams.N,
		StopWords:        cfg.DefaultParams.StopWords,
		Stream:           false,
		FrequencyPenalty: cfg.DefaultParams.FrequencyPenalty,
		PresencePenalty:  cfg.DefaultParams.PresencePenalty,
		LogitBias:        cfg.DefaultParams.LogitBias,
		User:             cfg.DefaultParams.User,
		Seed:             cfg.DefaultParams.Seed,
		Tools:            cfg.DefaultParams.Tools,
		ToolChoice:       cfg.DefaultParams.ToolChoice,
		ResponseFormat:   cfg.DefaultParams.ResponseFormat,
	}
}

// Chat sends a chat request to the specified OpenAI model.
func (c *Client) Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createRequestSchema(request)
	chatRequest.Stream = false

	chatResponse, err := c.doChatRequest(ctx, chatRequest)

	if err != nil {
		return nil, err
	}

	if len(chatResponse.ModelResponse.Message.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

// createRequestSchema creates a new ChatRequest object based on the given request.
func (c *Client) createRequestSchema(request *schemas.ChatRequest) *ChatRequest {
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

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.ChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal openai chat request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create openai chat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(c.config.APIKey)))

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.logger.Debug(
		"Chat Request",
		zap.String("chatURL", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send openai chat request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.errMapper.Map(resp)
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(
			"Failed to read chat response",
			zap.Error(err),
			zap.ByteString("rawResponse", bodyBytes),
		)

		return nil, err
	}

	// Parse the response JSON
	var chatCompletion ChatCompletion

	c.logger.Debug(
		"Raw chat response",
		zap.ByteString("resp", bodyBytes),
	)

	err = json.Unmarshal(bodyBytes, &chatCompletion)
	if err != nil {
		c.logger.Error(
			"Failed to unmarshal chat response",
			zap.ByteString("rawResponse", bodyBytes),
			zap.Error(err),
		)

		return nil, err
	}

	// Map response to ChatResponse schema
	response := schemas.ChatResponse{
		ID:        chatCompletion.ID,
		Created:   chatCompletion.Created,
		Provider:  providerName,
		ModelName: chatCompletion.ModelName,
		Cached:    false,
		ModelResponse: schemas.ModelResponse{
			Metadata: map[string]string{
				"system_fingerprint": chatCompletion.SystemFingerprint,
			},
			Message: schemas.ChatMessage{
				Role:    chatCompletion.Choices[0].Message.Role,
				Content: chatCompletion.Choices[0].Message.Content,
			},
			TokenUsage: schemas.TokenUsage{
				PromptTokens:   chatCompletion.Usage.PromptTokens,
				ResponseTokens: chatCompletion.Usage.CompletionTokens,
				TotalTokens:    chatCompletion.Usage.TotalTokens,
			},
		},
	}

	return &response, nil
}
