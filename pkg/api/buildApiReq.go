// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
// 1. receive client payload
// 2. determine provider
// 3. build request body based on provider
// 4. send request to provider
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"glide/pkg/api/providers"
	"glide/pkg/api/providers/openai"
	"glide/pkg/api/providers/cohere"
	"glide/pkg/api/providers/anyscale"
	"glide/pkg/api/providers/google"
	"io"
	"log/slog"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

const provider = "google" // placeholder for testing until provider pool is implemented

func Router(c *app.RequestContext) (interface{}, error) {
	// this function takes the client request and returns the response from the provider
	slog.Info("Router Function Called")

	requestBody := c.Request.Body()

	// Check if the request body is empty
	if len(requestBody) == 0 {
		slog.Error("request body cannot be empty")
		return nil, errors.New("request body cannot be empty")
	}

	// Send the request to the provider
	resp, err := sendRequest(requestBody)
	if err != nil {
		slog.Error("Error with request: %v", err)
	}

	slog.Info("Response: " + fmt.Sprintf("%v", resp))

	return resp, err

}

func sendRequest(payload []byte) (interface{}, error) {

	// this function takes the client payload and returns the response from the provider

	slog.Info("sendRequest Function Called")

	requestDetails, provider, err := definePayload(payload, "chat")

	if err != nil {
		slog.Error("error defining payload: %v", err)
		return nil, err
	}

	// Create the full URL
	url := buildProviderURL(provider, requestDetails.ApiConfig)

	slog.Debug("provider URL: " + url)

	// Marshal the requestDetails.RequestBody struct into JSON
	body, err := json.Marshal(requestDetails.RequestBody)
	if err != nil {
		slog.Error("error marshalling request body: %v", err)
		return nil, err
	}

	slog.Debug("request Body: " + string(body))

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("Error creating request: %v", err)
		return nil, err
	}

	// Set the headers
	setRequestHeaders(req, requestDetails.ApiConfig.Headers)

	// Send the request using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body: %v", err)
		return nil, err
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		slog.Error("error decoding JSON response: %v", err)
		return nil, err
	}

	return responseMap, nil

}

// definePayload takes the client payload and returns the request body for the provider as a struct.
func definePayload(payload []byte, endpoint string) (providers.RequestDetails, string, error) {
	slog.Info("definePayload Function Called")

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payload, &payloadData); err != nil {
		slog.Error("Error unmarshalling payload: %v", err)
		return providers.RequestDetails{}, "", err
	}

	// This is where the provider 
	provider, err := extractProvider(payloadData)
	if err != nil {
		return providers.RequestDetails{}, "", err
	}

	params, apiKey, err := extractEndpointDetails(payloadData, provider)
	if err != nil {
		return providers.RequestDetails{}, "", err
	}

	defaultConfig, finalApiConfig, err := buildApiConfig(provider, apiKey, endpoint)
	if err != nil {
		return providers.RequestDetails{}, "", err
	}

	mergedConfig, err := mergeConfigs(defaultConfig, params)
	if err != nil {
		return providers.RequestDetails{}, "", err
	}

	requestDetails := providers.RequestDetails{
		RequestBody: mergedConfig,
		ApiConfig:   finalApiConfig,
	}

	return requestDetails, provider, nil
}

func buildApiConfig(provider string, api_key string, endpoint string) (interface{}, providers.ProviderDefinedApiConfig, error) {

	// TODO: CLEAN THIS UP
	slog.Info("buildApiConfig function Called")

	var defaultConfig interface{}
	var apiConfig providers.ProviderApiConfig
	var finalApiConfig providers.ProviderDefinedApiConfig

	switch provider {
	case "openai":
		defaultConfig = openai.OpenAiChatDefaultConfig()
		apiConfig = openai.OpenAiApiConfig(api_key)
	case "cohere":
	    defaultConfig = cohere.CohereChatDefaultConfig()
	    apiConfig = cohere.CohereApiConfig(api_key)
	case "anyscale":
		defaultConfig = anyscale.AnyscaleChatDefaultConfig()
		apiConfig = anyscale.AnyscaleApiConfig(api_key)
	case "google":
		defaultConfig = google.GoogleChatDefaultConfig()
		apiConfig = google.GoogleApiConfig(api_key)
	default:
		return nil, providers.ProviderDefinedApiConfig{}, errors.New("invalid provider")
	}

	switch provider {
	case "google":
		finalApiConfig.BaseURL = apiConfig.BaseURL + apiConfig.Chat + "?key=" + api_key
		finalApiConfig.Headers = apiConfig.Headers(api_key)
	default:
		finalApiConfig.BaseURL = apiConfig.BaseURL
		finalApiConfig.Headers = apiConfig.Headers(api_key)

	}

	switch endpoint {
	case "chat":
		finalApiConfig.Endpoint = apiConfig.Chat
	case "complete":
		finalApiConfig.Endpoint = apiConfig.Complete
	default:
		return nil, providers.ProviderDefinedApiConfig{}, errors.New("invalid endpoint")
	}

	return defaultConfig, finalApiConfig, nil
}

// Helper functions below:

// buildProviderURL constructs the URL for the provider request.
func buildProviderURL(provider string, apiConfig providers.ProviderDefinedApiConfig) string {
	if provider == "google" {
		return apiConfig.BaseURL
	}
	return apiConfig.BaseURL + apiConfig.Endpoint
}

// setRequestHeaders sets the headers for the HTTP request.
func setRequestHeaders(req *http.Request, headers map[string][]string) {
	for key, values := range headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
}

// extractProvider extracts the provider from the payload data.
func extractProvider(payloadData map[string]interface{}) (string, error) {
	endpoints, ok := payloadData["endpoints"].([]interface{})
	if !ok {
		return "", errors.New("endpoints not found or invalid format")
	}

	for _, endpoint := range endpoints {
		endpointMap, ok := endpoint.(map[string]interface{})
		if !ok {
			continue
		}
		if provider, ok := endpointMap["provider"].(string); ok {
			return provider, nil
		}
	}
	return "", errors.New("provider not found in endpoints")
}

// extractEndpointDetails extracts details like parameters and API key for a given provider.
func extractEndpointDetails(payloadData map[string]interface{}, provider string) (map[string]interface{}, string, error) {
	endpoints, ok := payloadData["endpoints"].([]interface{})
	if !ok {
		return nil, "", errors.New("endpoints not found or invalid format")
	}

	for _, endpoint := range endpoints {
		endpointMap, ok := endpoint.(map[string]interface{})
		if !ok || endpointMap["provider"] != provider {
			continue
		}
		params, _ := endpointMap["params"].(map[string]interface{})
		apiKey, _ := endpointMap["api_key"].(string)
		return params, apiKey, nil
	}
	return nil, "", errors.New("endpoint details not found for provider")
}

// mergeConfigs merges the default configuration with the parameters provided.
func mergeConfigs(defaultConfig interface{}, params map[string]interface{}) (interface{}, error) {
	defaultConfigMap, err := structToMap(defaultConfig)
	if err != nil {
		return nil, err
	}

	for key, value := range params {
		defaultConfigMap[key] = value
	}

	return mapToStruct(defaultConfigMap, defaultConfig)
}

// structToMap converts a struct to a map.
func structToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &resultMap); err != nil {
		return nil, err
	}

	return resultMap, nil
}

// mapToStruct converts a map to a struct based on the provided struct type.
func mapToStruct(dataMap map[string]interface{}, structType interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}

	resultStruct := structType
	if err := json.Unmarshal(jsonData, &resultStruct); err != nil {
		return nil, err
	}

	return resultStruct, nil
}

