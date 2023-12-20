package openai

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOpenAIClient(t *testing.T) {
	// Initialize the OpenAI client

	poolName := "default"
	modelName := "gpt-3.5-turbo"

	payload := map[string]interface{}{
		"message": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": "tell me a joke",
			},
		},
		"messageHistory": []string{"Hello there", "How are you?", "I'm good, how about you?"},
	}

	payloadBytes, _ := json.Marshal(payload)

	c := &Client{}

	resp, _ := c.Run(poolName, modelName, payloadBytes)

	respJSON, _ := json.Marshal(resp)

	fmt.Println(string(respJSON))
}
