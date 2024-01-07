package routing

import "glide/pkg/providers"

const (
	Priority Strategy = "priority"
)

// PriorityRouting routes request to the first healthy model defined in the routing config
//
//	Priority of models are defined as position of the model on the list
//	(e.g. the first model definition has the highest priority, then the second model definition and so on)
type PriorityRouting struct {
	models *[]providers.LanguageModel
}

func NewPriorityRouting(models *[]providers.LanguageModel) *PriorityRouting {
	return &PriorityRouting{
		models: models,
	}
}

func (r *PriorityRouting) Iterator() LangModelIterator {
	return PriorityIterator{
		idx:    0,
		models: r.models,
	}
}

type PriorityIterator struct {
	idx    int
	models *[]providers.LanguageModel
}

func (r *PriorityIterator) Next() (providers.LanguageModel, error) {
	models := *r.models

	for r.idx < len(models) {
		model := models[r.idx]

		//if model.Healthy() {
		//	// TODO: check if it's healthy
		//	return model, nil
		//}

		// otherwise, try to pick the next model on the list
		r.idx++
	}

	return nil, ErrNoHealthyModels
}
