package pkg

import (
	"net/http"
)

type ProviderConfigs map[string]interface{}

type RequestDetails struct {
	RequestBody interface{}
	ApiConfig ProviderDefinedApiConfig
}

type ProviderApiConfig struct {
	BaseURL  string
	Headers  func(string) http.Header
	Complete string
	Chat     string
	Embed    string
}

type ProviderDefinedApiConfig struct {
	BaseURL  string
	Headers  http.Header
	Endpoint string
}