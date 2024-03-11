package health

import (
	"errors"
	"glide/pkg/providers/clients"
)

type HealthTracker struct {
	errBudget *TokenBucket
	rateLimit *RateLimitTracker
}

func NewHealthTracker(budget ErrorBudget) *HealthTracker {
	return &HealthTracker{
		rateLimit: NewRateLimitTracker(),
		errBudget: NewTokenBucket(budget.TimePerTokenMicro(), budget.Budget()),
	}
}

func (t *HealthTracker) Healthy() bool {
	return !t.rateLimit.Limited() && t.errBudget.HasTokens()
}

func (t *HealthTracker) TrackErr(err error) {
	var rateLimitErr *clients.RateLimitError

	if errors.As(err, &rateLimitErr) {
		t.rateLimit.SetLimited(rateLimitErr.UntilReset())

		return
	}

	_ = t.errBudget.Take(1)
}
