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
	"log"
	"log/slog"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"
)

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
		fmt.Println(err)
	}
	return resp, err

}

func sendRequest(payload []byte) (interface{}, error) {

	// this function takes the client payload and returns the response from the provider

	slog.Info("sendRequest Function Called")

	requestDetails, err := definePayload(payload, "chat")

	if err != nil {
		println("Error defining payload: %v", err)
		return nil, err
	}

	// Create the full URL
	url := requestDetails.ApiConfig.BaseURL + requestDetails.ApiConfig.Endpoint

	slog.Info("Provider URL: " + url)

	// Marshal the requestDetails.RequestBody struct into JSON
	body, err := json.Marshal(requestDetails.RequestBody)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return nil, err
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	// If there was an error with creating the request, handle it
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	// Set the headers
	for key, values := range requestDetails.ApiConfig.Headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Send the request using http Client
	client := &http.Client{}
	return client.Do(req)
}

func definePayload(payload []byte, endpoint string) (providers.RequestDetails, error) {

	// this function takes the client payload and returns the request body for the provider as a struct

	slog.Info("definePayload Function Called")

	// Define a map to hold the JSON data
	var payload_data map[string]interface{}

	// Parse the JSON
	err := json.Unmarshal([]byte(payload), &payload_data)
	if err != nil {
		// Handle error
		slog.Error("Error unmarshalling payload: %v", err)
	}

	jsonString, _ := json.Marshal(payload_data)

	slog.Debug("Payload: " + string(jsonString))

	endpointsInterface := payload_data["endpoints"].([]interface{})
	var endpoints []map[string]interface{}

	for _, endpoint := range endpointsInterface {
		endpointMap, ok := endpoint.(map[string]interface{})
		if !ok {
			slog.Error("Invalid endpoint format")
			continue
		}
		endpoints = append(endpoints, endpointMap)
	}

	jsonString, _ = json.Marshal(endpoints)

	slog.Info("Endpoint: " + string(jsonString))

	providerList := make([]string, len(endpoints))
	for i, endpoint := range endpoints {
		provider, ok := endpoint["provider"].(string)
		if !ok {
			// Handle error
			fmt.Println("Provider not found")
		}

		providerList[i] = provider
	}

	slog.Info("Provider List: " + fmt.Sprintf("%v", providerList))

	// TODO: Send the providerList to the provider pool to get the provider selection. Mode list can be used as well. Mode is the routing strategy.
	//modeList := payload_data["mode"].([]interface{})


	provider := "openai" // placeholder until provider pool is implemented

	var params map[string]interface{}

	var api_key string

	for _, endpoint := range endpoints {
		if endpoint["provider"] == provider {
			params, _ = endpoint["params"].(map[string]interface{})
			api_key, _ = endpoint["api_key"].(string)
			break
		}
	}

	var defaultConfig interface{}                         // Assuming defaultConfig is a struct
	var finalApiConfig providers.ProviderDefinedApiConfig // Assuming finalApiConfig is a struct

	defaultConfig, finalApiConfig, _ = buildApiConfig(provider, api_key, endpoint)

	defaultConfigJson, _ := json.Marshal(defaultConfig)

	paramsJson, _ := json.Marshal(params)

	var defaultConfigMap map[string]interface{}
	json.Unmarshal(defaultConfigJson, &defaultConfigMap)

	var paramsMap map[string]interface{}
	json.Unmarshal(paramsJson, &paramsMap)

	for key, value := range paramsMap {
		if _, ok := defaultConfigMap[key]; ok {
			defaultConfigMap[key] = value
		}
	}

	updatedConfigJson, _ := json.Marshal(defaultConfigMap)

	err = json.Unmarshal(updatedConfigJson, &defaultConfig)
	if err != nil {
		slog.Error("Error occurred during unmarshalling. %v", err)
	}

	slog.Info("Default Config: " + fmt.Sprintf("%v", defaultConfig))
	slog.Info("Final API Config: " + fmt.Sprintf("%v", finalApiConfig))

	// Validate the struct
	validate := validator.New()
	err = validate.Struct(defaultConfig)
	if err != nil {
		slog.Error("Validation failed: ", err)
		return providers.RequestDetails{}, err
	}

	// Convert the struct to JSON
	defaultConfig, err = json.Marshal(defaultConfig)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	var requestDetails providers.RequestDetails = providers.RequestDetails{RequestBody: defaultConfig, ApiConfig: finalApiConfig}

	return requestDetails, nil
}

func buildApiConfig(provider string, api_key string, endpoint string) (interface{}, providers.ProviderDefinedApiConfig, error) {
	
	slog.Info("buildApiConfig Function Called")
	
	var defaultConfig interface{}
	var apiConfig providers.ProviderApiConfig
	var finalApiConfig providers.ProviderDefinedApiConfig

	switch provider {
	case "openai":
		defaultConfig = openai.OpenAiChatDefaultConfig()
		apiConfig = openai.OpenAiApiConfig(api_key)
	//case "cohere":
	//  defaultConfig = cohere.CohereChatDefaultConfig()
	//apiConfig = cohere.CohereAiApiConfig(api_key)
	default:
		return nil, providers.ProviderDefinedApiConfig{}, errors.New("invalid provider")
	}

	finalApiConfig.BaseURL = apiConfig.BaseURL
	finalApiConfig.Headers = apiConfig.Headers(api_key)

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
