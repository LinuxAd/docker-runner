package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

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
	s.Container.ContainerName = s.Container.ContainerName + "-" + s.ID
	ctx := context.Background()

	r, err := docker.NewRunner(ctx)
	if err != nil {
		return err
	}

	for i := 0; i <= s.TargetCount; i++ {
		s.addContainer(ctx, *r)
	}

	Running = append(Running, s)

	return nil
}

func generateID() string {
	return uuid.NewV4().String()
}

func CheckRunning() {

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, s := range Running {
		wg.Add(1)
		go runCalculator(&wg, &mu, s, i)
	}
	wg.Wait()
}

func runCalculator(wg *sync.WaitGroup, mu *sync.Mutex, s *Service, i int) {
	defer wg.Done()

	ctx := context.Background()

	r, err := docker.NewRunner(ctx)
	if err != nil {
		log.Printf("error creating docker runner: %v", err)
	}

	a, err := r.CheckRunning(ctx, s.Container)
	if err != nil {
		log.Printf("error checking running containers for ImageName %v: %v", s.Container.ImageName, err)
	}
	s.ActualCount = len(a)
	if s.ActualCount == 0 {
		r.Pull(ctx, s.Container)
	}

	diff := s.ActualCount - s.TargetCount

	switch {
	case diff < 0:
		log.Printf("not enough running of %s, diff: %v\n", s.Name, diff)
		ctx := context.Background()
		s.addContainer(ctx, *r)
		mu.Lock()
		Running[i].ActualCount = s.ActualCount + 1
		mu.Unlock()
	case diff > 0:
		log.Println("too many running!")
		s.removeContainer(ctx, *r, a[0].ID)
	case diff == 0:
		fmt.Println("just the right amount")
	default:
		fmt.Println("looks like the maths is a bit off")
	}
}

func (s *Service) addContainer(ctx context.Context, runner docker.Runner) {
	c := s.Container
	c.ContainerName = fmt.Sprintf("%v-%v", s.Container.ContainerName, s.ActualCount+1)
	err := runner.Run(ctx, c)
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *Service) removeContainer(ctx context.Context, runner docker.Runner, id string) {
	if err := runner.Kill(ctx, id); err != nil {
		log.Println(err)
	}
}
