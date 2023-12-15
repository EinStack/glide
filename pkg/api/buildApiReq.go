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
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"
	"glide/pkg/api/providers"
	"glide/pkg/api/providers/openai"
	"log"
	"net/http"
	"reflect"
)

func LlmRouter(c *app.RequestContext) (interface{}, error) {

	// this function takes the client request and returns the response from the provider

	requestBody := c.Request.Body()

	if requestBody == nil {
		fmt.Println("Request body is nil")
	}

	fmt.Println(requestBody)

	// Send the request to the provider
	resp, err := sendRequest(requestBody)
	if err != nil {
		fmt.Println(err)
	}
	return resp, err

}

func sendRequest(payload []byte) (interface{}, error) {

	// this function takes the client payload and returns the response from the provider

	requestDetails, err := definePayload(payload, "chat")

	if err != nil {
		println("Error defining payload: %v", err)
		return nil, err
	}

	// Create the full URL
	url := requestDetails.ApiConfig.BaseURL + requestDetails.ApiConfig.Endpoint

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

	// Define a map to hold the JSON data
	var payload_data map[string]interface{}

	// Parse the JSON
	err := json.Unmarshal([]byte(payload), &payload_data)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	endpoints, ok := payload_data["endpoints"].([]interface{})
	if !ok {
		// Handle error
		fmt.Println("Endpoints not found")
	}

	providerList := make([]string, len(endpoints))
	for i, endpoint := range endpoints {
		endpointMap, ok := endpoint.(map[string]interface{})
		if !ok {
			// Handle error
			fmt.Println("Endpoint is not a map")
		}

		provider, ok := endpointMap["provider"].(string)
		if !ok {
			// Handle error
			fmt.Println("Provider not found")
		}

		providerList[i] = provider
	}

	// TODO: Send the providerList to the provider pool to get the provider selection. Mode list can be used as well. Mode is the routing strategy.
	//modeList := payload_data["mode"].([]interface{})

	// TODO: the following is inefficient. Needs updating.

	provider := "openai" // placeholder until provider pool is implemented

	endpointsMap := payload_data["endpoints"].([]map[string]interface{})

	var params map[string]interface{}

	var api_key string

	for _, endpoint := range endpointsMap {
		if endpoint["provider"] == provider {
			params = endpoint["params"].(map[string]interface{})
			api_key = endpoint["api_key"].(string)
			fmt.Println(params)
			break
		}
	}

	var defaultConfig interface{}                   // Assuming defaultConfig is a struct
	var finalApiConfig providers.ProviderDefinedApiConfig // Assuming finalApiConfig is a struct

	defaultConfig, finalApiConfig, _ = buildApiConfig(provider, api_key, endpoint)

	// Use reflect to set the value in defaultConfig based on client payload
	v := reflect.ValueOf(defaultConfig).Elem()
	for key, value := range params {
		field := v.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Int:
				if val, ok := value.(int); ok {
					field.SetInt(int64(val))
				}
			case reflect.String:
				if val, ok := value.(string); ok {
					field.SetString(val)
				}
			}
		}
	}

	// Validate the struct
	validate := validator.New()
	err = validate.Struct(defaultConfig)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
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
