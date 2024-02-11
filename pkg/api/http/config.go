package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"glide/pkg/version"
	"time"
)

type ServerConfig struct {
	Host               string         `yaml:"host"`
	Port               int            `yaml:"port"`
	ReadTimeout        *time.Duration `yaml:"read_timeout"`
	WriteTimeout       *time.Duration `yaml:"write_timeout"`
	IdleTimeout        *time.Duration `yaml:"idle_timeout"`
	MaxRequestBodySize *int           `yaml:"max_request_body_size"`
}

func DefaultServerConfig() *ServerConfig {
	maxReqBodySizeBytes := 4 * 1024 * 1024 // 4Mb
	readTimeout := 3 * time.Minute
	writeTimeout := 3 * time.Minute
	idleTimeout := 3 * time.Minute

	return &ServerConfig{
		Host:               "127.0.0.1",
		Port:               9099,
		IdleTimeout:        &idleTimeout,
		ReadTimeout:        &readTimeout,
		WriteTimeout:       &writeTimeout,
		MaxRequestBodySize: &maxReqBodySizeBytes,
	}
}

func (cfg *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

func (cfg *ServerConfig) ToServer() *fiber.App {
	// More configs are listed on https://docs.gofiber.io/api/fiber
	// TODO: Consider alternative JSON marshallers that provides better performance over the standard marshaller
	serverConfig := fiber.Config{
		AppName:                      "glide",
		DisableDefaultDate:           true,
		ServerHeader:                 fmt.Sprintf("glide/%v", version.Version),
		StreamRequestBody:            true,
		Immutable:                    false,
		DisablePreParseMultipartForm: true,
		EnablePrintRoutes:            false,
	}

	if cfg.IdleTimeout != nil {
		serverConfig.IdleTimeout = *cfg.IdleTimeout
	}

	if cfg.ReadTimeout != nil {
		serverConfig.ReadTimeout = *cfg.ReadTimeout
	}

	if cfg.WriteTimeout != nil {
		serverConfig.WriteTimeout = *cfg.WriteTimeout
	}

	if cfg.MaxRequestBodySize != nil {
		serverConfig.BodyLimit = *cfg.MaxRequestBodySize
	}

	return fiber.New(serverConfig)
}
