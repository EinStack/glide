package main

import (
	"log"

	"glide/pkg/cmd"
)

//	@title			Glide Gateway
//	@version		1.0
//	@description	API documentation for Glide, an open-source lightweight high-performance model gateway

//	@contact.name	Glide Community
//	@contact.url	https://github.com/modelgateway/glide

//	@license.name	Apache 2.0
//	@license.url	https://github.com/modelgateway/glide/blob/develop/LICENSE

// @host		localhost:9099
// @BasePath	/
// @schemes	http
func main() {
	cli := cmd.NewCLI()

	if err := cli.Execute(); err != nil {
		log.Fatalf("Glide has finished with error: %v", err)
	}
}
