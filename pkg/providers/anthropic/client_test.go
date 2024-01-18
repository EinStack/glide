package anthropic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"

	"glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
)

func TestOpenAIClient_ChatRequest(t *testing.T) {
	// Anthropic Messages API: https://docs.anthropic.com/claude/reference/messages_post
	AnthropicMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	AnthropicServer := httptest.NewServer(AnthropicMock)
	defer AnthropicServer.Close()

	ctx := context.Background()
	providerCfg := DefaultConfig()
	clientCfg := clients.DefaultClientConfig()

	providerCfg.BaseURL = AnthropicServer.URL

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "human",
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, "msg_013Zva2CMHLNnXjNJJKqJ2EF", response.ID)
}
