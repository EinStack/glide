// Package api contains functions to build and send API requests to various providers.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"glide/pkg/api/providers"
	"glide/pkg/api/providers/anyscale"
	"glide/pkg/api/providers/cohere"
	"glide/pkg/api/providers/google"
	"glide/pkg/api/providers/openai"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Router handles the client request and returns the response from the provider.
func Router(c *app.RequestContext) (interface{}, error) {
	logrus.Info("Router Function Called")

	requestBody, err := io.ReadAll(c.Request.Body())
	if err != nil {
		logrus.Error("failed to read request body:", err)
		return nil, err
	}

	if len(requestBody) == 0 {
		errMsg := "request body cannot be empty"
		logrus.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	resp, err := sendRequest(requestBody)
	if err != nil {
		logrus.Errorf("Error with request: %v", err)
		return nil, err
	}

	logrus.Infof("Response: %+v", resp)
	return resp, nil
}

// sendRequest takes the client payload and returns the response from the provider.
func sendRequest(payload []byte) (interface{}, error) {
	logrus.Info("sendRequest Function Called")

	requestDetails, provider, err := definePayload(payload, "chat")
	if err != nil {
		logrus.Errorf("error defining payload: %v", err)
		return nil, err
	}

	url := buildProviderURL(provider, requestDetails.ApiConfig)
	logrus.Infof("Provider URL: %s", url)

	body, err := json.Marshal(requestDetails.RequestBody)
	if err != nil {
		logrus.Errorf("error marshalling request body: %v", err)
		return nil, err
	}

	logrus.Infof("Request Body: %s", string(body))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("Error creating request: %v", err)
		return nil, err
	}

	setRequestHeaders(req, requestDetails.ApiConfig.Headers)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error reading response body: %v", err)
		return nil, err
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		logrus.Errorf("Error decoding JSON response: %v", err)
		return nil, err
	}

	return responseMap, nil
}

// definePayload takes the client payload and returns the request body for the provider as a struct.
func definePayload(payload []byte, endpoint string) (providers.RequestDetails, string, error) {
	logrus.Info("definePayload Function Called")

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payload, &payloadData); err != nil {
		logrus.Errorf("Error unmarshalling payload: %v", err)
		return providers.RequestDetails{}, "", err
	}

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

// buildApiConfig constructs the API configuration based on the provider and API key.
func buildApiConfig(provider string, apiKey string, endpoint string) (interface{}, providers.ProviderDefinedApiConfig, error) {
	logrus.Info("buildApiConfig function Called")

	var defaultConfig interface{}
	var apiConfig providers.ProviderApiConfig
	var finalApiConfig providers.ProviderDefinedApiConfig

	switch provider {
	case "openai":
		defaultConfig = openai.OpenAiChatDefaultConfig()
		apiConfig = openai.OpenAiApiConfig(apiKey)
	case "cohere":
		defaultConfig = cohere.CohereChatDefaultConfig()
		apiConfig = cohere.CohereApiConfig(apiKey)
	case "anyscale":
		defaultConfig = anyscale.AnyscaleChatDefaultConfig()
		apiConfig = anyscale.AnyscaleApiConfig(apiKey)
	case "google":
		defaultConfig = google.GoogleChatDefaultConfig()
		apiConfig = google.GoogleApiConfig(apiKey)
	default:
		return nil, providers.ProviderDefinedApiConfig{}, errors.New("invalid provider")
	}

	finalApiConfig = buildFinalApiConfig(apiConfig, apiKey, endpoint)
	return defaultConfig, finalApiConfig, nil
}

// validateStruct validates a struct using the validator package.
func validateStruct(config interface{}) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("invalid validation error: %v", err)
		}
		return fmt.Errorf("validation failed for the configuration: %v", err)
	}
	return nil
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

// buildFinalApiConfig constructs the final API configuration based on the base configuration, API key, and endpoint.
func buildFinalApiConfig(apiConfig providers.ProviderApiConfig, apiKey string, endpoint string) providers.ProviderDefinedApiConfig {
	var finalApiConfig providers.ProviderDefinedApiConfig
	if endpoint == "chat" && apiConfig.Chat != "" {
		finalApiConfig.Endpoint = apiConfig.Chat
	} else if endpoint == "complete" && apiConfig.Complete != "" {
		finalApiConfig.Endpoint = apiConfig.Complete
	}

	finalApiConfig.BaseURL = apiConfig.BaseURL
	finalApiConfig.Headers = apiConfig.Headers(apiKey)

	if endpoint == "chat" && apiConfig.ProviderName == "google" {
		finalApiConfig.BaseURL += "?key=" + apiKey
	}

	return finalApiConfig
}