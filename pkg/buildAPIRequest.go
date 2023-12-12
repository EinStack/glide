// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
package pkg

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"fmt"
	"glide/pkg/providers"
	"glide/pkg/providers/openai"
)

type ProviderConfigs = pkg.ProviderConfigs

// Initialize configList

var configList = map[string]interface{}{
    "openai": openai.OpenAIConfig,
}

// Create a new validator instance
var validate *validator.Validate = validator.New()


func BuildAPIRequest(provider string, params map[string]string, mode string) (interface{}, error) {
    // provider is the name of the provider, e.g. "openai", params is the map of parameters from the client, 
	// mode is the mode of the provider, e.g. "chat", configList is the list of provider configurations 


	var providerConfig map[string]interface{}
	if config, ok := configList[provider].(ProviderConfigs); ok {
    	if modeConfig, ok := config[mode].(map[string]interface{}); ok {
        providerConfig = modeConfig
    }
}

    // If the provider is not supported, return an error
    if providerConfig == nil {
        return nil, errors.New("unsupported provider")
    }
	

	// Build the providerConfig map by iterating over the keys in the providerConfig map and checking if the key exists in the params map

	for key := range providerConfig {
		if value, exists := params[key]; exists {
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
