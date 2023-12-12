package pkg

import (
    "glide/pkg/providers"
)

type ProviderConfigs = pkg.ProviderConfigs

var OpenAIConfig = ProviderConfigs{
    "api":                 OpenAIAPIConfig,
    "chat":        OpenAiChatDefaultConfig,
}
