package main

import (
	"glide/pkg/cmd"
	"log"
)

func main() {
	cli := cmd.NewCLI()

	if err := cli.Execute(); err != nil {
		log.Fatalf("glide run finished with error: %v", err)
	}
}
