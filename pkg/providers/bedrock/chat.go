package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	//"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is an Bedrock-specific request schema
type ChatRequest struct {
	Messages             string               `json:"inputText"`
	TextGenerationConfig TextGenerationConfig `json:"textGenerationConfig"`
}

type TextGenerationConfig struct {
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		TextGenerationConfig: TextGenerationConfig{
			MaxTokenCount: cfg.DefaultParams.MaxTokens,
			StopSequences: cfg.DefaultParams.StopSequence,
			Temperature:   cfg.DefaultParams.Temperature,
			TopP:          cfg.DefaultParams.TopP,
		},
	}
}

func NewChatMessagesFromUnifiedRequest(request *schemas.UnifiedChatRequest) string {
	// message history not yet supported for AWS models
	message := fmt.Sprintf("Role: %s, Content: %s", request.Message.Role, request.Message.Content)

	return message
}

// Chat sends a chat request to the specified bedrock model.
func (c *Client) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)

	chatResponse, err := c.doChatRequest(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.ModelResponse.Message.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

func (c *Client) createChatRequestSchema(request *schemas.UnifiedChatRequest) *ChatRequest {
	// TODO: consider using objectpool to optimize memory allocation
	chatRequest := c.chatRequestTemplate // hoping to get a copy of the template
	chatRequest.Messages = NewChatMessagesFromUnifiedRequest(request)

	return chatRequest
}

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.UnifiedChatResponse, error) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal chat request payload: %w", err)
	}

	cfg, _ := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.config.AccessKey, c.config.SecretKey, "")),
		config.WithRegion(c.config.AWSRegion),
	)

	client := bedrockruntime.NewFromConfig(cfg)

	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.config.Model),
		ContentType: aws.String("application/json"),
		Body:        rawPayload,
	})
	if err != nil {
		c.telemetry.Logger.Error("Error: Couldn't invoke model. Here's why: %v\n", zap.Error(err))
		return nil, err
	}

	var bedrockCompletion schemas.BedrockChatCompletion

	err = json.Unmarshal(result.Body, &bedrockCompletion)
	if err != nil {
		c.telemetry.Logger.Error("failed to parse bedrock chat response", zap.Error(err))
		return nil, err
	}

	response := schemas.UnifiedChatResponse{
		ID:       uuid.NewString(),
		Created:  int(time.Now().Unix()),
		Provider: "aws-bedrock",
		Model:    c.config.Model,
		Cached:   false,
		ModelResponse: schemas.ProviderResponse{
			SystemID: map[string]string{
				"system_fingerprint": "none",
			},
			Message: schemas.ChatMessage{
				Role:    "assistant",
				Content: bedrockCompletion.Results[0].OutputText,
				Name:    "",
			},
			TokenUsage: schemas.TokenUsage{
				PromptTokens:   float64(bedrockCompletion.Results[0].TokenCount),
				ResponseTokens: -1,
				TotalTokens:    float64(bedrockCompletion.Results[0].TokenCount),
			},
		},
	}

	return &response, nil
}
