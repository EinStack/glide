// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
// 1. receive client payload
// 2. determine provider
// 3. build request body based on provider
// 4. send request to provider
package pkg

import (
	//"errors"
	//"github.com/go-playground/validator/v10"
	"fmt"
	"glide/pkg/providers"
	"glide/pkg/providers/openai"
	"encoding/json"
	//"net/http"
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

	// TODO: use mode and providerList to determine which provider to use
	//modeList := payload_data["mode"].([]interface{})

	provider := "openai"

	// select the predefined config for the provider
	var providerConfig map[string]interface{}
	if config, ok := configList[provider].(pkg.ProviderConfigs); ok { // this pulls the config in index.go
    	if modeConfig, ok := config["chat"].(map[string]interface{}); ok { // this pulls the specific config for the endpoint
        providerConfig = modeConfig
    }
}

	// Build the providerConfig map by iterating over the keys in the providerConfig map and checking if the key exists in the params map

	for key := range providerConfig {
		if value, exists := payload_data[key]; exists {
			providerConfig[key] = value
		}
	}

	// If everything is fine, return the providerConfig and nil error
	println(providerConfig)
    return providerConfig, nil
}

