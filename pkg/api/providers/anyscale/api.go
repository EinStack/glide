package anyscale

import (
	"fmt"
	"net/http"

	"glide/pkg/api/providers"
)


func AnyscaleApiConfig(APIKey string) providers.ProviderApiConfig {
	return providers.ProviderApiConfig{
		BaseURL: "https://api.endpoints.anyscale.com/v1",
		Headers: func(APIKey string) http.Header {
			headers := make(http.Header)
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", APIKey))
			headers.Set("Content-Type", "application/json")
			return headers
		},
		Complete: "/completions",
		Chat:     "/chat/completions",
		Embed:    "/embeddings",
	}
}
