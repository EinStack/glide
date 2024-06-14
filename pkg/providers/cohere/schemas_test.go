package cohere

import (
	"testing"

	"github.com/EinStack/glide/pkg/api/schemas"
	"github.com/stretchr/testify/require"
)

func TestChatRequest_ApplyParams(t *testing.T) {
	tests := []struct {
		name     string
		chatReq  ChatRequest
		params   *schemas.ChatParams
		expected ChatRequest
	}{
		{
			name:    "should set role to default USER when role is empty string",
			chatReq: ChatRequest{},
			params: &schemas.ChatParams{
				Messages: []schemas.ChatMessage{
					{Role: "", Content: "Hello"},
					{Role: schemas.RoleAssistant, Content: "Hi there!"},
				},
			},
			expected: ChatRequest{
				Message: "Hi there!",
				ChatHistory: []schemas.ChatMessage{
					{Role: "USER", Content: "Hello"},
				},
			},
		},
		{
			name:    "should set role to default USER when role is RoleUser",
			chatReq: ChatRequest{},
			params: &schemas.ChatParams{
				Messages: []schemas.ChatMessage{
					{Role: schemas.RoleUser, Content: "Hello"},
					{Role: schemas.RoleAssistant, Content: "Hi there!"},
				},
			},
			expected: ChatRequest{
				Message: "Hi there!",
				ChatHistory: []schemas.ChatMessage{
					{Role: "USER", Content: "Hello"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.chatReq.ApplyParams(tt.params)
			require.Equal(t, tt.expected, tt.chatReq)
		})
	}
}
