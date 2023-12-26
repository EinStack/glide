package providers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type GatewayConfig struct {
	Gateway PoolsConfig `yaml:"gateway" validate:"required"`
}
type PoolsConfig struct {
	Pools []Pool `yaml:"pools" validate:"required"`
}

type Pool struct {
	Name      string     `yaml:"pool" validate:"required"`
	Balancing string     `yaml:"balancing" validate:"required"`
	Providers []Provider `yaml:"providers" validate:"required"`
}

type Provider struct {
	Name          string                 `yaml:"name" validate:"required"`
	Model         string                 `yaml:"model"`
	APIKey        string                 `yaml:"api_key" validate:"required"`
	TimeoutMs     int                    `yaml:"timeout_ms,omitempty"`
	DefaultParams map[string]interface{} `yaml:"default_params,omitempty"`
}

type ProviderVars struct {
	Name        string `yaml:"name"`
	ChatBaseURL string `yaml:"chatBaseURL"`
}

type RequestBody struct {
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	MessageHistory []string `json:"messageHistory"`
}

// Variables

var HTTPClient = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 2,
	},
}

type UnifiedAPIData struct {
	Provider       string            `json:"provider"`
	Model          string            `json:"model"`
	APIKey         string            `json:"api_key"`
	Params         map[string]interface{} `json:"params"`
	Message        string            `json:"message"`
	MessageHistory []string          `json:"messageHistory"`
}


// Helper Functions

func FindPoolByName(pools []Pool, name string) (*Pool, error) {
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
// - providers: a slice of Provider, the list of providers to search in.
// - providerName: a string, the name of the provider to search for.
// - modelName: a string, the name of the model to search for.
//
// Returns:
// - *Provider: a pointer to the found provider.
// - error: an error indicating whether a provider was found or not.
func FindProviderByModel(providers []Provider, providerName string, modelName string) (*Provider, error) {
	for i := range providers {
		provider := &providers[i]
		if provider.Name == providerName && provider.Model == modelName {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("no provider found in config for model: %s", modelName)
}

func ReadProviderVars(filePath string) ([]ProviderVars, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute file path: %w", err)
	}

	// Validate that the absolute path is a file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("provided path is a directory, not a file")
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider vars file: %w", err)
	}

	var provVars []ProviderVars
	if err := yaml.Unmarshal(data, &provVars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider vars data: %w", err)
	}

	return provVars, nil
}

func GetDefaultBaseURL(provVars []ProviderVars, providerName string) (string, error) {
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

func ReadConfig(filePath string) (GatewayConfig, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return GatewayConfig{}, fmt.Errorf("failed to get absolute file path: %w", err)
	}

	// Validate that the absolute path is a file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return GatewayConfig{}, fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.IsDir() {
		return GatewayConfig{}, fmt.Errorf("provided path is a directory, not a file")
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		slog.Error("Error:", err)
		return GatewayConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config GatewayConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		slog.Error("Error:", err)
		return GatewayConfig{}, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return config, nil
}
