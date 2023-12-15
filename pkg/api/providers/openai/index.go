package openai

import (
    "glide/pkg/api/providers"
)


// TODO: this needs to be imported into buildAPIRequest.go
var OpenAIConfig = providers.ProviderConfigs{
    "api":                 OpenAiApiConfig,
    "chat":        OpenAiChatDefaultConfig,
}
