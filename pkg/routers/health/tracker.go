package health

import (
	"errors"

	"glide/pkg/providers/clients"
)

// Tracker tracks errors and general health of model provider
type Tracker struct {
	errBudget *TokenBucket
	rateLimit *RateLimitTracker
}

func NewTracker(budget ErrorBudget) *Tracker {
	return &Tracker{
		rateLimit: NewRateLimitTracker(),
		errBudget: NewTokenBucket(budget.TimePerTokenMicro(), budget.Budget()),
	}
}

func (t *Tracker) Healthy() bool {
	return !t.rateLimit.Limited() && t.errBudget.HasTokens()
}

func (t *Tracker) TrackErr(err error) {
	var rateLimitErr *clients.RateLimitError

	if errors.As(err, &rateLimitErr) {
		t.rateLimit.SetLimited(rateLimitErr.UntilReset())

		return
	}

	_ = t.errBudget.Take(1)
}
