package dockerutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Run uses docker run to actually run a container
// expose hostport 6443 as default connect port
func (c *ContainerCmd) Run() error {
	args := []string{
		"run",
		"--privileged",
	}
	if c.Detach {
		args = append(args, "-d")
	}
	if err := checkDir(filepath.Join(k3sServerFile, c.ID)); err != nil {
		return fmt.Errorf("kubeserver path failed")
	}
	if err := checkDir(filepath.Join(kubeCfgFolder, c.ID)); err != nil {
		return fmt.Errorf("kubeconfig path failed")
	}
	args = append(args,
		"-e", "K3S_KUBECONFIG_OUTPUT="+filepath.Join(kubeCfgFolder, c.ID, "kubeconfig.yaml"),
		"-e", "K3S_KUBECONFIG_MODE=666",
		"-v", "/lib/modules:/lib/modules",
		"-v", filepath.Join(k3sServerFile, c.ID)+":/var/lib/rancher/k3s",
		"--name", c.ID)
	// args = append(args, c.Args...)
	for _, port := range c.Args {
		portStr := port+":"+port
		args = append(args, "-p", portStr)
	}
	args = append(args, c.Image)
	cmd := exec.Command(c.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debug(fmt.Sprintf("begin run container %s", c.ID))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
