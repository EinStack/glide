package main

import (
	"github.com/EinStack/glide/pkg/cmd"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

//	@title			Glide
//	@version		0.0.1
//	@description	API documentation for Glide, an open-source lightweight high-performance model gateway

//	@contact.name	EinStack Community
//	@contact.url	https://github.com/EinStack/glide/
//  @contact.email  contact@einstack.ai

//	@license.name	Apache 2.0
//	@license.url	https://github.com/EinStack/glide/blob/develop/LICENSE

// @externalDocs.description  Documentation
// @externalDocs.url          https://glide.einstack.ai/

// @host		localhost:9099
// @BasePath	/
// @schemes	http
func main() {
	cli := cmd.NewCLI()

	if err := cli.Execute(); err != nil {
		logger.Fatal("💥Glide has finished with error: %v", zap.Error(err))
	}
}
