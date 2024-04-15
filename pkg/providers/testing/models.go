package testing

import (
	"time"

	"glide/pkg/config/fields"

	"glide/pkg/providers"
	"glide/pkg/routers/latency"
)

// LangModelMock
type LangModelMock struct {
	modelID     string
	healthy     bool
	chatLatency *latency.MovingAverage
	weight      int
}

func NewLangModelMock(ID string, healthy bool, avgLatency float64, weight int) LangModelMock {
	chatLatency := latency.NewMovingAverage(0.06, 3)

	if avgLatency > 0.0 {
		chatLatency.Set(avgLatency)
	}

	return LangModelMock{
		modelID:     ID,
		healthy:     healthy,
		chatLatency: chatLatency,
		weight:      weight,
	}
}

func (m LangModelMock) ID() string {
	return m.modelID
}

func (m LangModelMock) Healthy() bool {
	return m.healthy
}

func (m *LangModelMock) ChatLatency() *latency.MovingAverage {
	return m.chatLatency
}

func (m LangModelMock) LatencyUpdateInterval() *fields.Duration {
	updateInterval := 30 * time.Second

	return (*fields.Duration)(&updateInterval)
}

func (m LangModelMock) Weight() int {
	return m.weight
}

func ChatMockLatency(model providers.Model) *latency.MovingAverage {
	return model.(LangModelMock).chatLatency
}
