package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"glide/pkg/api/schemas"

	"github.com/stretchr/testify/require"
	"glide/pkg/telemetry"
)

func TestOpenAIClient_ChatRequest(t *testing.T) {
	// OpenAI Chat API: https://platform.openai.com/docs/api-reference/chat/create
	openAIMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPayload, _ := io.ReadAll(r.Body)

		var data interface{}
		// Parse the JSON body
		err := json.Unmarshal(rawPayload, &data)

		if err != nil {
			t.Errorf("error decoding payload (%q): %v", string(rawPayload), err)
		}

		chatResponse, err := os.ReadFile(filepath.Clean("./testdata/chat.success.json"))
		if err != nil {
			t.Errorf("error reading openai chat mock response: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(chatResponse)

		if err != nil {
			t.Errorf("error on sending chat response: %v", err)
		}
	})

	openAIServer := httptest.NewServer(openAIMock)
	defer openAIServer.Close()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.BaseURL = openAIServer.URL

	client, err := NewClient(&cfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "human",
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, "chatcmpl-123", response.ID)
}
