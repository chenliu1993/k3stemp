package dockerutils

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Run uses docker run to actually run a container
func (c *ContainerCmd) Run() error {
	args := []string{
		"run",
		"--privileged",
	}
	for _, env := range c.Env {
		args = append(args, "-e", env)
	}
	if err := checkDir(k3sServerFiles); err != nil {
		return fmt.Errorf("kubeconfig path failed")
	}

	args = append(args, "-p", "6444:6443",
		"-v", "/lib/modules:/lib/modules",
		"-v", k3sServerFiles+":/var/lib/rancher/k3s",
		"--name", c.ID)
	args = append(args, c.Args...)
	args = append(args, c.Image)
	cmd := exec.Command(c.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debug("begin run container")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
