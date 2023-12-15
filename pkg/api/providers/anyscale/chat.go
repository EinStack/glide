package anyscale

type AnyscaleProviderConfig struct {
	Model            string           `json:"model" validate:"required,lowercase"`
	Messages         string           `json:"messages" validate:"required"` // does this need to be updated to []string?
	MaxTokens        int              `json:"max_tokens,omitempty" validate:"omitempty,gte=0"`
	Temperature      int              `json:"temperature,omitempty" validate:"omitempty,gte=0,lte=2"`
	TopP             int              `json:"top_p,omitempty" validate:"omitempty,gte=0,lte=1"`
	N                int              `json:"n,omitempty" validate:"omitempty,gte=1"`
	Stream           bool             `json:"stream,omitempty" validate:"omitempty,boolean"`
	Stop             interface{}      `json:"stop,omitempty" validate:"omitempty"`
	PresencePenalty  int              `json:"presence_penalty,omitempty" validate:"omitempty,gte=-2,lte=2"`
	FrequencyPenalty int              `json:"frequency_penalty,omitempty" validate:"omitempty,gte=-2,lte=2"`
	LogitBias        *map[int]float64 `json:"logit_bias,omitempty" validate:"omitempty"`
	User             interface{}      `json:"user,omitempty" validate:"omitempty"`
	Seed             interface{}      `json:"seed,omitempty" validate:"omitempty,gte=0"`
	Tools            []string         `json:"tools,omitempty" validate:"omitempty"`
	ToolChoice       interface{}      `json:"tool_choice,omitempty" validate:"omitempty"`
	ResponseFormat   interface{}      `json:"response_format,omitempty" validate:"omitempty"`
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
func AnyscaleChatDefaultConfig() AnyscaleProviderConfig {
	return AnyscaleProviderConfig{
		Model:            "meta-llama/Llama-2-7b-chat-hf",
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