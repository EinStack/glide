package providers

type GatewayConfig struct {
	Gateway PoolsConfig `yaml:"gateway" validate:"required"`
}
type PoolsConfig struct {
	Pools []Pool `yaml:"pools" validate:"required"`
}

type Pool struct {
	Name      string     `yaml:"pool" validate:"required"`
	Balancing string     `yaml:"balancing" validate:"required"`
	Providers []Provider `yaml:"providers" validate:"required"`
}

type Provider struct {
	Provider      string                 `yaml:"provider" validate:"required"`
	Model         string                 `yaml:"model"`
	APIKey        string                 `yaml:"api_key" validate:"required"`
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
