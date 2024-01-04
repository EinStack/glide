package octoml

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

	"glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
)

func TestOctoMLClient_ChatRequest(t *testing.T) {
	// OctoML Chat API: https://docs.octoai.cloud/docs/text-gen-api-docs
	octoMLMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPayload, _ := io.ReadAll(r.Body)

		var data interface{}
		// Parse the JSON body
		err := json.Unmarshal(rawPayload, &data)
		if err != nil {
			t.Errorf("error decoding payload (%q): %v", string(rawPayload), err)
		}

		chatResponse, err := os.ReadFile(filepath.Clean("./testdata/chat.success.json"))
		if err != nil {
			t.Errorf("error reading octoml chat mock response: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(chatResponse)

		if err != nil {
			t.Errorf("error on sending chat response: %v", err)
		}
	})

	octoMLServer := httptest.NewServer(octoMLMock)
	defer octoMLServer.Close()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.BaseURL = octoMLServer.URL

	client, err := NewClient(cfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "human",
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, cfg.Model, response.Model)
	require.Equal(t, "cmpl-8ea213aece0747aca6d0608b02b57196", response.ID)
}

func TestOctoMLClient_Chat_Error(t *testing.T) {
	// Set up the test case
	// Create a mock API server that returns an error
	octoMLMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return an error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	// Create a mock API server
	octoMLServer := httptest.NewServer(octoMLMock)
	defer octoMLServer.Close()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.BaseURL = octoMLServer.URL

	// Create a new client with the mock API server
	client, err := NewClient(cfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	// Create a chat request
	request := schemas.UnifiedChatRequest{
		Message: schemas.ChatMessage{
			Role:    "human",
			Content: "What's the biggest animal?",
		},
	}

	// Call the Chat function
	_, err = client.Chat(ctx, &request)

	// Check the error
	require.Error(t, err)
	require.Contains(t, err.Error(), "provider is not available")
}

func TestDoChatRequest_ErrorResponse(t *testing.T) {
	// Create a mock HTTP server that returns a non-OK status code
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer mockServer.Close()

	// Create a new client with the mock server URL
	client := &Client{
		httpClient: http.DefaultClient,
		chatURL:    mockServer.URL,
		config:     &Config{APIKey: "dummy_key"},
		telemetry:  telemetry.NewTelemetryMock(),
	}

	// Create a chat request payload
	payload := &ChatRequest{
		Model:    "dummy_model",
		Messages: []ChatMessage{{Role: "human", Content: "Hello"}},
	}

	// Call the doChatRequest function
	_, err := client.doChatRequest(context.Background(), payload)

	require.Error(t, err)
	require.Contains(t, err.Error(), "provider is not available")
}
