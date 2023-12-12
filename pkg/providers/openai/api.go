package pkg

import (
	"fmt"
	"net/http"
)


// provides the base URL and headers for the OpenAI API
type ProviderAPIConfig struct {
	BaseURL  string
	Headers  func(string) http.Header
	Complete string
	Chat     string
	Embed    string
}

func OpenAIAPIConfig(APIKey string) *ProviderAPIConfig {
	return &ProviderAPIConfig{
		BaseURL: "https://api.openai.com/v1",
		Headers: func(APIKey string) http.Header {
			headers := make(http.Header)
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", APIKey))
			return headers
		},
		Complete: "/completions",
		Chat:     "/chat/completions",
		Embed:    "/embeddings",
	}
}
