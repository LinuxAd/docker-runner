package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type Container struct {
	ImageName string
	Command   string
}

type Runner struct {
	ctx    context.Context
	client *client.Client
}
