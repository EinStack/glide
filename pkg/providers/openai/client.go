package openai

import (
	"github.com/EinStack/glide/pkg/providers"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/EinStack/glide/pkg/telemetry"

	"github.com/EinStack/glide/pkg/providers/clients"
)

const (
	providerName = "openai"
)

// Client is a client for accessing OpenAI API
type Client struct {
	baseURL             string
	chatURL             string
	chatRequestTemplate *ChatRequest
	errMapper           *ErrorMapper
	finishReasonMapper  *FinishReasonMapper
	config              *Config
	httpClient          *http.Client
	tel                 *telemetry.Telemetry
	logger              *zap.Logger
}

// ensure interfaces are implemented at compilation
var _ providers.LangProvider = (*Client)(nil)

// NewClient creates a new OpenAI client for the OpenAI API.
func NewClient(providerConfig *Config, clientConfig *clients.ClientConfig, tel *telemetry.Telemetry) (*Client, error) {
	chatURL, err := url.JoinPath(providerConfig.BaseURL, providerConfig.ChatEndpoint)
	if err != nil {
		return nil, err
	}

	logger := tel.L().With(
		zap.String("provider", providerName),
	)

	c := &Client{
		baseURL:             providerConfig.BaseURL,
		chatURL:             chatURL,
		config:              providerConfig,
		chatRequestTemplate: NewChatRequestFromConfig(providerConfig),
		finishReasonMapper:  NewFinishReasonMapper(tel),
		errMapper:           NewErrorMapper(tel),
		httpClient: &http.Client{
			Timeout: time.Duration(*clientConfig.Timeout),
			Transport: &http.Transport{
				MaxIdleConns:        *clientConfig.MaxIdleConns,
				MaxIdleConnsPerHost: *clientConfig.MaxIdleConnsPerHost,
			},
		},
		tel:    tel,
		logger: logger,
	}

	return c, nil
}

func (c *Client) Provider() string {
	return providerName
}

func (c *Client) ModelName() string {
	return c.config.ModelName
}
