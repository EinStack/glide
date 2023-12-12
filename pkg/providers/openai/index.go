package openai

import (
    "glide/pkg/providers"
)


// TODO: this needs to be imported into buildAPIRequest.go
var OpenAIConfig = pkg.ProviderConfigs{
    "api":                 OpenAIAPIConfig,
    "chat":        OpenAiChatDefaultConfig,
}
