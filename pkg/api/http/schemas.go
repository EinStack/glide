package http

import "glide/pkg/routers"

type ErrorSchema struct {
	Message string `json:"message"`
}

type HealthSchema struct {
	Healthy bool `json:"healthy"`
}

type RouterListSchema struct {
	Routers []*routers.LangRouterConfig `json:"routers"`
}
