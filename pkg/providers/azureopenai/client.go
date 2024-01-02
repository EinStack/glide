package azureopenai

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"glide/pkg/telemetry"
)

// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling

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
func NewClient(cfg *Config, tel *telemetry.Telemetry) (*Client, error) {
	chatURL := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", cfg.BaseURL, cfg.Model, cfg.APIVersion)

	fmt.Println("chatURL", chatURL)

	c := &Client{
		baseURL:             cfg.BaseURL,
		chatURL:             chatURL,
		config:              cfg,
		chatRequestTemplate: NewChatRequestFromConfig(cfg),
		httpClient: &http.Client{
			// TODO: use values from the config
			Timeout: time.Second * 30,
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
