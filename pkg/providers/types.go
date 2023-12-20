package providers

type GatewayConfig struct {
	Gateway PoolsConfig `yaml:"gateway"`
}
type PoolsConfig struct {
	Pools []Pool `yaml:"pools"`
}

type Pool struct {
	Name      string     `yaml:"pool"`
	Balancing string     `yaml:"balancing"`
	Providers []Provider `yaml:"providers"`
}

type Provider struct {
	Provider      string                 `yaml:"provider"`
	Model         string                 `yaml:"model"`
	ApiKey        string                 `yaml:"api_key"`
	TimeoutMs     int                    `yaml:"timeout_ms,omitempty"`
	DefaultParams map[string]interface{} `yaml:"default_params,omitempty"`
}

type RequestBody struct {
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	MessageHistory []string `json:"messageHistory"`
}
