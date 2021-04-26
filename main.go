package main

import (
	"context"
	"log"

	"github.com/LinuxAd/docker-runner/docker"
)

func main() {

	c, err := docker.NewRunner(
		docker.Container{
			ImageName: "nginx:latest",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	err = c.Run(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}
