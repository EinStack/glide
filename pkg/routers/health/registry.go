package health

// ProviderHealthRegistry holds information about provider API health across routers
type ProviderHealthRegistry struct{}

func NewHealthRegistry() *ProviderHealthRegistry {
	return &ProviderHealthRegistry{}
}
