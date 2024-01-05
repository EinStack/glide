package octoml

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"glide/pkg/telemetry"
)

// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling

const (
	providerName = "octoml"
)

// ErrEmptyResponse is returned when the OctoML API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// Client is a client for accessing OctoML API
type Client struct {
	baseURL             string
	chatURL             string
	chatRequestTemplate *ChatRequest
	config              *Config
	httpClient          *http.Client
	telemetry           *telemetry.Telemetry
}

// NewClient creates a new OctoML client for the OctoML API.
func NewClient(cfg *Config, tel *telemetry.Telemetry) (*Client, error) {
	chatURL, err := url.JoinPath(cfg.BaseURL, cfg.ChatEndpoint)
	if err != nil {
		return nil, err
	}

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
