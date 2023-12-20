// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling

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
	baseURL      string
	organization string
	apiType      APIType
	httpClient   *client.Client

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion string
}

func (c *Client) Run(poolName string, modelName string, payload []byte) (*ChatResponse, error) {
	c, err := c.NewClient(poolName, modelName)
	if err != nil {
		slog.Error("Error:" + err.Error())
		return nil, err
	}

	// Create a new chat request

	slog.Info("creating chat request")

	chatRequest := c.CreateChatRequest(payload)

	slog.Info("chat request created")

	// Send the chat request

	slog.Info("sending chat request")

	resp, err := c.CreateChatResponse(context.Background(), chatRequest)

	return resp, err
}

func (c *Client) NewClient(poolName string, modelName string) (*Client, error) {
	// Returns a []*Client of OpenAI
	// modelName is determined by the model pool
	// poolName is determined by the route the request came from

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
		organization: defaultOrganization,
		apiType:      APITypeOpenAI,
		httpClient:   HTTPClient(),
	}

	return client, nil
}

func HTTPClient() *client.Client {
	c, err := client.NewClient()
	if err != nil {
		slog.Error(err.Error())
	}
	return c
}
