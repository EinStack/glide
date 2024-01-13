package routing

import "glide/pkg/providers"

const (
	LeastLatency Strategy = "least_latency"
)

// LeastLatencyRouting routes requests to the model that responses the fastest
// At the beginning, we try to send requests to all models to find out the quickest one.
// After that, we use the that model for some time. But we don't want to stick to that model forever (as some
// other model latency may improve over time overperform the best one),
// so we need to send some traffic to other models from time to time to update their latency stats
type LeastLatencyRouting struct {
	models []*providers.LangModel
}

func NewLeastLatencyRouting(models []*providers.LangModel) *LeastLatencyRouting {
	return &LeastLatencyRouting{
		models: models,
	}
}

func (r *LeastLatencyRouting) Iterator() LangModelIterator {
	return r
}

func (r LeastLatencyRouting) Next() (*providers.LangModel, error) {
	models := r.models

	for _, model := range models {
		if !model.Healthy() {
			// cannot do much with unavailable model
			continue
		}

		return model, nil
	}

	// responseTime := time.Since(startTime)
	//	h.avgResponseTime = lb.alpha*float64(responseTime/time.Millisecond) + (1.0-lb.alpha)*h.avgResponseTime

	//return nil, ErrNoHealthyModels
}
