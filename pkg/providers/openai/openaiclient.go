// TODO: Explore resource pooling
// TODO: Optimize Type use
// TODO: Explore Hertz TLS & resource pooling
// OpenAI package provide a set of functions to interact with the OpenAI API.
package openai

import (
	"errors"

	"glide/pkg/providers"
)

const (
	providerName   = "openai"
	defaultBaseURL = "https://api.openai.com/v1"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var (
	ErrEmptyResponse = errors.New("empty response")
)

// OpenAiClient creates a new client for the OpenAI API.
//
// Parameters:
// - poolName: The name of the pool to connect to.
// - modelName: The name of the model to use.
//
// Returns:
// - *Client: A pointer to the created client.
// - error: An error if the client creation failed.
func Client() (*ProviderClient, error) {
	// Create a new client
	c := &ProviderClient{
		BaseURL:    defaultBaseURL,
		HTTPClient: providers.HTTPClient,
	}

	return c, nil
}
