package openai

type OpenAiProviderConfig struct {
	Model            string           `json:"model" validate:"required,lowercase"`
	Messages         string           `json:"messages" validate:"required"` // does this need to be updated to []string?
	MaxTokens        int              `json:"max_tokens" validate:"omitempty,gte=0"`
	Temperature      int              `json:"temperature" validate:"omitempty,gte=0,lte=2"`
	TopP             int              `json:"top_p" validate:"omitempty,gte=0,lte=1"`
	N                int              `json:"n" validate:"omitempty,gte=1"`
	Stream           bool             `json:"stream" validate:"omitempty, boolean"`
	Stop             interface{}      `json:"stop"`
	PresencePenalty  int              `json:"presence_penalty" validate:"omitempty,gte=-2,lte=2"`
	FrequencyPenalty int              `json:"frequency_penalty" validate:"omitempty,gte=-2,lte=2"`
	LogitBias        *map[int]float64 `json:"logit_bias" validate:"omitempty"`
	User             interface{}      `json:"user"`
	Seed             interface{}      `json:"seed" validate:"omitempty,gte=0"`
	Tools            []string         `json:"tools"`
	ToolChoice       interface{}      `json:"tool_choice"`
	ResponseFormat   interface{}      `json:"response_format"`
}

var defaultMessage = `[
	{
	  "role": "system",
	  "content": "You are a helpful assistant."
	},
	{
	  "role": "user",
	  "content": "Hello!"
	}
  ]`

// Provide the request body for OpenAI's ChatCompletion API
func OpenAiChatDefaultConfig() OpenAiProviderConfig {
	return OpenAiProviderConfig{
		Model:            "gpt-3.5-turbo",
		Messages:         defaultMessage,
		MaxTokens:        100,
		Temperature:      1,
		TopP:             1,
		N:                1,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             nil,
		Seed:             nil,
		Tools:            nil,
		ToolChoice:       nil,
		ResponseFormat:   nil,
	}
}
