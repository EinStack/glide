package openai

import (
    "glide/pkg/api/providers"
)

func OpenAiFullConfig() providers.ProviderConfigsAll {
    return providers.ProviderConfigsAll {
    Api: OpenAiApiConfig,
    Chat: OpenAiChatDefaultConfig,
    }
}