package config

import "encoding"

// Secret is a string that is marshaled in an opaque way, so we are not leaking sensitive information
type Secret string

const maskedSecret = "[REDACTED]"

var _ encoding.TextMarshaler = Secret("")

// MarshalText marshals the string as `[REDACTED]`.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(maskedSecret), nil
}
