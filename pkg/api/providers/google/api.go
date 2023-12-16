package google

import (
	"net/http"

	"glide/pkg/api/providers"
)


func GoogleApiConfig(APIKey string) providers.ProviderApiConfig {
	return providers.ProviderApiConfig{
		BaseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro",
		Headers: func(APIKey string) http.Header {
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			return headers
		},
		Complete: "/completions",
		Chat:     ":generateContent",
		Embed:    "/embeddings",
	}
}
