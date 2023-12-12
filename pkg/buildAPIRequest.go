package pkg

import (
	"errors"
)

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
	

	// TODO: Next need to build the request based on the params from the client
	// First for each param in param check if present. If yes then add it to the request.
	// If not & the param is required, return a default value from the provider config

	for key := range providerConfig {
		if value, exists := params[key]; exists {
			providerConfig[key] = value
		}
	}

	
}

    // For now, return providerConfig and nil error to satisfy the function signature
    return providerConfig, nil

	
}