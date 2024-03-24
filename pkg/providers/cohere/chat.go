package cohere

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"
	"go.uber.org/zap"
)


// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		Model:             cfg.Model,
		Temperature:       cfg.DefaultParams.Temperature,
		PreambleOverride:  cfg.DefaultParams.PreambleOverride,
		ChatHistory:       cfg.DefaultParams.ChatHistory,
		ConversationID:    cfg.DefaultParams.ConversationID,
		PromptTruncation:  cfg.DefaultParams.PromptTruncation,
		Connectors:        cfg.DefaultParams.Connectors,
		SearchQueriesOnly: cfg.DefaultParams.SearchQueriesOnly,
		CitiationQuality:  cfg.DefaultParams.CitiationQuality,
		Stream:            false,
	}
}

// Chat sends a chat request to the specified cohere model.
func (c *Client) Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createRequestSchema(request)

	chatResponse, err := c.doChatRequest(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.ModelResponse.Message.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

func (c *Client) createRequestSchema(request *schemas.ChatRequest) *ChatRequest {
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

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.ChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal cohere chat request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create cohere chat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.tel.Logger.Debug(
		"cohere chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send cohere chat request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.tel.Logger.Error("failed to read cohere chat response", zap.Error(err))
		}

		c.tel.Logger.Error(
			"cohere chat request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.Any("headers", resp.Header),
		)

		if resp.StatusCode != http.StatusOK {
			return c.handleErrorResponse(resp)
		}

		// Server & client errors result in the same error to keep gateway resilient
		return nil, clients.ErrProviderUnavailable
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.tel.Logger.Error("failed to read cohere chat response", zap.Error(err))
		return nil, err
	}

	// Parse the response JSON
	var responseJSON map[string]interface{}

	err = json.Unmarshal(bodyBytes, &responseJSON)
	if err != nil {
		c.tel.Logger.Error("failed to parse cohere chat response", zap.Error(err))
		return nil, err
	}

	// Parse the response JSON
	var cohereCompletion ChatCompletion

	err = json.Unmarshal(bodyBytes, &cohereCompletion)
	if err != nil {
		c.tel.Logger.Error("failed to parse cohere chat response", zap.Error(err))
		return nil, err
	}

	// Map response to ChatResponse schema
	response := schemas.ChatResponse{
		ID:        cohereCompletion.ResponseID,
		Created:   int(time.Now().UTC().Unix()), // Cohere doesn't provide this
		Provider:  providerName,
		ModelName: c.config.Model,
		Cached:    false,
		ModelResponse: schemas.ModelResponse{
			SystemID: map[string]string{
				"generationId": cohereCompletion.GenerationID,
				"responseId":   cohereCompletion.ResponseID,
			},
			Message: schemas.ChatMessage{
				Role:    "model",
				Content: cohereCompletion.Text,
				Name:    "",
			},
			TokenUsage: schemas.TokenUsage{
				PromptTokens:   cohereCompletion.TokenCount.PromptTokens,
				ResponseTokens: cohereCompletion.TokenCount.ResponseTokens,
				TotalTokens:    cohereCompletion.TokenCount.TotalTokens,
			},
		},
	}

	return &response, nil
}

func (c *Client) handleErrorResponse(resp *http.Response) (*schemas.ChatResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.tel.Logger.Error("failed to read cohere chat response", zap.Error(err))
		return nil, err
	}

	c.tel.Logger.Error(
		"cohere chat request failed",
		zap.Int("status_code", resp.StatusCode),
		zap.String("response", string(bodyBytes)),
		zap.Any("headers", resp.Header),
	)

	if resp.StatusCode == http.StatusTooManyRequests {
		cooldownDelay, err := c.getCooldownDelay(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cooldown delay from headers: %w", err)
		}

		return nil, clients.NewRateLimitError(&cooldownDelay)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, clients.ErrUnauthorized
	}

	return nil, clients.ErrProviderUnavailable
}

func (c *Client) getCooldownDelay(resp *http.Response) (time.Duration, error) {
	retryAfter := resp.Header.Get("Retry-After")

	cooldownDelay, err := time.ParseDuration(retryAfter)
	if err != nil {
		return 0, fmt.Errorf("failed to parse cooldown delay from headers: %w", err)
	}

	return cooldownDelay, nil
}
