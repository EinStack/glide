package openai

import (
	"context"
	"encoding/json"
	"fmt"
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

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "user", // should be ['system', 'assistant', 'user', 'function']
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, "chatcmpl-123", response.ID)
}

func TestOpenAIClient_StreamChatRequest(t *testing.T) {
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

	// providerCfg.BaseURL = openAIServer.URL

	providerCfg.APIKey = ""

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "user", // should be ['system', 'assistant', 'user', 'function']
		Content: "What's the biggest animal?",
	}}

	// Assume you have a function to call StreamChat and receive the response channel
	responseChannel, _ := client.StreamChat(ctx, &request)

	go func() {
		// Iterate over the response channel to print each response
		for chatResponse := range responseChannel {
			fmt.Println("stream response:", chatResponse)
		}
		close(responseChannel)
	}()
}
