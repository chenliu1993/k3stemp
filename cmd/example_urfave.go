package cmd

import (
	"context"
	"os"
	"time"

	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gitlab.eng.vmware.com/vkd/run-crx-cli/pkg/oci"
	"gitlab.eng.vmware.com/vkd/run-crx-cli/pkg/runcrxutils"
	"gitlab.eng.vmware.com/vkd/run-crx-cli/pkg/virtsdk"
	"gitlab.eng.vmware.com/vkd/run-crx-cli/pkg/virtsdk/specutils"
	"gitlab.eng.vmware.com/vkd/run-crx-cli/pkg/virtsdk/store"
)

var CreateCLICommand = cli.Command{
	Name:  "create",
	Usage: "Create a container",
	ArgsUsage: `<container-id>

   <container-id> is your name for the instance of the container that you
   are starting. The name you provide for the container instance must be unique
   on your host.`,
	Description: `The create command creates an instance of a container for a bundle. The
   bundle is a directory with a specification file named "` + "specConfig" + `" and a
   root filesystem.
   The specification file includes an args parameter. The args parameter is
   used to specify command(s) that get run when the container is started.
   To change the command(s) that get executed on start, edit the args
   parameter of the spec.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		// cli.StringFlag{
		// 	Name:  "console",
		// 	Value: "",
		// 	Usage: "path to a pseudo terminal",
		// },
		// cli.StringFlag{
		// 	Name:  "console-socket",
		// 	Value: "",
		// 	Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		// },
		// cli.StringFlag{
		// 	Name:  "pid-file",
		// 	Value: "",
		// 	Usage: "specify the file to write the process id to",
		// },
		// cli.BoolFlag{
		// 	Name:  "no-pivot",
		// 	Usage: "warning: this flag is meaningless to kata-runtime, just defined in order to be compatible with docker in ramdisk",
		// },
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return create(ctx, context.Args().First(),
			context.String("bundle"),
			context.String("console-socket"),
			context.String("pid-file"),
			true, //context.Bool("detach"),
			context.Bool("systemd-cgroup"),
		)
	},
}

// TODO
// --console-socket support
func create(ctx context.Context, containerID, bundlePath, consoleSocket, pidFilePath string, detach, systemdCgroup bool) error {
	log.Debug(time.Now().String())
	if bundlePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		log.Debug("Defaulting bundle path to current directory " + cwd)
		bundlePath = cwd
	}

	ociSpec, err := oci.ParseConfigJSON(bundlePath)
	if err != nil {
		return err
	}

	containerType, err := ociSpec.ContainerType()
	if err != nil {
		return err
	}

	if bundlePath, err = runcrxutils.ValidCreateParams(ctx, containerID, bundlePath); err != nil {
		return err
	}

	var sandbox virtsdk.VCSandbox

	var sandboxID = ""

	switch containerType {
	case virtsdk.PodSandbox:
		log.Debug("PodSandbox type")
		sandbox, err = runcrxutils.CreateSandbox(ctx, ociSpec, containerID)
		if err != nil {
			return err
		}
		sandboxID = sandbox.ID()
		fallthrough
	case virtsdk.PodContainer:
		log.Debug("PodContainer type")
		oci, err := specutils.LoadSpec(bundlePath)
		if err != nil {
			log.Fatal("load spec: ", err)
		}
		oci.Mounts = append(oci.Mounts, specs.Mount{
			Destination: "/etc/resolv.conf",
			Type:        "bind",
			Source:      "/etc/resolv.conf",
			Options:     []string{"rbind", "ro"},
		})
		newSandbox, _, err := runcrxutils.CreateContainer(ctx, ociSpec, sandboxID, containerID, oci, bundlePath, consoleSocket, true)
		if err != nil {
			log.Debug(err)
			errDel := sandbox.Delete()
			if errDel != nil {
				log.Debug(errDel)
			}
			return err
		}
		sandbox = newSandbox
	}
	err = store.AddContainerIDMapping(ctx, containerID, sandbox.ID())
	if err != nil {
		return err
	}

	log.Debug(time.Now().String())
	return nil
}
