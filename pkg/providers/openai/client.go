package openai

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
	providerName = "openai"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// Client is a client for accessing OpenAI API
type Client struct {
	baseURL    string
	chatURL    string
	config     *Config
	httpClient *http.Client
	telemetry  *telemetry.Telemetry
}

// NewClient creates a new OpenAI client for the OpenAI API.
func NewClient(cfg *Config, tel *telemetry.Telemetry) (*Client, error) {
	// Create a new client
	c := &Client{
		baseURL: cfg.BaseURL,
		chatURL: fmt.Sprintf("%s%s", cfg.BaseURL, cfg.ChatEndpoint),
		config:  cfg,
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
