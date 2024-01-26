package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"bufio"

	"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"
	"go.uber.org/zap"
)

type ErrorAccumulator interface {
	Write(p []byte) error
	Bytes() []byte
}

// ChatRequest is an OpenAI-specific request schema
type ChatRequest struct {
	Model            string           `json:"model"`
	Messages         []schemas.ChatMessage    `json:"messages"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
	MaxTokens        int              `json:"max_tokens,omitempty"`
	N                int              `json:"n,omitempty"`
	StopWords        []string         `json:"stop,omitempty"`
	Stream           bool             `json:"stream,omitempty"`
	FrequencyPenalty int              `json:"frequency_penalty,omitempty"`
	PresencePenalty  int              `json:"presence_penalty,omitempty"`
	LogitBias        *map[int]float64 `json:"logit_bias,omitempty"`
	User             *string          `json:"user,omitempty"`
	Seed             *int             `json:"seed,omitempty"`
	Tools            []string         `json:"tools,omitempty"`
	ToolChoice       interface{}      `json:"tool_choice,omitempty"`
	ResponseFormat   interface{}      `json:"response_format,omitempty"`
}

// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		Model:            cfg.Model,
		Temperature:      cfg.DefaultParams.Temperature,
		TopP:             cfg.DefaultParams.TopP,
		MaxTokens:        cfg.DefaultParams.MaxTokens,
		N:                cfg.DefaultParams.N,
		StopWords:        cfg.DefaultParams.StopWords,
		Stream:           cfg.DefaultParams.Stream,
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

func NewChatMessagesFromUnifiedRequest(request *schemas.UnifiedChatRequest) []schemas.ChatMessage {
	messages := make([]schemas.ChatMessage, 0, len(request.MessageHistory)+1)

	// Add items from messageHistory first and the new chat message last
	for _, message := range request.MessageHistory {
		messages = append(messages, schemas.ChatMessage{Role: message.Role, Content: message.Content})
	}

	messages = append(messages, schemas.ChatMessage{Role: request.Message.Role, Content: request.Message.Content})

	return messages
}

// Chat sends a chat request to the specified OpenAI model.
func (c *Client) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)

	if chatRequest.Stream {
		// Create channels for receiving responses and errors
	   responseChannel := make(chan *schemas.UnifiedChatResponse, 5)
	   errChannel := make(chan error, 5)

	   fmt.Println("Starting streaming chat request")

	   defer close(responseChannel)
       defer close(errChannel)

	   c.doStreamingChatRequest(ctx, chatRequest, responseChannel, errChannel)

	   // Handle streaming responses and errors
	   for {
		   select {
		   case chatResponse := <-responseChannel:
			   // Process the streaming response
			   fmt.Println("Received response:", chatResponse)
			   //return chatResponse, nil
		   case err := <-errChannel:
			   // Handle the error
			   fmt.Println("Received error:", err)
			default:
				fmt.Println("DEFAULT")
		   }
	   }
   }

	chatResponse, err := c.doChatRequest(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.ModelResponse.Message.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

func (c *Client) createChatRequestSchema(request *schemas.UnifiedChatRequest) *ChatRequest {
	// TODO: consider using objectpool to optimize memory allocation
	chatRequest := c.chatRequestTemplate // hoping to get a copy of the template
	chatRequest.Messages = NewChatMessagesFromUnifiedRequest(request)

	return chatRequest
}

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.UnifiedChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal openai chat request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
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

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.telemetry.Logger.Error("failed to read openai chat response", zap.Error(err))
		}

		c.telemetry.Logger.Error(
			"openai chat request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.Any("headers", resp.Header),
		)

		if resp.StatusCode == http.StatusTooManyRequests {
			// Read the value of the "Retry-After" header to get the cooldown delay
			retryAfter := resp.Header.Get("Retry-After")

			// Parse the value to get the duration
			cooldownDelay, err := time.ParseDuration(retryAfter)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cooldown delay from headers: %w", err)
			}

			return nil, clients.NewRateLimitError(&cooldownDelay)
		}

		// Server & client errors result in the same error to keep gateway resilient
		return nil, clients.ErrProviderUnavailable
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.telemetry.Logger.Error("failed to read openai chat response", zap.Error(err))
		return nil, err
	}

	//fmt.Println(string(bodyBytes))

	var openAICompletion schemas.OpenAIChatCompletion

	err = json.Unmarshal(bodyBytes, &openAICompletion)
	if err != nil {
		c.telemetry.Logger.Error("failed to parse openai chat response", zap.Error(err))
		return nil, err
	}

	// Map response to UnifiedChatResponse schema
	response := schemas.UnifiedChatResponse{
		ID:       openAICompletion.ID,
		Created:  openAICompletion.Created,
		Provider: providerName,
		Model:    openAICompletion.Model,
		Cached:   false,
		ModelResponse: schemas.ProviderResponse{
			ResponseID: map[string]string{
				"system_fingerprint": openAICompletion.SystemFingerprint,
			},
			Message: schemas.ChatMessage{
				Role:    openAICompletion.Choices[0].Message.Role,
				Content: openAICompletion.Choices[0].Message.Content,
				Name:    "",
			},
			TokenCount: schemas.TokenCount{
				PromptTokens:   openAICompletion.Usage.PromptTokens,
				ResponseTokens: openAICompletion.Usage.CompletionTokens,
				TotalTokens:    openAICompletion.Usage.TotalTokens,
			},
		},
	}

	return &response, nil
}

func (c *Client) doStreamingChatRequest(ctx context.Context, payload *ChatRequest, responseChannel chan *schemas.UnifiedChatResponse, errChannel chan error) {
   defer close(responseChannel)
   defer close(errChannel)
	
	// build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		errChannel <- fmt.Errorf("unable to marshal openai chat request payload: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		errChannel <- fmt.Errorf("unable to create openai chat request: %w", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+string(c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	c.telemetry.Logger.Debug(
		"openai chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		errChannel <- fmt.Errorf("failed to send openai chat request: %w", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			errChannel <- fmt.Errorf("failed to read openai chat response: %w", err)
			c.telemetry.Logger.Error("failed to read openai chat response", zap.Error(err))
		}

		c.telemetry.Logger.Error(
			"openai chat request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.Any("headers", resp.Header),
		)

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			cooldownDelay, err := time.ParseDuration(retryAfter)
			if err != nil {
				errChannel <- fmt.Errorf("failed to parse cooldown delay from headers: %w", err)
				return
			}

			errChannel <- clients.NewRateLimitError(&cooldownDelay)
			return
		}

		// Server & client errors result in the same error to keep gateway resilient
		errChannel <- clients.ErrProviderUnavailable
		return
	}

	var (
		headerData  = []byte("data: ")
		errorPrefix = []byte(`data: {"error":`)
		hasErrorPrefix     bool
	)

	// Read the response body into a byte slice
	//bodyBytes, _ := io.ReadAll(resp.Body)

	//fmt.Println(string(bodyBytes))

	// Create a scanner to read from the response body
	// TODO: NEed to figure out why it only returns the first line
	reader := bufio.NewReader(resp.Body)
	
	for {
		r, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// Handle EOF error
				errChannel <- fmt.Errorf("EOF: %w", err)
			}
			errChannel <- fmt.Errorf("failed to read openai chat response: %w", err)
			c.telemetry.Logger.Error("failed to read openai chat response", zap.Error(err))
			continue
		}

		// Apply the processing steps to each chunk
		noSpaceLine := bytes.TrimSpace(r)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)

		var openAICompletion schemas.OpenAIChatStreamCompletion

		decoder := json.NewDecoder(bytes.NewReader(noPrefixLine))
		if err := decoder.Decode(&openAICompletion); err != nil {
			if err == io.EOF {
				// Handle EOF error
				return
			}
			errChannel <- fmt.Errorf("failed to parse openai chat response: %w", err)
			c.telemetry.Logger.Error("failed to parse openai chat response", zap.Error(err))
			continue
		}

		response := schemas.UnifiedChatResponse{
			ID:       openAICompletion.ID,
			Created:  openAICompletion.Created,
			Provider: providerName,
			Model:    openAICompletion.Model,
			Cached:   false,
			ModelResponse: schemas.ProviderResponse{
				ResponseID: map[string]string{
					"system_fingerprint": openAICompletion.SystemFingerprint,
				},
				Message: schemas.ChatMessage{
					Role:    openAICompletion.StreamChoice[0].Delta.Role,
					Content: openAICompletion.StreamChoice[0].Delta.Content,
					Name:    "",
				},
				TokenCount: schemas.TokenCount{
					PromptTokens:   openAICompletion.Usage.PromptTokens,
					ResponseTokens: openAICompletion.Usage.CompletionTokens,
					TotalTokens:    openAICompletion.Usage.TotalTokens,
				},
			},
		}
		responseChannel <- &response
	}
	// ... (remaining code)
}

