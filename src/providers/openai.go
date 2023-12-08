package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Define your Go struct types here

	ChatGptRequest struct {
		Prompt     string `json:"prompt"`
		MaxTokens  uint   `json:"max_tokens"`
	}

	ChatGptResponse struct {
		ID        string  `json:"id"`
		Object    string  `json:"object"`
		Created   uint64  `json:"created"`
		Model     string  `json:"model"`
		Choices   []Choice `json:"choices"`
		Usage     []Usage   `json:"usage"`
	}

	GptResponse struct {
		Data []Choice `json:"data"`
	}

	Choice struct {
		Index         uint    `json:"index"`
		Message       Message `json:"message"`
		FinishReason  string  `json:"finish_reason"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Usage struct {
		PromptTokens      uint32 `json:"prompt_tokens"`
		CompletionTokens  uint32 `json:"completion_tokens"`
		TotalTokens       uint32 `json:"total_tokens"`
	}
)

// Function to interact with ChatGPT
func chatWithGPT(input string) (map[string]interface{}, error) {
	fmt.Printf("input: %s\n", input)

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	// Set your OpenAI API key
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_KEY not set")
	}

	fmt.Print("Running OpenAI Chat")

	// Set up the HTTP client
	client := http.Client{}

	fmt.Printf("Request Payload: %s\n", input)

	// Make the API request
	openaiEndpoint := os.Getenv("OPENAI_ENDPOINT")
	if openaiEndpoint == "" {
		return nil, fmt.Errorf("OPENAI_ENDPOINT not set")
	}

	reqBody, err := json.Marshal(map[string]string{"input": input})
	if err != nil {
		return nil, fmt.Errorf("Error encoding request body: %w", err)
	}

	req, err := http.NewRequest("POST", openaiEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making API request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	fmt.Printf("OpenAI Response: %s\n", body)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: %w", err)
	}

	return response, nil
}

func main() {
	// Example usage
	input := "Hello, ChatGPT!"
	response, err := chatWithGPT(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %+v\n", response)
}
