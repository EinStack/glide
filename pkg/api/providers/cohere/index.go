package cohere

import (
    "glide/pkg/api/providers"
)

var CohereConfig = providers.ProviderConfigs{
    "api":                 CohereApiConfig,
    "chat":        CohereChatDefaultConfig,
}

