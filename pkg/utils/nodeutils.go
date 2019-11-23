package utils

import (
	"fmt"
	"github.com/google/uuid"
	// "os"
	"os/exec"
	"bytes"
	"bufio"
	"path/filepath"
	"io/ioutil"
	"strings"
	log "github.com/sirupsen/logrus"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
)

// This file contains some node related functions like get server's ip and token


// func Init() {
// 	file, err := os.Open("/dev/urandom")
//     if err != nil {
//             panic(fmt.Sprintf("Failed to open urandom: %v", err))
//     }
//     uuid.SetRand(file)
// }
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


// ExecOutput save the output to strings
func ExecOutput(cmd exec.Cmd) (lines []string,err  error) {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err = cmd.Run()
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines,err
}

// GetClusterNames and returns it 
func GetClusterNames(clusterName string) (lines []string, err error) {
	// For now, only supports one server, so server name will be based on th cluster name
	cmd := exec.Command(
		"docker",
		"ps",
		"-q",         // quiet output for parsing
		"-a",         // show stopped nodes
		"--no-trunc", // don't truncate
		// filter for nodes with the cluster label
		"--filter", fmt.Sprintf("label=%s=%s", clusterconfig.ClusterLabelKey, clusterName),
		// format to include the cluster name
		"--format", `{{.Names}}`,
	)
	lines, err = ExecOutput(*cmd)
	if err != nil {
		return nil, err
	}
	// currentlt only supports one server
	if len(lines) != 1 {
		return nil, fmt.Errorf("k3scli don't support multiserver now...")
	}
	return lines, nil
}

// Generate container a unique container name 
func GenCtrName() string {
	return uuid.New().String()
}