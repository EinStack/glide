package fields

import (
	"crypto/md5"
	"encoding"
	"encoding/hex"
)

// Secret is a string that is marshaled in an opaque way, so we are not leaking sensitive information
type Secret string

const maskedSecret = "[REDACTED]"

var _ encoding.TextMarshaler = Secret("")

// MarshalText marshals the secret as `[REDACTED]`.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(maskedSecret), nil
}

// Hash generates a digest of the secret to be used instead of the actual secret value
func (s Secret) Hash() string {
	hash := md5.Sum([]byte(s))

	return hex.EncodeToString(hash[:])
}
