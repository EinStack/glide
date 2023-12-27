package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"glide/pkg/providers"
)

const (
	defaultChatModel = "gpt-3.5-turbo"
	defaultEndpoint  = "/chat/completions"
)

// Client is a client for the OpenAI API.
type ProviderClient struct {
	BaseURL    string       `json:"baseURL"`
	HTTPClient *http.Client `json:"httpClient"`
}

// ChatRequest is a request to complete a chat completion..
type ChatRequest struct {
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

// Chat sends a chat request to the specified OpenAI model.
//
// Parameters:
// - payload: The user payload for the chat request.
// Returns:
// - *ChatResponse: a pointer to a ChatResponse
// - error: An error if the request failed.
func (c *ProviderClient) Chat(u *providers.UnifiedAPIData) (*ChatResponse, error) {
	// Create a new chat request

	slog.Info("creating chat request")

	chatRequest := CreateChatRequest(u)

	slog.Info("chat request created")

	// Send the chat request

	slog.Info("sending chat request")

	resp, err := CreateChatResponse(context.Background(), chatRequest, u)

	return resp, err
}

func CreateChatRequest(u *providers.UnifiedAPIData) *ChatRequest {
	slog.Info("creating chatRequest from payload")

	var messages []map[string]string

	// Add items from messageHistory first
	messages = append(messages, u.MessageHistory...)

	// Add msg variable last
	messages = append(messages, u.Message)

	// iterate throughunifiedData.Params and add them to the request otherwise leave the default value

	chatRequest := &ChatRequest{
		Model:            u.Model,
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
	// TODO: refactor to avoid reflection(?)
	defaultParams := u.Params

	chatRequestValue := reflect.ValueOf(chatRequest).Elem()
	chatRequestType := chatRequestValue.Type()

	for i := 0; i < chatRequestValue.NumField(); i++ {
		jsonTags := strings.Split(chatRequestType.Field(i).Tag.Get("json"), ",")
		jsonTag := jsonTags[0]
		if value, ok := defaultParams[jsonTag]; ok {
			fieldValue := chatRequestValue.Field(i)
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return chatRequest
}

// CreateChatResponse creates chat Response.
func CreateChatResponse(ctx context.Context, r *ChatRequest, u *providers.UnifiedAPIData) (*ChatResponse, error) {
	_ = ctx // keep this for future use

	resp, err := createChatHTTP(r, u)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func createChatHTTP(payload *ChatRequest, u *providers.UnifiedAPIData) (*ChatResponse, error) {
	slog.Info("running createChatHttp")

	if payload.StreamingFunc != nil {
		payload.Stream = true
	}
	// Build request payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Build request
	if defaultBaseURL == "" {
		slog.Error("baseURL not set")
		return nil, errors.New("baseURL not set")
	}

	reqBody := bytes.NewBuffer(payloadBytes)
	req, err := http.NewRequest("POST", buildURL(defaultEndpoint), reqBody)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	fmt.Println(reqBody.String())

	req.Header.Set("Authorization", "Bearer "+u.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := providers.HTTPClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	slog.Info(fmt.Sprintf("%d", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(err.Error())
		}
		bodyString := string(bodyBytes)
		slog.Warn(bodyString)
	}

	// Parse response
	var response ChatResponse
	return &response, json.NewDecoder(resp.Body).Decode(&response)
}

func buildURL(suffix string) string {
	slog.Info("request url: " + fmt.Sprintf("%s%s", defaultBaseURL, suffix))

	// open ai implement:
	return fmt.Sprintf("%s%s", defaultBaseURL, suffix)
}
