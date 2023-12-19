package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"log/slog"

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
	token        string
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

// Doer performs a HTTP request.
type Doer interface {
	Do(req protocol.Request) (protocol.Response, error)
}

// New returns a new OpenAI client.
func New(token string, model string, baseURL string, organization string,
	apiType APIType, apiVersion string, httpClient Doer, embeddingsModel string,
	opts ...Option,
) (*Client, error) {
	c := &Client{
		token:           token,
		Model:           model,
		embeddingsModel: embeddingsModel,
		baseURL:         baseURL,
		organization:    organization,
		apiType:         apiType,
		apiVersion:      apiVersion,
		httpClient:      HertzClient(),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func HertzClient() *client.Client {

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
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else {
		req.Header.Set("api-key", c.token)
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