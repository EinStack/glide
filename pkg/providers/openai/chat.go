package openai

type ProviderConfig struct {
	Model            ConfigItem          `json:"model"`
	Messages         ConfigItem          `json:"messages"`
	Functions        ConfigItem          `json:"functions"`
	FunctionCall     ConfigItem          `json:"function_call"`
	MaxTokens        NumericConfigItem   `json:"max_tokens"`
	Temperature      NumericConfigItem   `json:"temperature"`
	TopP             NumericConfigItem   `json:"top_p"`
	N                NumericConfigItem   `json:"n"`
	Stream           BoolConfigItem      `json:"stream"`
	Stop             ConfigItem          `json:"stop"`
	PresencePenalty  NumericConfigItem   `json:"presence_penalty"`
	FrequencyPenalty NumericConfigItem   `json:"frequency_penalty"`
	LogitBias        ConfigItem          `json:"logit_bias"`
	User             ConfigItem          `json:"user"`
	Seed             ConfigItem          `json:"seed"`
	Tools            ConfigItem          `json:"tools"`
	ToolChoice       ConfigItem          `json:"tool_choice"`
	ResponseFormat   ConfigItem          `json:"response_format"`
}

type ConfigItem struct {
	Param    string      `json:"param"`
	Required bool        `json:"required,omitempty"`
	Default  interface{} `json:"default,omitempty"`
	Min      interface{} `json:"min,omitempty"`
	Max      interface{} `json:"max,omitempty"`
}

type NumericConfigItem struct {
	Param   string    `json:"param"`
	Default float64   `json:"default,omitempty"`
	Min     float64   `json:"min,omitempty"`
	Max     float64   `json:"max,omitempty"`
}

type BoolConfigItem struct {
	Param   string `json:"param"`
	Default bool   `json:"default,omitempty"`
}

// DefaultProviderConfig returns an instance of ProviderConfig with default values.
func OpenAiDefaultConfig() ProviderConfig {
	return ProviderConfig{
		Model: ConfigItem{
			Param:    "model",
			Required: true,
			Default:  "gpt-3.5-turbo",
		},
		Messages: ConfigItem{
			Param:   "messages",
			Default: "",
		},
		Functions: ConfigItem{
			Param: "functions",
		},
		FunctionCall: ConfigItem{
			Param: "function_call",
		},
		MaxTokens: NumericConfigItem{
			Param:   "max_tokens",
			Default: 100,
			Min:     0,
		},
		Temperature: NumericConfigItem{
			Param:   "temperature",
			Default: 1,
			Min:     0,
			Max:     2,
		},
		TopP: NumericConfigItem{
			Param:   "top_p",
			Default: 1,
			Min:     0,
			Max:     1,
		},
		N: NumericConfigItem{
			Param:   "n",
			Default: 1,
		},
		Stream: BoolConfigItem{
			Param:   "stream",
			Default: false,
		},
		Stop: ConfigItem{
			Param: "stop",
		},
		PresencePenalty: NumericConfigItem{
			Param: "presence_penalty",
			Min:   -2,
			Max:   2,
		},
		FrequencyPenalty: NumericConfigItem{
			Param: "frequency_penalty",
			Min:   -2,
			Max:   2,
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
