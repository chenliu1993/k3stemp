package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/chenliu1993/k3scli/cmd"
)


var runtimeCommands = []cli.Command{
	cmds.RunCLICommand,
}

func beforeSubcommands(c *cli.Context) error {
	loglevel := c.GlobalString("log-level")
	level, err := log.ParseLevel(loglevel)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}

fund main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("exiting...")
		//os.Exit(1)
	}()
	ctx := context.Background()
	cliApp := cli.NewApp()
	cliApp.Commands = runtimeCommands
	cliApp.Before = beforeSubcommands
	cliApp.Metadata = map[string]interface{}{
		"context": ctx,
	}
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
		},
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
