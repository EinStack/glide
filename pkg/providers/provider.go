package providers

import (
	"context"
	"errors"

	"glide/pkg/providers/clients"
	"glide/pkg/routers/health"

	"glide/pkg/api/schemas"
)

// LangModelProvider defines an interface a provider should fulfill to be able to serve language chat requests
type LangModelProvider interface {
	Provider() string
	Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error)
}

type Model interface {
	ID() string
	Healthy() bool
	Weight() int
}

type LanguageModel interface {
	Model
	LangModelProvider
}

// LangModel
type LangModel struct {
	modelID     string
	weight      int
	client      LangModelProvider
	rateLimit   *health.RateLimitTracker
	errorBudget *health.TokenBucket // TODO: centralize provider API health tracking in the registry
}

func NewLangModel(modelID string, client LangModelProvider, budget health.ErrorBudget, weight int) *LangModel {
	return &LangModel{
		modelID:     modelID,
		weight:      weight,
		client:      client,
		rateLimit:   health.NewRateLimitTracker(),
		errorBudget: health.NewTokenBucket(budget.TimePerTokenMicro(), budget.Budget()),
	}
}

func (m *LangModel) ID() string {
	return m.modelID
}

func (m *LangModel) Provider() string {
	return m.client.Provider()
}

func (m *LangModel) Healthy() bool {
	return !m.rateLimit.Limited() && m.errorBudget.HasTokens()
}

func (m *LangModel) Weight() int {
	return m.weight
}

func (m *LangModel) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	resp, err := m.client.Chat(ctx, request)

	if err == nil {
		// successful response
		resp.ModelID = m.modelID

		return resp, err
	}

	var rle *clients.RateLimitError

	if errors.As(err, &rle) {
		m.rateLimit.SetLimited(rle.UntilReset())

		return resp, err
	}

	_ = m.errorBudget.Take(1)

	return resp, err
}
