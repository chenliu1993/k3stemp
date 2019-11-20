package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("exiting...")
		//os.Exit(1)
	}()

	app := cli.NewApp()
	app.Commands = runtimeCommands
	app.Before = beforeSubcommands
	// app.Metadata = map[string]interface{}{
	// 	"context": ctx,
	// }
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
