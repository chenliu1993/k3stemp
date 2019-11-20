package dockerutils

import (
	"fmt"
	"os"
)

const (
	kubeCfgFile    = "/tmp/k3s/config"
	k3sServerFiles = "/tmp/k3s/files"
)

// ContainerCmd used for wrapping an docker command
type ContainerCmd struct {
	ID      string // the container name or ID
	Command string
	Args    []string
	Env     []string
	Image   string
}

func checkDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("remove existing k3s files")
		}
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("create k3s files failed")
		}
	}
	return nil
}
