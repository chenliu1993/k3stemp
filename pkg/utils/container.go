package utils

import (
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	log "github.com/sirupsen/logrus"
)

// RunContainer used for wrap exec run
func RunContainer(containerID string, detach bool, image string) error {
	log.Debug("generating cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	ctrCmd.Detach = detach
	ctrCmd.Image = image
	return ctrCmd.Run()
}
