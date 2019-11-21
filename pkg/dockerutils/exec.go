package dockerutils

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Exec uses docker eec to actually exec a process in a container
func (c *ContainerCmd) Exec() error {
	args := []string{
		"exec",
	}
	if c.Detach {
		args = append(args, "-d")
	}else{
		args = append(args, "-it")
	}

	args = append(args, c.ID)
	args = append(args, c.Args...)
	cmd := exec.Command(c.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debug(fmt.Sprintf("begin exec process in container: %s", c.ID))
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
