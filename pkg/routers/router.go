package routers

import "glide/pkg/telemetry"

type Router struct {
	telemetry *telemetry.Telemetry
}

func NewRouter(tel *telemetry.Telemetry) (*Router, error) {
	return &Router{
		telemetry: tel,
	}, nil
}
