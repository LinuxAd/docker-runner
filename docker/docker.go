package docker

import "github.com/docker/docker/client"

type Container struct {
	ImageName string
}

func (c *Container) Run(client *client.Client) {

}
