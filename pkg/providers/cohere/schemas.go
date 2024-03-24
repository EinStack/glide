package cohere

// Cohere Chat Response
type ChatCompletion struct {
	Text          string                 `json:"text"`
	GenerationID  string                 `json:"generation_id"`
	ResponseID    string                 `json:"response_id"`
	TokenCount    TokenCount             `json:"token_count"`
	Citations     []Citation             `json:"citations"`
	Documents     []Documents            `json:"documents"`
	SearchQueries []SearchQuery          `json:"search_queries"`
	SearchResults []SearchResults        `json:"search_results"`
	Meta          Meta                   `json:"meta"`
	ToolInputs    map[string]interface{} `json:"tool_inputs"`
}

type TokenCount struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
	BilledTokens   int `json:"billed_tokens"`
}

type Meta struct {
	APIVersion struct {
		Version string `json:"version"`
	} `json:"api_version"`
	BilledUnits struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"billed_units"`
}

type Citation struct {
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Text       string   `json:"text"`
	DocumentID []string `json:"document_id"`
}

type Documents struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"` // TODO: This needs to be updated
}

type SearchQuery struct {
	Text         string `json:"text"`
	GenerationID string `json:"generation_id"`
}

type SearchResults struct {
	SearchQuery []SearchQueryObject  `json:"search_query"`
	Connectors  []ConnectorsResponse `json:"connectors"`
	DocumentID  []string             `json:"documentId"`
}

type SearchQueryObject struct {
	Text         string `json:"text"`
	GenerationID string `json:"generationId"`
}

type ConnectorsResponse struct {
	ID              string            `json:"id"`
	UserAccessToken string            `json:"user_access_token"`
	ContOnFail      string            `json:"continue_on_failure"`
	Options         map[string]string `json:"options"`
}
