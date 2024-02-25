package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"

	"glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
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
	providerCfg := DefaultConfig()
	clientCfg := clients.DefaultClientConfig()

	providerCfg.BaseURL = openAIServer.URL

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.ChatRequest{Message: schemas.ChatMessage{
		Role:    "user",
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, "chatcmpl-123", response.ID)
}

func TestOpenAIClient_RateLimit(t *testing.T) {
	openAIMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", strconv.FormatInt(int64(5*time.Minute), 10))
		w.WriteHeader(429)
	})

	openAIServer := httptest.NewServer(openAIMock)
	defer openAIServer.Close()

	ctx := context.Background()
	providerCfg := DefaultConfig()
	clientCfg := clients.DefaultClientConfig()

	providerCfg.BaseURL = openAIServer.URL

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.ChatRequest{Message: schemas.ChatMessage{
		Role:    "user",
		Content: "What's the biggest animal?",
	}}

	_, err = client.Chat(ctx, &request)

	require.Error(t, err)
	require.ErrorIs(t, err, clients.RateLimitError{})
}
