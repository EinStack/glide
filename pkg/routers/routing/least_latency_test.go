package routing

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"glide/pkg/providers"
)

func TestLeastLatencyRouting_Warmup(t *testing.T) {
	type Model struct {
		modelID string
		healthy bool
		latency float64
	}

	type TestCase struct {
		models           []Model
		expectedModelIDs []string
	}

	tests := map[string]TestCase{
		"all cold models":             {[]Model{{"first", true, 0.0}, {"second", true, 0.0}, {"third", true, 0.0}}, []string{"first", "second", "third"}},
		"all cold models & unhealthy": {[]Model{{"first", true, 0.0}, {"second", false, 0.0}, {"third", true, 0.0}}, []string{"first", "third", "first"}},
		"some models are warmed":      {[]Model{{"first", true, 100.0}, {"second", true, 0.0}, {"third", true, 120.0}}, []string{"second", "second", "second"}},
		"cold unhealthy model":        {[]Model{{"first", true, 120.0}, {"second", false, 0.0}, {"third", true, 100.0}}, []string{"third", "third", "third"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			models := make([]providers.Model, 0, len(tc.models))

			for _, model := range tc.models {
				models = append(models, providers.NewLangModelMock(model.modelID, model.healthy, model.latency))
			}

			routing := NewLeastLatencyRouting(models)
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

func TestLeastLatencyRouting_NoHealthyModels(t *testing.T) {
	tests := map[string][]float64{
		"all cold models unhealthy":    {0.0, 0.0, 0.0},
		"all warm models unhealthy":    {100.0, 120.0, 150.0},
		"cold & warm models unhealthy": {0.0, 120.0, 150.0},
	}

	for name, latencies := range tests {
		t.Run(name, func(t *testing.T) {
			models := make([]providers.Model, 0, len(latencies))

			for idx, latency := range latencies {
				models = append(models, providers.NewLangModelMock(strconv.Itoa(idx), false, latency))
			}

			routing := NewPriorityRouting(models)
			iterator := routing.Iterator()

			_, err := iterator.Next()
			require.Error(t, err)
		})
	}
}
