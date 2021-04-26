package main

import (
	"fmt"

	"github.com/LinuxAd/docker-runner/docker"
	"github.com/LinuxAd/docker-runner/runner"
)

func main() {
	svc := runner.Service{
		Name: "test service",
		Container: docker.Container{
			ImageName: "nginx:latest",
		},
		TargetCount: 2,
	}
	fmt.Println(svc)

}
