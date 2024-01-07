package health

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrNoTokens = errors.New("not enough tokens in the bucket")
)

// TokenBucket concurrency-safe implementation based on atomic operations for the max performance
//
//	We are not tracking the number of tokens directly,
//	but rather the time that has passed since the last token consumption
type TokenBucket struct {
	timePointer  uint64
	timePerToken uint64
	timePerBurst uint64
}

func NewTokenBucket(rate, burstSize uint64) *TokenBucket {
	return &TokenBucket{
		timePointer:  0,
		timePerToken: 1_000_000 / rate,
		timePerBurst: burstSize * (1_000_000 / rate),
	}
}

// Take one bucket from the bucket
func (b *TokenBucket) Take(tokens uint64) error {
	oldTime := atomic.LoadUint64(&b.timePointer)
	newTime := oldTime

	timeNeeded := tokens * b.timePerToken

	for {
		now := uint64(time.Now().UnixNano() / 1000)
		minTime := now - b.timePerBurst

		// Take into account burst size.
		if minTime > oldTime {
			newTime = minTime
		}

		// Now shift by the time needed.
		newTime += timeNeeded

		// Check if too many tokens.
		if newTime > now {
			return ErrNoTokens
		}

		if atomic.CompareAndSwapUint64(&b.timePointer, oldTime, newTime) {
			// consumed tokens
			return nil
		}

		// Otherwise load old value and try again.
		oldTime = atomic.LoadUint64(&b.timePointer)
		newTime = oldTime
	}
}

func (b *TokenBucket) HasTokens() bool {
	return b.Tokens() >= 1
}

// Tokens returns number of available tokens in the bucket
func (b *TokenBucket) Tokens() uint64 {
	timePointer := atomic.LoadUint64(&b.timePointer)
	now := uint64(time.Now().UnixNano() / 1000)
	minTime := now - b.timePerBurst

	newTime := timePointer

	// Take into account burst size.
	if minTime > timePointer {
		newTime = minTime
	}

	return newTime / b.timePerToken
}
