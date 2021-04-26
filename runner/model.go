package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/LinuxAd/docker-runner/docker"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Name        string
	Container   docker.Container
	TargetCount int
	ActualCount int
}

type Response struct {
	Services []Service `json:"services,omitempty"`
	Error    `json:"error,omitempty"`
}

type Error struct {
	Msg  string `json:"msg"`
	Body string `json:"body"`
}

var (
	Running []Service
)

func writeResponse(res Response, w *http.ResponseWriter) {
	bytes, err := json.Marshal(res)
	if err != nil {
		logrus.Error(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(*w, "%s\n", bytes)
}

func containerFromConfig(cont docker.Container) (docker.Container, error) {
	dock := docker.Container{
		ImageName: cont.ImageName,
	}

	if dock.ImageName == "" {
		return dock, errors.New("ImageName cannot be blank")
	}

	if cont.Command != "" {
		dock.Command = cont.Command
	}

	return dock, nil
}

func (s *Service) newService(runner docker.Runner) error {
	ctx := context.Background()

	container, err := containerFromConfig(s.Container)
	if err != nil {
		return err
	}

	err = runner.Pull(ctx, os.Stdout, container)
	if err != nil {
		return err
	}

	return nil
}

func AddToServices(s Service) {
	Running = append(Running, s)
}

func CheckRunning() {
	for _, s := range Running {
		diff := s.ActualCount - s.TargetCount
		switch {
		case diff < 0:
			fmt.Println("not enough running!")
		case diff > 0:
			fmt.Println("too many running!")
		case diff == 0:
			fmt.Println("just the right amount")
		default:
			fmt.Println("looks like the maths is a bit off")
		}
	}
}
