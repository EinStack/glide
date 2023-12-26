package openai

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"glide/pkg/providers"

	"github.com/joho/godotenv"
)

func TestOpenAIClient(t *testing.T) {
	// Initialize the OpenAI client

	err := godotenv.Load("/Users/max/code/Glide/.env")
	if err != nil {
		fmt.Println("Error loading .env file" + err.Error())
	}

	_ = t

	fakeData := providers.UnifiedAPIData{
		Provider: "openai",
		Model:    "gpt-3.5-turbo",
		APIKey:   os.Getenv("openai"),
		Params: map[string]interface{}{
			"temperature": 0.3,
		},
		Message: map[string]string{
			"role":    "user",
			"content": "Where was it played?",
		},
		MessageHistory: []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": "Who won the world series in 2020?",
			},
			{
				"role":    "assistant",
				"content": "The Los Angeles Dodgers won the World Series in 2020.",
			},
		},
	}

	c, err := Client(fakeData)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	resp, _ := c.Chat()

	respJSON, _ := json.Marshal(resp)

	fmt.Println(string(respJSON))
}
