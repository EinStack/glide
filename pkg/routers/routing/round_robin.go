package routing

import (
	"sync/atomic"

	"glide/pkg/providers"
)

const (
	RoundRobin Strategy = "round-robin"
)

// RoundRobinRouting routes request to the next model in the list in cycle
type RoundRobinRouting struct {
	idx    atomic.Uint64
	models []*providers.LangModel
}

func NewRoundRobinRouting(models []*providers.LangModel) *RoundRobinRouting {
	return &RoundRobinRouting{
		models: models,
	}
}

func (r *RoundRobinRouting) Iterator() LangModelIterator {
	return r
}

func (r *RoundRobinRouting) Next() (*providers.LangModel, error) {
	modelLen := len(r.models)

	// in order to avoid infinite loop in case of no healthy model is available,
	// we need to track whether we made a whole cycle around the model slice looking for a healthy model
	for i := 0; i < modelLen; i++ {
		idx := r.idx.Add(1)
		model := r.models[idx%uint64(modelLen)]

		if !model.Healthy() {
			continue
		}

		return model, nil
	}

	return nil, ErrNoHealthyModels
}
