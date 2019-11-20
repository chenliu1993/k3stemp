package cmd

import (
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
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
		cli.StringFlag{
			Name:  "env, e",
			Value: "",
			Usage: `environment used for docker run --env`,
		},
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return run(ctx, context.Args().First(),
			context.String("label"),
			context.String("env"),
		)
	},
}

func run(ctx cli.Context, containerID string, env []string) error {
	log.Debug("begin running container")
	if env == "" {
		log.Fatal("env not set")
	}
	return utils.RunContainer(containerID, env)
}
