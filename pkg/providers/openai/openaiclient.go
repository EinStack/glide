package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"log/slog"
	"gopkg.in/yaml.v2"
	"os"
	"json"
	"reflect"

	//"Glide/pkg/providers"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"

)

const (
	defaultBaseURL              = "https://api.openai.com/v1"
	defaultOrganization          = ""
	defaultFunctionCallBehavior = "auto"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
	requestBody struct {
		Message        []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		MessageHistory []string `json:"messageHistory"`
	}
)
type APIType string

const (
	APITypeOpenAI  APIType = "OPEN_AI"
	APITypeAzure   APIType = "AZURE"
	APITypeAzureAD APIType = "AZURE_AD"
)

type GatewayConfig struct {
	Pools []Pool `yaml:"pools"`
}

type Pool struct {
	Name       string      `yaml:"name"`
	Balancing  string      `yaml:"balancing"`
	Providers  []Provider  `yaml:"providers"`
}

type Provider struct {
	Name           string                 `yaml:"name"`
	Provider       string                 `yaml:"provider"`
	Model          string                 `yaml:"model"`
	ApiKey         string                 `yaml:"api_key"`
	TimeoutMs      int                    `yaml:"timeout_ms,omitempty"`
	DefaultParams  map[string]interface{} `yaml:"default_params,omitempty"`
}



// Client is a client for the OpenAI API.
type Client struct {
	Provider     Provider
	baseURL      string
	organization string
	apiType      APIType
	httpClient   *client.Client

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion      string
	embeddingsModel string
}

// Option is an option for the OpenAI client.
type Option func(*Client) error


func (c *Client) Init(poolName string, modelName string, providerName string) (*Client, error) {
	// Returns a []*Client of OpenAI

	// Read the YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		slog.Error("Failed to read file: %v", err)
	}

	// Unmarshal the YAML data into your struct
	var config GatewayConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("Failed to unmarshal YAML: %v", err)
	}

	// Find the pool with the specified name
	var selectedPool *Pool
	for _, pool := range config.Pools {
		if pool.Name == poolName {
			selectedPool = &pool
			break
		}
	}

	// Check if the pool was found
	if selectedPool == nil {
		slog.Error("pool '%s' not found", poolName)
	}

	// Find the OpenAI provider in the selected pool with the specified model
    var selectedProvider *Provider
    for _, provider := range selectedPool.Providers {
        if provider.Name == providerName && provider.Model == modelName {
            selectedProvider = &provider
            break
        }
    }

    // Check if the provider was found
    if selectedProvider == nil {
        slog.Error("provider for model '%s' not found in pool '%s'", modelName, poolName)
    }

	// Create clients for each OpenAI provider
	client := &Client{
			Provider:     *selectedProvider,
			organization: defaultOrganization,
			apiType:      APITypeOpenAI,
			httpClient:   HttpClient(),
		}

	return client, nil
	
}


func HttpClient() *client.Client {

	c, err := client.NewClient()
	if err != nil {
		slog.Error(err.Error())
	}
	return c

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

func IsAzure(apiType APIType) bool {
	return apiType == APITypeAzure || apiType == APITypeAzureAD
}

func (c *Client) setHeaders(req *protocol.Request) {
	req.Header.Set("Content-Type", "application/json")
	if c.apiType == APITypeOpenAI || c.apiType == APITypeAzureAD {
		req.Header.Set("Authorization", "Bearer "+c.Provider.ApiKey)
	} else {
		req.Header.Set("api-key", c.Provider.ApiKey)
	}
	if c.organization != "" {
		req.Header.Set("OpenAI-Organization", c.organization)
	}
}

func (c *Client) buildURL(suffix string, model string) string {
	if IsAzure(c.apiType) {
		return c.buildAzureURL(suffix, model)
	}

	// open ai implement:
	return fmt.Sprintf("%s%s", c.baseURL, suffix)
}

func (c *Client) buildAzureURL(suffix string, model string) string {
	baseURL := c.baseURL
	baseURL = strings.TrimRight(baseURL, "/")

	// azure example url:
	// /openai/deployments/{model}/chat/completions?api-version={api_version}
	return fmt.Sprintf("%s/openai/deployments/%s%s?api-version=%s",
		baseURL, model, suffix, c.apiVersion,
	)
}

func main() {

	c := &Client{}

	c, err := c.Init("pool1", "gpt-3.5-turbo", "openai")
	if err != nil {
		// Handle the error
	}


}	