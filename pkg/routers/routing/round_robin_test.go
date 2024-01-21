package routing

import (
	"testing"

	"github.com/stretchr/testify/require"
	"glide/pkg/providers"
)

func TestRoundRobinRouting_PickModelsSequentially(t *testing.T) {
	type Model struct {
		modelID string
		healthy bool
	}

	type TestCase struct {
		models           []Model
		expectedModelIDs []string
	}

	tests := map[string]TestCase{
		"all healthy":             {[]Model{{"first", true}, {"second", true}, {"third", true}}, []string{"first", "second", "third"}},
		"unhealthy in the middle": {[]Model{{"first", true}, {"second", false}, {"third", true}}, []string{"first", "third"}},
		"two unhealthy":           {[]Model{{"first", true}, {"second", false}, {"third", false}}, []string{"first"}},
		"first unhealthy":         {[]Model{{"first", false}, {"second", true}, {"third", true}}, []string{"second", "third"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			models := make([]providers.Model, 0, len(tc.models))

			for _, model := range tc.models {
				models = append(models, providers.NewLangModelMock(model.modelID, model.healthy, 100, 1))
			}

			routing := NewRoundRobinRouting(models)
			iterator := routing.Iterator()

			for i := 0; i < 3; i++ {
				// loop three times over the whole pool to check if we return back to the begging of the list
				for _, modelID := range tc.expectedModelIDs {
					model, err := iterator.Next()
					require.NoError(t, err)
					require.Equal(t, modelID, model.ID())
				}
			}
		})
	}
}

func TestRoundRobinRouting_NoHealthyModels(t *testing.T) {
	models := []providers.Model{
		providers.NewLangModelMock("first", false, 0, 1),
		providers.NewLangModelMock("second", false, 0, 1),
		providers.NewLangModelMock("third", false, 0, 1),
	}

	routing := NewRoundRobinRouting(models)
	iterator := routing.Iterator()

	_, err := iterator.Next()
	require.Error(t, err)
}
