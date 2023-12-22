// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling
// OpenAI package provide a set of functions to interact with the OpenAI API.
package openai

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

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
	requestBody      struct {
		Message []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		MessageHistory []string `json:"messageHistory"`
	}
)

var httpClient = &http.Client{
	Timeout: time.Second * 60,
	Transport: &http.Transport{
		MaxIdleConns:        90,
		MaxIdleConnsPerHost: 5,
	},
}

// Client is a client for the OpenAI API.
type ProviderClient struct {
	Provider   providers.Provider `validate:"required"`
	PoolName   string             `validate:"required"`
	baseURL    string             `validate:"required"`
	payload    []byte             `validate:"required"`
	httpClient *http.Client       `validate:"required"`
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
func Client(poolName string, modelName string, payload []byte) (*ProviderClient, error) {
	provVars, err := readProviderVars(providerVarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider vars: %w", err)
	}

	defaultBaseURL, err := getDefaultBaseURL(provVars, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get default base URL: %w", err)
	}

	config, err := readConfig(configPath) // TODO: replace with struct built in router/pool
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Find the pool with the specified name from global config. This may not be necessary if details are passed directly in struct
	selectedPool, err := findPoolByName(config.Gateway.Pools, poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Find the OpenAI provider params in the selected pool with the specified model. This may not be necessary if details are passed directly in struct
	selectedProvider, err := findProviderByModel(selectedPool.Providers, providerName, modelName)
	if err != nil {
		return nil, fmt.Errorf("provider error: %w", err)
	}

	// Create a new client
	c := &ProviderClient{
		Provider:   *selectedProvider,
		PoolName:   poolName,
		baseURL:    defaultBaseURL,
		payload:    payload,
		httpClient: httpClient,
	}

	v := validator.New()
	err = v.Struct(c)
	if err != nil {
		return nil, fmt.Errorf("failed to validate client: %w", err)
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

// findProviderByModel find provider params in the given config file by the specified provider name and model name.
//
// Parameters:
// - providers: a slice of providers.Provider, the list of providers to search in.
// - providerName: a string, the name of the provider to search for.
// - modelName: a string, the name of the model to search for.
//
// Returns:
// - *providers.Provider: a pointer to the found provider.
// - error: an error indicating whether a provider was found or not.
func findProviderByModel(providers []providers.Provider, providerName string, modelName string) (*providers.Provider, error) {
	for i := range providers {
		provider := &providers[i]
		if provider.Name == providerName && provider.Model == modelName {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("no provider found in config for model: %s", modelName)
}

func readProviderVars(filePath string) ([]providers.ProviderVars, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider vars file: %w", err)
	}

	var provVars []providers.ProviderVars
	if err := yaml.Unmarshal(data, &provVars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider vars data: %w", err)
	}

	return provVars, nil
}

func getDefaultBaseURL(provVars []providers.ProviderVars, providerName string) (string, error) {
	providerVarsMap := make(map[string]string)
	for _, providerVar := range provVars {
		providerVarsMap[providerVar.Name] = providerVar.ChatBaseURL
	}

	defaultBaseURL, ok := providerVarsMap[providerName]
	if !ok {
		return "", fmt.Errorf("default base URL not found for provider: %s", providerName)
	}

	return defaultBaseURL, nil
}

func readConfig(filePath string) (providers.GatewayConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Error:", err)
		return providers.GatewayConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config providers.GatewayConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		slog.Error("Error:", err)
		return providers.GatewayConfig{}, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return config, nil
}
