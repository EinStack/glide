package bedrock

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/EinStack/glide/pkg/telemetry"

	"github.com/EinStack/glide/pkg/providers/clients"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	providerName = "bedrock"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// Client is a client for accessing OpenAI API
type Client struct {
	baseURL             string
	bedrockClient       *bedrockruntime.Client
	chatURL             string
	chatRequestTemplate *ChatRequest
	config              *Config
	httpClient          *http.Client
	telemetry           *telemetry.Telemetry
}

// NewClient creates a new OpenAI client for the OpenAI API.
func NewClient(providerConfig *Config, clientConfig *clients.ClientConfig, tel *telemetry.Telemetry) (*Client, error) {
	chatURL, err := url.JoinPath(providerConfig.BaseURL, providerConfig.ChatEndpoint, providerConfig.Model, "/invoke")
	if err != nil {
		return nil, err
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), // Is this the right context?
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(providerConfig.AccessKey, providerConfig.SecretKey, "")),
		config.WithRegion(providerConfig.AWSRegion),
	)

	bedrockClient := bedrockruntime.NewFromConfig(cfg)

	c := &Client{
		baseURL:             providerConfig.BaseURL,
		bedrockClient:       bedrockClient,
		chatURL:             chatURL,
		config:              providerConfig,
		chatRequestTemplate: NewChatRequestFromConfig(providerConfig),
		httpClient: &http.Client{
			Timeout: time.Duration(*clientConfig.Timeout),
			// TODO: use values from the config
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
