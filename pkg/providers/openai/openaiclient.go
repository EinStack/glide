// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling
// OpenAI package provide a set of functions to interact with the OpenAI API.
package openai

import (
	"errors"
	"fmt"

	"glide/pkg/providers"

	"github.com/go-playground/validator/v10"
)

const (
	providerName    = "openai"
	providerVarPath = "/Users/max/code/Glide/pkg/providers/providerVars.yaml"
	configPath      = "/Users/max/code/Glide/config.yaml"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// OpenAiClient creates a new client for the OpenAI API.
//
// Parameters:
// - poolName: The name of the pool to connect to.
// - modelName: The name of the model to use.
//
// Returns:
// - *Client: A pointer to the created client.
// - error: An error if the client creation failed.
func Client(poolName string, modelName string, payload []byte) (*ProviderClient, error) {
	provVars, err := providers.ReadProviderVars(providerVarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider vars: %w", err)
	}

	defaultBaseURL, err := providers.GetDefaultBaseURL(provVars, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get default base URL: %w", err)
	}

	config, err := providers.ReadConfig(configPath) // TODO: replace with struct built in router/pool
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Find the pool with the specified name from global config. This may not be necessary if details are passed directly in struct
	selectedPool, err := providers.FindPoolByName(config.Gateway.Pools, poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Find the OpenAI provider params in the selected pool with the specified model. This may not be necessary if details are passed directly in struct
	selectedProvider, err := providers.FindProviderByModel(selectedPool.Providers, providerName, modelName)
	if err != nil {
		return nil, fmt.Errorf("provider error: %w", err)
	}

	// Create a new client
	c := &ProviderClient{
		Provider:   *selectedProvider,
		PoolName:   poolName,
		BaseURL:    defaultBaseURL,
		Payload:    payload,
		HTTPClient: providers.HTTPClient,
	}

	v := validator.New()
	err = v.Struct(c)
	if err != nil {
		return nil, fmt.Errorf("failed to validate client: %w", err)
	}

	return c, nil
}
