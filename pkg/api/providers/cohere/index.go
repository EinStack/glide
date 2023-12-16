package cohere

import (
    "glide/pkg/api/providers"
)

func CohereFullConfig() providers.ProviderConfigsAll {
    return providers.ProviderConfigsAll {
    Api: CohereApiConfig,
    Chat: CohereChatDefaultConfig,
    }
}

