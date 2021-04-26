package services

import "github.com/LinuxAd/docker-runner/docker"

type Service struct {
	Name        string
	Container   docker.Container
	TargetCount int
}

func (s *Service) newService()
