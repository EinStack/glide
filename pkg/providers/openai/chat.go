package openai

type ProviderConfig struct {
	Model            ConfigItem          `json:"model" validate:"required,lowercase"`
	Messages         ConfigItem          `json:"messages" validate:"required"`
	MaxTokens        ConfigItem   		 `json:"max_tokens" validate:"omitempty,gte=0"`
	Temperature      ConfigItem   		 `json:"temperature" validate:"omitempty,gte=0,lte=2"`
	TopP             ConfigItem   		 `json:"top_p" validate:"omitempty,gte=0,lte=1"`
	N                ConfigItem          `json:"n" validate:"omitempty,gte=1"`
	Stream           ConfigItem      	 `json:"stream" validate:"omitempty, boolean"`
	Stop             ConfigItem          `json:"stop"`
	PresencePenalty  ConfigItem   		 `json:"presence_penalty" validate:"omitempty,gte=-2,lte=2"`
	FrequencyPenalty ConfigItem   		 `json:"frequency_penalty" validate:"omitempty,gte=-2,lte=2"`
	LogitBias        ConfigItem  		 `json:"logit_bias" validate:"omitempty"`
	User             ConfigItem          `json:"user"`
	Seed             ConfigItem          `json:"seed" validate:"omitempty,gte=0"`
	Tools            ConfigItem          `json:"tools"`
	ToolChoice       ConfigItem          `json:"tool_choice"`
	ResponseFormat   ConfigItem          `json:"response_format"`
}

type ConfigItem struct {
	Param    string      `json:"param" validate:"required"`
	Default  interface{} `json:"default"`
}

// Provide the request body for OpenAI's ChatCompletion API
func OpenAiDefaultConfig() ProviderConfig {
	return ProviderConfig{
		Model: ConfigItem{
			Param:    "model",
			Default:  "gpt-3.5-turbo",
		},
		Messages: ConfigItem{
			Param:   "messages",
			Default: "",
		},
		MaxTokens: ConfigItem{
			Param:   "max_tokens",
			Default: 100,
		},
		Temperature: ConfigItem{
			Param:   "temperature",
			Default: 1,
		},
		TopP: ConfigItem{
			Param:   "top_p",
			Default: 1,
		},
		N: ConfigItem{
			Param:   "n",
			Default: 1,
		},
		Stream: ConfigItem{
			Param:   "stream",
			Default: false,
		},
		Stop: ConfigItem{
			Param: "stop",
		},
		PresencePenalty: ConfigItem{
			Param: "presence_penalty",
			Default: 0,
		},
		FrequencyPenalty: ConfigItem{
			Param: "frequency_penalty",
			Default: 0,
		},
		LogitBias: ConfigItem{
			Param: "logit_bias",
		},
		User: ConfigItem{
			Param: "user",
		},
		Seed: ConfigItem{
			Param: "seed",
		},
		Tools: ConfigItem{
			Param: "tools",
		},
		ToolChoice: ConfigItem{
			Param: "tool_choice",
		},
		ResponseFormat: ConfigItem{
			Param: "response_format",
		},
	}
}
