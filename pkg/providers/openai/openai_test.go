package openai

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
)

func TestOpenAIClient(t *testing.T) {
	// Initialize the OpenAI client

	_ = t



	c, err := Client(poolName, modelName, payloadBytes)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	resp, _ := c.Chat()

	respJSON, _ := json.Marshal(resp)

	fmt.Println(string(respJSON))
}
