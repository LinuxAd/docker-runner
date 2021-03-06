package docker

import (
	"github.com/docker/docker/client"
)

type Container struct {
	ImageName     string `json:"image_name"`
	ContainerName string `json:"container_name"`
	Command       string `json:"command,omitempty"`
}

type Runner struct {
	client *client.Client
}
