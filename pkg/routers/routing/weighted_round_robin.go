package routing

import (
	"sync"

	"glide/pkg/providers"
)

const (
	WeightedRoundRobin Strategy = "weighed-round-robin"
)

type Weighter struct {
	model         providers.Model
	currentWeight int
}

func (w *Weighter) Current() int {
	return w.currentWeight
}

func (w *Weighter) Weight() int {
	return w.model.Weight()
}

func (w *Weighter) Incr() {
	w.currentWeight += w.Weight()
}

func (w *Weighter) Decr(totalWeight int) {
	w.currentWeight -= totalWeight
}

type WRoundRobinRouting struct {
	mu      sync.Mutex
	weights []*Weighter
}

func NewWeightedRoundRobin(models []providers.Model) *WRoundRobinRouting {
	weights := make([]*Weighter, 0, len(models))

	for _, model := range models {
		weights = append(weights, &Weighter{
			model:         model,
			currentWeight: 0,
		})
	}

	return &WRoundRobinRouting{
		weights: weights,
	}
}

func (r *WRoundRobinRouting) Iterator() LangModelIterator {
	return r
}

func (r *WRoundRobinRouting) Next() (providers.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	totalWeight := 0

	var maxWeighter *Weighter

	for _, weighter := range r.weights {
		if !weighter.model.Healthy() {
			continue
		}

		weighter.Incr()
		totalWeight += weighter.Weight()

		if maxWeighter == nil {
			maxWeighter = weighter
			continue
		}

		if weighter.Current() > maxWeighter.Current() {
			maxWeighter = weighter
		}
	}

	if maxWeighter != nil {
		maxWeighter.Decr(totalWeight)

		return maxWeighter.model, nil
	}

	return nil, ErrNoHealthyModels
}
