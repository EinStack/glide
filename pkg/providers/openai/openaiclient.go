// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling
// OpenAI package provide a set of functions to interact with the OpenAI API.
package openai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"

	"glide/pkg/providers"

	"github.com/cloudwego/hertz/pkg/app/client"
)

const (
	defaultBaseURL      = "https://api.openai.com/v1"
	defaultOrganization = ""
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
	requestBody      struct {
		Message []struct {
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

// Client is a client for the OpenAI API.
type Client struct {
	Provider     providers.Provider
	PoolName     string
	baseURL      string
	payload      []byte
	organization string
	apiType      APIType
	httpClient   *client.Client

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion string
}

// OpenAiClient creates a new client for the OpenAI API.
//
// Parameters:
// - poolName: The name of the pool to connect to.
// - modelName: The name of the model to use.
//
// Returns:
// - *Client: A pointer to the created client.
// - error: An error if the client creation failed.
func OpenAiClient(poolName string, modelName string, payload []byte) (*Client, error) {
	providerName := "openai"

	// Read the YAML file
	data, err := os.ReadFile("/Users/max/code/Glide/config.yaml")
	if err != nil {
		slog.Error("Failed to read file: %v", err)
		return nil, err
	}

	slog.Info("config found")

	// Unmarshal the YAML data into your struct
	var config providers.GatewayConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("Failed to unmarshal YAML: %v", err)
	}

	fmt.Println(config)

	// Find the pool with the specified name
	var selectedPool *providers.Pool
	for i := range config.Gateway.Pools {
		pool := &config.Gateway.Pools[i]
		if pool.Name == poolName {
			selectedPool = pool
			break
		}
	}

	// Check if the pool was found
	if selectedPool == nil {
		slog.Error("pool not found")
		return nil, fmt.Errorf("pool not found: %s", poolName)
	}

	// Find the OpenAI provider in the selected pool with the specified model
	var selectedProvider *providers.Provider
	for i := range selectedPool.Providers {
		provider := &selectedPool.Providers[i]
		if provider.Provider == providerName && provider.Model == modelName {
			selectedProvider = provider
			break
		}
	}

	// Check if the provider was found
	if selectedProvider == nil {
		slog.Error("provider for model '%s' not found in pool '%s'", modelName, poolName)
		return nil, fmt.Errorf("provider for model '%s' not found in pool '%s'", modelName, poolName)
	}

	// Create clients for each OpenAI provider
	client := &Client{
		Provider:     *selectedProvider,
		PoolName:     poolName,
		baseURL:      defaultBaseURL,
		payload:      payload,
		organization: defaultOrganization,
		httpClient:   HTTPClient(),
	}

	return client, nil
}

// Chat sends a chat request to the specified OpenAI model.
//
// Parameters:
// - payload: The user payload for the chat request.
// Returns:
// - *ChatResponse: a pointer to a ChatResponse
// - error: An error if the request failed.
func (c *Client) Chat() (*ChatResponse, error) {
	// Create a new chat request

	slog.Info("creating chat request")

	chatRequest := c.CreateChatRequest(c.payload)

	slog.Info("chat request created")

	// Send the chat request

	slog.Info("sending chat request")

	resp, err := c.CreateChatResponse(context.Background(), chatRequest)

	return resp, err
}

// HTTPClient returns a new Hertz HTTP client.
//
// It creates a new client using the client.NewClient() function and returns the client.
// If an error occurs during the creation of the client, it logs the error using slog.Error().
// The function returns the created client or nil if an error occurred.
func HTTPClient() *client.Client {
	c, err := client.NewClient()
	if err != nil {
		slog.Error(err.Error())
	}
	return c
}
