package routing

import (
	"testing"

	"github.com/stretchr/testify/require"
	"glide/pkg/providers"
)

func TestPriorityRouting_PickModelsInOrder(t *testing.T) {
	type Model struct {
		modelID string
		healthy bool
	}

	type TestCase struct {
		models           []Model
		expectedModelIDs []string
	}

	tests := map[string]TestCase{
		"all healthy":         {[]Model{{"first", true}, {"second", true}, {"third", true}}, []string{"first", "first", "first"}},
		"first unhealthy":     {[]Model{{"first", false}, {"second", true}, {"third", true}}, []string{"second", "second", "second"}},
		"first two unhealthy": {[]Model{{"first", false}, {"second", false}, {"third", true}}, []string{"third", "third", "third"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			models := make([]providers.Model, 0, len(tc.models))

			for _, model := range tc.models {
				models = append(models, providers.NewLangModelMock(model.modelID, model.healthy, 1))
			}

			routing := NewPriority(models)
			iterator := routing.Iterator()

			// loop three times over the whole pool to check if we return back to the begging of the list
			for _, modelID := range tc.expectedModelIDs {
				model, err := iterator.Next()
				require.NoError(t, err)
				require.Equal(t, modelID, model.ID())
			}
		})
	}
}

func TestPriorityRouting_NoHealthyModels(t *testing.T) {
	models := []providers.Model{
		providers.NewLangModelMock("first", false, 1),
		providers.NewLangModelMock("second", false, 1),
		providers.NewLangModelMock("third", false, 1),
	}

	routing := NewPriority(models)
	iterator := routing.Iterator()

	_, err := iterator.Next()
	require.Error(t, err)
}
