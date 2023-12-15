package pkg

import (
	"net/http"
)

type ProviderConfigs map[string]interface{}

type RequestDetails struct {
	RequestBody interface{}
	ApiConfig ProviderApiConfig
}

type ProviderApiConfig struct {
	BaseURL  string
	Headers  func(string) http.Header
	Complete string
	Chat     string
	Embed    string
}