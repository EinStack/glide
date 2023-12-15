// this file contains the BuildAPIRequest function which takes in the provider name, params map, and mode and returns the providerConfig map and error
// The providerConfig map can be used to build the API request to the provider
// 1. receive client payload
// 2. determine provider
// 3. build request body based on provider
// 4. send request to provider
package pkg

import (
	//"errors"
	"github.com/go-playground/validator/v10"
	"fmt"
	//"glide/pkg/providers"
	"glide/pkg/providers/openai"
	"encoding/json"
	//"net/http"
	"reflect"
)

func DefinePayload(payload []byte) (interface{}, error) {

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

    // TODO: the following is inefficient. Needs updating.
    endpointsMap := payload_data["endpoints"].([]map[string]interface{})

	var params map[string]interface{} 

    for _, endpoint := range endpointsMap {
        if endpoint["provider"] == provider {
            params := endpoint["params"].(map[string]interface{})
            fmt.Println(params)
            break
        }
    }

    var defaultConfig interface{} // Assuming defaultConfig is a struct

    if provider == "openai" {
        defaultConfig = openai.OpenAiChatDefaultConfig() // this is a struct
    } else if provider == "cohere" {
        defaultConfig = openai.OpenAiChatDefaultConfig() //TODO: change this to cohere
    }

    // Use reflect to set the value in defaultConfig
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
        return nil, err
    }

    return defaultConfig, nil
}