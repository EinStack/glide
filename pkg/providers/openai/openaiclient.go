package openai

import (
	"errors"
	"fmt"
	"strings"
	"log/slog"
	"gopkg.in/yaml.v2"
	"os"

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
		slog.Error("pool not found")
		return nil, fmt.Errorf("pool not found: %s", poolName)
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
		return nil, fmt.Errorf("provider for model '%s' not found in pool '%s'", modelName, poolName)
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
