package providers

import "errors"

var ErrProviderUnavailable = errors.New("provider is not available")

// ModelProvider defines an interface all model providers should support
type ModelProvider interface {
	Provider() string
}
