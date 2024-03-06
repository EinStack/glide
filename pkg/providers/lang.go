package providers

import (
	"context"
	"errors"
	"time"

	"glide/pkg/api/schemas"
	"glide/pkg/providers/clients"
	"glide/pkg/routers/health"
	"glide/pkg/routers/latency"
)

// LangProvider defines an interface a provider should fulfill to be able to serve language chat requests
type LangProvider interface {
	ModelProvider

	SupportChatStream() bool

	Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error)
	ChatStream(ctx context.Context, request *schemas.ChatRequest, responseC chan<- schemas.ChatResponse) error
}

type LangModel interface {
	LangProvider
	Model
}

// LanguageModel wraps provider client and expend it with health & latency tracking
//
//	The model health is assumed to be independent of model actions (e.g. chat & chatStream)
//	The latency is assumed to be action-specific (e.g. streaming chat chunks are much low latency than the full chat action)
type LanguageModel struct {
	modelID               string
	weight                int
	client                LangProvider
	rateLimit             *health.RateLimitTracker
	errBudget             *health.TokenBucket
	chatLatency           *latency.MovingAverage
	chatStreamLatency     *latency.MovingAverage
	latencyUpdateInterval *time.Duration
}

func NewLangModel(modelID string, client LangProvider, budget health.ErrorBudget, latencyConfig latency.Config, weight int) *LanguageModel {
	return &LanguageModel{
		modelID:               modelID,
		client:                client,
		rateLimit:             health.NewRateLimitTracker(),
		errBudget:             health.NewTokenBucket(budget.TimePerTokenMicro(), budget.Budget()),
		chatLatency:           latency.NewMovingAverage(latencyConfig.Decay, latencyConfig.WarmupSamples),
		chatStreamLatency:     latency.NewMovingAverage(latencyConfig.Decay, latencyConfig.WarmupSamples),
		latencyUpdateInterval: latencyConfig.UpdateInterval,
		weight:                weight,
	}
}

func (m LanguageModel) ID() string {
	return m.modelID
}

func (m LanguageModel) Weight() int {
	return m.weight
}

func (m LanguageModel) LatencyUpdateInterval() *time.Duration {
	return m.latencyUpdateInterval
}

func (m LanguageModel) ChatLatency() *latency.MovingAverage {
	return m.chatLatency
}

func (m LanguageModel) ChatStreamLatency() *latency.MovingAverage {
	return m.chatStreamLatency
}

func (m LanguageModel) Healthy() bool {
	return !m.rateLimit.Limited() && m.errBudget.HasTokens()
}

func (m *LanguageModel) Chat(ctx context.Context, request *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	startedAt := time.Now()
	resp, err := m.client.Chat(ctx, request)

	if err == nil {
		// record latency per token to normalize measurements
		m.chatLatency.Add(float64(time.Since(startedAt)) / resp.ModelResponse.TokenUsage.ResponseTokens)

		// successful response
		resp.ModelID = m.modelID

		return resp, err
	}

	var rateLimitErr *clients.RateLimitError

	if errors.As(err, &rateLimitErr) {
		m.rateLimit.SetLimited(rateLimitErr.UntilReset())

		return resp, err
	}

	_ = m.errBudget.Take(1)

	return resp, err
}

func (m *LanguageModel) ChatStream(ctx context.Context, request *schemas.ChatRequest, responseC chan<- schemas.ChatResponse) error {
	err := m.client.ChatStream(ctx, request, responseC)

	if err == nil {
		return err
	}

	var rateLimitErr *clients.RateLimitError

	if errors.As(err, &rateLimitErr) {
		m.rateLimit.SetLimited(rateLimitErr.UntilReset())

		return err
	}

	_ = m.errBudget.Take(1)

	return err
}

func (m *LanguageModel) SupportChatStream() bool {
	return m.client.SupportChatStream()
}

func (m *LanguageModel) Provider() string {
	return m.client.Provider()
}

func ChatLatency(model Model) *latency.MovingAverage {
	return model.(LanguageModel).chatLatency
}

func ChatStreamLatency(model Model) *latency.MovingAverage {
	return model.(LanguageModel).chatStreamLatency
}
