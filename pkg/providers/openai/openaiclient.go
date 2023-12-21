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
	providerName = "openai"
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
	
	// Read the YAML file
	data, err := os.ReadFile("/Users/max/code/Glide/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	slog.Info("config loaded")

	// Unmarshal the YAML data into your struct
	var config providers.GatewayConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	// Find the pool with the specified name
	selectedPool, err := findPoolByName(config.Gateway.Pools, poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Find the OpenAI provider params in the selected pool with the specified model
	selectedProvider, err := findProviderByModel(selectedPool.Providers, providerName, modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to find provider: %w", err)
	}

	// Create a new client
	c := &Client{
		Provider:   *selectedProvider,
		PoolName:   poolName,
		baseURL:    "", // Set the appropriate base URL
		payload:    payload,
		httpClient: HTTPClient(),
	}

	return c, nil
}

func findPoolByName(pools []providers.Pool, name string) (*providers.Pool, error) {
	for i := range pools {
		pool := &pools[i]
		if pool.Name == name {
			return pool, nil
		}
	}

	return nil, fmt.Errorf("pool not found: %s", name)
}

func findProviderByModel(providers []providers.Provider, providerName string, modelName string) (*providers.Provider, error) {
	for i := range providers {
		provider := &providers[i]
		if provider.Provider == providerName && provider.Model == modelName {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("provider not found: %s", modelName)
}
