package strategy

import "glide/pkg/providers"

const (
	Priority RoutingStrategy = "priority"
)

// PriorityRouting routes request to the first healthy model defined in the routing config
//
//	Priority of models are defined as position of the model on the list
//	(e.g. the first model definition has the highest priority, then the second model definition and so on)
type PriorityRouting struct {
}

func NewPriorityRouting(models []providers.LanguageModel) *PriorityRouting {

	return &PriorityRouting{}
}

func (r *PriorityRouting) Next() (*providers.LanguageModel, error) {

}
