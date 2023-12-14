// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
package pkg

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"fmt"
	"glide/pkg/providers"
	"glide/pkg/providers/openai"
	"encoding/json"
)

type ProviderConfigs = pkg.ProviderConfigs

// Initialize configList

var configList = map[string]interface{}{
    "openai": openai.OpenAIConfig,
}

// Create a new validator instance
var validate *validator.Validate = validator.New()


func BuildChatAPIRequest(payload map[string]string) (interface{}, error) {
    // API route api.glide.com/v1/chat
	// {"provider": "openai", "params": {"model": "gpt-3.5-turbo", "messages": "Hello, how are you?"}}

	// Sample JSON
	jsonStr := `{"provider": "openai", "params": {"model": "gpt-3.5-turbo", "messages": "Hello, how are you?"}}`

	// Define a map to hold the JSON data
	var data map[string]interface{}

	// Parse the JSON
	err := json.Unmarshal([]byte(jsonStr), &data)
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
	err := validate.Struct(providerConfig)
    if err != nil {
        // Handle validation error
        return nil, fmt.Errorf("validation error: %v", err)
    }
	// If everything is fine, return the providerConfig and nil error
    return providerConfig, nil
}
