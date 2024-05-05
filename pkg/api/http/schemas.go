package http

import "github.com/EinStack/glide/pkg/routers"

type HealthSchema struct {
	Healthy bool `json:"healthy"`
}

type RouterListSchema struct {
	Routers []*routers.LangRouterConfig `json:"routers"`
}
