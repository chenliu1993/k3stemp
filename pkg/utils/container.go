package utils

import (
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	log "github.com/sirupsen/logrus"
)

// RunContainer used for wrap exec run
func RunContainer(containerID string, env []string, image string) error {
	log.Debug("generating cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	ctrCmd.Env = env
	ctrCmd.Image = image
	return ctrCmd.Run()
}
