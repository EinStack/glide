package anyscale

import (
    "glide/pkg/api/providers"
)

func CohereFullConfig() providers.ProviderConfigsAll {
    return providers.ProviderConfigsAll {
    Api: AnyscaleApiConfig,
    Chat: AnyscaleChatDefaultConfig,
    }
}

