package fields

import "time"

type Duration time.Duration

// MarshalText serializes Durations in a human-friendly way (it's shown in nanoseconds by default)
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}
