package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewRunner(ctx context.Context) (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	c := Runner{
		client: cli,
	}
	return &c, err

}

func (r *Runner) createContainer(ctx context.Context, cont Container) (container.ContainerCreateCreatedBody, error) {

	resp, err := r.client.ContainerCreate(ctx, &container.Config{
		Image: cont.ImageName,
		Tty:   false,
	}, nil, nil, nil, cont.ContainerName)

	return resp, err
}

func (r *Runner) Pull(ctx context.Context, cont Container) error {
	_, err := r.client.ImagePull(
		ctx, cont.ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return err
}

func (r *Runner) Run(ctx context.Context, cont Container) error {
	log.Printf("running container '%s'", cont.ImageName)
	resp, err := r.createContainer(ctx, cont)
	if err != nil {
		return err
	}
	r.Pull(ctx, cont)
	return r.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}

func (r *Runner) CheckRunning(ctx context.Context, cont Container) ([]types.ImageSummary, error) {
	return r.client.ImageList(ctx, types.ImageListOptions{})
}
