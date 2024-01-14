// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Glide Community",
            "url": "https://github.com/modelgateway/glide"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://github.com/modelgateway/glide/blob/develop/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/health/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Gateway Health",
                "operationId": "glide-health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.HealthSchema"
                        }
                    }
                }
            }
        },
        "/v1/language/": {
            "get": {
                "description": "Retrieve list of configured language routers and their configurations",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Language"
                ],
                "summary": "Language Router List",
                "operationId": "glide-language-routers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.RouterListSchema"
                        }
                    }
                }
            }
        },
        "/v1/language/{router}/chat": {
            "post": {
                "description": "Talk to different LLMs Chat API via unified endpoint",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Language"
                ],
                "summary": "Language Chat",
                "operationId": "glide-language-chat",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Router ID",
                        "name": "router",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Data",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.UnifiedChatRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.UnifiedChatResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorSchema"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "azureopenai.Config": {
            "type": "object",
            "required": [
                "apiVersion",
                "baseUrl",
                "model"
            ],
            "properties": {
                "apiVersion": {
                    "description": "The API version to use for this operation. This follows the YYYY-MM-DD format (e.g 2023-05-15)",
                    "type": "string"
                },
                "baseUrl": {
                    "description": "The name of your Azure OpenAI Resource (e.g https://glide-test.openai.azure.com/)",
                    "type": "string"
                },
                "chatEndpoint": {
                    "type": "string"
                },
                "defaultParams": {
                    "$ref": "#/definitions/azureopenai.Params"
                },
                "model": {
                    "description": "The name of your model deployment. You're required to first deploy a model before you can make calls (e.g. glide-gpt-35)",
                    "type": "string"
                }
            }
        },
        "azureopenai.Params": {
            "type": "object",
            "properties": {
                "frequency_penalty": {
                    "type": "integer"
                },
                "logit_bias": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "max_tokens": {
                    "type": "integer"
                },
                "n": {
                    "type": "integer"
                },
                "presence_penalty": {
                    "type": "integer"
                },
                "response_format": {
                    "description": "TODO: should this be a part of the chat request API?"
                },
                "seed": {
                    "type": "integer"
                },
                "stop": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "temperature": {
                    "type": "number"
                },
                "tool_choice": {},
                "tools": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "top_p": {
                    "type": "number"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "clients.ClientConfig": {
            "type": "object",
            "properties": {
                "timeout": {
                    "type": "integer"
                }
            }
        },
        "cohere.ChatHistory": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "cohere.Config": {
            "type": "object",
            "required": [
                "baseUrl",
                "chatEndpoint",
                "model"
            ],
            "properties": {
                "baseUrl": {
                    "type": "string"
                },
                "chatEndpoint": {
                    "type": "string"
                },
                "defaultParams": {
                    "$ref": "#/definitions/cohere.Params"
                },
                "model": {
                    "type": "string"
                }
            }
        },
        "cohere.Params": {
            "type": "object",
            "properties": {
                "chat_history": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/cohere.ChatHistory"
                    }
                },
                "citiation_quality": {
                    "type": "string"
                },
                "connectors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "conversation_id": {
                    "type": "string"
                },
                "preamble_override": {
                    "type": "string"
                },
                "prompt_truncation": {
                    "type": "string"
                },
                "search_queries_only": {
                    "type": "boolean"
                },
                "stream": {
                    "description": "unsupported right now",
                    "type": "boolean"
                },
                "temperature": {
                    "type": "number"
                }
            }
        },
        "http.ErrorSchema": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "http.HealthSchema": {
            "type": "object",
            "properties": {
                "healthy": {
                    "type": "boolean"
                }
            }
        },
        "http.RouterListSchema": {
            "type": "object",
            "properties": {
                "routers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/routers.LangRouterConfig"
                    }
                }
            }
        },
        "octoml.Config": {
            "type": "object",
            "required": [
                "baseUrl",
                "chatEndpoint",
                "model"
            ],
            "properties": {
                "baseUrl": {
                    "type": "string"
                },
                "chatEndpoint": {
                    "type": "string"
                },
                "defaultParams": {
                    "$ref": "#/definitions/octoml.Params"
                },
                "model": {
                    "type": "string"
                }
            }
        },
        "octoml.Params": {
            "type": "object",
            "properties": {
                "frequency_penalty": {
                    "type": "integer"
                },
                "max_tokens": {
                    "type": "integer"
                },
                "presence_penalty": {
                    "type": "integer"
                },
                "stop": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "temperature": {
                    "type": "number"
                },
                "top_p": {
                    "type": "number"
                }
            }
        },
        "openai.Config": {
            "type": "object",
            "required": [
                "baseUrl",
                "chatEndpoint",
                "model"
            ],
            "properties": {
                "baseUrl": {
                    "type": "string"
                },
                "chatEndpoint": {
                    "type": "string"
                },
                "defaultParams": {
                    "$ref": "#/definitions/openai.Params"
                },
                "model": {
                    "type": "string"
                }
            }
        },
        "openai.Params": {
            "type": "object",
            "properties": {
                "frequency_penalty": {
                    "type": "integer"
                },
                "logit_bias": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "max_tokens": {
                    "type": "integer"
                },
                "n": {
                    "type": "integer"
                },
                "presence_penalty": {
                    "type": "integer"
                },
                "response_format": {
                    "description": "TODO: should this be a part of the chat request API?"
                },
                "seed": {
                    "type": "integer"
                },
                "stop": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "temperature": {
                    "type": "number"
                },
                "tool_choice": {},
                "tools": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "top_p": {
                    "type": "number"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "providers.LangModelConfig": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "azureopenai": {
                    "$ref": "#/definitions/azureopenai.Config"
                },
                "client": {
                    "$ref": "#/definitions/clients.ClientConfig"
                },
                "cohere": {
                    "$ref": "#/definitions/cohere.Config"
                },
                "enabled": {
                    "description": "Is the model enabled?",
                    "type": "boolean"
                },
                "error_budget": {
                    "type": "string"
                },
                "id": {
                    "description": "Model instance ID (unique in scope of the router)",
                    "type": "string"
                },
                "octoml": {
                    "$ref": "#/definitions/octoml.Config"
                },
                "openai": {
                    "$ref": "#/definitions/openai.Config"
                }
            }
        },
        "retry.ExpRetryConfig": {
            "type": "object",
            "properties": {
                "base_multiplier": {
                    "type": "integer"
                },
                "max_delay": {
                    "type": "integer"
                },
                "max_retries": {
                    "type": "integer"
                },
                "min_delay": {
                    "type": "integer"
                }
            }
        },
        "routers.LangRouterConfig": {
            "type": "object",
            "required": [
                "models",
                "routers"
            ],
            "properties": {
                "enabled": {
                    "description": "Is router enabled?",
                    "type": "boolean"
                },
                "models": {
                    "description": "the list of models that could handle requests",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/providers.LangModelConfig"
                    }
                },
                "retry": {
                    "description": "retry when no healthy model is available to router",
                    "allOf": [
                        {
                            "$ref": "#/definitions/retry.ExpRetryConfig"
                        }
                    ]
                },
                "routers": {
                    "description": "Unique router ID",
                    "type": "string"
                },
                "strategy": {
                    "description": "strategy on picking the next model to serve the request",
                    "allOf": [
                        {
                            "$ref": "#/definitions/routing.Strategy"
                        }
                    ]
                }
            }
        },
        "routing.Strategy": {
            "type": "string",
            "enum": [
                "priority",
                "round-robin"
            ],
            "x-enum-varnames": [
                "Priority",
                "RoundRobin"
            ]
        },
        "schemas.ChatMessage": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "The content of the message.",
                    "type": "string"
                },
                "name": {
                    "description": "The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,\nwith a maximum length of 64 characters.",
                    "type": "string"
                },
                "role": {
                    "description": "The role of the author of this message. One of system, user, or assistant.",
                    "type": "string"
                }
            }
        },
        "schemas.ProviderResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/schemas.ChatMessage"
                },
                "responseId": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "tokenCount": {
                    "$ref": "#/definitions/schemas.TokenCount"
                }
            }
        },
        "schemas.TokenCount": {
            "type": "object",
            "properties": {
                "promptTokens": {
                    "type": "number"
                },
                "responseTokens": {
                    "type": "number"
                },
                "totalTokens": {
                    "type": "number"
                }
            }
        },
        "schemas.UnifiedChatRequest": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/schemas.ChatMessage"
                },
                "messageHistory": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemas.ChatMessage"
                    }
                }
            }
        },
        "schemas.UnifiedChatResponse": {
            "type": "object",
            "properties": {
                "cached": {
                    "type": "boolean"
                },
                "created": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "modelResponse": {
                    "$ref": "#/definitions/schemas.ProviderResponse"
                },
                "model_id": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "router": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9099",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Glide Gateway",
	Description:      "API documentation for Glide, an open-source lightweight high-performance model gateway",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
