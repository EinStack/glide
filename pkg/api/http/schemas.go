package http

type ErrorSchema struct {
	Message string `json:"message"`
}

type HealthSchema struct {
	Healthy bool `json:"healthy"`
}
