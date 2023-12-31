package providers

// ModelProvider defines an interface all model providers should support
type ModelProvider interface {
	Provider() string
}
