package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"net/http"

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	defaultChatModel = "gpt-3.5-turbo"
)

// ChatRequest is a request to complete a chat completion..
type ChatRequest struct {
	Model            string           `json:"model" validate:"required,lowercase"`
	Messages         []*ChatMessage   `json:"messages" validate:"required"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty" validate:"omitempty,gte=0,lte=1"`
	MaxTokens        int              `json:"max_tokens,omitempty" validate:"omitempty,gte=0"`
	N                int              `json:"n,omitempty" validate:"omitempty,gte=1"`
	StopWords        []string         `json:"stop,omitempty"`
	Stream           bool             `json:"stream,omitempty" validate:"omitempty, boolean"`
	FrequencyPenalty int              `json:"frequency_penalty,omitempty"`
	PresencePenalty  int              `json:"presence_penalty,omitempty"`
	LogitBias        *map[int]float64 `json:"logit_bias,omitempty" validate:"omitempty"`
	User             interface{}      `json:"user,omitempty"`
	Seed             interface{}      `json:"seed,omitempty" validate:"omitempty,gte=0"`
	Tools            []string         `json:"tools,omitempty"`
	ToolChoice       interface{}      `json:"tool_choice,omitempty"`
	ResponseFormat   interface{}      `json:"response_format,omitempty"`

	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct {
	// The role of the author of this message. One of system, user, or assistant.
	Role string `json:"role"`
	// The content of the message.
	Content string `json:"content"`
	// The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,
	// with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`

}

// ChatChoice is a choice in a chat response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatUsage is the usage of a chat completion request.
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (c *Client) CreateChatRequest(message []byte) *ChatRequest {


	err := json.Unmarshal(message, &requestBody)
	if err != nil {
		slog.Error("Error:", err)
		return nil
	}

	var messages []*ChatMessage
	for _, msg := range requestBody.Message {
		chatMsg := &ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
		if msg.Role == "user" {
			chatMsg.Content += " " + strings.Join(requestBody.MessageHistory, " ")
		}
		messages = append(messages, chatMsg)
	}

	// iterate through self.Provider.DefaultParams and add them to the request otherwise leave the default value
	
	chatRequest := &ChatRequest{
		Model:            c.Provider.Model,
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

	// Use reflection to dynamically assign default parameter values
	defaultParams := c.Provider.DefaultParams
	v := reflect.ValueOf(chatRequest).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		defaultValue, ok := defaultParams[fieldName]
		if ok && defaultValue != nil {
			fieldValue := v.FieldByName(fieldName)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(defaultValue))
			}
		}
	}

	return chatRequest
}


// ChatResponse is a response to a chat request.
type ChatResponse struct {
	ID      string        `json:"id,omitempty"`
	Created float64       `json:"created,omitempty"`
	Choices []*ChatChoice `json:"choices,omitempty"`
	Model   string        `json:"model,omitempty"`
	Object  string        `json:"object,omitempty"`
	Usage   struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

// StreamedChatResponsePayload is a chunk from the stream.
type StreamedChatResponsePayload struct {
	ID      string  `json:"id,omitempty"`
	Created float64 `json:"created,omitempty"`
	Model   string  `json:"model,omitempty"`
	Object  string  `json:"object,omitempty"`
	Choices []struct {
		Index float64 `json:"index,omitempty"`
		Delta struct {
			Role         string        `json:"role,omitempty"`
			Content      string        `json:"content,omitempty"`
		} `json:"delta,omitempty"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
}

// CreateChatResponse creates chat Response.
func (c *Client) CreateChatResponse(ctx context.Context, r *ChatRequest) (*ChatResponse, error) {
	if r.Model == "" {
		if c.Provider.Model == "" {
			r.Model = defaultChatModel
		} else {
			r.Model = c.Provider.Model
		}
	}
	
	resp, err := c.createChat(ctx, r)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}


func (c *Client) createChat(ctx context.Context, payload *ChatRequest) (*ChatResponse, error) {
	if payload.StreamingFunc != nil {
		payload.Stream = true
	}
	// Build request payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Build request
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}

	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI(c.buildURL("/chat/completions", c.Provider.Model))
	req.SetBody(payloadBytes)

	c.setHeaders(req) // sets additional headers

	// Send request
	err = c.httpClient.Do(ctx, req, res) //*client.Client
	if err != nil {
		return nil, err
	}

	defer res.ConnectionClose() // replaced r.Body.Close()

	if res.StatusCode() != http.StatusOK {
		msg := fmt.Sprintf("API returned unexpected status code: %d", res.StatusCode())

		return nil, fmt.Errorf("%s: %s", msg, err.Error()) // nolint:goerr113
	}
	if payload.StreamingFunc != nil {
		return parseStreamingChatResponse(ctx, res, payload)
	}
	// Parse response
	var response ChatResponse
	return &response, json.NewDecoder(bytes.NewReader(res.Body())).Decode(&response)
}

func parseStreamingChatResponse(ctx context.Context, r *protocol.Response, payload *ChatRequest) (*ChatResponse, error) {
	scanner := bufio.NewScanner(bytes.NewReader(r.Body()))
	responseChan := make(chan StreamedChatResponsePayload)
	go func() {
		defer close(responseChan)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data:") {
				slog.Warn("unexpected line:" + line)
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var streamPayload StreamedChatResponsePayload
			err := json.NewDecoder(bytes.NewReader([]byte(data))).Decode(&streamPayload)
			if err != nil {
				slog.Error("failed to decode stream payload: %v", err)
			}
			responseChan <- streamPayload
		}
		if err := scanner.Err(); err != nil {
			slog.Error("issue scanning response:", err)
		}
	}()

	// Parse response
	response := ChatResponse{
		Choices: []*ChatChoice{
			{},
		},
	}

	for streamResponse := range responseChan {
		if len(streamResponse.Choices) == 0 {
			continue
		}
		chunk := []byte(streamResponse.Choices[0].Delta.Content)
		response.Choices[0].Message.Content += streamResponse.Choices[0].Delta.Content
		response.Choices[0].FinishReason = streamResponse.Choices[0].FinishReason

		if payload.StreamingFunc != nil {
			err := payload.StreamingFunc(ctx, chunk)
			if err != nil {
				return nil, fmt.Errorf("streaming func returned an error: %w", err)
			}
		}
	}
	return &response, nil
}