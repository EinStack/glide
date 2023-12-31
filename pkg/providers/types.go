package providers

import (
	"net/http"
	"time"
)

type ProviderVars struct {
	Name        string `yaml:"name"`
	ChatBaseURL string `yaml:"chatBaseURL"`
}

type RequestBody struct {
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	MessageHistory []string `json:"messageHistory"`
}

// Variables

var HTTPClient = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 2,
	},
}

type UnifiedAPIData struct {
	Model          string                 `json:"model"`
	APIKey         string                 `json:"api_key"`
	Params         map[string]interface{} `json:"params"`
	Message        map[string]string      `json:"message"`
	MessageHistory []map[string]string    `json:"messageHistory"`
}
