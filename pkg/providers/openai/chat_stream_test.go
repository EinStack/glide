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
	"glide/pkg/providers/clients"
	"glide/pkg/telemetry"
)

func TestOpenAIClient_ChatStreamSupported(t *testing.T) {
	providerCfg := DefaultConfig()
	clientCfg := clients.DefaultClientConfig()

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	require.True(t, client.SupportChatStream())
}

func TestOpenAIClient_ChatStreamRequest(t *testing.T) {
	tests := map[string]string{
		"success stream": "./testdata/chat_stream.success.txt",
		"success stream, but no last done message": "./testdata/chat_stream.nodone.txt",
		"success stream, but with empty event":     "./testdata/chat_stream.empty.txt",
	}

	for name, streamFile := range tests {
		t.Run(name, func(t *testing.T) {
			openAIMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawPayload, _ := io.ReadAll(r.Body)

				var data interface{}
				// Parse the JSON body
				err := json.Unmarshal(rawPayload, &data)
				if err != nil {
					t.Errorf("error decoding payload (%q): %v", string(rawPayload), err)
				}

				chatResponse, err := os.ReadFile(filepath.Clean(streamFile))
				if err != nil {
					t.Errorf("error reading openai chat mock response: %v", err)
				}

				w.Header().Set("Content-Type", "text/event-stream")

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

			req := schemas.ChatRequest{Message: schemas.ChatMessage{
				Role:    "user",
				Content: "What's the capital of the United Kingdom?",
			}}

			resultC := client.ChatStream(ctx, &req)

			for chunkResult := range resultC {
				require.NoError(t, chunkResult.Error())
				require.NotNil(t, chunkResult.Chunk().ModelResponse.Message.Content)
			}
		})
	}
}
