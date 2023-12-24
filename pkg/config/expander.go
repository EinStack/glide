package config

const directivePattern = `[A-Za-z][A-Za-z0-9+.-]+`

type RawConfig = map[string]interface{}

// Expander finds special directives like ${env:ENV_VAR} in the config file and fill them with actual values
type Expander struct{}

func (e *Expander) Expand(rawConfig RawConfig) (RawConfig, error) {
	for configKey, configValue := range rawConfig {
		expandedValue, err := e.expandConfigValue(configValue)

		if err != nil {
			return nil, err
		}

		rawConfig[configKey] = expandedValue
	}

	return rawConfig, nil
}

func (e *Expander) expandConfigValue(value any) (any, error) {
	switch v := value.(type) {
	case string:
		return e.expandString(v)
	case []any:
		return e.expandSlice(v)
	case map[string]any:
		return e.expandMap(v)
	default:
		// could be int or float
		return value, nil
	}
}

func (e *Expander) expandString(value string) (string, error) {
	// TODO: implement
	return value, nil
}

func (e *Expander) expandSlice(value []any) ([]any, error) {
	slice := make([]any, 0, len(value))

	for _, vint := range value {
		val, err := e.expandConfigValue(vint)
		if err != nil {
			return nil, err
		}

		slice = append(slice, val)
	}

	return slice, nil
}

func (e Expander) expandMap(value map[string]any) (map[string]any, error) {
	newMap := map[string]any{}

	for mapKey, mapValue := range value {
		val, err := e.expandConfigValue(mapValue)
		if err != nil {
			return nil, err
		}

		newMap[mapKey] = val
	}

	return newMap, nil
}
