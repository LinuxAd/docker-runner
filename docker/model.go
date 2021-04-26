package docker

import (
	"github.com/docker/docker/client"
)

type Container struct {
	ImageName string
	Command   string
}

type Runner struct {
	client *client.Client
}
