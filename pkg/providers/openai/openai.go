package openai

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
)

type OpenAiClient struct {
	apiKey   string
	baseURL  string
	http     *http.Client
}

func NewOpenAiClient(apiKey string) *OpenAiClient {
	return &OpenAiClient{
		apiKey:   apiKey,
		baseURL:  "https://api.openai.com/v1",
		http:     http.DefaultClient,
	}
}

func (c *OpenAiClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *OpenAiClient) SetHTTPOpenAiClient(httpOpenAiClient *http.Client) {
	c.http = httpOpenAiClient
}

func (c *OpenAiClient) GetAPIKey() string {
	return c.apiKey
}

func (c *OpenAiClient) Get(endpoint string) (string, error) {
	// Implement the logic to make a GET request to the OpenAI API

	return "", nil
}

func (c *OpenAiClient) Post(endpoint string, payload []byte) (string, error) {
	// Implement the logic to make a POST request to the OpenAI API

	// Create the full URL
	url := c.baseURL + endpoint

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request using http Client
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// Add more methods to interact with OpenAI API

func main() {
	// Example usage of the OpenAI OpenAiClient
	OpenAiClient := NewOpenAiClient("YOUR_API_KEY")
	
	// Call methods on the OpenAiClient to interact with the OpenAI API
	// For example:
	payrload := []byte(`{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello!"}]}`)
	response, err := OpenAiClient.Post("/chat", payrload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Println("Response:", response)
}