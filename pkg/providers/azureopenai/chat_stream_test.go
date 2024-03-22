package azureopenai

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

	"glide/pkg/providers/clients"
	"glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
)

func TestAzureOpenAIClient_ChatStreamSupported(t *testing.T) {
	providerCfg := DefaultConfig()
	clientCfg := clients.DefaultClientConfig()

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	require.True(t, client.SupportChatStream())
}

func TestAzureOpenAIClient_ChatStreamRequest(t *testing.T) {
	tests := map[string]string{
		"success stream": "./testdata/chat_stream.success.txt",
	}

	for name, streamFile := range tests {
		t.Run(name, func(t *testing.T) {
			AzureOpenAIMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawPayload, _ := io.ReadAll(r.Body)

				var data interface{}
				// Parse the JSON body
				err := json.Unmarshal(rawPayload, &data)
				if err != nil {
					t.Errorf("error decoding payload (%q): %v", string(rawPayload), err)
				}

				chatResponse, err := os.ReadFile(filepath.Clean(streamFile))
				if err != nil {
					t.Errorf("error reading azureopenai chat mock response: %v", err)
				}

				w.Header().Set("Content-Type", "text/event-stream")

				_, err = w.Write(chatResponse)
				if err != nil {
					t.Errorf("error on sending chat response: %v", err)
				}
			})

			AzureopenAIServer := httptest.NewServer(AzureOpenAIMock)
			defer AzureopenAIServer.Close()

			ctx := context.Background()
			providerCfg := DefaultConfig()
			clientCfg := clients.DefaultClientConfig()

			providerCfg.BaseURL = AzureopenAIServer.URL

			client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
			require.NoError(t, err)

			req := schemas.NewChatStreamFromStr("What's the capital of the United Kingdom?")
			stream, err := client.ChatStream(ctx, req)
			require.NoError(t, err)

			err = stream.Open()
			require.NoError(t, err)

			for {
				chunk, err := stream.Recv()

				if err == io.EOF {
					return
				}

				require.NoError(t, err)
				require.NotNil(t, chunk)
			}
		})
	}
}

func TestAzureOpenAIClient_ChatStreamRequestInterrupted(t *testing.T) {
	tests := map[string]string{
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

			req := schemas.NewChatStreamFromStr("What's the capital of the United Kingdom?")

			stream, err := client.ChatStream(ctx, req)
			require.NoError(t, err)

			err = stream.Open()
			require.NoError(t, err)

			for {
				chunk, err := stream.Recv()
				if err != nil {
					require.ErrorIs(t, err, clients.ErrProviderUnavailable)
					return
				}

				require.NoError(t, err)
				require.NotNil(t, chunk)
			}
		})
	}
}
