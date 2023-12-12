package pkg

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// Create a new validator instance
var validate *validator.Validate = validator.New()

type ProviderConfigs map[string]interface{} // TODO: import from types.go

func BuildAPIRequest(provider string, params map[string]string, mode string, configList map[string]interface{}) (interface{}, error) {
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
        return nil, err
    }

	
}

    // For now, return providerConfig and nil error to satisfy the function signature
    return providerConfig, nil

	
}