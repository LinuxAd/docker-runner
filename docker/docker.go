package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewRunner(container Container) (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	c := Runner{
		Container: container,
		client:    cli,
	}
	return &c, err

}

func (r *Runner) createContainer(ctx context.Context) (container.ContainerCreateCreatedBody, error) {
	resp, err := r.client.ContainerCreate(ctx, &container.Config{
		Image: r.ImageName,
		Tty:   false,
	}, nil, nil, nil, r.ImageName+"_docker-runner-managed")

	return resp, err
}

func (r *Runner) Pull(ctx context.Context, w io.Writer) error {
	out, err := r.client.ImagePull(
		ctx, r.ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, err = io.Copy(w, out)
	return err
}

func (r *Runner) Run(ctx context.Context) error {
	resp, err := r.createContainer(ctx)
	if err != nil {
		return err
	}
	return r.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}
