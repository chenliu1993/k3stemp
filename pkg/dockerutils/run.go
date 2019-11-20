package dockerutils

import (
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
	args = append(args, c.Args...)
	args = append(args, c.ID)
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
