package routing

import (
	"glide/pkg/routers/health"
)

const (
	Priority Strategy = "priority"
)

// PriorityRouting routes request to the first healthy model defined in the routing config
//
//	Priority of models are defined as position of the model on the list
//	(e.g. the first model definition has the highest priority, then the second model definition and so on)
type PriorityRouting struct {
	models *[]health.LangModelHealthTracker
}

func NewPriorityRouting(models *[]health.LangModelHealthTracker) *PriorityRouting {
	return &PriorityRouting{
		models: models,
	}
}

func (r *PriorityRouting) Iterator() LangModelIterator {
	iterator := PriorityIterator{
		idx:    0,
		models: r.models,
	}

	return iterator
}

type PriorityIterator struct {
	idx    int
	models *[]health.LangModelHealthTracker
}

func (r PriorityIterator) Next() (*health.LangModelHealthTracker, error) {
	models := *r.models

	for r.idx < len(models) {
		model := models[r.idx]

		r.idx++

		if model.Healthy() {
			return &model, nil
		}

		// otherwise, try to pick the next model on the list
	}

	return nil, ErrNoHealthyModels
}
