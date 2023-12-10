package pkg

import (
	"errors"
)

type ProviderConfigs map[string]interface{} // TODO: import from types.go

func BuildAPIRequest(provider string, params map[string]string, mode string, config_list map[string]interface{}) (interface{}, error) {
    
	var providerConfig interface{}
    if config, ok := config_list[provider].(ProviderConfigs); ok {
        providerConfig = config[mode]
    }

    // If the provider is not supported, return an error
    if providerConfig == nil {
        return nil, errors.New("unsupported provider")
    }

    // For now, return providerConfig and nil error to satisfy the function signature
    return providerConfig, nil

	// TODO: Next need to build the request based on the params from the client
	// First check if the param is present. If yes then add it to the request.
	// If not & the param is required, return a default value from the provider config
}