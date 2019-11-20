package cmd

import (
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// RunCommand wraps docker run for k3scli
var RunCommand = cli.Command{
	Name:  "run",
	Usage: "run a k3sbase/k3snode container",
	ArgsUsage: `<container-id> is your name for the instance of the container that you
	are starting. The name you provide for the container instance must be unique
	on your host.`,
	Description: `The run command allows you to start a new k3sbase/k3snode container`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "label, l",
			Value: "",
			Usage: `label used for docker run --label used for distinguishing from server to worker (primary)`,
		},
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: `environment used for docker run --env`,
		},
		cli.StringFlag{
			Name:  "image",
			Usage: `image used`,
		},
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return run(ctx, context.Args().First(),
			context.String("label"),
			context.StringSlice("env"),
			context.String("image"),
		)
	},
}

func run(ctx context.Context, containerID, label string, env []string, image string) error {
	log.Debug("begin running container")
	if env == nil {
		log.Fatal("env not set")
	}
	if image == "" {
		log.Fatal("k3s image not set")
	}
	return utils.RunContainer(containerID, env, image)
}
