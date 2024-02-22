package anthropic

// Anthropic Chat Response
type ChatCompletion struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Model        string    `json:"model"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	StopReason   string    `json:"stop_reason"`
	StopSequence string    `json:"stop_sequence"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
