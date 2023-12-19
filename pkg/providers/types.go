package provider
type GatewayConfig struct {
	Pools []Pool `yaml:"pools"`
}

type Pool struct {
	Name       string      `yaml:"name"`
	Balancing  string      `yaml:"balancing"`
	Providers  []Provider  `yaml:"providers"`
}

type Provider struct {
	Name           string                 `yaml:"name"`
	Provider       string                 `yaml:"provider"`
	Model          string                 `yaml:"model"`
	APIKey         string                 `yaml:"api_key"`
	TimeoutMs      int                    `yaml:"timeout_ms,omitempty"`
	DefaultParams  map[string]interface{} `yaml:"default_params,omitempty"`
}