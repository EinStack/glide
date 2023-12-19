package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"log/slog"
	"gopkg.in/yaml.v2"
	"os"

	"Glide/pkg/providers"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"

)

const (
	defaultBaseURL              = "https://api.openai.com/v1"
	defaultFunctionCallBehavior = "auto"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var ErrEmptyResponse = errors.New("empty response")

type APIType string

const (
	APITypeOpenAI  APIType = "OPEN_AI"
	APITypeAzure   APIType = "AZURE"
	APITypeAzureAD APIType = "AZURE_AD"
)

// Client is a client for the OpenAI API.
type Client struct {
	apiKey        string
	Model        string
	baseURL      string
	organization string
	apiType      APIType
	httpClient   *client.Client

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion      string
	embeddingsModel string
}

// Option is an option for the OpenAI client.
type Option func(*Client) error


func Init(c *Client) (*client.Client, error) {
	// initializes the client

	// Read the YAML file
	data, err := os.ReadFile("path/to/file.yaml")
	if err != nil {
		slog.Error("Failed to read file: %v", err)
	}

	// Unmarshal the YAML data into your struct
	var config GatewayConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("Failed to unmarshal YAML: %v", err)
	}

}


// New returns a new OpenAI client.
func New(apiKey string, model string, baseURL string, organization string,
	apiType APIType, apiVersion string, httpClient *client.Client, embeddingsModel string,
	opts ...Option,
) (*Client, error) {
	c := &Client{
		apiKey:           apiKey,
		Model:           model,
		embeddingsModel: embeddingsModel,
		baseURL:         baseURL,
		organization:    organization,
		apiType:         apiType,
		apiVersion:      apiVersion,
		httpClient:      HttpClient(),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func HttpClient() *client.Client {

	c, err := client.NewClient()
	if err != nil {
		slog.Error(err.Error())
	}
	return c

}

// Completion is a completion.
type Completion struct {
	Text string `json:"text"`
}


// CreateChat creates chat request.
func (c *Client) CreateChat(ctx context.Context, r *ChatRequest) (*ChatResponse, error) {
	if r.Model == "" {
		if c.Model == "" {
			r.Model = defaultChatModel
		} else {
			r.Model = c.Model
		}
	}
	
	resp, err := c.createChat(ctx, r)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func IsAzure(apiType APIType) bool {
	return apiType == APITypeAzure || apiType == APITypeAzureAD
}

func (c *Client) setHeaders(req *protocol.Request) {
	req.Header.Set("Content-Type", "application/json")
	if c.apiType == APITypeOpenAI || c.apiType == APITypeAzureAD {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	} else {
		req.Header.Set("api-key", c.apiKey)
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