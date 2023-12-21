// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling
// OpenAI package provide a set of functions to interact with the OpenAI API.
package openai

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"

	"glide/pkg/providers"

	"github.com/cloudwego/hertz/pkg/app/client"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
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

// Client is a client for the OpenAI API.
type Client struct {
	Provider   providers.Provider `validate:"required"`
	PoolName   string             `validate:"required"`
	baseURL    string             `validate:"required"`
	payload    []byte             `validate:"required"`
	httpClient *client.Client     `validate:"required"`
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
		return nil, err
	}

	slog.Info("config loaded")

	// Unmarshal the YAML data into your struct
	var config providers.GatewayConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

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
		slog.Error("double check the config.yaml for errors")
		return nil, fmt.Errorf("provider for model '%s' not found in pool '%s'", modelName, poolName)
	}

	// Create clients for each OpenAI provider
	client := &Client{
		Provider:   *selectedProvider,
		PoolName:   poolName,
		baseURL:    defaultBaseURL,
		payload:    payload,
		httpClient: HTTPClient(),
	}

	return client, nil
}
