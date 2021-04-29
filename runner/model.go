package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/LinuxAd/docker-runner/docker"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Name        string           `json:"name"`
	ID          string           `json:"id"`
	Container   docker.Container `json:"container"`
	TargetCount int              `json:"target_count"`
	ActualCount int              `json:"actual_count"`
}

type Response struct {
	Services []*Service `json:"services,omitempty"`
	Error    `json:"error,omitempty"`
}

type Error struct {
	Msg  string `json:"msg"`
	Body string `json:"body"`
}

var (
	Running []*Service
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

func (s *Service) newService() error {
	s.ID = generateID()
	Running = append(Running, s)
	return nil
}

func generateID() string {
	return uuid.NewV4().String()
}

func CheckRunning() {
	for _, s := range Running {
		diff := s.ActualCount - s.TargetCount
		switch {
		case diff < 0:
			fmt.Printf("not enough running of %s, diff: %v\n", s.Name, diff)
			s.addContainers(math.Abs(float64(diff)))
		case diff > 0:
			fmt.Println("too many running!")
		case diff == 0:
			fmt.Println("just the right amount")
		default:
			fmt.Println("looks like the maths is a bit off")
		}
	}
}

func (s *Service) addContainers(count float64) {
	log.Printf("adding %v of %v", count, s.Container.ImageName)
	c := int(count)
	for i := 0; i == c; i++ {
		ctx := context.Background()
		r, err := docker.NewRunner(ctx)
		if err != nil {
			log.Println(err)
		}
		err = r.Run(ctx, s.Container)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
