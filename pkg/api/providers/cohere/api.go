package cohere

import (
	"fmt"
	"net/http"

	"glide/pkg/api/providers"
)

// CohereAPIConfig represents the API configuration for Cohere
func CohereApiConfig(APIKey string) providers.ProviderApiConfig {
	return providers.ProviderApiConfig{
		BaseURL: "https://api.cohere.ai/v1",
		Headers: func(APIKey string) http.Header {
			headers := make(http.Header)
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", APIKey))
			headers.Set("content-type", "application/json")
			headers.Set("accept", "application/json")
			return headers
		},
		Complete: "/generate",
		Chat: "/chat",
		Embed: "/embed",
	}
}