package providers

import (
	"context"

	"glide/pkg/routers/latency"

	"glide/pkg/api/schemas"
)

type ResponseMock struct {
	Msg string
	Err *error
}

func (m *ResponseMock) Resp() *schemas.UnifiedChatResponse {
	return &schemas.UnifiedChatResponse{
		ID: "rsp0001",
		ModelResponse: schemas.ProviderResponse{
			ResponseID: map[string]string{
				"ID": "0001",
			},
			Message: schemas.ChatMessage{
				Content: m.Msg,
			},
		},
	}
}

type ProviderMock struct {
	idx       int
	responses []ResponseMock
}

func NewProviderMock(responses []ResponseMock) *ProviderMock {
	return &ProviderMock{
		idx:       0,
		responses: responses,
	}
}

func (c *ProviderMock) Chat(_ context.Context, _ *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	response := c.responses[c.idx]
	c.idx++

	if response.Err != nil {
		return nil, *response.Err
	}

	return response.Resp(), nil
}

func (c *ProviderMock) Provider() string {
	return "provider_mock"
}

type LangModelMock struct {
	modelID string
	healthy bool
	latency *latency.MovingAverage
}

func NewLangModelMock(ID string, healthy bool, avgLatency float64) *LangModelMock {
	movingAverage := latency.NewMovingAverage(0.06, 3)
	movingAverage.Set(avgLatency)

	return &LangModelMock{
		modelID: ID,
		healthy: healthy,
		latency: movingAverage,
	}
}

func (m *LangModelMock) ID() string {
	return m.modelID
}

func (m *LangModelMock) Healthy() bool {
	return m.healthy
}

func (m *LangModelMock) Latency() *latency.MovingAverage {
	return m.latency
}
