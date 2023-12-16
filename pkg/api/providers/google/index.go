package google

import (
    "glide/pkg/api/providers"
)

func GoogleFullConfig() providers.ProviderConfigsAll {
    return providers.ProviderConfigsAll {
    Api: GoogleApiConfig,
    Chat: GoogleChatDefaultConfig,
    }
}