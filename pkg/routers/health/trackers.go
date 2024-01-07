package health

import (
	"context"

	"glide/pkg/api/schemas"
	"glide/pkg/providers"
)

type HealthTracker interface {
	Healthy() bool
}

// LangModelHealthTracker decorates the LangModel struct to add health/availability tracking
type LangModelHealthTracker struct {
	Model       providers.LanguageModel
	rateLimit   *RateLimitTracker
	errorBudget *TokenBucket // TODO: centralize provider API health tracking in the registry
}

func NewLangModelHealthTracker(model providers.LanguageModel) *LangModelHealthTracker {
	return &LangModelHealthTracker{
		Model:       model,
		rateLimit:   NewRateLimitTracker(),
		errorBudget: NewTokenBucket(1, 10), // TODO: set from configs
	}
}

func (t *LangModelHealthTracker) Healthy() bool {
	return !t.rateLimit.Limited() && t.errorBudget.HasTokens()
}

func (t *LangModelHealthTracker) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	resp, err := t.Model.Chat(ctx, request)

	if err == nil {
		// successful response
		return resp, err
	}

	// TODO: track all availability issues

	return resp, err
}
