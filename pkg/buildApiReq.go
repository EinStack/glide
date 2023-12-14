// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
// 1. receive client payload
// 2. determine provider
// 3. build request body based on provider
// 4. send request to provider
package pkg

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"fmt"
	"glide/pkg/providers"
	"glide/pkg/providers/openai"
	"encoding/json"
	"net/http"
	"time"
)


// Initialize configList

var configList = map[string]interface{}{
    "openai": openai.OpenAIConfig,
}

// Sample JSON
//var jsonStr = `{"provider": "openai", "params": {"model": "gpt-3.5-turbo", "messages": "Hello, how are you?"}}`


func DefinePayload(payload []byte) (interface{}, error) {
    // API route api.glide.com/v1/chat
	// {"provider": "openai", "params": {"model": "gpt-3.5-turbo", "messages": "Hello, how are you?"}}

	// Define a map to hold the JSON data
	var payload_data map[string]interface{}

	// Parse the JSON
	err := json.Unmarshal([]byte(payload), &payload_data)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	// Extract the provider
	provider, ok := data["provider"].(string)
	if !ok {
		// Handle error
		fmt.Println("Provider not found")
	}

	// select the predefined config for the provider
	var providerConfig map[string]interface{}
	if config, ok := configList[provider].(ProviderConfigs); ok {
    	if modeConfig, ok := config["chat"].(map[string]interface{}); ok {
        providerConfig = modeConfig
    }
}

    // If the provider is not supported, return an error
    if providerConfig == nil {
        return nil, errors.New("unsupported provider")
    }
	

	// Build the providerConfig map by iterating over the keys in the providerConfig map and checking if the key exists in the params map

	for key := range providerConfig {
		if value, exists := payload[key]; exists {
			providerConfig[key] = value
		}
	}

	// Validate the providerConfig map using the validator package
	validate := validator.New()
    err := validate.Struct(payload)
    if err != nil {
        // Handle validation error
        return nil, fmt.Errorf("validation error: %v", err)
    }

	// If everything is fine, return the providerConfig and nil error
    return providerConfig, nil
}

// SendRequest sends an HTTP request to a specific URL path with given headers and payload.
func SendRequest(config pkg.providerConfig) (map[string]interface{}, error) {
	client := &http.Client{
	}

	// TODO: convert the pr
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(result["error"].(string))
	}

	return result, nil
}
