package cmd

import (
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)
// Join command combines run and join.
// Used to join a k3s worker node to a server

const (
        BASE_VERSION = "0.10"
        BASE_IMAGE = "cliu2/k3sbase:"+BASE_VERSION
)

var (
        defaultPorts = []string{"6443"}
)

var JoinCommand = cli.Command{
	Name:  "join",
        Usage: "join a k3sbase container to a existing a server",
        ArgsUsage: `join <--detach> --server <SERVER-IP> --token <TOKEN> to <worker-container-id> <server-container-id`,
        Description: `The join command run a k3sbase container and join it to an existing k3snode server container`,
	Flags: []cli.Flag{
                &cli.StringFlag{
                        Name:  "server-ip",
                        Value: "",
                        Usage: `server container ip`,
                },
                &cli.BoolFlag{
                        Name:  "detach, d",
                        Usage: `run in detach mode or not`,
                },
                &cli.StringFlag{
                        Name:  "token, t",
                        Usage: `server token resides in /var/lib/rancher/k3s/server/node-token on server container`,
                },
                // &cli.StringFlag{
                //         Name:  "server-name",
                //         Usage: `server token resides in /var/lib/rancher/k3s/server/node-token on server container`,
                // },
        },
        Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return join(ctx, context.Args().First(),
			context.String("server-ip"),
                        context.String("token"),
                        context.Bool("detach"),
		)
        },
}

func join(ctx context.Context, containerID, serverIP, token string, detach bool) error {
        log.Debug("Begin join server node, first checking args")
        if serverIP == "" {
                log.Fatal("no server ip provided")
        }
        if token == "" {
                log.Fatal("no server token provided")
        }
        // First run a worker container
        log.Debug("run worker container")
        err := run(ctx, containerID, "worker", true, BASE_IMAGE, defaultPorts)
        if err != nil {
                log.Debug(err)
                return err
        }
        // Second join to server container
        err = utils.Join(containerID, serverIP, token, detach)
        if err != nil {
                log.Debug(err)
                return err
        }
        return nil
}

	

