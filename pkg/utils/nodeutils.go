package utils

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"strings"
	log "github.com/sirupsen/logrus"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
)

// This file contains some node related functions like get server's ip and token

// GetServerToken get server token content
func GetServerToken(containerID string) (string, error) {
	log.Debug("read token out from k3s server files")
	// token place 
	token := filepath.Join(docker.K3sServerFile, containerID, "server", "token")
	bytes, err := ioutil.ReadFile(token)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	tokenStr := strings.Replace(string(bytes), "\n", "", -1)
	fmt.Print(tokenStr)
	return string(tokenStr), nil
}

// GetServerIP get server internal IP through docker inspect
func GetServerIP(containerID string) (string, error) {
	log.Debug("get server ip from docker inspect")
	ip, err := InspectContainerIP(containerID)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	// remove the unneccessary '
	ip = ip[1:len(ip)-2]
	server := "https://"+ip+":6443"
	return server, nil
}