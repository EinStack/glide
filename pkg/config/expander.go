package config

// Expander finds special directives like ${env:ENV_VAR} in the config file and fill them with actual values
type Expander struct{}

func (e *Expander) Expand(content []byte) ([]byte, error) {
	expandedContent := string(content)

	// TODO: Expand env vars & files

	return []byte(expandedContent), nil
}
