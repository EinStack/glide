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
	Required bool        `json:"required" validate:"omitempty,boolean"`
	Default  interface{} `json:"default"`
}

// Provide the request body for OpenAI's ChatCompletion API
var OpenAiChatDefaultConfig = ProviderConfig {
		Model: ConfigItem{
			Param:    "model",
			Required: true,
			Default:  "gpt-3.5-turbo",
		},
		Messages: ConfigItem{
			Param:   "messages",
			Required: true,
			Default: "",
		},
		MaxTokens: ConfigItem{
			Param:   "max_tokens",
			Required: false,
			Default: 100,
		},
		Temperature: ConfigItem{
			Param:   "temperature",
			Required: false,
			Default: 1,
		},
		TopP: ConfigItem{
			Param:   "top_p",
			Required: false,
			Default: 1,
		},
		N: ConfigItem{
			Param:   "n",
			Required: false,
			Default: 1,
		},
		Stream: ConfigItem{
			Param:   "stream",
			Required: false,
			Default: false,
		},
		Stop: ConfigItem{
			Param: "stop",
			Required: false,
		},
		PresencePenalty: ConfigItem{
			Param: "presence_penalty",
			Required: false,
		},
		FrequencyPenalty: ConfigItem{
			Param: "frequency_penalty",
			Required: false,
		},
		LogitBias: ConfigItem{
			Param: "logit_bias",
			Required: false,
		},
		User: ConfigItem{
			Param: "user",
			Required: false,
		},
		Seed: ConfigItem{
			Param: "seed",
			Required: false,
		},
		Tools: ConfigItem{
			Param: "tools",
			Required: false,
		},
		ToolChoice: ConfigItem{
			Param: "tool_choice",
			Required: false,
		},
		ResponseFormat: ConfigItem{
			Param: "response_format",
			Required: false,
		},
	}
