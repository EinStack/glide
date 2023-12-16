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
	"github.com/go-playground/validator/v10"
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
	var url string
	switch provider {
	case "google":
		url = requestDetails.ApiConfig.BaseURL
	default:
		url = requestDetails.ApiConfig.BaseURL + requestDetails.ApiConfig.Endpoint
	}

	slog.Info("rovider URL: " + url)

	// Marshal the requestDetails.RequestBody struct into JSON
	body, err := json.Marshal(requestDetails.RequestBody)
	if err != nil {
		slog.Error("error marshalling request body: %v", err)
		return nil, err
	}

	slog.Info("request Body: " + string(body))

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	// If there was an error with creating the request, handle it
	if err != nil {
		slog.Error("Error creating request: %v", err)
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

	resp, err := client.Do(req)

	if err != nil {
		slog.Error("Error sending request: ", err)
	}

	// retrieve payload from response
	body, _ = io.ReadAll(resp.Body)

	// Unmarshal the JSON response into a map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		slog.Error("Error decoding JSON response:", err)
	}

	return responseMap, nil

}

func definePayload(payload []byte, endpoint string) (providers.RequestDetails, string, error) {

	// this function takes the client payload and returns the request body for the provider as a struct

	slog.Info("definePayload Function Called")

	// Define a map to hold the JSON data. Should these be declared at the top?
	var defaultConfig interface{}                         // Assuming defaultConfig is a struct
	var finalApiConfig providers.ProviderDefinedApiConfig // Assuming finalApiConfig is a struct
	var payload_data map[string]interface{}
	var params map[string]interface{}
	var api_key string
	var endpoints []map[string]interface{}
	var defaultConfigMap map[string]interface{}
	var paramsMap map[string]interface{}

	// Parse the JSON
	err := json.Unmarshal([]byte(payload), &payload_data)
	if err != nil {
		// Handle error
		slog.Error("Error unmarshalling payload: %v", err)
	}

	jsonString, _ := json.Marshal(payload_data)

	slog.Debug("Payload: " + string(jsonString))

	endpointsInterface := payload_data["endpoints"].([]interface{})

	for _, endpoint := range endpointsInterface {
		endpointMap, ok := endpoint.(map[string]interface{})
		if !ok {
			slog.Error("Invalid endpoint format")
			continue
		}
		endpoints = append(endpoints, endpointMap)
	}

	jsonString, _ = json.Marshal(endpoints)

	slog.Debug("Endpoint: " + string(jsonString))

	providerList := make([]string, len(endpoints))
	for i, endpoint := range endpoints {
		provider, ok := endpoint["provider"].(string)
		if !ok {
			// Handle error
			fmt.Println("provider not found")
		}

		providerList[i] = provider
	}

	slog.Info("provider list: " + fmt.Sprintf("%v", providerList))

	// TODO: Send the providerList to the provider pool to get the provider selection. Mode list can be used as well. Mode is the routing strategy.
	//modeList := payload_data["mode"].([]interface{})
	// provider is used in sendRequest function to build Google URL. This might need to change or provider should be a const

	for _, endpoint := range endpoints {
		if endpoint["provider"] == provider {
			params, _ = endpoint["params"].(map[string]interface{})
			api_key, _ = endpoint["api_key"].(string)
			break
		}
	}

	// retrieve the default configs for the provider
	defaultConfig, finalApiConfig, _ = buildApiConfig(provider, api_key, endpoint)

	// convert the defaultConfig and params to maps
	defaultConfigJson, _ := json.Marshal(defaultConfig)

	paramsJson, _ := json.Marshal(params)


	json.Unmarshal(defaultConfigJson, &defaultConfigMap)

	json.Unmarshal(paramsJson, &paramsMap)

	for key, value := range paramsMap {
		if _, ok := defaultConfigMap[key]; ok {
			defaultConfigMap[key] = value
		}
	}

	updatedConfigJson, _ := json.Marshal(defaultConfigMap)

	err = json.Unmarshal(updatedConfigJson, &defaultConfig)
	if err != nil {
		slog.Error("error occurred during unmarshalling. %v", err)
	}

	//err = validateStruct(defaultConfig)

	//if err != nil {
	//	slog.Error("error occurred during validation.", err)
	//	return providers.RequestDetails{}, err
	//}

	slog.Info("default Config: " + fmt.Sprintf("%v", defaultConfig))

	var requestDetails providers.RequestDetails = providers.RequestDetails{RequestBody: defaultConfig, ApiConfig: finalApiConfig}

	slog.Debug("requestDetails: " + fmt.Sprintf("%v", requestDetails))

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


func validateStruct(config interface{}) error {
    // Validate a struct
	validate := validator.New()
    err := validate.Struct(config)
    if err != nil {
        if invalidErr, ok := err.(*validator.InvalidValidationError); ok {
            return fmt.Errorf("invalid validation error: %v", invalidErr)
        }

        for _, err := range err.(validator.ValidationErrors) {
            fmt.Println(err.Namespace())
            fmt.Println(err.Field())
            fmt.Println(err.StructNamespace())
            fmt.Println(err.StructField())
            fmt.Println(err.Tag())
            fmt.Println(err.ActualTag())
            fmt.Println(err.Kind())
            fmt.Println(err.Type())
            fmt.Println(err.Value())
            fmt.Println(err.Param())
            fmt.Println()
        }

        return fmt.Errorf("validation failed for the configuration")
    }

    return nil
}
