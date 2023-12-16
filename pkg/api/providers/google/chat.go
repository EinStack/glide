package google

// Content represents the structure of the 'contents' array in the request body.
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents the structure of the 'parts' array in the 'contents' array.
type Part struct {
	Text string `json:"text"`
}

// SafetySetting represents the structure of the 'safetySettings' array in the request body.
type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GenerationConfig represents the structure of the 'generationConfig' object in the request body.
type GenerationConfig struct {
	StopSequences   []string `json:"stopSequences"`
	Temperature     float64  `json:"temperature"`
	MaxOutputTokens int      `json:"maxOutputTokens"`
	TopP            float64  `json:"topP"`
	TopK            int      `json:"topK"`
}

type GoogleProviderConfig struct {
	Model            string           `json:"model" validate:"required,lowercase"`
	Contents         []Content        `json:"contents"`
	SafetySettings   []SafetySetting  `json:"safetySettings"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

// Provide the request body for OpenAI's ChatCompletion API
func GoogleChatDefaultConfig() GoogleProviderConfig {
	return GoogleProviderConfig{
		Model: "gemini-pro",
		Contents: []Content{
			{
				Parts: []Part{
					{Text: "Write a story about a magic backpack."},
				},
			},
		},
		SafetySettings: []SafetySetting{
			{
				Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
				Threshold: "BLOCK_ONLY_HIGH",
			},
		},
		GenerationConfig: GenerationConfig{
			StopSequences:   []string{"Title"},
			Temperature:     1.0,
			MaxOutputTokens: 800,
			TopP:            0.8,
			TopK:            10,
		},
	}
}
