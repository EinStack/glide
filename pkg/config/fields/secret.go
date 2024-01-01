package fields

import "encoding"

// Secret is a string that is marshaled in an opaque way, so we are not leaking sensitive information
type Secret string

const maskedSecret = "***"

var _ encoding.TextMarshaler = Secret("")

// MarshalText marshals the secret as `[REDACTED]`.
func (s Secret) MarshalText() ([]byte, error) {
	if len(s) == 0 {
		return []byte(maskedSecret), nil
	}

	return []byte(maskedSecret + s[:3]), nil
}
