package azureopenai

import (
	"errors"
	"fmt"
	"net/http"

	"glide/pkg/providers/clients"
	"glide/pkg/telemetry"
)


const (
	providerName = "azureopenai"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// Client is a client for accessing Azure OpenAI API
type Client struct {
	baseURL             string // The name of your Azure OpenAI Resource (e.g https://glide-test.openai.azure.com/)
	chatURL             string
	chatRequestTemplate *ChatRequest
	config              *Config
	httpClient          *http.Client
	telemetry           *telemetry.Telemetry
}

// NewClient creates a new Azure OpenAI client for the OpenAI API.
func NewClient(providerConfig *Config, clientConfig *clients.ClientConfig, tel *telemetry.Telemetry)  (*Client, error) {
	chatURL := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", providerConfig.BaseURL, providerConfig.Model, providerConfig.APIVersion)

	fmt.Println("chatURL", chatURL)

	c := &Client{
		baseURL:             providerConfig.BaseURL,
		chatURL:             chatURL,
		config:              providerConfig,
		chatRequestTemplate: NewChatRequestFromConfig(providerConfig),
		httpClient: &http.Client{
			// TODO: use values from the config
			Timeout: *clientConfig.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 2,
			},
		},
		telemetry: tel,
	}

	return c, nil
}

func (c *Client) Provider() string {
	return providerName
}
