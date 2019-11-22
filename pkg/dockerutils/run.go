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
	ctrFiles := filepath.Join(K3sServerFile, c.ID)
	if err := checkDir(ctrFiles); err != nil {
		return fmt.Errorf("kubeserver path failed")
	}
	ctrCfg := filepath.Join(KubeCfgFolder, c.ID)
	if err := checkDir(ctrCfg); err != nil {
		return fmt.Errorf("kubeconfig path failed")
	}
	for _, port := range c.Args {
		portStr := port+":"+port
		args = append(args, "-p", portStr)
	}
	args = append(args,
		"-e", "K3S_KUBECONFIG_OUTPUT="+filepath.Join(ctrCfg, "kubeconfig.yaml"),
		"-e", "K3S_KUBECONFIG_MODE=666",
		"-v", "/lib/modules:/lib/modules",
		"-v", ctrFiles+":/var/lib/rancher/k3s",
		"--name", c.ID)
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