package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"glide/pkg/api/schemas"
	"go.uber.org/zap"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is an Anthropic-specific request schema
type ChatRequest struct {
	Model         string        `json:"model"`
	Messages      []ChatMessage `json:"messages"`
	System        string        `json:"system,omitempty"`
	Temperature   float64       `json:"temperature,omitempty"`
	TopP          float64       `json:"top_p,omitempty"`
	TopK          int           `json:"top_k,omitempty"`
	MaxTokens     int           `json:"max_tokens,omitempty"`
	Stream        bool          `json:"stream,omitempty"`
	Metadata      *string       `json:"metadata,omitempty"`
	StopSequences []string      `json:"stop_sequences,omitempty"`
}

// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		Model:         cfg.Model,
		System:        cfg.DefaultParams.System,
		Temperature:   cfg.DefaultParams.Temperature,
		TopP:          cfg.DefaultParams.TopP,
		TopK:          cfg.DefaultParams.TopK,
		MaxTokens:     cfg.DefaultParams.MaxTokens,
		Metadata:      cfg.DefaultParams.Metadata,
		StopSequences: cfg.DefaultParams.StopSequences,
		Stream:        false, // unsupported right now
	}
}

func NewChatMessagesFromUnifiedRequest(request *schemas.ChatRequest) []ChatMessage {
	messages := make([]ChatMessage, 0, len(request.MessageHistory)+1)

	// Add items from messageHistory first and the new chat message last
	for _, message := range request.MessageHistory {
		messages = append(messages, ChatMessage{Role: message.Role, Content: message.Content})
	}

	messages = append(messages, ChatMessage{Role: request.Message.Role, Content: request.Message.Content})

	return messages
}

// Chat sends a chat request to the specified anthropic model.
//
//	Ref: https://docs.anthropic.com/claude/reference/messages_post
func (c *Client) Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)

	chatResponse, err := c.doChatRequest(ctx, chatRequest)

	if err != nil {
		return nil, err
	}

	return chatResponse, nil
}

func (c *Client) createChatRequestSchema(request *schemas.ChatRequest) *ChatRequest {
	// TODO: consider using objectpool to optimize memory allocation
	chatRequest := c.chatRequestTemplate // hoping to get a copy of the template
	chatRequest.Messages = NewChatMessagesFromUnifiedRequest(request)

	return chatRequest
}

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.ChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal anthropic chat request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create anthropic chat request: %w", err)
	}

	req.Header.Set("x-api-key", string(c.config.APIKey)) // must be in lower case
	req.Header.Set("anthropic-version")
	req.Header.Set("Content-Type", "application/json")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.tel.L().Debug(
		"Anthropic chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send anthropic chat request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.errMapper.Map(resp)
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.tel.L().Error("Failed to read anthropic chat response", zap.Error(err))
		return nil, err
	}

	// Parse the response JSON
	var anthropicResponse ChatCompletion

	err = json.Unmarshal(bodyBytes, &anthropicResponse)
	if err != nil {
		c.tel.L().Error("Failed to parse anthropic chat response", zap.Error(err))
		return nil, err
	}

	if len(anthropicResponse.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	completion := anthropicResponse.Content[0]
	usage := anthropicResponse.Usage

	// Map response to ChatResponse schema
	response := schemas.ChatResponse{
		ID:        anthropicResponse.ID,
		Created:   int(time.Now().UTC().Unix()), // not provided by anthropic
		Provider:  providerName,
		ModelName: anthropicResponse.Model,
		Cached:    false,
		ModelResponse: schemas.ModelResponse{
			SystemID: map[string]string{},
			Message: schemas.ChatMessage{
				Role:    completion.Type,
				Content: completion.Text,
			},
			TokenUsage: schemas.TokenUsage{
				PromptTokens:   usage.InputTokens,
				ResponseTokens: usage.OutputTokens,
				TotalTokens:    usage.InputTokens + usage.OutputTokens,
			},
		},
	}

	return &response, nil
}
