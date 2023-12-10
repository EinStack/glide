package pkg

type ProviderConfigs map[string]interface{} // TODO: import from types.go

var OpenAIConfig = ProviderConfigs{
    "api":                 OpenAIAPIConfig,
    "chat":        OpenAiChatDefaultConfig,
}
