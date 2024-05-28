package schemas

// ChatRequest defines Glide's Chat Request Schema unified across all language models
type ChatRequest struct {
	Message        ChatMessage                     `json:"message" validate:"required"`
	MessageHistory []ChatMessage                   `json:"message_history"`
	OverrideParams *map[string]ModelParamsOverride `json:"override_params,omitempty"`
}

func (r *ChatRequest) ModelParams(modelNameOrID string) *ModelParamsOverride {
	if r.OverrideParams == nil {
		return nil
	}

	if override, found := (*r.OverrideParams)[modelNameOrID]; found {
		return &override
	}

	return nil
}

// ModelParamsOverride allows to redefine chat message and model params based on the model ID
//
//	Glide provides an abstraction around concreate models and this is a way to be able to provide model-specific params if needed.
//	The override is going to be applied if Glide picks the referenced there (it may pick another model to serve a given request)
type ModelParamsOverride struct {
	// TODO: should be just string?
	Message ChatMessage `json:"message,omitempty"`
	// TODO(185): Add an ability to override model params
}

// ChatParams represents a chat request params that overrides the default model params from configs
type ChatParams struct {
	Messages []ChatMessage
	// TODO(185): set other params
}

// ChatParams returns a specific chat request params account for model-specific overrides.
func (r *ChatRequest) ChatParams(modelID string, modelName string) *ChatParams {
	params := &ChatParams{
		Messages: make([]ChatMessage, 0, len(r.MessageHistory)+1),
	}

	reqMessage := r.Message

	if override := r.ModelParams(modelName); override != nil {
		reqMessage = override.Message
		// TODO(185): set other params
	}

	if override := r.ModelParams(modelID); override != nil {
		reqMessage = override.Message
		// TODO(185): set other params
	}

	params.Messages = append(params.Messages, r.MessageHistory...)
	params.Messages = append(params.Messages, reqMessage)

	return params
}

func NewChatFromStr(message string) *ChatRequest {
	return &ChatRequest{
		Message: ChatMessage{
			"user",
			message,
		},
	}
}

// ChatResponse defines Glide's Chat Response Schema unified across all language models
type ChatResponse struct {
	ID            string        `json:"id,omitempty"`
	Created       int           `json:"created_at,omitempty"`
	Provider      string        `json:"provider_id,omitempty"`
	RouterID      string        `json:"router_id,omitempty"`
	ModelID       string        `json:"model_id,omitempty"`
	ModelName     string        `json:"model_name,omitempty"`
	Cached        bool          `json:"cached,omitempty"`
	ModelResponse ModelResponse `json:"model_response,omitempty"`
}

// ModelResponse is the unified response from the provider.
type ModelResponse struct {
	Metadata   map[string]string `json:"metadata,omitempty"`
	Message    ChatMessage       `json:"message"`
	TokenUsage TokenUsage        `json:"token_usage"`
}

type TokenUsage struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct {
	// The role of the author of this message. One of system, user, or assistant.
	Role string `json:"role" validate:"required"`
	// The content of the message.
	Content string `json:"content" validate:"required"`
}
