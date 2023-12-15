package openai

import (
    "glide/pkg/api/providers"
)

var OpenAIConfig = providers.ProviderConfigs{
    "api":                 OpenAiApiConfig,
    "chat":        OpenAiChatDefaultConfig,
}

