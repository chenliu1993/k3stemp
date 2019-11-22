package utils

import (
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	log "github.com/sirupsen/logrus"
)

// RunServerContainer used for wrap exec run
func RunContainer(containerID string, detach bool, image string, ports []string) error {
	log.Debug("generating docker run cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	ctrCmd.Args = ports
	ctrCmd.Detach = detach
	ctrCmd.Image = image
	return ctrCmd.Run()
}

func Join(containerID, serverIP, token string, detach bool) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID: containerID,
		Command: "docker",
	}
	// Has to be true, because k3scli now it is not a input tty
	ctrCmd.Detach = detach
	// k3s agent --server https://myserver:6443 --token ${NODE_TOKEN}
	// since container IP uses a differenet network namespace, here join may leads to failed
	// needs to handle
	ctrCmd.Args = []string{
		"k3s", "agent",
		"--server", "https://"+serverIP+":6443",
		"--token", token,
	}
	return ctrCmd.Exec()
}
