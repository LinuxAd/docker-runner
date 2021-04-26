package main

import (
	"log"

	"github.com/LinuxAd/docker-runner/runner"
)

func main() {
	log.Println("app initialising")
	a := runner.App{}
	a.Init()
	log.Println("app initialised")
	a.Run(":8080")
}
