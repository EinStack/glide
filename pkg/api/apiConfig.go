package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Provider string

const (
	OPENAI                                 Provider = "openai"
	ANTHROPIC                              Provider = "anthropic"
	COHERE                                 Provider = "cohere"
	AI21LABS                               Provider = "ai21labs"
	MLFLOW_MODEL_SERVING                   Provider = "mlflow-model-serving"
	MOSAICML                               Provider = "mosaicml"
	HUGGINGFACE_TEXT_GENERATION_INFERENCE  Provider = "huggingface-text-generation-inference"
	PALM                                   Provider = "palm"
	BEDROCK                                Provider = "bedrock"
)

type RouteType string

const (
	LLM_V1_COMPLETIONS  RouteType = "llm/v1/completions"
	LLM_V1_CHAT         RouteType = "llm/v1/chat"
	LLM_V1_EMBEDDINGS   RouteType = "llm/v1/embeddings"
)

type CohereConfig struct {
	CohereAPIKey string `json:"cohere_api_key" validate:"required"`
}

type AI21LabsConfig struct {
	AI21LabsAPIKey string `json:"ai21labs_api_key" validate:"required"`
}

type MosaicMLConfig struct {
	MosaicMLAPIKey string `json:"mosaicml_api_key" validate:"required"`
	MosaicMLAPIBase string `json:"mosaicml_api_base"`
}

type OpenAIAPIType string

const (
	OPENAI_API_TYPE  OpenAIAPIType = "openai"
	AZURE            OpenAIAPIType = "azure"
	AZUREAD          OpenAIAPIType = "azuread"
)

type OpenAIConfig struct {
	OpenAIAPIKey string `json:"openai_api_key" validate:"required"`
	OpenAIAPIType OpenAIAPIType `json:"openai_api_type"`
	OpenAIAPIBase string `json:"openai_api_base"`
}

type AnthropicConfig struct {
	AnthropicAPIKey string `json:"anthropic_api_key" validate:"required"`
}

type PaLMConfig struct {
	PalmAPIKey string `json:"palm_api_key" validate:"required"`
}

type MlflowModelServingConfig struct {
	ModelServerURL string `json:"model_server_url" validate:"required"`
}

type HuggingFaceTextGenerationInferenceConfig struct {
	HFServerURL string `json:"hf_server_url" validate:"required"`
}

type AWSBaseConfig struct {
	AWSRegion string `json:"aws_region"`
}

type AWSRole struct {
	AWSBaseConfig
	AWSRoleARN string `json:"aws_role_arn" validate:"required"`
	SessionLengthSeconds int `json:"session_length_seconds" validate:"required"`
}

type AWSIdAndKey struct {
	AWSBaseConfig
	AWSAccessKeyID string `json:"aws_access_key_id" validate:"required"`
	AWSSecretAccessKey string `json:"aws_secret_access_key" validate:"required"`
	AWSSessionToken string `json:"aws_session_token"`
}

type AWSBedrockConfig struct {
	AWSConfig interface{} `json:"aws_config" validate:"required"`
}

type ModelInfo struct {
	Name string `json:"name"`
	Provider Provider `json:"provider" validate:"required"`
}

type Model struct {
	Name string `json:"name"`
	Provider Provider `json:"provider" validate:"required"`
	Config interface{} `json:"config"`
}

type RouteConfig struct {
	Name string `json:"name" validate:"required"`
	RouteType RouteType `json:"route_type" validate:"required"`
	Model Model `json:"model" validate:"required"`
}

type RouteModelInfo struct {
	Name string `json:"name"`
	Provider string `json:"provider" validate:"required"`
}

type Route struct {
	Name string `json:"name" validate:"required"`
	RouteType string `json:"route_type" validate:"required"`
	Model RouteModelInfo `json:"model" validate:"required"`
	RouteURL string `json:"route_url" validate:"required"`
}

type Limit struct {
	Calls int `json:"calls" validate:"required"`
	Key string `json:"key"`
	RenewalPeriod string `json:"renewal_period" validate:"required"`
}

type GatewayConfig struct {
	Routes []RouteConfig `json:"routes" validate:"required"`
}

type LimitsConfig struct {
	Limits []Limit `json:"limits"`
}

func LoadRouteConfig(path string) (GatewayConfig, error) {
	var config GatewayConfig
	data, err := os.ReadFile(path) //TODO: double check implementation of os
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func SaveRouteConfig(config GatewayConfig, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644) //TODO: double check implementation of os
	if err != nil {
		return err
	}
	return nil
}

func ValidateConfig(configPath string) (GatewayConfig, error) {
	var config GatewayConfig
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, errors.New(fmt.Sprintf("%s does not exist", configPath))
	}
	config, err := LoadRouteConfig(configPath)
	if err != nil {
		return config, err
	}
	return config, nil
}
